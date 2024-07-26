import argparse

from test_available_cameras import query_and_open_available_cameras
from follow import Follower
from utils import sACNSenderWrapper

def validate_mode(value):
    choices = {
        't': 'test_available_cameras',
        'r': 'run',
        'c': 'calibrate',
        'test_available_cameras': 'test_available_cameras',
        'run': 'run',
        'calibrate': 'calibrate'
    }
    if value in choices:
        return choices[value]
    else:
        raise argparse.ArgumentTypeError(
            f"Invalid mode '{value}'. Choose from '[t]est_available_cameras', '[r]un', or '[c]alibrate'.")

def main():
    parser = argparse.ArgumentParser(prog = "Follow Spot", description="A program for calibrating a spot light's pan and tilt on a camera feed, and then follow the mouse on the same feed.")

    parser.add_argument("-m", "--mode",              dest = "mode",    nargs = 1, required = True,  action = "store", type = validate_mode, help = "Mode of operation: [t]est_available_cameras, [r]un, [c]alibrate")
    parser.add_argument("-i", "--videoCaptureIndex", dest = "vcindex", nargs = 1, required = False, action = "store", type = int,           help = "The index of the video stream to use.")

    vals = parser.parse_args()

    match vals.mode[0]:
        case "calibrate" | "run":
            if vals.vcindex is None:
                parser.error("When in calibrate or run mode you must set --videoCaptureIndex")

            with sACNSenderWrapper() as sender:
                follow = Follower(vcindex = vals.vcindex[0], sender = sender)

                if vals.mode[0] == "calibrate":
                    follow.start_calibrate()
                elif vals.mode[0] == "run":
                    follow.start_run()

        case "test_available_cameras":
            query_and_open_available_cameras()
        case _:
            parser.error(f"Unkown --mode {vals.mode}")

if __name__ == "__main__":
    main()