package sensors

import "math/rand"

type SensorData struct {
	Pitch     float64
	Roll      float64
	Yaw       float64
	Altitude  float64
	EngineRPM float64
}

type MovingAverage struct {
	values []float64
	size   int
	index  int
	sum    float64
}

func NewMovingAverage(size int) *MovingAverage {
	return &MovingAverage{
		values: make([]float64, size),
		size:   size,
	}
}

func (ma *MovingAverage) Add(value float64) float64 {
	ma.sum -= ma.values[ma.index]
	ma.values[ma.index] = value
	ma.sum += value
	ma.index = (ma.index + 1) % ma.size
	return ma.sum / float64(ma.size)
}

var pitchMA = NewMovingAverage(10)
var rollMA = NewMovingAverage(10)
var yawMA = NewMovingAverage(10)
var altitudeMA = NewMovingAverage(10)
var engineRPMMA = NewMovingAverage(10)

func InitSensors() error {
	// Initialize actual sensors here
	return nil
}

func ReadSensorData() SensorData {
	return SensorData{
		Pitch:     pitchMA.Add(rand.Float64() * 10),
		Roll:      rollMA.Add(rand.Float64() * 10),
		Yaw:       yawMA.Add(rand.Float64() * 10),
		Altitude:  altitudeMA.Add(rand.Float64() * 100),
		EngineRPM: engineRPMMA.Add(rand.Float64() * 5000),
	}
}
