import math

class AutopilotSystem:
    def __init__(self, gps_module):
        self.gps = gps_module
        self.target_latitude = None
        self.target_longitude = None

    def set_target(self, latitude, longitude):
        self.target_latitude = latitude
        self.target_longitude = longitude

    def calculate_distance(self, lat1, lon1, lat2, lon2):
        # Haversine formula to calculate distance between two lat/lon points
        R = 6371000  # Radius of the Earth in meters
        phi1 = math.radians(lat1)
        phi2 = math.radians(lat2)
        delta_phi = math.radians(lat2 - lat1)
        delta_lambda = math.radians(lon2 - lon1)

        a = math.sin(delta_phi / 2) ** 2 + math.cos(phi1) * math.cos(phi2) * math.sin(delta_lambda / 2) ** 2
        c = 2 * math.atan2(math.sqrt(a), math.sqrt(1 - a))
        distance = R * c
        return distance

    def navigate(self):
        if self.gps.latitude is None or self.gps.longitude is None:
            print("Waiting for GPS fix...")
            return

        distance = self.calculate_distance(self.gps.latitude, self.gps.longitude, self.target_latitude, self.target_longitude)
        print(f"Distance to target: {distance} meters")

        if distance < 10:  # Target reached
            print("Target reached!")
        else:
            self.adjust_course()

    def adjust_course(self):
        # Simplified course adjustment logic
        print("Adjusting course...")
        # Here you would implement the logic to control the vehicle's actuators

if __name__ == "__main__":
    gps = GPSModule()
    autopilot = AutopilotSystem(gps)
    autopilot.set_target(30.2672, -97.7431)  # Example target coordinates

    while True:
        gps.read_data()
        autopilot.navigate()
        time.sleep(1)
