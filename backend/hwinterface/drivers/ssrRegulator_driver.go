package drivers

import (
	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
)

// SSRRegulatorDriver represents a ssr regulator
type SSRRegulatorDriver struct {
	pin        string
	name       string
	connection gpio.DigitalWriter
	high       bool
	gobot.Commander
}

// NewSSRRegulator return a new SSRRegulator given a DigitalWriter and pin.
// It is actually the same as a led used in PWM mode, nothing more, then the circuit should implement a conversion of voltage required (in my case, a RCRC with opAmp does the trick)
// Adds the following API Commands:
//
//	"Brightness" - See SSRRegulator.Brightness
//	"Off" - See SSRRegulator.Off
func NewSSRRegulator(a gpio.DigitalWriter, pin string) *SSRRegulatorDriver {
	l := &SSRRegulatorDriver{
		name:       gobot.DefaultName("SSRRegulator"),
		pin:        pin,
		connection: a,
		high:       false,
		Commander:  gobot.NewCommander(),
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

// Start implements the Driver interface
func (l *SSRRegulatorDriver) Start() (err error) { return }

// Halt implements the Driver interface
func (l *SSRRegulatorDriver) Halt() (err error) { return }

// Name returns the SSRRegulators name
func (l *SSRRegulatorDriver) Name() string { return l.name }

// SetName sets the SSRRegulators name
func (l *SSRRegulatorDriver) SetName(n string) { l.name = n }

// Pin returns the SSRRegulators pin
func (l *SSRRegulatorDriver) Pin() string { return l.pin }

// Connection returns the SSRRegulators Connection
func (l *SSRRegulatorDriver) Connection() gobot.Connection {
	return l.connection.(gobot.Connection)
}

// Off sets the regulator to 0.
func (l *SSRRegulatorDriver) Off() (err error) {
	if err = l.connection.DigitalWrite(l.Pin(), 0); err != nil {
		return
	}
	l.high = false
	return
}

// SetPower sets the ssrRegulator to the specified level of power-voltage
func (l *SSRRegulatorDriver) SetPower(level byte) (err error) {
	if writer, ok := l.connection.(gpio.PwmWriter); ok {
		return writer.PwmWrite(l.Pin(), level)
	}
	return gpio.ErrPwmWriteUnsupported
}
