<script lang="ts">
    import { onMount } from "svelte";
    import { writable, get } from "svelte/store";
    import Fixtures from "./Fixtures.svelte";
    import type { Fixture, CalibrationPoint } from "./types";

    interface MousePos {
        x: number;
        y: number;
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
    let fixtures = writable<Fixture[]>([]);

    let calibrationPoints = writable<CalibrationPoint[]>([
        { uid: "aoeu", x: 0.5, y: 0.5 },
        { uid: "2aoeu", x: 0.1, y: 0.5 },
        { uid: "3aoeu", x: 0.9, y: 0.5 },
        { uid: "4aoeu", x: 0.5, y: 0.1 },
        { uid: "5aoeu", x: 0.5, y: 0.9 },
    ]);

    let showSettingsButtons = false;
    let hideAllSettings = false;

    let showMousePosition = false;
    let showCalibrationPoints = false;

    const toggleShowMousePosition = () => {
        showMousePosition = !showMousePosition;
    };

    const toggleShowCalibrationPoints = () => {
        showCalibrationPoints = !showCalibrationPoints;
    };

    const toggleShowFixtures = () => {
        showFixtures = !showFixtures;
    };

    const toggleShowSettingsDropdown = () => {
        showSettingsButtons = !showSettingsButtons;
    };

    onMount(() => {
        getStream().then(getDevices).then(gotDevices);
    });

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
                height: { ideal: 1080 }
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

<main>
    <div class="content" on:mousemove={handleMouseMove}>
        <video autoplay bind:this={videoElement} />
        {#if showCalibrationPoints}
            <div
                class="video-overlay"
                style="top: {videoStartY}px; left: {videoStartX}px; width: {videoRenderedWidth}px; height: {videoRenderedHeight}px;"
            >
                {#each $calibrationPoints as calibrationPoint, index}
                    <div
                        class="calibration-point"
                        style="
                            top: {calibrationPoint.y * 100}%;
                            left: {calibrationPoint.x * 100}%;
                        "
                    ></div>
                {/each}
            </div>
        {/if}
    </div>
    <button
        class="settings-button {hideAllSettings ? 'hidden' : ''}"
        on:click={toggleShowSettingsDropdown}
    >
        <svg
            xmlns="http://www.w3.org/2000/svg"
            x="0px"
            y="0px"
            width="40px"
            height="40px"
            viewBox="0 0 50 50"
            fill="#ffffff"
        >
            <path
                d="M47.16,21.221l-5.91-0.966c-0.346-1.186-0.819-2.326-1.411-3.405l3.45-4.917c0.279-0.397,0.231-0.938-0.112-1.282 l-3.889-3.887c-0.347-0.346-0.893-0.391-1.291-0.104l-4.843,3.481c-1.089-0.602-2.239-1.08-3.432-1.427l-1.031-5.886 C28.607,2.35,28.192,2,27.706,2h-5.5c-0.49,0-0.908,0.355-0.987,0.839l-0.956,5.854c-1.2,0.345-2.352,0.818-3.437,1.412l-4.83-3.45 c-0.399-0.285-0.942-0.239-1.289,0.106L6.82,10.648c-0.343,0.343-0.391,0.883-0.112,1.28l3.399,4.863 c-0.605,1.095-1.087,2.254-1.438,3.46l-5.831,0.971c-0.482,0.08-0.836,0.498-0.836,0.986v5.5c0,0.485,0.348,0.9,0.825,0.985 l5.831,1.034c0.349,1.203,0.831,2.362,1.438,3.46l-3.441,4.813c-0.284,0.397-0.239,0.942,0.106,1.289l3.888,3.891 c0.343,0.343,0.884,0.391,1.281,0.112l4.87-3.411c1.093,0.601,2.248,1.078,3.445,1.424l0.976,5.861C21.3,47.647,21.717,48,22.206,48 h5.5c0.485,0,0.9-0.348,0.984-0.825l1.045-5.89c1.199-0.353,2.348-0.833,3.43-1.435l4.905,3.441 c0.398,0.281,0.938,0.232,1.282-0.111l3.888-3.891c0.346-0.347,0.391-0.894,0.104-1.292l-3.498-4.857 c0.593-1.08,1.064-2.222,1.407-3.408l5.918-1.039c0.479-0.084,0.827-0.5,0.827-0.985v-5.5C47.999,21.718,47.644,21.3,47.16,21.221z M25,32c-3.866,0-7-3.134-7-7c0-3.866,3.134-7,7-7s7,3.134,7,7C32,28.866,28.866,32,25,32z"
            ></path>
        </svg>
    </button>
    <div class="settings {!showSettingsButtons || hideAllSettings ? 'hidden' : ''}">
        <select bind:this={videoSelect} on:change={getStream}>
            {#each $deviceInfos as deviceInfo, index}
                {#if deviceInfo.kind === "videoinput"}
                    <option value={deviceInfo.deviceId}
                        >{deviceInfo.label || `Camera ${index + 1}`}</option
                    >
                {/if}
            {/each}
        </select>
        <button on:click={toggleShowFixtures}> Fixtures </button>
        {#if showFixtures}
            <div class="overlay" on:click={toggleShowFixtures}>
                <div on:click|stopPropagation>
                    <Fixtures bind:fixtures />
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
        <button
            on:click={() => {
                showCalibrationPoints = true;
            }}>Add Calibration Point</button
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
    </div>
</main>

<style>
    main {
        height: 100vh;
    }

    video {
        width: 100%;
        height: 100%;
        background-color: rgb(2, 12, 24);
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
        background-color: rgba(123, 147, 183, 0.2);
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
        background: rgba(0, 0, 0, 0.5);
        z-index: 1000;
    }

    .video-overlay {
        position: absolute;
        z-index: 2;
    }

    .calibration-point {
        position: absolute;
        width: 20px;
        height: 20px;
        background-color: red;
        border-radius: 50%;
        transform: translate(-50%, -50%);
    }

    .info-box {
        position: absolute;
        pointer-events: none;
        top: 0;
        right: 0;
        background-color: rgba(0, 0, 0, 0.4);
        padding: 10px;
        border-radius: 10px;
    }
</style>
