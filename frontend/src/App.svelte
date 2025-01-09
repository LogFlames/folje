<script lang="ts">
    import { onMount } from "svelte";
    import { writable, get } from "svelte/store";
    import Fixtures from "./Fixtures.svelte";
    import type { Fixture, CalibrationPoint } from "./types";
    import { v4 as uuidv4 } from "uuid";
    import { calcPan, calcTilt } from "./utils";

    interface MousePos {
        x: number;
        y: number;
    }

    interface CalibratingFixture {
        fixture_id: string;
        calibration_point_id: string;
    }

    let videoElement;
    let videoSelect;

    let videoStartX;
    let videoStartY;
    let videoRenderedWidth;
    let videoRenderedHeight;

    let deviceInfos = writable<MediaDeviceInfo[]>([]);
    let stream = writable<MediaStream>();

    let mousePos = writable<MousePos>({ x: 0, y: 0 });

    let showFixtures = false;
    let fixtures = writable<{ [id: string]: Fixture }>({});

    let calibrationPointCounter = writable<number>(0);

    function getNewCalibrationName() {
        calibrationPointCounter.update((value) => {
            return value + 1;
        });
        return `Point ${get(calibrationPointCounter)}`;
    }

    let calibrationPoints = writable<{ [id: string]: CalibrationPoint }>({
        "1aoeu": { id: "1aoeu", name: getNewCalibrationName(), x: 0.5, y: 0.5 },
        "2aoeu": { id: "2aoeu", name: getNewCalibrationName(), x: 0.1, y: 0.5 },
        "3aoeu": { id: "3aoeu", name: getNewCalibrationName(), x: 0.9, y: 0.5 },
        "4aoeu": { id: "4aoeu", name: getNewCalibrationName(), x: 0.5, y: 0.1 },
        "5aoeu": { id: "5aoeu", name: getNewCalibrationName(), x: 0.5, y: 0.9 },
    });

    let showSettingsMenu = false;
    let hideAllSettings = false;

    let showMousePosition = false;
    let showCalibrationPoints = false;
    let addingCalibrationPoint = false;
    let removingCalibrationPoint = false;

    let calibrateForOnePointSelectCalibrationPoint = false;

    let fixturesToCalibrate = writable<string[]>([]);
    let calibrationPointsToCalibrate = writable<string[]>([]);

    let currentlyCalibrating = writable<CalibratingFixture | null>(null);

    const toggleShowMousePosition = () => {
        showMousePosition = !showMousePosition;
    };

    const toggleShowCalibrationPoints = () => {
        showCalibrationPoints = !showCalibrationPoints;
    };

    const toggleShowFixtures = () => {
        showFixtures = !showFixtures;
    };

    const toggleShowSettingsMenu = () => {
        showSettingsMenu = !showSettingsMenu;
    };

    onMount(() => {
        getStream().then(getDevices).then(gotDevices);
    });

    onMount(() => {
        window.addEventListener("keyup", (event) => {
            if (event.preventDefault) {
                return;
            }

            handleKeyup(event);
        });
    });

    function showNotification(message: string) {}

    function addCalibrationPoint() {
        hideAllSettings = true;
        showCalibrationPoints = true;
        addingCalibrationPoint = true;
    }

    function removeCalibrationPoint() {
        hideAllSettings = true;
        showCalibrationPoints = true;
        removingCalibrationPoint = true;
    }

    function handleKeyup(event: KeyboardEvent) {
        if (event.repeat) {
            return;
        }

        if (event.key === "Escape") {
            if (addingCalibrationPoint) {
                showNotification("Canceling adding calibration point");
                addingCalibrationPoint = false;
                hideAllSettings = false;
                fixturesToCalibrate.set([]);
                calibrationPointsToCalibrate.set([]);
                currentlyCalibrating.set(null);
            } else if (removingCalibrationPoint) {
                removingCalibrationPoint = false;
                hideAllSettings = false;
            } else if (
                get(currentlyCalibrating) !== null &&
                !calibrateForOnePointSelectCalibrationPoint
            ) {
                moveToNextFixtureOrCalibrationPointOrCancel();
            } else if (calibrateForOnePointSelectCalibrationPoint) {
                calibrateForOnePointSelectCalibrationPoint = false;
                hideAllSettings = false;
                currentlyCalibrating.set(null);
                fixturesToCalibrate.set([]);
                calibrationPointsToCalibrate.set([]);
            } else if (showFixtures) {
                showFixtures = false;
            } else if (showSettingsMenu) {
                showSettingsMenu = false;
            } else {
                showSettingsMenu = true;
            }
        }

        event.preventDefault();
    }

    function handleKeydown(event: KeyboardEvent) {
        if (event.repeat) {
            return;
        }

        if (event.key == "Space") {
        }
    }

    function handleClickOnCalibrationPoint(event: MouseEvent, id: string) {
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

    function calibrateFixtureForOnePoint(fixture_id: string) {
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

    function calibrateFixtureForAllPoints(fixture_id: string) {
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

    function handleClickOnVideo() {
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
        } else if (
            get(currentlyCalibrating) !== null &&
            !calibrateForOnePointSelectCalibrationPoint
        ) {
            fixtures.update((fixtures) => {
                let fixture = fixtures[get(currentlyCalibrating).fixture_id];
                let pan = calcPan(fixture, get(mousePos).x);
                let tilt = calcTilt(fixture, get(mousePos).y);

                fixture.calibration[
                    get(currentlyCalibrating).calibration_point_id
                ] = {
                    id: get(currentlyCalibrating).calibration_point_id,
                    pan: pan,
                    tilt: tilt,
                };

                fixtures[get(currentlyCalibrating).fixture_id] = fixture;

                return fixtures;
            });

            moveToNextFixtureOrCalibrationPointOrCancel();
        } else if (removingCalibrationPoint) {
            removingCalibrationPoint = false;
            hideAllSettings = false;
        }
    }

    function getDevices() {
        return navigator.mediaDevices.enumerateDevices();
    }

    function gotDevices(p_deviceInfos) {
        console.log(p_deviceInfos);
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
    }

    function handleError(error) {
        console.error("Error: ", error);
    }

    function handleMouseMove(event) {
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

        const x = Math.max(
            Math.min((event.clientX - videoStartX) / videoRenderedWidth, 1),
            0,
        );
        const y = Math.max(
            Math.min((event.clientY - videoStartY) / videoRenderedHeight, 1),
            0,
        );

        mousePos.set({ x, y });
    }
</script>

<svelte:window on:keyup={handleKeyup} on:keydown={handleKeydown} />

<main>
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div
        class="content"
        on:mousemove={handleMouseMove}
        on:click={handleClickOnVideo}
    >
        <!-- svelte-ignore a11y-media-has-caption -->
        <video autoplay bind:this={videoElement} />
        <div
            class="video-overlay"
            style="top: {videoStartY}px; left: {videoStartX}px; width: {videoRenderedWidth}px; height: {videoRenderedHeight}px;"
        >
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
        <select bind:this={videoSelect} on:change={getStream}>
            {#each $deviceInfos as deviceInfo, index (deviceInfo.deviceId)}
                {#if deviceInfo.kind === "videoinput"}
                    <option value={deviceInfo.deviceId}
                        >{deviceInfo.label || `Camera ${index + 1}`}</option
                    >
                {/if}
            {/each}
        </select>
        <button on:click={toggleShowFixtures}> Fixtures </button>
        {#if showFixtures}
            <!-- svelte-ignore a11y-click-events-have-key-events -->
            <div class="overlay" on:click={toggleShowFixtures}>
                <!-- svelte-ignore a11y-click-events-have-key-events -->
                <div on:click|stopPropagation>
                    <Fixtures
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
                    />
                </div>
            </div>
        {/if}
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
        <button on:click={addCalibrationPoint}>Add Calibration Point</button>
        <button on:click={removeCalibrationPoint}
            >Remove Calibration Point</button
        >
    </div>
    <div class="info-box">
        {#if showMousePosition}
            <div>
                x: {$mousePos.x}
            </div>
            <div>
                y: {$mousePos.y}
            </div>
        {/if}
        {#if showCalibrationPoints}
            <div>Showing Calibration Points</div>
        {/if}
        {#if addingCalibrationPoint}
            <div>Adding calibration point</div>
        {/if}
        {#if removingCalibrationPoint}
            <div>Removing calibration point</div>
        {/if}
        {#if $currentlyCalibrating !== null && !calibrateForOnePointSelectCalibrationPoint}
            <div>
                Calibrating {$fixtures[$currentlyCalibrating.fixture_id].name} on
                point {$calibrationPoints[
                    $currentlyCalibrating.calibration_point_id
                ].name}
            </div>
            <div>
                Pan: {Math.floor(
                    calcPan(
                        $fixtures[$currentlyCalibrating.fixture_id],
                        $mousePos.x,
                    ),
                )}
            </div>
            <div>
                Tilt: {Math.floor(
                    calcTilt(
                        $fixtures[$currentlyCalibrating.fixture_id],
                        $mousePos.y,
                    ),
                )}
            </div>
        {/if}
        {#if calibrateForOnePointSelectCalibrationPoint}
            <div>
                Calibrating {$fixtures[$currentlyCalibrating.fixture_id].name},
                select calibration point.
            </div>
        {/if}
        {#if $fixturesToCalibrate.length > 0}
            <div>
                Fixtures to calibrate:
                <br />
                <span class="small">
                    {#each $fixturesToCalibrate as fixture}
                        {$fixtures[fixture].name},
                    {/each}
                </span>
            </div>
        {/if}
        {#if $calibrationPointsToCalibrate.length > 0}
            <div>
                Calibration points to calibrate for:
                <br />
                <span class="small">
                    {#each $calibrationPointsToCalibrate as calibrationPoint}
                        {$calibrationPoints[calibrationPoint].name},
                    {/each}
                </span>
            </div>
        {/if}
    </div>
    {#if addingCalibrationPoint || removingCalibrationPoint || ($currentlyCalibrating !== null && !calibrateForOnePointSelectCalibrationPoint) || calibrateForOnePointSelectCalibrationPoint}
        <div class="tooltip">
            {#if addingCalibrationPoint}
                <div>
                    Click on the camera-feed to create a calibration point
                    there.
                    <br />
                    Press ESC to cancel.
                </div>
            {/if}
            {#if removingCalibrationPoint}
                <div>
                    Click on a calibration point to remove it and all
                    calibrations to that point on the fixtures.
                    <br />
                    Press ESC to cancel.
                </div>
            {/if}
            {#if $currentlyCalibrating !== null && !calibrateForOnePointSelectCalibrationPoint}
                <div>
                    Calibrating a fixture for the green calibration point.
                    <br />
                    Press ESC to not calibrate this fixture. (NOTE if you have more
                    fixtures to be calibrated you will move on to the next one)
                </div>
            {/if}
            {#if calibrateForOnePointSelectCalibrationPoint}
                <div>
                    Select the point to calibrate the fixture for by clicking on
                    it.
                    <br />
                    Press ESC to cancel.
                </div>
            {/if}
        </div>
    {/if}
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
        border-radius: 8px;
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
        pointer-events: none;
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

    .info-box {
        position: absolute;
        pointer-events: none;
        top: 0;
        right: 0;
        background-color: var(--overlay);
        padding: 10px;
        border-radius: 10px;
        text-align: right;
    }

    .tooltip {
        position: absolute;
        pointer-events: none;
        bottom: 10px;
        left: 50%;
        width: 90%;
        transform: translate(-50%, 0);
        border-radius: 20px;
        background-color: var(--overlay);
        padding: 10px;
    }

    .info-box > div {
        margin-bottom: 6px;
    }

    span.small {
        font-size: 10px;
    }
</style>
