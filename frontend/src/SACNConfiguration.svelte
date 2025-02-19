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
        <table class="sacn-settings-table">
            <tr>
                <td>IP Address:</td>
                <td>
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
                </td>
            </tr>
            <tr>
                <button on:click={refreshIPAdresses}> Refresh </button>
            </tr>
            <tr>
                <td>FPS:</td>
                <td>
                    <input
                        type="number"
                        min="1"
                        step="1"
                        bind:value={$sacnConfig.fps}
                        on:change={sacnConfigUpdated}
                    />
                </td>
            </tr>
            <tr>
                <td>Multicast:</td>
                <td>
                    <input
                        type="checkbox"
                        bind:checked={$sacnConfig.multicast}
                        on:change={sacnConfigUpdated}
                    />
                </td>
            </tr>
            <tr>
                <td>Destinations:</td>
                <td>
                    {#each $sacnConfig.destinations as destination, index}
                        <div>
                            <input
                                type="text"
                                bind:value={$sacnConfig.destinations[index]}
                                on:change={sacnConfigUpdated}
                            />
                            <button
                                on:click={() => {
                                    removeDestination(index);
                                }}
                                class="sacn-destination-remove-button">x</button
                            >
                        </div>
                    {/each}
                    <button
                        class="sacn-add-destination-button"
                        on:click={addDestination}>Add</button
                    >
                </td>
            </tr>
        </table>
        <div class="sacn-settings-separator"></div>
        {#if sacnConfigDirty}
            <div>
                <button class="sacn-cancel-button" on:click={cancelSACNConfig}
                    >Cancel</button
                >
                <button class="sacn-apply-button" on:click={applySACNConfig}
                    >Apply</button
                >
            </div>
        {/if}
    </div>
</div>

<style>
    .sacn-settings-separator {
        margin-top: 20px;
    }

    .sacn-settings-list > div > button {
        margin-top: 10px;
    }

    .sacn-destination-remove-button {
        background-color: var(--main-red-color);
        padding: 2px;
    }

    .sacn-settings-table {
        text-align: left;
    }

    .sacn-settings-table > tr > td:first-child {
        width: 140px;
        align-content: start;
    }

    .sacn-add-destination-button {
        background-color: var(--main-button-color);
        padding: 4px;
    }

    .sacn-cancel-button {
        background-color: var(--main-red-color);
    }
</style>
