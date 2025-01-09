# Följe

This is a project to use a camera and moving light fixture to create a follow light. After setting it up you will have a camera-feed of the stage where you can follow a person with your mouse and PAN/TILT instructions will be sent over sACN to any dmx-controller supporting it.

It started as a python-script written in 6h for Kårspexet 23/24 by Elias Lundell, but has since been rewritten and expanded, and is now being rewritten in Go with [wails](https://github.com/wailsapp/wails). For the original project see [./python-project/](https://github.com/LogFlames/folje/tree/main/python-project).

## Todo in rewrite

- Add sACN output
- Add documentation and instructinos in readme

## Build

This is build using [wails](https://wails.io/) and [Go](https://go.dev/). To build the project, [install wails](https://wails.io/docs/gettingstarted/installation) and run `wails build`.
