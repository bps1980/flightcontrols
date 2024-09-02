package pid

import "time"

type PIDController struct {
    Kp, Ki, Kd    float64
    previousError float64
    integral      float64
    previousTime  time.Time
}

func NewPIDController(kp, ki, kd float64) *PIDController {
    return &PIDController{
        Kp: kp,
        Ki: ki,
        Kd: kd,
        previousTime: time.Now(),
    }
}

func (pid *PIDController) Update(setpoint, measured float64) float64 {
    currentTime := time.Now()
    deltaTime := currentTime.Sub(pid.previousTime).Seconds()
    pid.previousTime = currentTime

    error := setpoint - measured
    pid.integral += error * deltaTime
    derivative := (error - pid.previousError) / deltaTime
    output := pid.Kp*error + pid.Ki*pid.integral + pid.Kd*derivative
    pid.previousError = error
    return output
}

