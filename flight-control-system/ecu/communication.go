package ecu

import (
    "fmt"
    "flight-control-system/shared"
)

type ECU struct {
    // Define ECU-specific fields, e.g., connection parameters
}

func InitECU() (*ECU, error) {
    // Initialize ECU connection here
    fmt.Println("ECU initialized")
    return &ECU{}, nil
}

func (ecu *ECU) SetThrottle(throttle float64) {
    // Send throttle command to ECU
    fmt.Printf("Setting throttle to %.2f\n", throttle)
}

func (ecu *ECU) SetPropPitch(propPitch float64) {
    // Send propeller pitch command to ECU
    fmt.Printf("Setting prop pitch to %.2f\n", propPitch)
}

func (ecu *ECU) SetLeanness(leanness float64) {
    // Send leanness command to ECU
    fmt.Printf("Setting leanness to %.2f\n", leanness)
}

func (ecu *ECU) UpdateControlSignals(control shared.ControlSignals) {
    ecu.SetThrottle(control.Throttle)
    ecu.SetPropPitch(control.PropPitch)
    ecu.SetLeanness(control.Leanness)
}
