<script lang="ts">
    import {
        SvelteToast,
        toast,
        type SvelteToastOptions,
    } from "@zerodevx/svelte-toast";
    import { onMount } from "svelte";
    import { get, writable } from "svelte/store";
    import { v4 as uuidv4 } from "uuid";
    import * as App from "../wailsjs/go/main/App";
    import { main } from "../wailsjs/go/models";
    import FixtureConfiguration from "./FixtureConfiguration.svelte";
    import Info from "./Info.svelte";
    import SACNConfiguration from "./SACNConfiguration.svelte";
    import Config from "./Config.svelte";
    import type {
        CalibratingFixture,
        CalibrationPoint,
        CalibrationPoints,
        Fixture,
        Fixtures,
        MousePos,
        Point,
        SACNConfig,
    } from "./types";
    import {
        calcPan,
        calcTilt,
        convertCalibrationPointsToGo,
        convertFixturesToGo,
        convexHull,
    } from "./utils";

    let videoElement: HTMLVideoElement;
    let videoSelect: HTMLSelectElement;

    let videoStartX: number;
    let videoStartY: number;
    let videoRenderedWidth: number;
    let videoRenderedHeight: number;

    let deviceInfos = writable<MediaDeviceInfo[]>([]);
    let stream = writable<MediaStream>();

    let lockMousePos = false;
    let mousePos = writable<MousePos>({ x: 0, y: 0 });
    let mouseDragStart = writable<MousePos>(null);

    let fixtures = writable<Fixtures>({});
    let allFixturesCalibrated = writable<boolean>(true);
    let calibrationPointOutline = writable<Point[]>([]);
    let calibrationPointCounter = writable<number>(0);
    let calibrationPoints = writable<CalibrationPoints>({});

    let addingCalibrationPoint = false;
    let removingCalibrationPoint = false;
    let calibrateForOnePointSelectCalibrationPoint = false;

    let showMousePosition = false;
    let showCalibrationPoints = false;
    let showFixtureConfiguration = false;
    let showSACNConfiguration = false;
    let showSettingsMenu = false;
    let hideAllSettings = false;

    let fixturesToCalibrate = writable<string[]>([]);
    let calibrationPointsToCalibrate = writable<string[]>([]);
    let currentlyCalibrating = writable<CalibratingFixture | null>(null);

    let sacnConfig = writable<SACNConfig>(null);
    let sacnConfigDirty = false;

    fixtures.subscribe((fixtures) => {
        checkAllFixturesCalibrated(fixtures, get(calibrationPoints));

        let goFixtures: { [id: string]: main.Fixture } = convertFixturesToGo(
            fixtures,
            get(calibrationPoints),
        );
        App.SetFixtures(goFixtures);
    });

    calibrationPoints.subscribe((calibrationPoints) => {
        checkAllFixturesCalibrated(get(fixtures), calibrationPoints);
        calculateCalibrationPointOutline(calibrationPoints);

        let goCalibrationPoints: { [id: string]: main.CalibrationPoint } =
            convertCalibrationPointsToGo(calibrationPoints);
        App.SetCalibrationPoints(goCalibrationPoints);
    });

    onMount(() => {
        App.GetSACNConfig().then((sacnConfigFromApp) => {
            sacnConfig.set({
                ipAddress: sacnConfigFromApp.IpAddress,
                possibleIdAddresses: sacnConfigFromApp.PossibleIpAddresses,
                fps: sacnConfigFromApp.Fps,
                multicast: sacnConfigFromApp.Multicast,
                destinations: sacnConfigFromApp.Destinations,
            });
        });
    });

    onMount(() => {
        getStream().then(getDevices).then(gotDevices);
    });

    const toggleShowMousePosition = () => {
        showMousePosition = !showMousePosition;
    };

    const toggleShowCalibrationPoints = () => {
        showCalibrationPoints = !showCalibrationPoints;
    };

    const toggleShowFixtureConfiguration = () => {
        showFixtureConfiguration = !showFixtureConfiguration;
    };

    const toggleShowSACNConfiguration = () => {
        showSACNConfiguration = !showSACNConfiguration;
    };

    const toggleShowSettingsMenu = () => {
        showSettingsMenu = !showSettingsMenu;
    };

    function getNewCalibrationName() {
        calibrationPointCounter.update((value) => {
            return value + 1;
        });
        return `Point ${get(calibrationPointCounter)}`;
    }

    function checkAllFixturesCalibrated(
        fixtures: { [id: string]: Fixture },
        calibrationPoints: { [id: string]: CalibrationPoint },
    ) {
        for (let fixture of Object.values(fixtures)) {
            if (
                Object.keys(calibrationPoints).filter(
                    (calibration_point_id) =>
                        !Object.keys(fixture.calibration).includes(
                            calibration_point_id,
                        ),
                ).length !== 0
            ) {
                allFixturesCalibrated.set(false);
                return;
            }
        }

        allFixturesCalibrated.set(true);
    }

    function calculateCalibrationPointOutline(calibrationPoints: {
        [id: string]: CalibrationPoint;
    }) {
        calibrationPointOutline.set(
            convexHull(Object.values(calibrationPoints)),
        );
    }

    function showNotification(message: string, time_ms: number = 5000) {
        let options: SvelteToastOptions = {
            duration: time_ms,
        };
        toast.push(message, options);
    }

    function unlockMouse(event: MouseEvent) {
        lockMousePos = false;
        handleMouseMove(event);
        showNotification("Unlocked mouse");
    }

    function addCalibrationPoint() {
        hideAllSettings = true;
        showCalibrationPoints = true;
        addingCalibrationPoint = true;
    }

    function removeCalibrationPoint() {
        if (Object.keys(get(calibrationPoints)).length === 0) {
            App.AlertDialog("No calibration points", "Nothing to remove, as there are no calibration points.");
            return;
        }

        hideAllSettings = true;
        showCalibrationPoints = true;
        removingCalibrationPoint = true;
    }

    function calibrateFixtureForOnePoint(fixture_id: string) {
        if (Object.keys(get(calibrationPoints)).length === 0) {
            App.AlertDialog("No calibration points", "You need to add calibration points first.");
            return;
        }

        hideAllSettings = true;
        showCalibrationPoints = true;
        calibrateForOnePointSelectCalibrationPoint = true;

        currentlyCalibrating.set({
            fixture_id: fixture_id,
            calibration_point_id: null,
        });

        fixturesToCalibrate.set([]);
        calibrationPointsToCalibrate.set([]);
    }

    function calibrateFixtureForMissingPoints(
        fixture_id: string,
        calibration_points_missing: string[],
    ) {
        if (calibration_points_missing.length === 0) {
            App.AlertDialog("No calibration points missing", `The fixture '${get(fixtures)[fixture_id].name}' has all calibration points.`);
            return;
        }

        hideAllSettings = true;
        showCalibrationPoints = true;

        currentlyCalibrating.set({
            fixture_id: fixture_id,
            calibration_point_id: calibration_points_missing.pop(),
        });

        calibrationPointsToCalibrate.set(calibration_points_missing);
    }

    function calibrateFixtureForAllPoints(fixture_id: string) {
        if (Object.keys(get(calibrationPoints)).length === 0) {
            App.AlertDialog("No calibration points", "You need to add calibration points first.");
            return;
        }

        hideAllSettings = true;
        showCalibrationPoints = true;

        currentlyCalibrating.set({
            fixture_id: fixture_id,
            calibration_point_id: Object.keys(get(calibrationPoints))[0],
        });

        calibrationPointsToCalibrate.set(
            Object.keys(get(calibrationPoints)).filter(
                (calibration_point_id) =>
                    calibration_point_id !==
                    get(currentlyCalibrating).calibration_point_id,
            ),
        );
    }

    function handleKeyup(event: KeyboardEvent) {
        if (event.key === "Shift" || event.key === "Escape") {
            event.preventDefault();
        }

        if (event.repeat) {
            return;
        }

        if (event.key === "Escape") {
            if (addingCalibrationPoint) {
                showNotification("Cancelled adding calibration point");
                addingCalibrationPoint = false;
                hideAllSettings = false;
                fixturesToCalibrate.set([]);
                calibrationPointsToCalibrate.set([]);
                currentlyCalibrating.set(null);
            } else if (removingCalibrationPoint) {
                showNotification("Cancelled removing calibration point");
                removingCalibrationPoint = false;
                hideAllSettings = false;
            } else if (
                get(currentlyCalibrating) !== null &&
                !calibrateForOnePointSelectCalibrationPoint
            ) {
                showNotification("Skipping calibrating fixture");
                moveToNextFixtureOrCalibrationPointOrCancel();
            } else if (calibrateForOnePointSelectCalibrationPoint) {
                showNotification("Cancelled selecting point for calibration");
                calibrateForOnePointSelectCalibrationPoint = false;
                hideAllSettings = false;
                currentlyCalibrating.set(null);
                fixturesToCalibrate.set([]);
                calibrationPointsToCalibrate.set([]);
            } else if (showFixtureConfiguration) {
                showFixtureConfiguration = false;
            } else if (showSACNConfiguration) {
                showSACNConfiguration = false;
            } else if (showSettingsMenu) {
                showSettingsMenu = false;
            } else {
                showSettingsMenu = true;
            }
        } else if (event.key === "Shift") {
            mouseDragStart.set(null);
        }
    }

    function handleKeydown(event: KeyboardEvent) {
        if (event.key === "Shift" || event.key === "Escape") {
            event.preventDefault();
        }

        if (event.repeat) {
            return;
        }

        if (event.key === "Shift") {
            mouseDragStart.set(get(mousePos));
        }
    }

    function handleClickOnCalibrationPoint(event: MouseEvent, id: string) {
        if (lockMousePos) {
            unlockMouse(event);
            return;
        }

        if (removingCalibrationPoint) {
            calibrationPoints.update((calibrationPoints) => {
                delete calibrationPoints[id];
                return calibrationPoints;
            });

            fixtures.update((fixtures) => {
                for (let fixture_id in fixtures) {
                    delete fixtures[fixture_id].calibration[id];
                }

                return fixtures;
            });

            hideAllSettings = false;
            removingCalibrationPoint = false;
            event.stopPropagation();
        } else if (calibrateForOnePointSelectCalibrationPoint) {
            currentlyCalibrating.update((currentlyCalibrating) => {
                return {
                    fixture_id: currentlyCalibrating.fixture_id,
                    calibration_point_id: id,
                };
            });

            calibrateForOnePointSelectCalibrationPoint = false;
            event.stopPropagation();
        }

        handleMouseMove(event);
    }

    function handleClickOnVideo(event: MouseEvent) {
        if (lockMousePos) {
            unlockMouse(event);
            return;
        }

        if (addingCalibrationPoint) {
            let newId = uuidv4();

            calibrationPoints.update((calibrationPoints) => {
                calibrationPoints[newId] = {
                    id: newId,
                    name: getNewCalibrationName(),
                    x: get(mousePos).x,
                    y: get(mousePos).y,
                };
                return calibrationPoints;
            });

            addingCalibrationPoint = false;

            if (Object.keys(get(fixtures)).length === 0) {
                hideAllSettings = false;
                currentlyCalibrating.set(null);
                fixturesToCalibrate.set([]);
            } else {
                showCalibrationPoints = false;

                currentlyCalibrating.set({
                    fixture_id: Object.keys(get(fixtures))[0],
                    calibration_point_id: newId,
                });

                fixturesToCalibrate.set(
                    Object.keys(get(fixtures)).filter(
                        (fixture_id) =>
                            fixture_id !== get(currentlyCalibrating).fixture_id,
                    ),
                );

                calibrationPointsToCalibrate.set([]);
            }
        } else if (get(currentlyCalibrating) !== null) {
            if (calibrateForOnePointSelectCalibrationPoint) {
                showNotification(
                    "Click on calibration point to select it. ESC to cancel.",
                );
            } else {
                fixtures.update((fixtures) => {
                    let fixture =
                        fixtures[get(currentlyCalibrating).fixture_id];
                    let calibrationPoint =
                        get(calibrationPoints)[
                            get(currentlyCalibrating).calibration_point_id
                        ];

                    let pan = calcPan(
                        fixture,
                        get(mousePos),
                        get(mouseDragStart),
                    );
                    let tilt = calcTilt(
                        fixture,
                        get(mousePos),
                        get(mouseDragStart),
                    );

                    fixture.calibration[
                        get(currentlyCalibrating).calibration_point_id
                    ] = {
                        id: get(currentlyCalibrating).calibration_point_id,
                        pan: pan,
                        tilt: tilt,
                    };

                    fixtures[get(currentlyCalibrating).fixture_id] = fixture;

                    showNotification(
                        `Calibrated '${fixture.name}' at '${calibrationPoint.name}' (x: ${calibrationPoint.x.toFixed(4)}, y: ${calibrationPoint.y.toFixed(4)}) with pan: ${Math.floor(pan)}, tilt: ${Math.floor(tilt)}.`,
                        10000,
                    );

                    return fixtures;
                });

                moveToNextFixtureOrCalibrationPointOrCancel();
            }
        } else if (removingCalibrationPoint) {
            removingCalibrationPoint = false;
            hideAllSettings = false;
        } else {
            handleMouseMove(event);
            lockMousePos = true;
            showNotification("Locked mouse");
        }

        handleMouseMove(event);
    }

    function moveToNextFixtureOrCalibrationPointOrCancel() {
        if (
            get(fixturesToCalibrate).length === 0 &&
            get(calibrationPointsToCalibrate).length === 0
        ) {
            currentlyCalibrating.set(null);
            hideAllSettings = false;
            showCalibrationPoints = true;
        } else if (get(calibrationPointsToCalibrate).length !== 0) {
            currentlyCalibrating.update((currentlyCalibrating) => {
                return {
                    fixture_id: currentlyCalibrating.fixture_id,
                    calibration_point_id: get(calibrationPointsToCalibrate)[0],
                };
            });
            calibrationPointsToCalibrate.update(
                (calibrationPointsToCalibrate) => {
                    return calibrationPointsToCalibrate.filter(
                        (calibration_point_id) =>
                            calibration_point_id !==
                            get(currentlyCalibrating).calibration_point_id,
                    );
                },
            );
        } else if (get(fixturesToCalibrate).length !== 0) {
            currentlyCalibrating.update((currentlyCalibrating) => {
                return {
                    fixture_id: get(fixturesToCalibrate)[0],
                    calibration_point_id:
                        currentlyCalibrating.calibration_point_id,
                };
            });
            fixturesToCalibrate.update((fixturesToCalibrate) => {
                return fixturesToCalibrate.filter(
                    (fixture_id) =>
                        fixture_id !== get(currentlyCalibrating).fixture_id,
                );
            });
        }
    }

    function calculateVideoSize() {
        const videoWidth = videoElement.videoWidth;
        const videoHeight = videoElement.videoHeight;
        const videoAspectRatio = videoWidth / videoHeight;

        const elementWidth = videoElement.clientWidth;
        const elementHeight = videoElement.clientHeight;
        const elementAspectRatio = elementWidth / elementHeight;

        let offsetX, offsetY;

        if (elementAspectRatio > videoAspectRatio) {
            // Element is wider than video, scale by height
            videoRenderedHeight = elementHeight;
            videoRenderedWidth = videoRenderedHeight * videoAspectRatio;
            offsetX = (elementWidth - videoRenderedWidth) / 2; // Center horizontally
            offsetY = 0;
        } else {
            // Element is taller than video, scale by width
            videoRenderedWidth = elementWidth;
            videoRenderedHeight = videoRenderedWidth / videoAspectRatio;
            offsetX = 0;
            offsetY = (elementHeight - videoRenderedHeight) / 2; // Center vertically
        }

        const rect = videoElement.getBoundingClientRect();

        videoStartX = rect.left + offsetX;
        videoStartY = rect.top + offsetY;
    }

    function handleMouseMove(event: MouseEvent) {
        if (lockMousePos) {
            return;
        }

        calculateVideoSize();

        const x = Math.max(
            Math.min((event.clientX - videoStartX) / videoRenderedWidth, 1),
            0,
        );
        const y = Math.max(
            Math.min((event.clientY - videoStartY) / videoRenderedHeight, 1),
            0,
        );

        mousePos.set({ x, y });

        if (
            get(currentlyCalibrating) !== null &&
            !calibrateForOnePointSelectCalibrationPoint
        ) {
            let fixture = get(fixtures)[get(currentlyCalibrating).fixture_id];
            let pan = calcPan(fixture, get(mousePos), get(mouseDragStart));
            let tilt = calcTilt(fixture, get(mousePos), get(mouseDragStart));
            App.SetPanTiltForFixture(
                get(currentlyCalibrating).fixture_id,
                Math.floor(pan),
                Math.floor(tilt),
            );
        } else {
            App.SetMouseForAllFixtures(get(mousePos).x, get(mousePos).y);
        }
    }

    function getDevices() {
        return navigator.mediaDevices.enumerateDevices();
    }

    function gotDevices(p_deviceInfos) {
        deviceInfos.set(p_deviceInfos);
    }

    function getStream() {
        if (get(stream)) {
            get(stream)
                .getTracks()
                .forEach((track) => {
                    track.stop();
                });
        }

        const videoSource = videoSelect.value;
        const constraints = {
            video: {
                deviceId: videoSource ? { exact: videoSource } : undefined,
                width: { ideal: 1920 },
                height: { ideal: 1080 },
            },
        };

        return navigator.mediaDevices
            .getUserMedia(constraints)
            .then(gotStream)
            .catch(handleError);
    }

    function gotStream(p_stream) {
        stream.set(p_stream);
        videoSelect.selectedIndex = [...videoSelect.options].findIndex(
            (option) => option.text === p_stream.getVideoTracks()[0].label,
        );
        videoElement.srcObject = p_stream;

        setTimeout(() => {
            calculateVideoSize();
        }, 100);
    }

    function handleError(error) {
        console.error("Error: ", error);
    }
