from ecu_communication import ECUCommunication
from sensors_integration import get_sensor_data
from control_logic import calculate_target_rpm
from safety_checks import check_safety_conditions, emergency_shutdown
import time

def control_loop():
    ecu_comm = ECUCommunication(port="COM1", baud_rate=115200)
    ecu_comm.initialize()
    
    while True:
        sensor_data = get_sensor_data()
        current_rpm = sensor_data['rpm']
        target_rpm = calculate_target_rpm(current_rpm)
        ecu_comm.set_engine_rpm(target_rpm)
        check_safety_conditions(sensor_data)
        time.sleep(0.1)

if __name__ == "__main__":
    try:
        control_loop()
    except KeyboardInterrupt:
        emergency_shutdown()
