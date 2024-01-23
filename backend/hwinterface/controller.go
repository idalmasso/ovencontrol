package hwinterface

import (
	"log/slog"
	"sync"
	"time"

	"github.com/idalmasso/ovencontrol/backend/config"
	"github.com/idalmasso/ovencontrol/backend/hwinterface/drivers/spi"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/raspi"
)

type piController struct {
	buttonInput *gpio.ButtonDriver
	//ledOk             *gpio.LedDriver
	ledPower            *gpio.LedDriver
	analogInput         *spi.MAX31856Driver
	mutex               *sync.RWMutex
	buttonPressFunc     func()
	actualProcessName   string
	actualPercentual    float64
	maxPower            float64
	internalArea        float64
	thermalCapacity     float64
	thermalConductivity float64
	weight              float64
	insulationWidth     float64
	logger              Logger
}
type Logger interface {
	Info(msc string, args ...any)
	spi.Logger
}

func (c *piController) GetPercentual() float64 {
	return c.actualPercentual
}
func (c *piController) GetMaxPower() float64 {
	return c.maxPower
}
func (c *piController) SetPercentual(f float64) {
	c.actualPercentual = f
}
func (c *piController) InitStartProgram() {
	c.logger.Info("Init start program")
}
func (c *piController) EndProgram() {
	c.logger.Info("End program")
}
func (d *piController) InitConfig(c config.Config) {

	d.actualPercentual = 0
	for _, v := range c.Oven.InsultationWidths {
		d.insulationWidth += v
	}
	d.internalArea = c.Oven.Height * c.Oven.Length * c.Oven.Width
	d.maxPower = c.Oven.MaxPower
	d.thermalCapacity = c.Oven.ThermalCapacity
	d.thermalConductivity = calculateConducibility(c.Oven.InsultationWidths, c.Oven.ThermalConductivities)
	d.weight = c.Oven.Weight
}

func (d *piController) SetLogger(logger Logger) {
	d.logger = logger
	d.analogInput.SetLogger(logger)
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

func (c *piController) SetOnButtonPress(callback func()) {
	c.buttonPressFunc = callback
}
func (c *piController) buttonPressed(interface{}) {
	c.buttonPressFunc()
}
func NewController() *piController {
	r := raspi.NewAdaptor()
	r.Connect()

	buttonInput := gpio.NewButtonDriver(r, "15", time.Duration(10*time.Millisecond))
	buttonInput.Start()
	//ledOk := gpio.NewLedDriver(r, "13")
	//ledOk.Start()
	//ledOk.On()
	ledPower := gpio.NewLedDriver(r, "11")
	ledPower.Start()

	err := ledPower.Brightness(0)
	if err != nil {
		slog.Error("setLedPowerBrigthness", "error", err)
	}

	analogInput := spi.NewMAX31856Driver(r)
	analogInput.Start()
	pi := piController{buttonInput: buttonInput, analogInput: analogInput, ledPower: ledPower}
	buttonInput.On(gpio.ButtonRelease, pi.buttonPressed)
	return &pi
}
