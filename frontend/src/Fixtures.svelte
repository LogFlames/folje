<script lang="ts">
    import { type Writable, writable } from 'svelte/store';
    import type { Fixture } from './types';
    
    export let fixtures: Writable<Fixture[]>;
    
    let selectedIndex = null;
    
    function removeFixture(index) {
        fixtures.update((fixtures) => {
            return fixtures.filter((_, i) => i !== index);
        });
    }
    
    function addFixture() {
        fixtures.update((fixtures) => {
            return [
            ...fixtures,
            {
                name: '',
                universe: 1,
                panAddress: 0,
                finePanAddress: 0,
                tiltAddress: 0,
                fineTiltAddress: 0,
                minPan: 0,
                maxPan: 65535,
                minTilt: 0,
                maxTilt: 65535
            }
            ];
        });
    }
    
    function selectIndex(index) {
        if (selectedIndex === index) {
            selectedIndex = null;
        } else {
            selectedIndex = index;
        }
    }
</script>

<style>
    .fixture-button {
        display: block;
        padding: 10px;
        margin: 10px;
        width: 200px;
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
        width: 220px;
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
    
    .remove-fixture-button {
        background-color: red;
        margin: 15px 0;
    }
</style>

<div class="overlay-content">
    <div class="side-by-side-holder">
        <div class="side-by-side-left">
            <div class="fixture-list-scroll">
                {#each $fixtures as fixture, index}
                <button class="fixture-button {index === selectedIndex ? "selected-fixture-button" : ""}" on:click={() => {selectIndex(index)}}>
                    Fixture {index + 1}: {fixture.name || "Unnamed"}
                </button>
                {/each}
            </div>
            <button on:click={addFixture}>
                Add Fixture
            </button>
        </div>
        {#if selectedIndex !== null && $fixtures[selectedIndex] !== undefined}
        <div class="side-by-side-right">
            <div class="fixture">
                <h2>
                    Fixture {selectedIndex + 1}: {$fixtures[selectedIndex].name || 'Unnamed'}
                </h2>
                
                <div>
                    <label>
                        Name:
                        <input type="text" bind:value={$fixtures[selectedIndex].name} placeholder="Enter fixture name" />
                    </label>
                </div>
                <div>
                    <label>
                        Universe:
                        <input type="number" bind:value={$fixtures[selectedIndex].universe} min="1" />
                    </label>
                </div>
                <div>
                    <label>
                        Pan Address:
                        <input type="number" bind:value={$fixtures[selectedIndex].panAddress} />
                    </label>
                </div>
                <div>
                    <label>
                        Fine Pan Address:
                        <input type="number" bind:value={$fixtures[selectedIndex].finePanAddress} />
                    </label>
                </div>
                <div>
                    <label>
                        Tilt Address:
                        <input type="number" bind:value={$fixtures[selectedIndex].tiltAddress} />
                    </label>
                </div>
                <div>
                    <label>
                        Fine Tilt Address:
                        <input type="number" bind:value={$fixtures[selectedIndex].fineTiltAddress} />
                    </label>
                </div>
                <div>
                    <label>
                        Min Pan:
                        <input type="number" bind:value={$fixtures[selectedIndex].minPan} />
                    </label>
                </div>
                <div>
                    <label>
                        Max Pan:
                        <input type="number" bind:value={$fixtures[selectedIndex].maxPan} />
                    </label>
                </div>
                <div>
                    <label>
                        Min Tilt:
                        <input type="number" bind:value={$fixtures[selectedIndex].minTilt} />
                    </label>
                </div>
                <div>
                    <label>
                        Max Tilt:
                        <input type="number" bind:value={$fixtures[selectedIndex].maxTilt} />
                    </label>
                </div>
                <button>Calibrate</button>
                <button class="remove-fixture-button" on:click={() => removeFixture(selectedIndex)}>Remove Fixture</button>
            </div>
        </div>
        {/if}
    </div>
</div>