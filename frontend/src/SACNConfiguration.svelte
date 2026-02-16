<script lang="ts">
    import { get, type Writable } from "svelte/store";
    import * as App from "../wailsjs/go/main/App";
    import { main } from "../wailsjs/go/models";
    import type { SACNConfig } from "./types";

    export let sacnConfig: Writable<SACNConfig>;
    export let sacnConfigDirty: boolean;

    function sacnConfigUpdated() {
        sacnConfigDirty = true;
    }

    function applySACNConfig() {
        sacnConfigDirty = false;
        let sacnConfigToApply = get(sacnConfig);
        App.SetSACNConfig(
            new main.SACNConfig({
                IpAddress: sacnConfigToApply.ipAddress,
                PossibleIpAddresses: sacnConfigToApply.possibleIdAddresses,
                Fps: sacnConfigToApply.fps,
                Multicast: sacnConfigToApply.multicast,
                Destinations: sacnConfigToApply.destinations,
            }),
        );
    }

    function cancelSACNConfig() {
        sacnConfigDirty = false;

        App.GetSACNConfig().then((sacnConfigFromApp) => {
            sacnConfig.set({
                ipAddress: sacnConfigFromApp.IpAddress,
                possibleIdAddresses: sacnConfigFromApp.PossibleIpAddresses,
                fps: sacnConfigFromApp.Fps,
                multicast: sacnConfigFromApp.Multicast,
                destinations: sacnConfigFromApp.Destinations,
            });
        });
    }

    function removeDestination(index: number) {
        sacnConfig.update((sacnConfig) => {
            sacnConfig.destinations.splice(index, 1);
            return sacnConfig;
        });

        sacnConfigDirty = true;
    }

    function addDestination() {
        sacnConfig.update((sacnConfig) => {
            sacnConfig.destinations.push("");
            return sacnConfig;
        });

        sacnConfigUpdated();
    }

    function refreshIPAdresses() {
        App.GetSACNConfig().then((sacnConfigFromApp) => {
            sacnConfig.update((oldSacnConfig) => {
                return {
                    ipAddress: sacnConfigFromApp.IpAddress,
                    possibleIdAddresses: sacnConfigFromApp.PossibleIpAddresses,
                    fps: oldSacnConfig.fps,
                    multicast: oldSacnConfig.multicast,
                    destinations: sacnConfigFromApp.Destinations,
                };
            });
        });

        applySACNConfig();
    }
</script>

<div class="overlay-content">
    <div class="sacn-settings-list">
        <div class="sacn-row">
            <span class="sacn-label">IP Address:</span>
            <select
                bind:value={$sacnConfig.ipAddress}
                on:change={sacnConfigUpdated}
            >
                {#each $sacnConfig.possibleIdAddresses as possibleIdAddress}
                    <option value={possibleIdAddress}
                        >{possibleIdAddress}</option
                    >
                {/each}
            </select>
        </div>
        <div class="sacn-row">
            <span class="sacn-label"></span>
            <button on:click={refreshIPAdresses}>Refresh</button>
        </div>
        <div class="sacn-row">
            <span class="sacn-label">FPS:</span>
            <input
                type="number"
                min="1"
                step="1"
                bind:value={$sacnConfig.fps}
                on:change={sacnConfigUpdated}
            />
        </div>
        <div class="sacn-row">
            <span class="sacn-label">Multicast:</span>
            <input
                type="checkbox"
                bind:checked={$sacnConfig.multicast}
                on:change={sacnConfigUpdated}
            />
        </div>
        <div class="sacn-row sacn-destinations">
            <span class="sacn-label">Destinations:</span>
            <div class="sacn-destination-list">
                {#each $sacnConfig.destinations as destination, index}
                    <div class="sacn-destination-row">
                        <input
                            type="text"
                            bind:value={$sacnConfig.destinations[index]}
                            on:change={sacnConfigUpdated}
                        />
                        <button
                            class="sacn-destination-remove-button btn-danger"
                            on:click={() => {
                                removeDestination(index);
                            }}>x</button
                        >
                    </div>
                {/each}
                <button on:click={addDestination}>Add</button>
            </div>
        </div>
        <div class="sacn-settings-separator"></div>
        {#if sacnConfigDirty}
            <div class="sacn-actions">
                <button class="btn-danger" on:click={cancelSACNConfig}
                    >Cancel</button
                >
                <button class="btn-primary" on:click={applySACNConfig}
                    >Apply</button
                >
            </div>
        {/if}
    </div>
</div>

<style>
    .sacn-settings-list {
        text-align: left;
    }

    .sacn-row {
        display: flex;
        align-items: center;
        gap: 12px;
        margin-bottom: 12px;
    }

    .sacn-label {
        width: 120px;
        flex-shrink: 0;
        color: var(--text-secondary);
        font-size: 13px;
    }

    .sacn-destinations {
        align-items: flex-start;
    }

    .sacn-destination-list {
        display: flex;
        flex-direction: column;
        gap: 6px;
    }

    .sacn-destination-row {
        display: flex;
        align-items: center;
        gap: 6px;
    }

    .sacn-destination-remove-button {
        padding: 4px 8px;
        font-size: 12px;
    }

    .sacn-settings-separator {
        margin-top: 20px;
    }

    .sacn-actions {
        display: flex;
        gap: 10px;
        margin-top: 12px;
    }
</style>
