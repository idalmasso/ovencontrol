package hwinterface

import (
	"fmt"
	"sync"
	"time"

	"github.com/golang/glog"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)

type piController struct {
	processing        bool
	buttonInput       *gpio.ButtonDriver
	ledOk             *gpio.LedDriver
	ledPower          *gpio.LedDriver
	analogInput       *i2c.ADS1x15Driver
	mutex             sync.RWMutex
	buttonPressFunc   func()
	actualProcessName string
}

// StartProcess actually starts the real process of making photo 360
func (c *piController) StartProcess() error {
	if glog.V(3) {
		glog.Infoln("piController - StartProcess called")
	}
	if !c.canSetStartProcess() {
		return ProcessingError{Operation: "Start Process"}
	}
	go func() {
		defer c.setProcessing(false)

		t := time.Now()
		var i byte
		c.actualProcessName = fmt.Sprintf("%04d%02d%02d%02d%02d%02d", t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second())
		for {
			val, err := c.analogInput.AnalogRead("0")
			if err != nil {
				glog.Errorln("ERR", err)
			}
			glog.Infoln("Value: ", val)
			glog.Infoln("Setting let to Value: ", i)
			c.ledPower.Brightness(i)
			i = (i + 1) % byte(255)
			time.Sleep(time.Second)
		}
	}()

	return nil
}

// StopProcess should stop the process at any time
func (c *piController) StopProcess() error {
	if glog.V(3) {
		glog.Infoln("piController - StopProcess called")
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.processing {
		c.ledOk.On()
		c.processing = false
		c.actualProcessName = ""
	}
	return nil
}

// Return true if the machine is actually doing a work and so can be stopped but cannot start another one
func (c *piController) IsWorking() bool {
	return c.isProcessing()
}

func (c *piController) isProcessing() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.processing
}

func (c *piController) canSetStartProcess() bool {
	if glog.V(4) {
		glog.Infoln("piController -  canSetStartProcess canStartProcess")
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.processing {
		return false
	} else {
		c.ledOk.Off()
		c.processing = true

		if glog.V(4) {
			glog.Infoln("piController - canSetStartProcess start processing")
		}
		return true
	}
}

func (c *piController) GetActualProcessName() string { return c.actualProcessName }
func (c *piController) setProcessing(value bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.processing && !value {
		if glog.V(4) {
			glog.Infoln("piController - setProcessing stop processing")
		}
		c.actualProcessName = ""
	} else if !c.processing && value {
		if glog.V(4) {
			glog.Infoln("piController - setProcessing start processing")
		}
	}
	c.processing = value
	if c.processing {
		c.ledOk.Off()
	} else {
		c.ledOk.On()
	}
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
	ledPower := gpio.NewLedDriver(r, "13")
	ledPower.Start()
	ledPower.Brightness(0)
	analogInput := i2c.NewADS1015Driver(r)
	analogInput.WithAddress(0x42)
	pi := piController{buttonInput: buttonInput, analogInput: analogInput, ledPower: ledPower}
	buttonInput.On(gpio.ButtonRelease, pi.buttonPressed)
	return &pi
}
