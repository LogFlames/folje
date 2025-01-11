<script lang="ts">
    import { type Writable } from "svelte/store";
    import type {
        CalibratingFixture,
        CalibrationPoints,
        Fixtures,
        Point,
    } from "./types";
    import { calcPan, calcTilt } from "./utils";

    export let addingCalibrationPoint: boolean;
    export let allFixturesCalibrated: Writable<boolean>;
    export let calibrateForOnePointSelectCalibrationPoint: boolean;
    export let calibrationPoints: Writable<CalibrationPoints>;
    export let calibrationPointsToCalibrate: Writable<string[]>;
    export let currentlyCalibrating: Writable<CalibratingFixture | null>;
    export let fixtures: Writable<Fixtures>;
    export let fixturesToCalibrate: Writable<string[]>;
    export let lockMousePos: boolean;
    export let mouseDragStart: Writable<Point | null>;
    export let mousePos: Writable<Point>;
    export let removingCalibrationPoint: boolean;
    export let showCalibrationPoints: boolean;
    export let showMousePosition: boolean;
</script>

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
            Calibrating {$fixtures[$currentlyCalibrating.fixture_id].name} on point
            {$calibrationPoints[$currentlyCalibrating.calibration_point_id]
                .name}
        </div>
        <div>
            Pan: {Math.floor(
                calcPan(
                    $fixtures[$currentlyCalibrating.fixture_id],
                    $mousePos,
                    $mouseDragStart,
                ),
            )}
        </div>
        <div>
            Tilt: {Math.floor(
                calcTilt(
                    $fixtures[$currentlyCalibrating.fixture_id],
                    $mousePos,
                    $mouseDragStart,
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
    {#if !$allFixturesCalibrated}
        <div>&lt;!&gt; There are uncalibrated fixtures &lt;!&gt;</div>
    {/if}
    {#if lockMousePos}
        <div>Mouse postion locked</div>
    {/if}
</div>
{#if addingCalibrationPoint || removingCalibrationPoint || ($currentlyCalibrating !== null && !calibrateForOnePointSelectCalibrationPoint) || calibrateForOnePointSelectCalibrationPoint || Object.keys($fixtures).length === 0 || Object.keys($calibrationPoints).length === 0}
    <div class="tooltip">
        {#if addingCalibrationPoint}
            <div>
                Click on the camera-feed to create a calibration point there.
                <br />
                Press ESC to cancel.
            </div>
        {/if}
        {#if removingCalibrationPoint}
            <div>
                Click on a calibration point to remove it and all calibrations
                to that point on the fixtures.
                <br />
                Press ESC to cancel.
            </div>
        {/if}
        {#if $currentlyCalibrating !== null && !calibrateForOnePointSelectCalibrationPoint}
            <div>
                Calibrating a fixture for the green calibration point. Click to
                lock pan/tilt.
                <br />
                To get finer control of pan and tilt, press and hold Space. You must
                keep pressing space until you have locked pan/tilt by clicking.
                <br />
                Press ESC to not calibrate this fixture. (NOTE if you have more fixtures
                to be calibrated you will move on to the next one)
            </div>
        {/if}
        {#if calibrateForOnePointSelectCalibrationPoint}
            <div>
                Select the point to calibrate the fixture for by clicking on it.
                <br />
                Press ESC to cancel.
            </div>
        {/if}
        {#if !$allFixturesCalibrated}
            <div>
                There are fixtures which are not calibrated for all calibration
                points. These will not get any pan/tilt data.
            </div>
        {/if}
        {#if Object.keys($fixtures).length === 0}
            <div>
                No fixtures. Either load a configuration or add them by going
                into Settings (top left cog) &gt; Fixtures &gt; Add fixture.
            </div>
        {/if}
        {#if Object.keys($calibrationPoints).length === 0}
            <div>
                No calibration points. Either load a configuration or add them
                by going into Settings (top left cog) &gt; Add calibration
                point.
            </div>
        {/if}
    </div>
{/if}

<style>
    .info-box {
        position: absolute;
        pointer-events: none;
        top: 0;
        right: 0;
        background-color: var(--overlay);
        padding: 8px;
        border-radius: 6px;
        text-align: right;
    }

    .tooltip {
        position: absolute;
        pointer-events: none;
        bottom: 10px;
        left: 50%;
        transform: translate(-50%, 0);
        width: 90%;
        border-radius: 20px;
        background-color: var(--overlay);
        padding: 10px;
    }

    .tooltip > div {
        margin-top: 8px;
    }

    .info-box > div {
        margin-bottom: 6px;
    }
</style>