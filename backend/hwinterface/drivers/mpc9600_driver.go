package drivers

import (
	"time"

	"gobot.io/x/gobot/v2/drivers/i2c"
)

const mpc9600DefaultAddress = 0x40

// Created to work with 16bit
type MPC9600Driver struct {
	*Driver
}

var mpc9600ChannelSelection = map[int]uint16{
	0: 0x08,
	1: 0x0C,
	2: 0x09,
	3: 0x0D,
	4: 0x0A,
	5: 0x0E,
	6: 0x0B,
	7: 0x0F,
}

// NewADS1015Driver creates a new driver for the MPC9600 (8-bit ADC)
func NewMPC9600Driver(a Connector, options ...func(i2c.Config)) *MPC9600Driver {

	d := &MPC9600Driver{
		Driver: NewDriver(a, "MPC9600", mpc9600DefaultAddress),
	}

	for _, option := range options {
		option(d)
	}

	return d
}

// ReadTemperature returns value from analog reading of specified pin using the default values.
func (d *MPC9600Driver) ReadTemperature() (value int, err error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	value, err = d.rawRead()

	return
}

func (d *MPC9600Driver) rawRead() (data int, err error) {
	var config uint16
	// Go out of power-down mode for conversion.
	config = 0x84
	// Specify mux value.
	//mux := mpc9600ChannelSelection[channel]
	//config |= (((channel<<2 | channel>>1) & 0x07) << 4)

	// Send the config value to start the ADC conversion.
	if err = d.writeWordBigEndian(0x00, config); err != nil {
		return
	}

	// Wait for the ADC sample to finish based on the sample rate plus a
	// small offset to be sure (0.1 millisecond).
	delay := time.Duration(10000) * time.Microsecond
	time.Sleep(delay)
	// Retrieve the result.
	udata, err := d.readWordBigEndian(0x00)
	if err != nil {
		return
	}

	// Handle negative values as two's complement
	return int(twosComplement16Bit(udata)), nil
}

func (d *MPC9600Driver) writeWordBigEndian(reg uint8, val uint16) error {
	return d.connection.WriteByte(byte(val))
}

func (d *MPC9600Driver) readWordBigEndian(reg uint8) (data uint16, err error) {
	if data, err := d.connection.ReadByte(); err != nil {
		return 0, err
	} else {
		return uint16(data), err
	}
}
