# Följe

A python script to control a front profile as a follow light, made in 6h for Kårspexet 23/24 by Elias Lundell.

## Install
`pip install -r requirements.txt`

## Find camera
Print available cameras with `python test_available_cam.py` and try to show their stream. Press `q` to exit the window.

Run the script with `python cam.py -i [INDEX]`

## Calibrate
Run with `python cam.py -i [INDEX] -c`. It loads `cal.txt` file if it exists.

`s` - save calibration points to `cal.txt`
`t` - enter and exit preview mode
`x` - remove the calibration point closest to the mouse
`c` - remove all calibration points
`q` - quit (must be pressed twice if calibration points are unsaved)

The calibration happens in two steps. First the mouse controls absolute Pan/Tilt based on 0-100% on the X and Y axis. When the light is in a desired position that has not been calibrated yet, press the LMB, the light will now freeze. Move the cursor to where the person is standing.


### Calibration tips

If the camera is more parallell to the stage, take a lighting doll and light their face, then click on where their feet are. Do not calibrate based on the center of the beam. It will be difficult when the person is standing 

### Calibration file `cal.txt`
The calibration file `cal.txt` contains 4 columns.
* pan - 16-bit (2 channels - for fixtures with fine pan)
* tilt - 16-bit (2 channels - for fixtures with fine tilt)
* mouse x position - normalized between 0 and 1
* mouse y position - normalized between 0 and 1 (and inverted)

## Running
Run with `python cam.py -i [INDEX]`. It loads `cal.txt`.
When the mouse is outside the calibration zone (outlined in green) it will not update the P/T information sent over sACN.

If you click LMB the cursor indicator will be turned into a red dot and the P/T will be locked to that position. You can freely move the mouse without affecting the light. To gain back control, click LMB again.

## Config
It currently sends pan, fine pan, tilt, fine tilt to channels 28-31 on universe 7. These must the changed in the code.
