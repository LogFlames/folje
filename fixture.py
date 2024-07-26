import os
import tomllib
from dataclasses import dataclass

@dataclass
class Fixture:
    name: str
    universe: int
    pan: int
    fpan: int
    tilt: int
    ftilt: int
    uid: str

def load_fixtures():
    if not os.path.exists("fixtures.toml"):
        raise Exception("Could not find fixtures.toml")

    with open("fixtures.toml", "r") as f:
        fixtures = []

        fixture_dict = tomllib.loads(f.read())

        if "fixture" not in fixture_dict:
            raise Exception("Found no fixture in fixtures.toml.")

        for fixture in fixture_dict['fixture']:
            fixtures.append(Fixture(**fixture))

    if len(fixtures) == 0:
        raise Exception("Error: Found no fixtures in fixtures.toml.")

    return fixtures