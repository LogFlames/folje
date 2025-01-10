<script lang="ts">
    import { type Writable } from "svelte/store";
    import type { CalibrationPoint, Fixture } from "./types";
    import { v4 as uuidv4 } from "uuid";
    import { createEventDispatcher } from "svelte";
    import * as App from "../wailsjs/go/main/App";

    export let fixtures: Writable<{ [id: string]: Fixture }>;
    export let calibrationPoints: Writable<{ [id: string]: CalibrationPoint }>;

    const dispatch = createEventDispatcher();

    let selectedId = null;

    function removeFixture(id) {
        App.ConfirmDialog(
            "Delete fixture",
            `Are you sure you want to delete '${$fixtures[id].name}'`,
        ).then((value) => {
            if (value === "Cancel") {
                return;
            }

            fixtures.update((fixtures) => {
                delete fixtures[id];
                return fixtures;
            });
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

    function fixtureUpdated() {
        fixtures.update((fixtures) => fixtures);
    }
</script>

<div class="overlay-content">
    <div class="side-by-side-holder">
        <div class="side-by-side-left">
            <div class="fixture-list-scroll">
                {#each Object.values($fixtures) as fixture, index (fixture.id)}
                    {@const calibrated =
                        Object.keys($calibrationPoints).filter(
                            (calibration_point_id) =>
                                !Object.keys(fixture.calibration).includes(
                                    calibration_point_id,
                                ),
                        ).length == 0}
                    <button
                        class="fixture-button {fixture.id === selectedId
                            ? 'selected-fixture-button'
                            : ''}"
                        on:click={() => {
                            selectId(fixture.id);
                        }}
                        title={calibrated
                            ? ""
                            : "Uncalibrated for atleast one point."}
                    >
                        Fixture {index + 1}: {fixture.name || "Unnamed"}
                        {calibrated ? "" : "<!>"}
                    </button>
                {/each}
            </div>
            <div class="fixture-list-separator"></div>
            <button on:click={addFixture}> Add Fixture </button>
        </div>
        {#if selectedId !== null && $fixtures[selectedId] !== undefined}
            {@const missingPoints = Object.keys($calibrationPoints).filter(
                (calibration_point_id) =>
                    !Object.keys($fixtures[selectedId].calibration).includes(
                        calibration_point_id,
                    ),
            )}
            <div class="side-by-side-right">
                <div class="fixture">
                    <h2>{$fixtures[selectedId].name || "Unnamed"}</h2>

                    <div>
                        <label>
                            Name:
                            <input
                                type="text"
                                bind:value={$fixtures[selectedId].name}
                                on:change={fixtureUpdated}
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
                                on:change={fixtureUpdated}
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
                                on:change={fixtureUpdated}
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
                                on:change={fixtureUpdated}
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
                                on:change={fixtureUpdated}
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
                                on:change={fixtureUpdated}
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
                                on:change={fixtureUpdated}
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
                                on:change={fixtureUpdated}
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
                                on:change={fixtureUpdated}
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
                                on:change={fixtureUpdated}
                                min="0"
                                max="65535"
                            />
                        </label>
                    </div>
                    <div class="fixture-list-separator"></div>
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
                    {#if missingPoints.length > 0}
                        <br />
                        <button
                            class="fixture-settings-button"
                            on:click={() => {
                                dispatch("calibrate_missing_points", {
                                    fixture_id: selectedId,
                                    calibration_points_missing: missingPoints,
                                });
                            }}
                        >
                            Calibrate for non-calibrated points ({missingPoints.length})
                        </button>
                    {/if}
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
        padding: 8px;
        margin: 10px;
        width: 280px;
    }

    .fixture {
        text-align: left;
    }

    .overlay-content {
        background: var(--secondary-bg-color);
        padding: 20px;
        border-radius: 6px;
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
        background-color: var(--main-button-active-color);
    }

    .fixture-settings-button {
        margin: 4px 0;
        padding: 10px;
    }

    .remove-fixture-button {
        background-color: var(--main-red-color);
    }

    .fixture-list-separator {
        margin-top: 25px;
    }
</style>
