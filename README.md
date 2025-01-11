# Följe

<img src="https://github.com/LogFlames/folje/blob/main/build/appicon.png?raw=true" width="256" height="256">

This is a project to use a camera and moving light fixture to create a follow light. After setting it up you will have a camera-feed of the stage where you can follow a person with your mouse and PAN/TILT instructions will be sent over sACN to any dmx-controller supporting it.

It started as a python-script written in 6h for Kårspexet 23/24 by Elias Lundell, but has since been rewritten and expanded, and is now being rewritten in Go with [wails](https://github.com/wailsapp/wails). For the original project see [d248f34
](https://github.com/LogFlames/folje/commit/d248f3438c96cdaaafaa230d976599d08036f53d).

## Configuration and calibration

### Fixtures

If you have a fixture without fine pan/tilt leave those addresses as `0`. Only addresses that are in the range `[1, 512]` will be used.

### Calibration points

### sACN configuration

## Build

This is build using [Go](https://go.dev/), [Wails](https://wails.io/) and [Svelte](https://svelte.dev/). To build the project, [install wails](https://wails.io/docs/gettingstarted/installation) and run `wails build`.
