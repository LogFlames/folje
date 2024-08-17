import math
import sacn
import cv2
from typing import Dict, List
from dataclasses import dataclass
from collections import namedtuple

from fixture import Fixture

PanTilt = namedtuple("PanTilt", "pan tilt")

class sACNSenderWrapper:
    def __init__(self):
        self.sender = sacn.sACNsender()
        self.senderStarted = False

        self.dmx_prep: Dict[int, List[int]] = {}

    def setup_universes(self, universes):
        for universe in universes:
            self.sender.activate_output(universe)
            self.sender[universe].multicast = True
            self.sender[universe].fps = 35

    def schedule_to_sacn(self, p_16bit: int, t_16bit: int, fixture: Fixture):
        if math.isnan(p_16bit) or math.isnan(t_16bit):
            return
        pan = int(p_16bit // 256)
        fpan = int(p_16bit % 256)
        tilt = int(t_16bit // 256)
        ftilt = int(t_16bit % 256)

        if fixture.universe not in self.dmx_prep:
            self.dmx_prep[fixture.universe] = [0] * 512

        self.dmx_prep[fixture.universe][fixture.pan - 1] = pan
        if fixture.fpan > 0:
            self.dmx_prep[fixture.universe][fixture.fpan - 1] = fpan
        self.dmx_prep[fixture.universe][fixture.tilt - 1] = tilt
        if fixture.ftilt > 0:
            self.dmx_prep[fixture.universe][fixture.ftilt - 1] = ftilt

    def send_to_sacn(self):
        for universe in self.dmx_prep:
            self.sender[universe].dmx_data = self.dmx_prep[universe]

        if not self.senderStarted:
            self.senderStarted = True
            self.sender.start()

    def __enter__(self):
        print("Opening sACN")
        return self

    def __exit__(self, type, value, traceback):
        print("Closing sACN")
        self.sender.stop()

@dataclass
class CalibrationPoint:
    x: float
    y: float
    pt: Dict[str, PanTilt]

def create_cap(index):
    cap = cv2.VideoCapture(index, apiPreference = cv2.CAP_DSHOW)
    # cap.set(cv2.CAP_PROP_FRAME_WIDTH, 1920)
    # cap.set(cv2.CAP_PROP_FRAME_HEIGHT, 1080)
    cv2.namedWindow("Camera", cv2.WINDOW_KEEPRATIO)
    return cap
