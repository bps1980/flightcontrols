def check_safety_conditions(sensor_data):
    max_temp = 100  # Example max temperature
    if sensor_data['temperature'] > max_temp:
        emergency_shutdown()

def emergency_shutdown():
    # Code to handle emergency shutdown
    print("Emergency Shutdown Activated")
