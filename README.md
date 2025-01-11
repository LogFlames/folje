# Följe

<img align="right" src="https://github.com/LogFlames/folje/blob/main/build/appicon.png?raw=true" width="256" height="256">

This is a project to use a camera and moving light fixture to create a follow light. After setting it up you will have a camera-feed of the stage where you can follow a person with your mouse and PAN/TILT instructions will be sent over sACN to any dmx-controller supporting it.

It started as a python-script written in 6h for Kårspexet 23/24 by Elias Lundell, but has since been rewritten and expanded, and then been rewritten again in Go with [wails](https://github.com/wailsapp/wails) and [Svelte](https://svelte.dev). For the original project see [d248f34
](https://github.com/LogFlames/folje/commit/d248f3438c96cdaaafaa230d976599d08036f53d).

- [Installation](#installation)
- [Usage](#usage)
  - [Configuration and calibration](#configuration-and-calibration)
    - [Fixtures](#fixtures)
    - [Calibration points](#calibration-points)
    - [sACN configuration](#sacn-configuration)
    - [Locking Position](#locking-position)
- [TODOs and known bugs](#todos-and-known-bugs)
  - [TODOs](#todos)
  - [Known Bugs](#known-bugs)
- [Build](#build)

## Installation

There exists automatic builds for Windows and MacOS. Download the [latest release](https://github.com/LogFlames/folje/releases/tag/latest). This is automatically built from the latest in the main branch and might be broken. Once the project has matured and is no longer in rapid development I will setup stable releases.

The Mac app is not signed, so you will need to disable Gatekeeper to run the app or [build](#build) it yourself.
```bash
xattr -d com.apple.quarantine Följe.app
```

## Usage

Demo video //TODO

### Configuration and calibration

When configuring and calibrating you can generally abort an operation by pressing ESC. Please also read the status info in the top-right corner as well as additional information in the lower middle of the screen.

When calibrating a fixture your mouse position will be converted to a pan/tilt value absolutely (top left is 0%/0%, bottom right is 100%/100%). If you wish finer control you can press and hold SPACE, your mouse movement will now be a small offset from where you started holding SPACE. Make sure to hold press when you click so that your desired pan/tilt is saved. 

A tip is to have a person walk around the stage and position the calibration points at their feet. This will make it easier to hit the correct depth as aiming for the middle of the spot will be inconsistent when things (backdrops, etc) are blocking the light. You can however allways trust the floor to not move. 

Around the outer edge of the calibration points there is a green outline. The fixtures can only track within this area, make sure to calibrate the entire stage/area you want to track later.

#### Fixtures

Only fixtures which are calibrated for all existing calibration ponits will get pan/tilt data. If you have uncalibrated fixtures there will be a warning in the status window and in the fixture settings you can choose to calibrate it only for the points which it is missing. 

If you have a fixture without fine pan/tilt leave those addresses as `0`. Only addresses that are in the range `[1, 512]` will be used.

- `name`: A name for the fixture. This is only for you to know which one it is. It is possible for multiple fixtures to have the same name.
- `universe`: Which DMX universe this fixtures pan/tilt should be sent to.
- `panAddress`: The DMX address for the pan channel.
- `finePanAddress`: The DMX address for the fine pan channel. If your fixture does not have fine pan leave this as 0.
- `tiltAddress`: The DMX address for the tilt channel.
- `fineTiltAddress`: The DMX address for the fine tilt channel. If your fixture does not have fine tilt leave this as 0.
- `minPan`, `maxPan`, `minTilt`, `maxTilt`: The range of the pan/tilt values. This is only used for calibration where the top left corner will be minPan/minTilt and the bottom right corner will be maxPan/maxTilt. Can make calibration easier if this range is as small as needed to cover the stage as you will get more precise control over the direction.

#### Calibration points

Add a calibration point by clicking on `Add calibration point` in the settings, then click on the video where you wish to place the calibration point. It will then go through all fixtures and calibrate them one by one to the new point. If you wish to skip calibrating fixtures for your new point press ESC for all of them.

Remove a calibration point by clicking on `Remove calibration ponit` in the settings and then click on one of the calibration points, you have to hit the quite small red dots. Abort by pressing ESC.

#### sACN configuration

- `ip address`: Följe will automatically detect all non-loopback ip addresses and lets you choose which of these to bind to, make sure choose the correct network interface that can communicate with you console/visualiser/etc.
- `fps`: When Följe is running it will send the latest calibration this many times a second. Make sure it is compatible with you console/reader. It will send updates indiscriminately of if any values have changed since last send.
- `multicast`: Wether to multicast, that is send the sACN packets to all ip addresses that are listening on the network you are connected to.
- `destinations`: If `multicast` is of you have to choose which IP Addresses to send the data to, this would be you console/visualiser/etc.

Nothing is changed until you hit `Apply`. When applying the settings the sACN sender will be stopped a new one will be created using new settings. There will be downtime in the packages sent. 

Theoretically the settings could be more fine-grained (multicast/destination per unvierse) but I never had the need to send different universes to different destinations. If you are interested in this feature please open an issue or PR, I will gladly merge it.

### Locking Position

When tracking someone you might want to move the mouse without having the fixtures follow (to change settings or interact with other programs). This can be done by clicking anywhere on the video. A red dot with a red ring around it will appear at the locked position. Clicking anywhere on the video will unlock it and it will resume following the mouse.

## TODOs and known bugs

### TODOs

- Increase efficiency of the locate point function for the linear interpolation
- Add save/load functionality

### Known Bugs

- When in 'calibrate all'-mode, after you click to confirm a calibration, the fixture moves up to 0, 0 until you move the mouse at which point it snaps back to absolute pan/tilt-mode. Should be solved by sending a 'App.SetPanTiltForFixture' directly when setting the new 'currentlyCalibrating' struct in the clickHandeler.



## Build

This is build using [Go](https://go.dev/), [Wails](https://wails.io/) and [Svelte](https://svelte.dev/). To build the project, install [go](https://go.dev/doc/install), [wails](https://wails.io/docs/gettingstarted/installation), [node](https://nodejs.org/en/download) and run `wails build` from the root of the project:
```bash
wails build
```
It will automatically install all dependencies for both the backand end frontend and build the project. The output will be in the `build/bin` folder.
