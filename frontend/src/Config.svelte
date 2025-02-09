<script lang="ts">
    import { get, type Writable } from "svelte/store";
    import * as App from "../wailsjs/go/main/App";
    import type { CalibrationPoint, Fixture, SACNConfig } from "./types";

    export let fixtures: Writable<{ [id: string]: Fixture }>;
    export let calibrationPoints: Writable<{ [id: string]: CalibrationPoint }>;
    export let sacnConfig: Writable<SACNConfig>;

    function loadConfig() {
        App.LoadFile().then(content => {
            let obj = JSON.parse(content);

            if (obj["fixtures"] !== undefined) {
                fixtures.set(obj["fixtures"]);
            }

            if (obj["calibrationPoints"] !== undefined) {
                calibrationPoints.set(obj["calibrationPoints"]);
            }

            App.AlertDialog("Loaded Config", "Loaded configuration from file.");
        }).catch(() => {
            App.AlertDialog("Load Config Error", "Error while trying to load configuration from file.");
        });
    }

    function saveConfig() {
        let content = JSON.stringify({
            fixtures: get(fixtures),
            calibrationPoints: get(calibrationPoints),
            date: String(new Date())
        });

        App.SaveFile(content).then(() => {
            App.AlertDialog("Save Config", "Saved configuration to file.");
        }).catch(() => {
            App.AlertDialog("Save Config Error", "Error while trying to save configuration to file.");
        })
    }
</script>

<div>
    <button on:click={loadConfig}>
        Load
    </button>
    <button on:click={saveConfig}>
        Save
    </button>
</div>