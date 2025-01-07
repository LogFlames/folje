<script lang="ts">
    import { type Writable, writable } from "svelte/store";
    import type { CalibrationPoint, Fixture } from "./types";
    import { v4 as uuidv4 } from "uuid";
    import { createEventDispatcher } from "svelte";

    export let fixtures: Writable<{ [id: string]: Fixture }>;
    export let calibrationPoints: Writable<{ [id: string]: CalibrationPoint }>;

    const dispatch = createEventDispatcher();

    let selectedId = null;

    function removeFixture(id) {
        fixtures.update((fixtures) => {
            delete fixtures[id];
            return fixtures;
        });
    }

    function addFixture() {
        let newId = uuidv4();

        fixtures.update((fixtures) => {
            fixtures[newId] = {
                id: newId,
                name: `fixture-${Object.keys(fixtures).length + 1}`,
                universe: 1,
                panAddress: 0,
                finePanAddress: 0,
                tiltAddress: 0,
                fineTiltAddress: 0,
                minPan: 0,
                maxPan: 65535,
                minTilt: 0,
                maxTilt: 65535,
                calibration: {},
            };
            return fixtures;
        });
    }

    function selectId(id) {
        if (selectedId === id) {
            selectedId = null;
        } else {
            selectedId = id;
        }
    }
</script>

<div class="overlay-content">
    <div class="side-by-side-holder">
        <div class="side-by-side-left">
            <div class="fixture-list-scroll">
                {#each Object.values($fixtures) as fixture, index (fixture.id)}
                    <button
                        class="fixture-button {fixture.id === selectedId
                            ? 'selected-fixture-button'
                            : ''}"
                        on:click={() => {
                            selectId(fixture.id);
                        }}
                    >
                        Fixture {index + 1}: {fixture.name || "Unnamed"}
                        <br />
                        {Object.keys($calibrationPoints).filter(
                            (calibration_point_id) =>
                                !Object.keys(fixture.calibration).includes(
                                    calibration_point_id,
                                ),
                        ).length > 0
                            ? "<!> uncalibrated <!>"
                            : ""}
                    </button>
                {/each}
            </div>
            <button on:click={addFixture}> Add Fixture </button>
        </div>
        {#if selectedId !== null && $fixtures[selectedId] !== undefined}
            <div class="side-by-side-right">
                <div class="fixture">
                    <h2>{$fixtures[selectedId].name || "Unnamed"}</h2>

                    <div>
                        <label>
                            Name:
                            <input
                                type="text"
                                bind:value={$fixtures[selectedId].name}
                                placeholder="Enter fixture name"
                            />
                        </label>
                    </div>
                    <div>
                        <label>
                            Universe:
                            <input
                                type="number"
                                bind:value={$fixtures[selectedId].universe}
                                min="1"
                                step="1"
                            />
                        </label>
                    </div>
                    <div>
                        <label>
                            Pan Address:
                            <input
                                type="number"
                                bind:value={$fixtures[selectedId].panAddress}
                                min="1"
                                max="512"
                            />
                        </label>
                    </div>
                    <div>
                        <label>
                            Fine Pan Address:
                            <input
                                type="number"
                                bind:value={$fixtures[selectedId]
                                    .finePanAddress}
                                min="1"
                                max="512"
                            />
                        </label>
                    </div>
                    <div>
                        <label>
                            Tilt Address:
                            <input
                                type="number"
                                bind:value={$fixtures[selectedId].tiltAddress}
                                min="1"
                                max="512"
                            />
                        </label>
                    </div>
                    <div>
                        <label>
                            Fine Tilt Address:
                            <input
                                type="number"
                                bind:value={$fixtures[selectedId]
                                    .fineTiltAddress}
                                min="1"
                                max="512"
                            />
                        </label>
                    </div>
                    <div>
                        <label>
                            Min Pan:
                            <input
                                type="number"
                                bind:value={$fixtures[selectedId].minPan}
                                min="0"
                                max="65535"
                            />
                        </label>
                    </div>
                    <div>
                        <label>
                            Max Pan:
                            <input
                                type="number"
                                bind:value={$fixtures[selectedId].maxPan}
                                min="0"
                                max="65535"
                            />
                        </label>
                    </div>
                    <div>
                        <label>
                            Min Tilt:
                            <input
                                type="number"
                                bind:value={$fixtures[selectedId].minTilt}
                                min="0"
                                max="65535"
                            />
                        </label>
                    </div>
                    <div>
                        <label>
                            Max Tilt:
                            <input
                                type="number"
                                bind:value={$fixtures[selectedId].maxTilt}
                                min="0"
                                max="65535"
                            />
                        </label>
                    </div>
                    <button
                        class="fixture-settings-button"
                        on:click={() => {
                            dispatch("calibrate_all_points", {
                                fixture_id: selectedId,
                            });
                        }}>Calibrate for all points</button
                    >
                    <br />
                    <button
                        class="fixture-settings-button"
                        on:click={() => {
                            dispatch("calibrate_one_point", {
                                fixture_id: selectedId,
                            });
                        }}>Calibrate for one point</button
                    >
                    <br />
                    <button
                        class="fixture-settings-button remove-fixture-button"
                        on:click={() => removeFixture(selectedId)}
                        >Remove Fixture</button
                    >
                </div>
            </div>
        {/if}
    </div>
</div>

<style>
    .fixture-button {
        display: block;
        padding: 10px;
        margin: 10px;
        width: 280px;
    }

    .fixture {
        text-align: left;
    }

    .overlay-content {
        background: #0f0a48;
        padding: 20px;
        border-radius: 8px;
        box-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
        text-align: center;
    }

    .side-by-side-holder {
        display: flex;
        justify-content: space-between;
        gap: 1rem;
    }

    .side-by-side-left {
        width: 300px;
        text-align: center;
    }

    .side-by-side-right {
        max-width: 40vw;
        padding: 20px;
        text-align: center;
    }

    .fixture-list-scroll {
        overflow-y: scroll;
        max-height: 60vh;
    }

    .selected-fixture-button {
        background-color: hsl(218, 56%, 25%);
    }

    .fixture-settings-button {
        margin: 10px 0;
    }

    .remove-fixture-button {
        background-color: red;
    }
</style>
