package utils

type FlightMode int

const (
    SensorMode FlightMode = iota
    PerformanceMode
    SuperPerformanceMode
)

func GetPitchSetpoint(mode FlightMode) float64 {
    switch mode {
    case SensorMode:
        return 10.0
    case PerformanceMode:
        return 20.0
    case SuperPerformanceMode:
        return 30.0
    default:
        return 0.0
    }
}

func GetLeannessSetting(mode FlightMode) float64 {
    switch mode {
    case SensorMode:
        return 0.8
    case PerformanceMode:
        return 1.0
    case SuperPerformanceMode:
        return 1.2
    default:
        return 1.0
    }
}

