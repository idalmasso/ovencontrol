package spi

import (
	"sync"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/spi"
)

const (
	// NotInitialized is the initial value for a bus/chip
	NotInitialized = -1
)

// Driver implements the interface gobot.Driver for SPI devices.
type Driver struct {
	name       string
	connector  spi.Connector
	connection spi.Connection
	afterStart func() error
	beforeHalt func() error
	spi.Config
	gobot.Commander
	mutex sync.Mutex
}

// NewDriver creates a new generic and basic SPI gobot driver.
func NewDriver(a spi.Connector, name string, options ...func(spi.Config)) *Driver {
	d := &Driver{
		name:       gobot.DefaultName(name),
		connector:  a,
		afterStart: func() error { return nil },
		beforeHalt: func() error { return nil },
		Config:     spi.NewConfig(),
		Commander:  gobot.NewCommander(),
	}
	for _, option := range options {
		option(d)
	}
	return d
}

// Name returns the name of the device.
func (d *Driver) Name() string { return d.name }

// SetName sets the name of the device.
func (d *Driver) SetName(n string) { d.name = n }

// Connection returns the Connection of the device.
func (d *Driver) Connection() gobot.Connection { return d.connector.(gobot.Connection) }

// Start initializes the driver.
func (d *Driver) Start() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	bus := d.GetBusNumberOrDefault(d.connector.SpiDefaultBusNumber())
	chip := d.GetChipNumberOrDefault(d.connector.SpiDefaultChipNumber())
	mode := d.GetModeOrDefault(d.connector.SpiDefaultMode())
	bits := d.GetBitCountOrDefault(d.connector.SpiDefaultBitCount())
	maxSpeed := d.GetSpeedOrDefault(d.connector.SpiDefaultMaxSpeed())

	var err error
	d.connection, err = d.connector.GetSpiConnection(bus, chip, mode, bits, maxSpeed)
	if err != nil {
		return err
	}
	return d.afterStart()
}

// Halt stops the driver.
func (d *Driver) Halt() (err error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if err := d.beforeHalt(); err != nil {
		return err
	}

	// currently there is nothing to do here for the driver, the connection is cached on adaptor side
	// and will be closed on adaptor Finalize()
	return nil
}
