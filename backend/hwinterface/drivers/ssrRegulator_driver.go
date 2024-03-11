package drivers

import (
	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
)

// SSRRegulatorDriver represents a ssr regulator
type SSRRegulatorDriver struct {
	*driver
	high bool
	gobot.Commander
}

// NewSSRRegulator return a new SSRRegulator given a DigitalWriter and pin.
// It is actually the same as a led used in PWM mode, nothing more, then the circuit should implement a conversion of voltage required (in my case, a RCRC with opAmp does the trick)
// Adds the following API Commands:
//
//	"Power" - See SSRRegulator.Power
//	"Off" - See SSRRegulator.Off
func NewSSRRegulator(a gpio.DigitalWriter, pin string, opts ...interface{}) *SSRRegulatorDriver {
	l := &SSRRegulatorDriver{
		driver:    newDriver(a.(gobot.Connection), "SSRRegulator", append(opts, withPin(pin))...),
		high:      false,
		Commander: gobot.NewCommander(),
	}

	l.AddCommand("Power", func(params map[string]interface{}) interface{} {
		level := byte(params["level"].(float64))
		return l.SetPower(level)
	})

	l.AddCommand("Off", func(params map[string]interface{}) interface{} {
		return l.Off()
	})

	return l
}

// State return true if the led is On and false if the led is Off
func (d *SSRRegulatorDriver) State() bool {
	return d.high
}

// On sets the led to a high state.
func (d *SSRRegulatorDriver) On() error {
	if err := d.digitalWrite(d.driverCfg.pin, 1); err != nil {
		return err
	}
	d.high = true
	return nil
}

// Off sets the led to a low state.
func (d *SSRRegulatorDriver) Off() error {
	if err := d.digitalWrite(d.driverCfg.pin, 0); err != nil {
		return err
	}
	d.high = false
	return nil
}

// Toggle sets the led to the opposite of it's current state
func (d *SSRRegulatorDriver) Toggle() error {
	if d.State() {
		return d.Off()
	}
	return d.On()
}

// Brightness sets the led to the specified level of power
func (d *SSRRegulatorDriver) SetPower(level byte) error {
	return d.pwmWrite(d.driverCfg.pin, level)
}
