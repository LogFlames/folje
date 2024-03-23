import cv2
import os
import argparse
import math
import alphashape
import sacn
import numpy as np
from scipy.interpolate import LinearNDInterpolator
from scipy.spatial import ConvexHull

mouse = {"x": 0.5, "y": 0.5, "pressed": False}
senderStarted = {"started": False}

class sACNSenderWrapper:
    def __init__(self):
        self.sender = sacn.sACNsender()

    def __enter__(self):
        print("Opening sACN")
        return self.sender

    def __exit__(self, type, value, traceback):
        print("Closing sACN")
        self.sender.stop()

def send_to_sacn(sender, p, t):
    if math.isnan(p) or math.isnan(t):
        print("nan")
        return
    pan = int(p // 256)
    fpan = int(p % 256)
    tilt = int(t // 256)
    ftilt = int(t % 256)
    sender[UNIVERSE].dmx_data = (0,) * 28 + (pan, fpan, tilt, ftilt)
    sender.flush([UNIVERSE])
    if not senderStarted["started"]:
        senderStarted["started"] = True
        sender.start()

def mouse_callback(event, x, y, *_):
    if event == cv2.EVENT_LBUTTONUP:
        mouse["pressed"] = True

    mouse["x"] = max(min(x / WIDTH, 1), 0)
    mouse["y"] = 1 - max(min(y / HEIGHT, 1), 0)

def run(sender, cal_data):
    cap = cv2.VideoCapture(vals.vcindex)
    cv2.namedWindow("Camera")
    cv2.setMouseCallback("Camera", mouse_callback)

    if not cap.isOpened():
        print("Error: Could not open camera.")
        return

    points = []
    ts = []
    ps = []
    for d in cal_data:
        points.append((d[2], d[3]))
        ps.append(d[0])
        ts.append(d[1])

    pinterp = LinearNDInterpolator(points, ps)
    tinterp = LinearNDInterpolator(points, ts)

    activeReading = True
    lastx = 0.5
    lasty = 0.5

    while True:
        ret, frame = cap.read()

        if not ret:
            print("Error: Could not read frame.")
            break

        frame = cv2.resize(frame, (WIDTH, HEIGHT))

        if activeReading:
            lastx = mouse["x"]
            lasty = mouse["y"]

            send_to_sacn(sender, pinterp(mouse["x"], mouse["y"]), tinterp(mouse["x"], mouse["y"]))
        else:
            frame = cv2.circle(frame, (int(lastx * WIDTH), int((1 - lasty) * HEIGHT)), radius = 6, color = (0, 0, 255), thickness=-1)

        if mouse["pressed"]:
            mouse["pressed"] = False
            activeReading = not activeReading

        hull = ConvexHull(points)
        for i in range(len(hull.vertices)):
            p1 = points[hull.vertices[i]]
            p2 = points[hull.vertices[(i + 1) % len(hull.vertices)]]

            frame = cv2.line(frame, (int(p1[0] * WIDTH), int((1 - p1[1]) * HEIGHT)), (int(p2[0] * WIDTH), int((1 - p2[1]) * HEIGHT)), color = (0, 255, 0), thickness = 3)

        cv2.imshow('Camera', frame)

        k = cv2.waitKey(1)
        if k & 0xFF == ord('q'):
            break

    cap.release()
    cv2.destroyAllWindows()

def calibrate(sender):
    cap = cv2.VideoCapture(vals.vcindex)
    cv2.namedWindow("Camera")
    cv2.setMouseCallback("Camera", mouse_callback)

    if not cap.isOpened():
        print("Error: Could not open camera.")
        return

    state = "pt"
    cal_data = []
    print("Pan/Tilt mode")

    p, t = 0, 0

    require_confirm_quit = True
    while True:
        ret, frame = cap.read()

        if not ret:
            print("Error: Could not read frame.")
            break

        frame = cv2.resize(frame, (WIDTH, HEIGHT))

        if state == "pt":
            p = int(mouse["x"] * (2**16-1))
            t = int(mouse["y"] * (2**16-1))
            send_to_sacn(sender, p, t)

            if mouse["pressed"]:
                mouse["pressed"] = False
                state = "m"
                print("Frozen P/T. Move mouse to spot")
        elif state == "m":
            if mouse["pressed"]:
                mouse["pressed"] = False
                cal_data.append([p, t, mouse["x"], mouse["y"]])
                state = "pt"
                require_confirm_quit = True
                print("Saved calibration point. P/T mode")

        for d in cal_data:
            frame = cv2.circle(frame, (int(d[2] * WIDTH), int((1 - d[3]) * HEIGHT)), color = (0, 0, 255), radius = 4, thickness=-1)

        if len(cal_data) > 2:
            points = []
            for d in cal_data:
                points.append((d[2], d[3]))
            hull = ConvexHull(points)
            for i in range(len(hull.vertices)):
                p1 = points[hull.vertices[i]]
                p2 = points[hull.vertices[(i + 1) % len(hull.vertices)]]

                frame = cv2.line(frame, (int(p1[0] * WIDTH), int((1 - p1[1]) * HEIGHT)), (int(p2[0] * WIDTH), int((1 - p2[1]) * HEIGHT)), color = (0, 255, 0), thickness = 3)

        cv2.imshow('Camera', frame)

        k = cv2.waitKey(1)
        if k & 0xFF == ord('s'):
            with open("cal.txt", 'w+') as f:
                s = ""
                for d in cal_data:
                    s += f"{d[0]} {d[1]} {d[2]} {d[3]}\n"
                f.write(s)
            print("Saved calibration file cal.txt")
            require_confirm_quit = False
        if k & 0xFF == ord('q'):
            if require_confirm_quit:
                require_confirm_quit = False
                print("You have unsaved calibration data, to quit, press q again.")
            else:
                break

    cap.release()
    cv2.destroyAllWindows()


parser = argparse.ArgumentParser(prog = "Cam", description="Follow Spot")
parser.add_argument( "-i", "-videoCaptureIndex", dest = "vcindex", help = "The index of the video stream to use.", required = True, action = "store", type = int)

parser.add_argument("-c", "--calibrate", dest = "calibrate", action = "store_true", help = "inititate calibration")

vals = parser.parse_args()

WIDTH = 1280
HEIGHT = 960
UNIVERSE = 7

with sACNSenderWrapper() as sender:
    sender.activate_output(UNIVERSE)
    sender[UNIVERSE].multicast = True
    sender[UNIVERSE].fps = 30

    if vals.calibrate:
        calibrate(sender)
    else:
        cal_data = []
        if not os.path.exists("cal.txt"):
            raise Exception("Must have a calibration file")

        with open("cal.txt", "r") as f:
            for line in f:
                p, t, x, y = line.strip().split()
                cal_data.append([int(p), int(t), float(x), float(y)])
        run(sender, cal_data)


