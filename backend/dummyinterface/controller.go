package dummyinterface

import (
	"math"
	"sync"
	"time"

	"github.com/idalmasso/ovencontrol/backend/config"
)

type DummyController struct {
	ovenTemperature, externalTemperature                                                  float64
	insulationWidth, thermalConductivity, internalArea, thermalCapacity, weight, maxPower float64
	timeSeconds                                                                           float64
	actualPercentual                                                                      float64
	actualDesiredProgramTemperature                                                       float64
	timeMultiplier                                                                        float64
	isWorking                                                                             bool
	mu                                                                                    *sync.Mutex
}

func (d DummyController) GetTemperature() float64 {
	return math.Round(d.ovenTemperature*100) / 100
}

func (d DummyController) GetTemperatureExpected() float64 {
	return math.Round(d.actualDesiredProgramTemperature*100) / 100
}

func (d *DummyController) IsWorking() bool {
	return d.isWorking
}

func (d *DummyController) TryStartTestRamp(temperature, timeMinutes float64) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.isWorking {
		return false
	}
	d.isWorking = true

	s := config.StepPoint{Temperature: temperature, TimeMinutes: timeMinutes}
	program := config.OvenProgram{Name: "TestProgram", Points: []config.StepPoint{s}}
	d.followOvenProgram(program)
	return true
}

func (d *DummyController) followOvenProgram(program config.OvenProgram) {
	go func(program config.OvenProgram) {
		d.timeSeconds = 0
		d.externalTemperature = 25
		d.ovenTemperature = 25
		d.actualPercentual = 0
		defer func() {
			d.mu.Lock()
			d.isWorking = false
			d.mu.Unlock()
		}()
		if len(program.Points) == 0 {
			return
		}
		firstPoint := program.Points[0]
		if firstPoint.Temperature > d.ovenTemperature {
			d.doRampUp(firstPoint)
		} else if firstPoint.Temperature == d.ovenTemperature {
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

func (d *DummyController) doRampUp(s config.StepPoint) {
	kp, ki, kd := 0.9, 0.001, 0.001
	d.actualDesiredProgramTemperature = d.ovenTemperature
	integral, previousError, derivative, lostPower := 0.0, 0.0, 0.0, 0.0
	desiredVariance := (s.Temperature - d.ovenTemperature) / s.TimeSeconds()
	for d.ovenTemperature < s.Temperature {
		dt := 1.0
		realDT := dt * d.timeMultiplier
		d.timeSeconds += realDT
		ovenPower := d.actualPercentual * d.maxPower

		lostPower = (d.thermalConductivity * d.internalArea *
			(d.ovenTemperature - d.externalTemperature)) / d.insulationWidth

		temperatureVariance := ((ovenPower - lostPower) / (d.weight * d.thermalCapacity)) * realDT
		d.ovenTemperature += temperatureVariance
		expectedVariance := desiredVariance * realDT
		d.actualDesiredProgramTemperature += expectedVariance
		if d.actualDesiredProgramTemperature > s.Temperature {
			d.actualDesiredProgramTemperature = s.Temperature
		}
		errorValue := expectedVariance - temperatureVariance

		integral = integral + errorValue*realDT
		derivative = (errorValue - previousError) / realDT
		d.actualPercentual += kp*errorValue + ki*integral + kd*derivative
		d.actualPercentual = min(d.actualPercentual, 1)
		d.actualPercentual = max(d.actualPercentual, 0)

		previousError = errorValue
		time.Sleep(time.Duration(dt) * time.Second)
	}
}
func (d *DummyController) maintainTemperature(s config.StepPoint) {
	kp, ki, kd := 0.01, 0.0001, 0.0001
	integral, previousError, derivative, lostPower := 0.0, 0.0, 0.0, 0.0
	lostPower = (d.thermalConductivity * d.internalArea *
		(d.ovenTemperature - d.externalTemperature)) / d.insulationWidth
	wantedPerc := lostPower / d.maxPower
	integral = wantedPerc / ki
	first := true
	d.actualDesiredProgramTemperature = s.Temperature
	for d.timeSeconds < s.TimeSeconds() {
		dt := 1.0
		realDT := dt * d.timeMultiplier
		d.timeSeconds += realDT
		ovenPower := d.actualPercentual * d.maxPower

		lostPower = (d.thermalConductivity * d.internalArea *
			(d.ovenTemperature - d.externalTemperature)) / d.insulationWidth

		temperatureVariance := ((ovenPower - lostPower) / (d.weight * d.thermalCapacity)) * realDT
		d.ovenTemperature += temperatureVariance

		errorValue := s.Temperature - d.ovenTemperature
		if first {
			first = false
			previousError = errorValue
		}
		integral = integral + errorValue*realDT
		derivative = (errorValue - previousError) / realDT
		d.actualPercentual = kp*errorValue + ki*integral + kd*derivative
		d.actualPercentual = min(d.actualPercentual, 1)
		d.actualPercentual = max(d.actualPercentual, 0)

		previousError = errorValue
		time.Sleep(time.Duration(dt) * time.Second)

	}
}
func (d *DummyController) doRampDown(s config.StepPoint) {
	kp, ki, kd := 0.9, 0.001, 0.001
	d.actualDesiredProgramTemperature = d.ovenTemperature
	integral, previousError, derivative, lostPower := 0.0, 0.0, 0.0, 0.0
	desiredVariance := (s.Temperature - d.ovenTemperature) / s.TimeSeconds()
	for d.ovenTemperature > s.Temperature {
		dt := 1.0
		realDT := dt * d.timeMultiplier
		d.timeSeconds += realDT
		ovenPower := d.actualPercentual * d.maxPower

		lostPower = (d.thermalConductivity * d.internalArea *
			(d.ovenTemperature - d.externalTemperature)) / d.insulationWidth

		temperatureVariance := ((ovenPower - lostPower) / (d.weight * d.thermalCapacity)) * realDT
		d.ovenTemperature += temperatureVariance
		expectedVariance := desiredVariance * realDT
		d.actualDesiredProgramTemperature += expectedVariance
		if d.actualDesiredProgramTemperature > s.Temperature {
			d.actualDesiredProgramTemperature = s.Temperature
		}
		errorValue := expectedVariance - temperatureVariance

		integral = integral + errorValue*realDT
		derivative = (errorValue - previousError) / realDT
		d.actualPercentual += kp*errorValue + ki*integral + kd*derivative
		d.actualPercentual = min(d.actualPercentual, 1)
		d.actualPercentual = max(d.actualPercentual, 0)

		previousError = errorValue
		time.Sleep(time.Duration(dt) * time.Second)

	}
}

func (d *DummyController) InitConfig(c config.Config) {
	d.actualDesiredProgramTemperature = 25
	d.actualPercentual = 0
	d.externalTemperature = 25
	for _, v := range c.Oven.InsultationWidths {
		d.insulationWidth += v
	}
	d.mu = &sync.Mutex{}
	d.internalArea = c.Oven.Height * c.Oven.Length * c.Oven.Width
	d.maxPower = c.Oven.MaxPower
	d.ovenTemperature = 25
	d.thermalCapacity = c.Oven.ThermalCapacity
	d.thermalConductivity = calculateConducibility(c.Oven.InsultationWidths, c.Oven.ThermalConductivities)
	d.timeMultiplier = 10
	d.weight = c.Oven.Weight
}

func calculateConducibility(lengths, conducibilities []float64) float64 {
	total := 0.0
	for _, l := range lengths {
		total += l
	}

	rTot := 0.0
	for idx := range lengths {
		rTot += lengths[idx] / conducibilities[idx]
	}

	return total / rTot
}
