package drivers

import (
	"log/slog"
	"time"

	"gobot.io/x/gobot/v2/drivers/i2c"
)

// Created to work with 16bit
type MPC9600Driver struct {
	*Driver
}

const (
	ADDRESS_DEVICE       = 0x60
	DEV_ADDR             = 0x40
	HOT_JUNC_TEMP        = 0x00
	DELTA_JUNC_TEMP      = 0x01
	COLD_JUNC_TEMP       = 0x02
	RAW_ADC              = 0x03
	SENSOR_STATUS        = 0x04
	THERMO_SENSOR_CONFIG = 0x05
	DEVICE_CONFIG        = 0x06
	ALERT1_CONFIG        = 0x08
	ALERT2_CONFIG        = 0x09
	ALERT3_CONFIG        = 0x0A
	ALERT4_CONFIG        = 0x0B
	ALERT1_HYSTERESIS    = 0x0C
	ALERT2_HYSTERESIS    = 0x0D
	ALERT3_HYSTERESIS    = 0x0E
	ALERT4_HYSTERESIS    = 0x0F
	ALERT1_LIMIT         = 0x10
	ALERT2_LIMIT         = 0x11
	ALERT3_LIMIT         = 0x12
	ALERT4_LIMIT         = 0x13
	DEVICE_ID            = 0x20
	TYPE_K               = 0b000
	TYPE_J               = 0b001
	TYPE_T               = 0b010
	TYPE_N               = 0b011
	TYPE_S               = 0b100
	TYPE_E               = 0b101
	TYPE_B               = 0b110
	TYPE_R               = 0b111
	RES_18_BIT           = 0b00
	RES_16_BIT           = 0b01
	RES_14_BIT           = 0b10
	RES_12_BIT           = 0b11
	NORMAL               = 0x00
	SHUTDOWN             = 0x01
	BURST                = 0x02
)

// NewADS1015Driver creates a new driver for the MPC9600 (8-bit ADC)
func NewMPC9600Driver(a Connector, options ...func(i2c.Config)) *MPC9600Driver {

	d := &MPC9600Driver{
		Driver: NewDriver(a, "MPC9600", ADDRESS_DEVICE),
	}

	for _, option := range options {
		option(d)
	}

	return d
}

// ReadTemperature returns value from analog reading of specified pin using the default values.
func (d *MPC9600Driver) AnalogRead() (value int, err error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	value, err = d.rawRead()

	return
}

func (d *MPC9600Driver) rawRead() (data int, err error) {
	// Send the config value to start the ADC conversion.
	if err = d.writeWordBigEndian(HOT_JUNC_TEMP); err != nil {
		slog.Error("Error writing HOT JUNC")
		return
	}

	// Wait for the ADC sample to finish based on the sample rate plus a
	// small offset to be sure (0.1 millisecond).
	delay := time.Duration(10000) * time.Microsecond
	time.Sleep(delay)

	// Retrieve the result.
	udata, err := d.readWordBigEndian()
	if err != nil {
		slog.Error("Error reading HOT JUNC")
		return
	}
	slog.Warn("READ HOT JUNC %d", udata)
	// Handle negative values as two's complement
	return int(twosComplement16Bit(udata)), nil
}

func (d *MPC9600Driver) writeWordBigEndian(val uint16) error {
	return d.connection.WriteByte(byte(val))
}

func (d *MPC9600Driver) readWordBigEndian() (data uint16, err error) {
	if data, err := d.connection.ReadByte(); err != nil {
		return 0, err
	} else {
		return uint16(data), err
	}
}
