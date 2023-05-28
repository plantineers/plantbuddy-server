#!/usr/bin/python3

"""
This script generates sensor data and saves it in the SENSOR_DATA table in the database.

The steps are as follows:
1. Read all controllers from the CONTROLLER table
2. For each plant group in the CONTROLLER table, read all data from SENSOR_RANGE
table to get the range of values for each sensor
3. Set a realistic start date time
4. Loop through all controllers until the current date time is reached (add 5 minutes on each iteration)):
    a. For each controller, loop through all sensors
    b. Generate a random value between the min and max values for that sensor
    c. Save the value in the SENSOR_DATA table
"""

from datetime import datetime, timedelta
from random import uniform
from sqlite3 import Cursor, connect
from sys import exit as _exit
from typing import Dict, List, Tuple


def connect_to_db(file: str) -> Cursor:
    """Connects to the database and returns a cursor"""
    return connect(file).cursor()


def read_controllers(cursor: Cursor) -> List[Tuple[str, int]]:
    """Reads all controllers from the CONTROLLER table and maps them to [(UUID, plant group)]"""
    return cursor.execute("SELECT UUID, PLANT_GROUP FROM CONTROLLER").fetchall()


def read_sensor_ranges(cursor: Cursor, plant_group: int) -> List[Tuple[str, int, int]]:
    """Reads all sensors for a given plant group from the SENSOR_RANGE table and maps them to [(sensor, min, max)]"""
    return cursor.execute("SELECT SENSOR, MIN, MAX FROM SENSOR_RANGE WHERE PLANT_GROUP = ?",
                          (plant_group,)).fetchall()


def generate_sensor_data(cursor: Cursor, controller: str, sensor: Tuple[str, int, int], time: datetime) -> None:
    """Inserts the generated sensor data into the SENSOR_DATA table"""
    value = round(uniform(sensor[1], sensor[2]), 2)
    cursor.execute("INSERT INTO SENSOR_DATA (CONTROLLER, SENSOR, VALUE, TIMESTAMP) VALUES (?, ?, ?, ?)",
                   (controller, sensor[0], value, time.isoformat()))


if __name__ == '__main__':
    cur = connect_to_db('buddy.sqlite')
    controllers = read_controllers(cur)

    full_controllers: Dict[str, List[Tuple[str, int, int]]] = {}
    for c, pg in controllers:
        full_controllers[c] = read_sensor_ranges(cur, pg)

    t = datetime(2023, 5, 20, 0, 0, 0)
    COUNTER = 0
    while t < datetime.now():
        for c, sr in full_controllers.items():
            for s in sr:
                generate_sensor_data(cur, c, s, t)
                COUNTER += 1
        cur.connection.commit()
        t = t + timedelta(minutes=5)

    print(f"Generated {COUNTER} sensor data entries")
else:
    print("This script cannot be imported")
    _exit(-1)
