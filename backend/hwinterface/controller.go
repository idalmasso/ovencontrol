package hwinterface

import (
	"encoding/binary"

	"github.com/idalmasso/ovencontrol/backend/commoninterface"
	"github.com/idalmasso/ovencontrol/backend/config"
	"github.com/idalmasso/ovencontrol/backend/hwinterface/drivers"
	"github.com/idalmasso/ovencontrol/backend/hwinterface/drivers/spi"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/raspi"
)

type piController struct {
	//buttonInput         *gpio.ButtonDriver
	ledOvenWorking      *gpio.LedDriver
	ssrPowerController  *drivers.SSRRegulatorDriver
	analogInput         *spi.MAX31856Driver
	buttonPressFunc     func()
	actualPercentual    float64
	maxPower            float64
	internalArea        float64
	thermalCapacity     float64
	thermalConductivity float64
	weight              float64
	insulationWidth     float64
	logger              commoninterface.Logger
}

func (c *piController) GetPercentual() float64 {
	return c.actualPercentual
}
func (c *piController) GetMaxPower() float64 {
	return c.maxPower
}
func (c *piController) SetPercentual(f float64) {
	c.actualPercentual = f
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], uint64(f*255))
	c.ssrPowerController.SetPower(b[0])
}
func (c *piController) InitStartProgram() {
	c.logger.Info("Init start program")
	c.ledOvenWorking.On()
}
func (c *piController) EndProgram() {
	c.logger.Info("End program")
	c.ledOvenWorking.Off()
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

func (d *piController) SetLogger(logger commoninterface.Logger) {
	d.logger = logger
	d.analogInput.SetLogger(logger)
}

func WithLogger(logger commoninterface.Logger) func(*piController) {
	return func(pc *piController) {
		pc.SetLogger(logger)
	}
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
func NewController(options ...func(*piController)) *piController {
	r := raspi.NewAdaptor()
	r.Connect()

	//buttonInput := gpio.NewButtonDriver(r, "15", time.Duration(10*time.Millisecond))
	//buttonInput.Start()
	//This is showing the server is on, I don't need to pass to the piController
	ledOk := gpio.NewLedDriver(r, "16")
	ledOk.Start()

	ledOvenWorking := gpio.NewLedDriver(r, "18")
	ledOvenWorking.Start()

	ssrRegulator := drivers.NewSSRRegulator(r, "11")
	ssrRegulator.Start()

	analogInput := spi.NewMAX31856Driver(r)
	analogInput.Start()
	pi := &piController{analogInput: analogInput, ssrPowerController: ssrRegulator}
	for _, o := range options {
		o(pi)
	}
	ledOk.On()
	ledOvenWorking.Off()
	pi.ssrPowerController.SetPower(0)
	//buttonInput.On(gpio.ButtonRelease, pi.buttonPressed)
	return pi
}
