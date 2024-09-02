import folium
import time
from folium.plugins import MarkerCluster
import random

# Create a map centered at a specific latitude and longitude
m = folium.Map(location=[30.2672, -97.7431], zoom_start=15)
marker_cluster = MarkerCluster().add_to(m)

# Function to simulate GPS data
def get_simulated_gps_data():
    # Simulate random GPS coordinates around the initial point
    lat = 30.2672 + random.uniform(-0.001, 0.001)
    lon = -97.7431 + random.uniform(-0.001, 0.001)
    return lat, lon

# Simulate updating the map with GPS data
for _ in range(10):  # Simulate 10 updates
    lat, lon = get_simulated_gps_data()
    folium.Marker(location=[lat, lon], popup=f'Location: {lat}, {lon}').add_to(marker_cluster)
    time.sleep(1)

# Save the map to an HTML file
m.save('gps_map.html')

print("Map has been updated and saved to gps_map.html")
