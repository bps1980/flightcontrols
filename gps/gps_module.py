import serial
import time
import threading
from flask_socketio import SocketIO, emit

class GPSModule:
    def __init__(self, port='/dev/ttyUSB0', baudrate=9600):
        self.ser = serial.Serial(port, baudrate, timeout=1)
        self.latitude = None
        self.longitude = None
        self.altitude = None
        self.speed = None
        self.heading = None
        self.vertical_speed = None
        self.pitch = None
        self.roll = None

    def read_data(self):
        line = self.ser.readline().decode('utf-8', errors='replace').strip()
        if line.startswith('$GPGGA'):  # NMEA sentence for GPS fix data
            self.parse_gpgga(line)

    def parse_gpgga(self, line):
        parts = line.split(',')
        if len(parts) < 6:
            return
        self.latitude = self.convert_to_decimal(parts[2], parts[3])
        self.longitude = self.convert_to_decimal(parts[4], parts[5])
        self.altitude = float(parts[9]) if parts[9] else None
        print(f"Latitude: {self.latitude}, Longitude: {self.longitude}, Altitude: {self.altitude}")

    def convert_to_decimal(self, value, direction):
        if not value or not direction:
            return None
        degrees = float(value[:2])
        minutes = float(value[2:])
        decimal = degrees + minutes / 60
        if direction in ['S', 'W']:
            decimal = -decimal
        return decimal

    def simulate_additional_data(self):
        # This function simulates other data (speed, heading, vertical speed, pitch, roll)
        # Replace this with actual data reading if available
        self.speed = random.uniform(50, 200)
        self.heading = random.uniform(0, 360)
        self.vertical_speed = random.uniform(-10, 10)
        self.pitch = random.uniform(-10, 10)
        self.roll = random.uniform(-30, 30)

def gps_reader_thread(gps, socketio):
    while True:
        gps.read_data()
        gps.simulate_additional_data()  # Simulate additional data
        if gps.latitude and gps.longitude:
            gps_data = {
                'latitude': gps.latitude,
                'longitude': gps.longitude,
                'altitude': gps.altitude,
                'speed': gps.speed,
                'heading': gps.heading,
                'verticalSpeed': gps.vertical_speed,
                'pitch': gps.pitch,
                'roll': gps.roll
            }
            socketio.emit('gps_update', gps_data)
        time.sleep(1)
