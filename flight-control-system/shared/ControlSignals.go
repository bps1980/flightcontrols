// shared/ControlSignals.go
package shared

type ControlSignals struct {
    Throttle       float64
    PropPitch      float64
    Leanness       float64
    MagnetoSwitch  []bool
    IgnitionSwitch []bool
}
