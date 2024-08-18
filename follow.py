import os
from typing import Dict, List

import cv2
from scipy.interpolate import LinearNDInterpolator
from scipy.spatial import ConvexHull

from fixture import Fixture, load_fixtures
from utils import CalibrationPoint, PanTilt, create_cap, sACNSenderWrapper

WIDTH = 1920
HEIGHT = 1080

def partial_pan(percent, fixture: Fixture | None):
    if fixture is None:
        return int(percent * (2**16 - 1))

    p = int(fixture.min_pan + (fixture.max_pan - fixture.min_pan) * percent)
    return p

def partial_tilt(percent, fixture: Fixture | None):
    if fixture is None:
        return int(percent * (2**16 - 1))

    t = int(fixture.min_tilt + (fixture.max_tilt - fixture.min_tilt) * percent)
    return t

class Follower:
    def __init__(self, vcindex: int, sender: sACNSenderWrapper):
        self.vcindex = vcindex
        self.sender = sender

        self.mouse = {"x": 0.5, "y": 0.5, "pressed": False}
        self.senderStarted = {"started": False}

        self.fixtures: List[Fixture] = load_fixtures()
        self.sender.setup_universes({fixture.universe for fixture in self.fixtures})

        self.load_calibration()

    def mouse_callback(self, event, x, y, *_):
        if event == cv2.EVENT_LBUTTONUP:
            self.mouse["pressed"] = True

        self.mouse["x"] = max(min(x / WIDTH, 1), 0)
        self.mouse["y"] = 1 - max(min(y / HEIGHT, 1), 0)

    def start_run(self):
        if len(self.calibration) < 3:
            raise Exception("Too few calibration points to run. It requires atleast three to function.")

        cap = create_cap(self.vcindex)
        cv2.setMouseCallback("Camera", self.mouse_callback)

        if not cap.isOpened():
            raise Exception("Error: Could not open camera.")

        pan_interps, tilt_interps, points = self.generate_ndinterps()
        hull = ConvexHull(points)

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
                lastx = self.mouse["x"]
                lasty = self.mouse["y"]

                for fixture in self.fixtures:
                    self.sender.schedule_to_sacn(pan_interps[fixture.uid](self.mouse["x"], self.mouse["y"]), tilt_interps[fixture.uid](self.mouse["x"], self.mouse["y"]), fixture)

                frame = cv2.circle(frame, (int(lastx * WIDTH), int((1 - lasty) * HEIGHT)), radius = 10, color = (0, 255, 0), thickness=2)
                frame = cv2.circle(frame, (int(lastx * WIDTH), int((1 - lasty) * HEIGHT)), radius = 2, color = (0, 255, 0), thickness=-1)
            else:
                frame = cv2.circle(frame, (int(lastx * WIDTH), int((1 - lasty) * HEIGHT)), radius = 8, color = (0, 0, 255), thickness=-1)

            if self.mouse["pressed"]:
                self.mouse["pressed"] = False
                activeReading = not activeReading

            for i in range(len(hull.vertices)):
                p1 = points[hull.vertices[i]]
                p2 = points[hull.vertices[(i + 1) % len(hull.vertices)]]

                frame = cv2.line(frame, (int(p1[0] * WIDTH), int((1 - p1[1]) * HEIGHT)), (int(p2[0] * WIDTH), int((1 - p2[1]) * HEIGHT)), color = (0, 255, 0), thickness = 2)

            cv2.imshow('Camera', frame)

            k = cv2.waitKey(1)
            if k & 0xFF == ord('q'):
                break

        cap.release()
        cv2.destroyAllWindows()

    def start_calibrate(self):
        cap = create_cap(self.vcindex)

        cv2.setMouseCallback("Camera", self.mouse_callback)

        state = "pan/tilt-absolute" # pan/tilt-absolute, mouse, track
        active_fixture = None

        last_unfinished_point: Dict[str, PanTilt] = {}
        unfinished_point: Dict[str, PanTilt] = {}

        if not cap.isOpened():
            print("Error: Could not open camera.")
            return

        require_confirm_quit = True
        while True:
            ret, frame = cap.read()

            if not ret:
                print("Error: Could not read frame.")
                break

            frame = cv2.resize(frame, (WIDTH, HEIGHT))

            if state == "pan/tilt-absolute":
                p = partial_pan(self.mouse["x"], active_fixture)
                t = partial_tilt(self.mouse["y"], active_fixture)

                if active_fixture is not None:
                    if active_fixture.uid in last_unfinished_point:
                        pan_percent = (last_unfinished_point[active_fixture.uid].pan - active_fixture.min_pan) / (active_fixture.max_pan - active_fixture.min_pan)
                        tilt_percent = (last_unfinished_point[active_fixture.uid].tilt - active_fixture.min_tilt) / (active_fixture.max_tilt - active_fixture.min_tilt)
                        frame = cv2.circle(frame, (int(pan_percent * WIDTH), int((1 - tilt_percent) * HEIGHT)), radius = 10, color = (255, 0, 0), thickness=2)

                if self.mouse["pressed"] and active_fixture is not None:
                    self.mouse["pressed"] = False
                    print(f"Freezing {active_fixture.uid}, start moving next fixture.")
                    unfinished_point[active_fixture.uid] = PanTilt(p, t)
                    active_fixture = None

                if active_fixture is None:
                    for fixture in self.fixtures:
                        if fixture.uid not in unfinished_point:
                            active_fixture = fixture
                            print(f"Currently moving: {fixture.uid}")
                            break
                    else:
                        print("State 'mouse': All spots are positioned. Move cursor to current spot location.")
                        state = "mouse"
                else:
                    self.sender.schedule_to_sacn(p, t, active_fixture)

            elif state == "mouse":
                if self.mouse["pressed"]:
                    self.mouse["pressed"] = False
                    self.calibration.append(CalibrationPoint(self.mouse["x"], self.mouse["y"], unfinished_point))
                    last_unfinished_point = {key: value for key, value in unfinished_point.items()}
                    unfinished_point = {}
                    
                    active_fixture = None
                    require_confirm_quit = True
                    state = "pan/tilt-absolute"
                    print("Saved calibration point. Going to 'pan/tilt-absolute' mode.")

            elif state == "track":
                pan_interps, tilt_interps, _ = self.generate_ndinterps()

                for fixture in self.fixtures:
                    self.sender.schedule_to_sacn(pan_interps[fixture.uid](self.mouse["x"], self.mouse["y"]), tilt_interps[fixture.uid](self.mouse["x"], self.mouse["y"]), fixture)

            for calibration_point in self.calibration:
                frame = cv2.circle(
                    frame, 
                    (
                        int(calibration_point.x * WIDTH), 
                        int((1 - calibration_point.y) * HEIGHT)
                    ), 
                    color = (0, 0, 255), 
                    radius = 4, 
                    thickness = -1)

            if len(self.calibration) > 2:
                points = []
                for calibration_point in self.calibration:
                    points.append((calibration_point.x, calibration_point.y))
                hull = ConvexHull(points)

                for i in range(len(hull.vertices)):
                    p1 = points[hull.vertices[i]]
                    p2 = points[hull.vertices[(i + 1) % len(hull.vertices)]]

                    frame = cv2.line(
                        frame, 
                        (
                            int(p1[0] * WIDTH), 
                            int((1 - p1[1]) * HEIGHT)
                        ), 
                        (
                            int(p2[0] * WIDTH), 
                            int((1 - p2[1]) * HEIGHT)
                        ), 
                        color = (0, 255, 0), 
                        thickness = 3)

            cv2.imshow("Camera", frame)

            k = cv2.waitKey(1)

            match chr(k & 0xFF):
                case 'c':
                    self.calibration = []
                    require_confirm_quit = True
                    state = "pan/tilt-absolute" # pan/tilt-absolute, mouse, track
                    unfinished_point = {}
                    active_fixture = None
                    print("Cleared calibration data. Moving to pan/tilt mode")
                case 'q':
                    if require_confirm_quit:
                        require_confirm_quit = False
                        print("ATTENTION! You have unsaved calibration data, to quit, press q again.")
                    else:
                        break
                case 'r':
                    if state == "pan/tilt-absolute":
                        unfinished_point = {}
                        active_fixture = None
                        print("Reset halfway calibration point.")
                    else:
                        print("Cannot reset halfway calibration when not in 'pan/tilt-absolute' mode.")
                case 's':
                    self.save_calibration()
                    require_confirm_quit = False
                    print("Saved calibration to cal.txt.")
                case 't':
                    if state == "pan/tilt-absolute":
                        if len(unfinished_point) == 0:
                            state = "track"
                            print("State 'track': For testing calibration.")
                        else:
                            print("Cannot go into track mode while halfway through a new calibration point. Press 'r' to reset halfway calibration.")
                    elif state == "track":
                        state = "pan/tilt-absolute"
                        active_fixture = None
                        unfinished_point = {}
                        print("State 'pan/tilt-absolute': Position all fixtures, one at a time.")
                    elif state == "mouse":
                        print("Cannot go into track-mode while in state 'mouse'")
                case 'u':
                    self.calibration = self.calibration[:-1]
                    print("Undid the last added calibration point.")
                case 'x':
                    smallestDist = 10000000000
                    smallestIndex = -1
                    for i in range(len(self.calibration)):
                        dist = ((self.mouse["x"] - self.calibration[i].x)**2 + (self.mouse["y"] - self.calibration[i].y)**2)**0.5
                        if dist < smallestDist:
                            smallestDist = dist
                            smallestIndex = i

                    if smallestIndex >= 0:
                        del self.calibration[smallestIndex]
                    require_confirm_quit = True
                    print("Removed calibration point closest to mouse.")

            self.sender.send_to_sacn()

        cap.release()
        cv2.destroyAllWindows()

    def generate_ndinterps(self):
        pans: Dict[str, int] = {}
        tilts: Dict[str, int] = {}
        points = []
        for point in self.calibration:
            points.append((point.x, point.y))
            for fixture in point.pt:
                if fixture not in pans:
                    pans[fixture] = []
                pans[fixture].append(point.pt[fixture].pan)

                if fixture not in tilts:
                    tilts[fixture] = []
                tilts[fixture].append(point.pt[fixture].tilt)

            for fixture in self.fixtures:
                if fixture.uid not in point.pt:
                    raise Exception(f"Error: Fixture with uid: {fixture.uid} is missing calibration on point.")

        pan_interps = {fixture.uid: LinearNDInterpolator(points = points, values = pans[fixture.uid]) for fixture in self.fixtures}
        tilt_interps = {fixture.uid: LinearNDInterpolator(points = points, values = tilts[fixture.uid]) for fixture in self.fixtures}

        return pan_interps, tilt_interps, points

    def load_calibration(self):
        self.calibration: List[CalibrationPoint] = []

        if os.path.exists("cal.txt"):
            with open("cal.txt", "r", encoding = "utf-8") as f:
                for line in f:
                    parts = line.split()
                    x = float(parts[0])
                    y = float(parts[1])
                    pt = {}
                    for i in range(2, len(parts), 3):
                        pt[parts[i]] = PanTilt(int(parts[i + 1]), int(parts[i + 2]))

                    self.calibration.append(CalibrationPoint(x, y, pt))

    def save_calibration(self):
        with open("cal.txt", "w", encoding = "utf-8") as f:
            content = ""
            for point in self.calibration:
                pt = " ".join([f"{fixture} {point.pt[fixture].pan} {point.pt[fixture].tilt}" for fixture in point.pt])
                content += f"{point.x} {point.y} {pt}\n"
            f.write(content)
