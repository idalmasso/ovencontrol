package hwinterface

import (
	"github.com/idalmasso/ovencontrol/backend/commoninterface"
	"github.com/idalmasso/ovencontrol/backend/config"
	"github.com/idalmasso/ovencontrol/backend/hwinterface/drivers"
	"github.com/idalmasso/ovencontrol/backend/hwinterface/drivers/spi"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/raspi"
)

type piController struct {
	ledOvenWorking                                        *gpio.LedDriver
	ssrPowerController                                    *drivers.SSRRegulatorDriver
	analogInput                                           *spi.MAX31856Driver
	ovenRelayPower, airCompressorPower, airCompressorOpen *gpio.RelayDriver
	gpio.RelayDriver
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

func NewController(options ...func(*piController)) *piController {
	r := raspi.NewAdaptor()
	r.Connect()

	//This is showing the server is on, I don't need to pass to the piController
	ledOk := gpio.NewLedDriver(r, "16")
	ledOk.Start()
	ledOk.On()
	ovenRelayPower := gpio.NewRelayDriver(r, "37", gpio.WithRelayInverted())
	airCompressorPower := gpio.NewRelayDriver(r, "35", gpio.WithRelayInverted())
	airCompressorOpen := gpio.NewRelayDriver(r, "33", gpio.WithRelayInverted())
	ledOvenWorking := gpio.NewLedDriver(r, "18")
	ssrRegulator := drivers.NewSSRRegulator(r, "38")
	analogInput := spi.NewMAX31856Driver(r, spi.WithAverageSample(4), spi.WithNoiseRejection(50), spi.WithThermocoupleType(spi.S))

	pi := &piController{analogInput: analogInput,
		ssrPowerController: ssrRegulator,
		ledOvenWorking:     ledOvenWorking,
		ovenRelayPower:     ovenRelayPower,
		airCompressorPower: airCompressorPower,
		airCompressorOpen:  airCompressorOpen}
	for _, o := range options {
		o(pi)
	}
	analogInput.Start()
	ledOvenWorking.Start()
	ssrRegulator.Start()
	ovenRelayPower.Start()
	airCompressorPower.Start()
	airCompressorOpen.Start()
	ovenRelayPower.Off()
	ledOvenWorking.Off()
	airCompressorPower.Off()
	airCompressorOpen.On()
	pi.ssrPowerController.SetPower(0)
	return pi
}
