package dummyinterface

import (
	"math"
	"time"

	"github.com/idalmasso/ovencontrol/backend/config"
)

type DummyController struct {
	ovenTemperature, externalTemperature                                                  float64
	insulationWidth, thermalConductivity, internalArea, thermalCapacity, weight, maxPower float64
	actualPercentual                                                                      float64
	timeMultiplier                                                                        float64
	isWorking                                                                             bool
}

func (d DummyController) GetTemperature() float64 {
	return math.Round(d.ovenTemperature*100) / 100
}

func (d *DummyController) IsWorking() bool {
	return d.isWorking
}

/*
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
*/
func (d *DummyController) InitConfig(c config.Config) {

	d.actualPercentual = 0
	d.externalTemperature = 25
	for _, v := range c.Oven.InsultationWidths {
		d.insulationWidth += v
	}
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

func (d *DummyController) GetPercentual() float64 {
	return d.actualPercentual
}
func (d *DummyController) GetMaxPower() float64 {
	return d.maxPower
}
func (d *DummyController) SetPercentual(percent float64) {
	d.actualPercentual = percent
}
func (d *DummyController) InitStartProgram() {
	if !d.isWorking {
		d.ovenTemperature = 0
		go func() {
			d.isWorking = true
			for d.isWorking {
				time.Sleep(time.Millisecond * 500)
				ovenPower := d.actualPercentual * d.maxPower
				lostPower := (d.thermalConductivity * d.internalArea *
					(d.ovenTemperature - d.externalTemperature)) / d.insulationWidth
				d.ovenTemperature += ((ovenPower - lostPower) / (d.weight * d.thermalCapacity)) * 0.5
			}
		}()
	}
}
func (d *DummyController) EndProgram() {
	d.isWorking = false
}
