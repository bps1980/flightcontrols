package actuators

import (
    "fmt"
    "flight-control-system/ecu"
    "flight-control-system/shared"
)

var ecuInstance *ecu.ECU

func InitActuators() error {
    var err error
    ecuInstance, err = ecu.InitECU()
    if err != nil {
        return fmt.Errorf("failed to initialize ECU: %v", err)
    }
    // Initialize other actuators here
    return nil
}

func UpdateActuators(control shared.ControlSignals) {
    for i, magneto := range control.MagnetoSwitch {
        fmt.Printf("Engine %d - Magneto: %v, Ignition: %v\n", i+1, magneto, control.IgnitionSwitch[i])
        if magneto && control.IgnitionSwitch[i] {
            fmt.Printf("Engine %d - Throttle: %.2f, PropPitch: %.2f, Leanness: %.2f\n", i+1, control.Throttle, control.PropPitch, control.Leanness)
            ecuInstance.UpdateControlSignals(control)
        } else {
            fmt.Printf("Engine %d is off\n", i+1)
        }
    }
}
