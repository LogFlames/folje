<script lang="ts">
    import { get, type Writable } from "svelte/store";
    import * as App from "../wailsjs/go/main/App";
    import type { CalibrationPoint, Fixture, SACNConfig } from "./types";

    export let fixtures: Writable<{ [id: string]: Fixture }>;
    export let calibrationPoints: Writable<{ [id: string]: CalibrationPoint }>;
    export let sacnConfig: Writable<SACNConfig>;

    function loadConfig() {
        App.LoadFile().then(content => {
            if (!content) return;

            let obj;
            try {
                obj = JSON.parse(content);
            } catch (err) {
                App.Log(`Failed to parse config file: ${err}`);
                App.AlertDialog("Load Config Error", "Config file contains invalid data.");
                return;
            }

            if (obj["fixtures"] !== undefined) {
                fixtures.set(obj["fixtures"]);
            }

            if (obj["calibrationPoints"] !== undefined) {
                calibrationPoints.set(obj["calibrationPoints"]);
            }

            // Restore sACN config if present
            if (obj["sacnConfig"] !== undefined) {
                sacnConfig.update((config) => {
                    if (config) {
                        const updatedConfig = {
                            ...config,
                            multicast: obj.sacnConfig.multicast ?? config.multicast,
                            destinations: obj.sacnConfig.destinations ?? config.destinations,
                            fps: obj.sacnConfig.fps ?? config.fps,
                        };
                        // Apply to backend
                        App.SetSACNConfig({
                            IpAddress: updatedConfig.ipAddress,
                            PossibleIpAddresses: updatedConfig.possibleIdAddresses,
                            Fps: updatedConfig.fps,
                            Multicast: updatedConfig.multicast,
                            Destinations: updatedConfig.destinations,
                        });
                        return updatedConfig;
                    }
                    return config;
                });
            }

            App.AlertDialog("Loaded Config", "Loaded configuration from file.");
        }).catch((err) => {
            App.Log(`Failed to load config file: ${err}`);
            App.AlertDialog("Load Config Error", "Error while trying to load configuration from file.");
        });
    }

    function saveConfig() {
        const currentSacnConfig = get(sacnConfig);
        let content = JSON.stringify({
            fixtures: get(fixtures),
            calibrationPoints: get(calibrationPoints),
            sacnConfig: currentSacnConfig ? {
                multicast: currentSacnConfig.multicast,
                destinations: currentSacnConfig.destinations,
                fps: currentSacnConfig.fps,
            } : undefined,
            date: String(new Date())
        });

        App.SaveFile(content).then((saved) => {
            if (saved) {
                App.AlertDialog("Save Config", "Saved configuration to file.");
            }
        }).catch((err) => {
            App.Log(`Failed to save config file: ${err}`);
            App.AlertDialog("Save Config Error", "Error while trying to save configuration to file.");
        })
    }
</script>

<div class="config-buttons">
    <button on:click={loadConfig}>
        Load
    </button>
    <button on:click={saveConfig}>
        Save
    </button>
</div>

<style>
    .config-buttons {
        display: flex;
        gap: 8px;
    }
</style>