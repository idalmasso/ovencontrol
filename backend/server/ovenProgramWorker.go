package server

import (
	"sync"
	"time"

	"github.com/idalmasso/ovencontrol/backend/config"
)

type Oven interface {
	GetTemperature() float64
	GetPercentual() float64
	GetMaxPower() float64
	SetPercentual(float64)
	InitStartProgram()
	EndProgram()
}
type OvenProgramWorker struct {
	timeSeconds                        float64
	oven                               Oven
	mu                                 *sync.Mutex
	isWorking                          bool
	kpRamp, kiRamp, kdRamp             float64
	kpMaintain, kiMaintain, kdMaintain float64
	stepTime, stepSave                 float64
	actualDesiredProgramTemperature    float64
}

func (d *OvenProgramWorker) followOvenProgram(program config.OvenProgram) {
	d.mu.Lock()
	if d.isWorking {
		return
	}
	d.isWorking = true
	d.mu.Unlock()
	go func(program config.OvenProgram) {
		d.timeSeconds = 0
		defer func() {
			d.mu.Lock()
			d.isWorking = false
			d.mu.Unlock()
		}()

		if len(program.Points) == 0 {
			return
		}
		d.oven.InitStartProgram()
		defer func() { d.oven.EndProgram() }()
		firstPoint := program.Points[0]
		if firstPoint.Temperature > d.oven.GetTemperature() {
			d.doRampUp(firstPoint)
		} else if firstPoint.Temperature == d.oven.GetTemperature() {
			d.maintainTemperature(firstPoint)
		} else {
			d.doRampDown(firstPoint)
		}
		lastTemp := firstPoint.Temperature
		for _, s := range program.Points[1:] {
			if s.Temperature > lastTemp {
				d.doRampUp(s)
			} else if s.Temperature == lastTemp {
				d.maintainTemperature(s)
			} else {
				d.doRampDown(s)
			}
			lastTemp = s.Temperature
		}

	}(program)
}

func (d *OvenProgramWorker) doRampUp(s config.StepPoint) {
	kp, ki, kd := 0.9, 0.001, 0.001
	d.actualDesiredProgramTemperature = d.oven.GetTemperature()
	integral, previousError, derivative, temperatureVariance := 0.0, 0.0, 0.0, 0.0
	desiredVariance := (s.Temperature - d.oven.GetTemperature()) / s.TimeSeconds()
	ovenTemperature := d.oven.GetTemperature()
	timeSave := 0.0
	for ovenTemperature < s.Temperature {
		d.timeSeconds += d.stepTime
		timeSave += d.stepTime
		newTemperature := d.oven.GetTemperature()
		temperatureVariance = newTemperature - ovenTemperature
		ovenTemperature = newTemperature
		expectedVariance := desiredVariance * d.stepTime
		d.actualDesiredProgramTemperature += expectedVariance
		if d.actualDesiredProgramTemperature > s.Temperature {
			d.actualDesiredProgramTemperature = s.Temperature
		}
		errorValue := expectedVariance - temperatureVariance

		integral = integral + errorValue*d.stepTime
		derivative = (errorValue - previousError) / d.stepTime
		actualPercentual := d.oven.GetPercentual()
		actualPercentual += kp*errorValue + ki*integral + kd*derivative
		actualPercentual = min(actualPercentual, 1)
		actualPercentual = max(actualPercentual, 0)
		d.oven.SetPercentual(actualPercentual)
		previousError = errorValue
		if timeSave > d.stepSave {
			d.Save()
		}
		time.Sleep(time.Duration(d.stepTime) * time.Second)
	}
}

func (d *OvenProgramWorker) maintainTemperature(s config.StepPoint) {
	kp, ki, kd := 0.01, 0.0001, 0.0001
	d.actualDesiredProgramTemperature = d.oven.GetTemperature()
	integral, previousError, derivative := 0.0, 0.0, 0.0
	d.actualDesiredProgramTemperature = s.Temperature
	timeSave := 0.0
	first := true
	for d.oven.GetTemperature() < s.Temperature {
		d.timeSeconds += d.stepTime
		timeSave += d.stepTime
		errorValue := s.Temperature - d.oven.GetTemperature()
		if first {
			first = false
			previousError = errorValue
		}
		integral = integral + errorValue*d.stepTime
		derivative = (errorValue - previousError) / d.stepTime
		actualPercentual := kp*errorValue + ki*integral + kd*derivative
		actualPercentual = min(actualPercentual, 1)
		actualPercentual = max(actualPercentual, 0)
		d.oven.SetPercentual(actualPercentual)
		previousError = errorValue
		if timeSave > d.stepSave {
			d.Save()
		}
		time.Sleep(time.Duration(d.stepTime) * time.Second)
	}
}
func (d *OvenProgramWorker) doRampDown(s config.StepPoint) {
	kp, ki, kd := 0.9, 0.001, 0.001
	d.actualDesiredProgramTemperature = d.oven.GetTemperature()
	integral, previousError, derivative, temperatureVariance := 0.0, 0.0, 0.0, 0.0
	desiredVariance := (s.Temperature - d.oven.GetTemperature()) / s.TimeSeconds()
	ovenTemperature := d.oven.GetTemperature()
	timeSave := 0.0
	for d.oven.GetTemperature() > s.Temperature {
		d.timeSeconds += d.stepTime
		timeSave += d.stepTime
		newTemperature := d.oven.GetTemperature()
		temperatureVariance = newTemperature - ovenTemperature
		ovenTemperature = newTemperature
		expectedVariance := desiredVariance * d.stepTime
		d.actualDesiredProgramTemperature += expectedVariance
		if d.actualDesiredProgramTemperature > s.Temperature {
			d.actualDesiredProgramTemperature = s.Temperature
		}
		errorValue := expectedVariance - temperatureVariance

		integral = integral + errorValue*d.stepTime
		derivative = (errorValue - previousError) / d.stepTime
		actualPercentual := d.oven.GetPercentual()
		actualPercentual += kp*errorValue + ki*integral + kd*derivative
		actualPercentual = min(actualPercentual, 1)
		actualPercentual = max(actualPercentual, 0)
		d.oven.SetPercentual(actualPercentual)
		previousError = errorValue
		if timeSave > d.stepSave {
			d.Save()
		}
		time.Sleep(time.Duration(d.stepTime) * time.Second)
	}
}
func (d *OvenProgramWorker) Save() {

}

func NewOvenProgramWorker(oven Oven, config config.Config) *OvenProgramWorker {
	o := OvenProgramWorker{oven: oven}
	o.mu = &sync.Mutex{}
	o.isWorking = false
	o.stepTime = 1
	o.stepSave = 5
	return &o
}
