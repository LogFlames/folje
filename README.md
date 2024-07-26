# Följe

A python script to control a front profile as a follow light, it started by being made in 6h for Kårspexet 23/24 by Elias Lundell, but has since been rewritten and expanded for more shows.

## Install
`pip install -r requirements.txt`

## Find camera
Print available cameras with `python main.py -m t` and try to show their stream. Press `q` to exit the window.

Run the script with `python main.py -m r -i [INDEX]`

## Calibrate
Run with `python main.py -m c -i [INDEX]`. It loads `cal.txt` file if it exists.

`c` - remove all calibration points
`q` - quit (must be pressed twice if calibration points are unsaved)
`r` - reset calibration of current point
`s` - save calibration points to `cal.txt`
`t` - enter and exit track mode (where you can preview the calibration so far)
`u` - undo/remove the last added calibration point
`x` - remove the calibration point closest to the mouse

The calibration happens in two steps. First the mouse controls absolute Pan/Tilt based on 0-100% on the X and Y axis. When the light is in a desired position that has not been calibrated yet, press the LMB, the light will now freeze. Do this for each of the fixtures you have configured in `fixtures.toml`. Move the cursor to where the person is standing.

### Calibration tips

If the camera is more parallell to the stage, take a lighting doll and light their face, then click on where their feet are. Do not calibrate based on the center of the beam. It will be difficult when the person is standing 

### Calibration file `cal.txt`
The calibration file `cal.txt` contains 2+3n columns separated by space.
* mouse x position - normalized between 0 and 1
* mouse y position - normalized between 0 and 1 (and inverted)
Then for each fixture in `fixtures.toml`:
* fixture uid - string with unique fixture name
* pan - 16-bit (2 channels - for fixtures with fine pan)
* tilt - 16-bit (2 channels - for fixtures with fine tilt)

## Running
Run with `python main.py -m r -i [INDEX]`. It loads `cal.txt`.
When the mouse is outside the calibration zone (outlined in green) it will not update the P/T information sent over sACN.

If you click LMB the cursor indicator will be turned into a red dot and the P/T will be locked to that position. You can freely move the mouse without affecting the light. To gain back control, click LMB again.

## Config
Add your desired fixtures in `fixtures.toml`. Define the channels for `pan`, `fpan` (fine pan), `tilt`, `ftilt` (fine tilt), the universe, give fixture a name as well as a unique string id that does not contain any spaces.

Currently only the `fixture-uid` is used, the `name` is ignored.
