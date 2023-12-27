package drivers

import (
	"fmt"
	"strconv"
	"time"

	"gobot.io/x/gobot/v2/drivers/i2c"
)

const ads7830DefaultAddress = 0x4b

// ADS1x15Driver is the Gobot driver for the ADS1015/ADS1115 ADC
// datasheet:
// https://www.ti.com/lit/gpn/ads1115
//
// reference implementations:
// * https://github.com/adafruit/Adafruit_Python_ADS1x15
// * https://github.com/Wh1teRabbitHU/ADS1115-Driver
type ADS7830Driver struct {
	*Driver
}

var ads7830ChannelSelection = map[int]uint16{
	0: 0x08,
	1: 0x0C,
	2: 0x09,
	3: 0x0D,
	4: 0x0A,
	5: 0x0E,
	6: 0x0B,
	7: 0x0F,
}

// NewADS1015Driver creates a new driver for the ADS7830 (8-bit ADC)
func NewADS7830Driver(a Connector, options ...func(i2c.Config)) *ADS7830Driver {

	d := &ADS7830Driver{
		Driver: NewDriver(a, "ADS7830", ads7830DefaultAddress),
	}

	for _, option := range options {
		option(d)
	}

	d.AddCommand("AnalogRead", func(params map[string]interface{}) interface{} {
		pin := params["pin"].(string)
		val, err := d.AnalogRead(pin)
		return map[string]interface{}{"val": val, "err": err}
	})

	return d
}

// AnalogRead returns value from analog reading of specified pin using the default values.
func (d *ADS7830Driver) AnalogRead(pin string) (value int, err error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	var channel uint64
	channel, err = strconv.ParseUint(pin, 10, 16)
	if err != nil {
		return
	}

	if err = d.checkChannel(uint16(channel)); err != nil {
		return
	}

	value, err = d.rawRead(uint16(channel))

	return
}

func (d *ADS7830Driver) rawRead(channel uint16) (data int, err error) {
	// Validate the passed in data rate (differs between ADS1015 and ADS1115).
	var config uint16
	// Go out of power-down mode for conversion.
	config = 0x84
	// Specify mux value.
	//mux := ads7830ChannelSelection[channel]
	config |= (((channel<<2 | channel>>1) & 0x07) << 4)

	// Send the config value to start the ADC conversion.
	if err = d.writeWordBigEndian(0x00, config); err != nil {
		return
	}

	// Wait for the ADC sample to finish based on the sample rate plus a
	// small offset to be sure (0.1 millisecond).
	delay := time.Duration(1000000/100) * time.Microsecond
	time.Sleep(delay)
	// Retrieve the result.
	udata, err := d.readWordBigEndian(0x00)
	if err != nil {
		return
	}

	// Handle negative values as two's complement
	return int(twosComplement16Bit(udata)), nil
}

func (d *ADS7830Driver) checkChannel(channel uint16) (err error) {
	if channel < 0 || channel > 7 {
		err = fmt.Errorf("Invalid channel (%d), must be between 0 and 7", channel)
	}
	return
}

func (d *ADS7830Driver) writeWordBigEndian(reg uint8, val uint16) error {
	return d.connection.WriteByte(byte(val))
}

func (d *ADS7830Driver) readWordBigEndian(reg uint8) (data uint16, err error) {
	if data, err := d.connection.ReadByte(); err != nil {
		return 0, err
	} else {
		return uint16(data), err
	}
}
