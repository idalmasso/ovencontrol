package dummyinterface

import (
	"math"
	"time"

	"github.com/idalmasso/ovencontrol/backend/commoninterface"
	"github.com/idalmasso/ovencontrol/backend/config"
)

type DummyController struct {
	ovenTemperature, externalTemperature                                                  float64
	insulationWidth, thermalConductivity, internalArea, thermalCapacity, weight, maxPower float64
	actualPercentual                                                                      float64
	timeMultiplier                                                                        float64
	isWorking                                                                             bool
	logger                                                                                commoninterface.Logger
}

func (d DummyController) GetTemperature() (float64, error) {
	return math.Round(d.ovenTemperature*100) / 100, nil
}

func (d *DummyController) IsWorking() bool {
	return d.isWorking
}

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

	go func() {
		for {
			time.Sleep(time.Millisecond * 500)
			ovenPower := d.actualPercentual * d.maxPower
			lostPower := (d.thermalConductivity * d.internalArea *
				(d.ovenTemperature - d.externalTemperature)) / d.insulationWidth
			d.ovenTemperature += ((ovenPower - lostPower) / (d.weight * d.thermalCapacity)) * 0.5
		}
	}()
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
func (d *DummyController) SetPercentual(percent float64) error {
	d.actualPercentual = percent
	return nil
}
func (d *DummyController) SetLogger(logger commoninterface.Logger) {
	d.logger = logger
}
func WithLogger(logger commoninterface.Logger) func(*DummyController) {
	return func(d *DummyController) {
		d.SetLogger(logger)
	}
}
func (d *DummyController) InitStartProgram() error {
	d.ovenTemperature = 0
	return nil
}
func (d *DummyController) EndProgram() error {
	d.isWorking = false
	return nil
}

func NewDummyController(options ...func(*DummyController)) *DummyController {
	d := &DummyController{}
	for _, o := range options {
		o(d)
	}
	return d
}
