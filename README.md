# Följe

<img src="https://github.com/LogFlames/folje/blob/main/build/appicon.png?raw=true" width="256" height="256">

This is a project to use a camera and moving light fixture to create a follow light. After setting it up you will have a camera-feed of the stage where you can follow a person with your mouse and PAN/TILT instructions will be sent over sACN to any dmx-controller supporting it.

It started as a python-script written in 6h for Kårspexet 23/24 by Elias Lundell, but has since been rewritten and expanded, and then been rewritten again in Go with [wails](https://github.com/wailsapp/wails) and [Svelte](https://svelte.dev). For the original project see [d248f34
](https://github.com/LogFlames/folje/commit/d248f3438c96cdaaafaa230d976599d08036f53d).

## Installation

Download the [latest release](https://github.com/LogFlames/folje/releases/tag/latest). This is automatically built from the latest in the main branch and might be broken. Once the project has matured I will setup semver releases.

The Mac app is not signed, so you will need to disable Gatekeeper to run the app or build it yourself.
```bash
xattr -d com.apple.quarantine Följe.app
```

## Usage

### Configuration and calibration

#### Fixtures

If you have a fixture without fine pan/tilt leave those addresses as `0`. Only addresses that are in the range `[1, 512]` will be used.

#### Calibration points

#### sACN configuration

## TODOs and known bugs

### TODOs

- Increase efficiency of the locate point function for the linear interpolation
- Add save/load functionality

## Known Bugs

- When in 'calibrate all'-mode, after you click to confirm a calibration, the fixture moves up to 0, 0 until you move the mouse at which point it snaps back to absolute pan/tilt-mode. Should be solved by sending a 'App.SetPanTiltForFixture' directly when setting the new 'currentlyCalibrating' struct in the clickHandeler.

## Build

This is build using [Go](https://go.dev/), [Wails](https://wails.io/) and [Svelte](https://svelte.dev/). To build the project, install [go](https://go.dev/doc/install), [wails](https://wails.io/docs/gettingstarted/installation), [node](https://nodejs.org/en/download) and run `wails build` from the root of the project.
