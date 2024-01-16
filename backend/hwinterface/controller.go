package hwinterface

import (
	"log/slog"
	"sync"
	"time"

	"github.com/golang/glog"
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
	slog.Info("Init start program")
}
func (c *piController) EndProgram() {
	slog.Info("End program")
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

// StartProcess actually starts the real process of making photo 360
/*func (c *piController) StartProcess() error {
	if glog.V(3) {
		glog.Infoln("piController - StartProcess called")
	}
	if !c.canSetStartProcess() {
		return ProcessingError{Operation: "Start Process"}
	}
	//go func() {
	defer c.setProcessing(false)

	t := time.Now()
	c.actualProcessName = fmt.Sprintf("%04d%02d%02d%02d%02d%02d", t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second())

	for temp := 0; temp < 1280; temp += 50 {
		value, err := c.analogInput.AnalogRead()
		if err != nil {
			glog.Errorln("ERR", err)
		}
		ovenTemperature := float64(value) / 255 * 1350
		rampTemperature := float64(temp)
		glog.Infoln("Forno temperatura: ", ovenTemperature)
		glog.Infoln("Rampa temperatura: ", rampTemperature)
		if rampTemperature > ovenTemperature {
			c.ledPower.Brightness(byte(255))
		} else {
			c.ledPower.Brightness(byte(0))
		}

		time.Sleep(time.Second)
	}
	//}()

	return nil
}
*/
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
		glog.Errorln(err)
	}

	analogInput := spi.NewMAX31856Driver(r, slog.Default())
	analogInput.Start()
	pi := piController{buttonInput: buttonInput, analogInput: analogInput, ledPower: ledPower}
	buttonInput.On(gpio.ButtonRelease, pi.buttonPressed)
	return &pi
}
