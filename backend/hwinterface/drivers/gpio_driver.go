package drivers

import (
	"fmt"
	"log"
	"sync"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
)

// configuration contains all changeable attributes of the driver.
type configuration struct {
	name string
	pin  string
}

// optionApplier needs to be implemented by each configurable option type
type optionApplier interface {
	apply(cfg *configuration)
}

// nameOption is the type for applying another name to the configuration
type nameOption string

// pinOption is the type for applying a pin to the configuration
type pinOption string

// Driver implements the interface gobot.Driver.
type driver struct {
	driverCfg  *configuration
	connection gobot.Adaptor
	afterStart func() error
	beforeHalt func() error
	gobot.Commander
	mutex *sync.Mutex // mutex often needed to ensure that write-read sequences are not interrupted
}

// newDriver creates a new generic and basic gpio gobot driver.
//
// Supported options:
//
//	"WithName"
//	"withPin"
func newDriver(a gobot.Adaptor, name string, opts ...interface{}) *driver {
	d := &driver{
		driverCfg:  &configuration{name: gobot.DefaultName(name)},
		connection: a,
		afterStart: func() error { return nil },
		beforeHalt: func() error { return nil },
		Commander:  gobot.NewCommander(),
		mutex:      &sync.Mutex{},
	}

	for _, opt := range opts {
		switch o := opt.(type) {
		case optionApplier:
			o.apply(d.driverCfg)
		default:
			panic(fmt.Sprintf("'%s' can not be applied on '%s'", opt, d.driverCfg.name))
		}
	}

	return d
}

// WithName is used to replace the default name of the driver.
func WithName(name string) optionApplier {
	return nameOption(name)
}

// withPin is used to add a pin to the driver. Only one pin can be linked.
// This option is not available outside gpio package.
func withPin(pin string) optionApplier {
	return pinOption(pin)
}

// Name returns the name of the gpio device.
func (d *driver) Name() string {
	return d.driverCfg.name
}

// SetName sets the name of the gpio device.
// Deprecated: Please use option [gpio.WithName] instead.
func (d *driver) SetName(name string) {
	WithName(name).apply(d.driverCfg)
}

// Pin returns the pin associated with the driver.
func (d *driver) Pin() string {
	return d.driverCfg.pin
}

// Connection returns the connection of the gpio device.
func (d *driver) Connection() gobot.Connection {
	if conn, ok := d.connection.(gobot.Connection); ok {
		return conn
	}

	log.Printf("%s has no gobot connection\n", d.driverCfg.name)
	return nil
}

// Start initializes the gpio device.
func (d *driver) Start() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// currently there is nothing to do here for the driver

	return d.afterStart()
}

// Halt halts the gpio device.
func (d *driver) Halt() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// currently there is nothing to do after halt for the driver

	return d.beforeHalt()
}

// digitalRead is a helper function with check that the connection implements DigitalReader
func (d *driver) digitalRead(pin string) (int, error) {
	if reader, ok := d.connection.(gpio.DigitalReader); ok {
		return reader.DigitalRead(pin)
	}

	return 0, gpio.ErrDigitalReadUnsupported
}

// digitalWrite is a helper function with check that the connection implements DigitalWriter
func (d *driver) digitalWrite(pin string, val byte) error {
	if writer, ok := d.connection.(gpio.DigitalWriter); ok {
		return writer.DigitalWrite(pin, val)
	}

	return gpio.ErrDigitalWriteUnsupported
}

// pwmWrite is a helper function with check that the connection implements PwmWriter
func (d *driver) pwmWrite(pin string, level byte) error {
	if writer, ok := d.connection.(gpio.PwmWriter); ok {
		return writer.PwmWrite(pin, level)
	}

	return gpio.ErrPwmWriteUnsupported
}

// servoWrite is a helper function with check that the connection implements ServoWriter
func (d *driver) servoWrite(pin string, level byte) error {
	if writer, ok := d.connection.(gpio.ServoWriter); ok {
		return writer.ServoWrite(pin, level)
	}

	return gpio.ErrServoWriteUnsupported
}

func (o nameOption) String() string {
	return "name option for digital drivers"
}

func (o pinOption) String() string {
	return "pin option for digital drivers"
}

// apply change the name in the configuration.
func (o nameOption) apply(c *configuration) {
	c.name = string(o)
}

// apply change the pins list of the configuration.
func (o pinOption) apply(c *configuration) {
	c.pin = string(o)
}