</script>

<svelte:window on:keyup={handleKeyup} on:keydown={handleKeydown} />

<SvelteToast />

<main>
    <div class="content" on:mousemove={handleMouseMove}>
        <!-- svelte-ignore a11y-media-has-caption -->
        <video autoplay bind:this={videoElement} />
        <!-- svelte-ignore a11y-click-events-have-key-events -->
        <div
            class="video-overlay"
            style="top: {videoStartY}px; left: {videoStartX}px; width: {videoRenderedWidth}px; height: {videoRenderedHeight}px;"
            on:click={handleClickOnVideo}
        >
            <svg class="video-cover-svg">
                {#each $calibrationPointOutline as point, index}
                    <line
                        class="outline-line"
                        x1="{point.x * 100}%"
                        y1="{point.y * 100}%"
                        x2="{$calibrationPointOutline[
                            (index + 1) % $calibrationPointOutline.length
                        ].x * 100}%"
                        y2="{$calibrationPointOutline[
                            (index + 1) % $calibrationPointOutline.length
                        ].y * 100}%"
                    ></line>
                {/each}

                {#if $mouseDragStart !== null}
                    <line
                        class="mouse-drag-line"
                        x1="{$mouseDragStart.x * 100}%"
                        y1="{$mouseDragStart.y * 100}%"
                        x2="{$mousePos.x * 100}%"
                        y2="{$mousePos.y * 100}%"
                    ></line>
                {/if}
            </svg>
            {#if showCalibrationPoints}
                {#each Object.values($calibrationPoints) as calibrationPoint (calibrationPoint.id)}
                    <div
                        class="calibration-point {$currentlyCalibrating &&
                        $currentlyCalibrating.calibration_point_id ===
                            calibrationPoint.id
                            ? 'active-calibration-point'
                            : ''}"
                        style="
                            top: {calibrationPoint.y * 100}%;
                            left: {calibrationPoint.x * 100}%;
                        "
                        on:click={(event) => {
                            handleClickOnCalibrationPoint(
                                event,
                                calibrationPoint.id,
                            );
                        }}
                    ></div>
                {/each}
            {/if}
            {#if lockMousePos}
                <div
                    class="lock-mouse-pos-div"
                    style="
                            top: {$mousePos.y * 100}%;
                            left: {$mousePos.x * 100}%;
                        "
                ></div>
            {/if}
            {#if $currentlyCalibrating !== null && !calibrateForOnePointSelectCalibrationPoint}
                <div
                    class="calibration-point active-calibration-point"
                    style="
                            top: {$calibrationPoints[
                        $currentlyCalibrating.calibration_point_id
                    ].y * 100}%;
                            left: {$calibrationPoints[
                        $currentlyCalibrating.calibration_point_id
                    ].x * 100}%;
                            pointer-events: none;
                        "
                ></div>
            {/if}
        </div>
    </div>
    <button
        class="settings-button {hideAllSettings ? 'hidden' : ''}"
        on:click={toggleShowSettingsMenu}
    >
        <svg
            xmlns="http://www.w3.org/2000/svg"
            x="0px"
            y="0px"
            width="40px"
            height="40px"
            viewBox="0 0 50 50"
            fill="var(--main-text-color)"
        >
            <path
                d="M47.16,21.221l-5.91-0.966c-0.346-1.186-0.819-2.326-1.411-3.405l3.45-4.917c0.279-0.397,0.231-0.938-0.112-1.282 l-3.889-3.887c-0.347-0.346-0.893-0.391-1.291-0.104l-4.843,3.481c-1.089-0.602-2.239-1.08-3.432-1.427l-1.031-5.886 C28.607,2.35,28.192,2,27.706,2h-5.5c-0.49,0-0.908,0.355-0.987,0.839l-0.956,5.854c-1.2,0.345-2.352,0.818-3.437,1.412l-4.83-3.45 c-0.399-0.285-0.942-0.239-1.289,0.106L6.82,10.648c-0.343,0.343-0.391,0.883-0.112,1.28l3.399,4.863 c-0.605,1.095-1.087,2.254-1.438,3.46l-5.831,0.971c-0.482,0.08-0.836,0.498-0.836,0.986v5.5c0,0.485,0.348,0.9,0.825,0.985 l5.831,1.034c0.349,1.203,0.831,2.362,1.438,3.46l-3.441,4.813c-0.284,0.397-0.239,0.942,0.106,1.289l3.888,3.891 c0.343,0.343,0.884,0.391,1.281,0.112l4.87-3.411c1.093,0.601,2.248,1.078,3.445,1.424l0.976,5.861C21.3,47.647,21.717,48,22.206,48 h5.5c0.485,0,0.9-0.348,0.984-0.825l1.045-5.89c1.199-0.353,2.348-0.833,3.43-1.435l4.905,3.441 c0.398,0.281,0.938,0.232,1.282-0.111l3.888-3.891c0.346-0.347,0.391-0.894,0.104-1.292l-3.498-4.857 c0.593-1.08,1.064-2.222,1.407-3.408l5.918-1.039c0.479-0.084,0.827-0.5,0.827-0.985v-5.5C47.999,21.718,47.644,21.3,47.16,21.221z M25,32c-3.866,0-7-3.134-7-7c0-3.866,3.134-7,7-7s7,3.134,7,7C32,28.866,28.866,32,25,32z"
            ></path>
        </svg>
    </button>
    <div
        class="settings {!showSettingsMenu || hideAllSettings ? 'hidden' : ''}"
    >
    <Config bind:fixtures bind:calibrationPoints bind:sacnConfig></Config>
        <select bind:this={videoSelect} on:change={getStream}>
            {#each $deviceInfos as deviceInfo, index}
                {#if deviceInfo.kind === "videoinput"}
                    <option value={deviceInfo.deviceId}
                        >{deviceInfo.label || `Camera ${index + 1}`}</option
                    >
                {/if}
            {/each}
        </select>
        <button on:click={toggleShowFixtureConfiguration}>
            Fixture Config
        </button>
        {#if showFixtureConfiguration}
            <!-- svelte-ignore a11y-click-events-have-key-events -->
            <div class="overlay" on:click={toggleShowFixtureConfiguration}>
                <!-- svelte-ignore a11y-click-events-have-key-events -->
                <div on:click|stopPropagation>
                    <FixtureConfiguration
                        bind:fixtures
                        bind:calibrationPoints
                        on:calibrate_all_points={(event) => {
                            calibrateFixtureForAllPoints(
                                event.detail.fixture_id,
                            );
                        }}
                        on:calibrate_one_point={(event) => {
                            calibrateFixtureForOnePoint(
                                event.detail.fixture_id,
                            );
                        }}
                        on:calibrate_missing_points={(event) => {
                            calibrateFixtureForMissingPoints(
                                event.detail.fixture_id,
                                event.detail.calibration_points_missing,
                            );
                        }}
                    />
                </div>
            </div>
        {/if}
        <button on:click={toggleShowSACNConfiguration}> sACN Config </button>
        {#if showSACNConfiguration}
            <!-- svelte-ignore a11y-click-events-have-key-events -->
            <div
                class="overlay"
                on:click={() => {
                    if (sacnConfigDirty) {
                        App.AlertDialog(
                            "Unsaved changes",
                            "You have unsaved changes. Either apply or cancel them before closing this window.",
                        );
                    } else {
                        toggleShowSACNConfiguration();
                    }
                }}
            >
                <div on:click|stopPropagation>
                    <SACNConfiguration bind:sacnConfig bind:sacnConfigDirty />
                </div>
            </div>
        {/if}
        <button on:click={addCalibrationPoint}> Add Calibration Point </button>
        <button on:click={removeCalibrationPoint}>
            Remove Calibration Point
        </button>
        <div class="settings-separator"></div>
        <button on:click={toggleShowMousePosition}
            >{showMousePosition
                ? "Hide Mouse Position"
                : "Show Mouse Position"}</button
        >
        <button on:click={toggleShowCalibrationPoints}>
            {showCalibrationPoints
                ? "Hide Calibration Points"
                : "Show Calibration Points"}
        </button>
    </div>
    <Info
        bind:addingCalibrationPoint
        bind:allFixturesCalibrated
        bind:calibrateForOnePointSelectCalibrationPoint
        bind:calibrationPoints
        bind:calibrationPointsToCalibrate
        bind:currentlyCalibrating
        bind:fixtures
        bind:fixturesToCalibrate
        bind:lockMousePos
        bind:mouseDragStart
        bind:mousePos
        bind:removingCalibrationPoint
        bind:showCalibrationPoints
        bind:showMousePosition
    />
</main>

<style>
    main {
        height: 100vh;
    }

    video {
        width: 100%;
        height: 100%;
        background-color: var(--main-bg-color);
        object-fit: contain;
    }

    .content {
        height: 100%;
        width: 100%;
    }

    .settings {
        position: fixed;
        display: flex;
        flex-direction: column;
        align-items: flex-start;
        gap: 10px;
        top: 60px;
        left: 10px;
        background-color: var(--main-bg-color-transparent);
        padding: 10px;
        border-radius: 6px;
        z-index: 3;
    }

    .settings-button {
        position: fixed;
        top: 10px;
        left: 10px;
        width: 40px;
        height: 40px;
        background-color: rgba(0, 0, 0, 0);
        padding: 0;
    }

    .overlay {
        position: fixed;
        top: 0;
        left: 0;
        width: 100vw;
        height: 100vh;
        display: flex;
        justify-content: center;
        align-items: center;
        background: var(--overlay);
        z-index: 1000;
    }

    .video-overlay {
        position: absolute;
        z-index: 2;
    }

    .calibration-point {
        position: absolute;
        width: 15px;
        height: 15px;
        background-color: var(--main-red-color-transparent);
        border-radius: 50%;
        transform: translate(-50%, -50%);
        pointer-events: auto;
    }

    .active-calibration-point {
        background-color: var(--main-green-color-transparent);
    }

    .video-cover-svg {
        width: 100%;
        height: 100%;
        pointer-events: none;
    }

    .video-cover-svg > line.outline-line {
        stroke: var(--main-green-color-transparent);
        stroke-width: 4px;
        stroke-linecap: round;
    }

    .video-cover-svg > line.mouse-drag-line {
        stroke: var(--main-red-color);
        stroke-width: 8px;
        stroke-linecap: round;
    }

    .lock-mouse-pos-div {
        pointer-events: none;
        width: 12px;
        height: 12px;
        background-color: red;
        border-radius: 50%;
        position: absolute;
        transform: translate(-50%, -50%);
    }

    .lock-mouse-pos-div::after {
        content: "";
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        width: 32px; /* Adjust to make the outer circle larger */
        height: 32px;
        border-style: solid;
        border-color: red;
        border-radius: 50%;
        box-sizing: border-box;
    }

    .settings-separator {
        margin-top: 12px;
    }
</style>
