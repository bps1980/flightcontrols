def calculate_target_rpm(current_rpm):
    desired_rpm = 2000  # Example desired RPM
    adjustment_factor = 10  # Example adjustment factor
    if current_rpm < desired_rpm:
        return current_rpm + adjustment_factor
    else:
        return current_rpm - adjustment_factor
