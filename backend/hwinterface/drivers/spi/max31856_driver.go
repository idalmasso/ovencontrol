package spi

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/idalmasso/ovencontrol/backend/commoninterface"
	"gobot.io/x/gobot/v2/drivers/spi"
)

// MAX31856Driver is a driver for the MAX31856 thermocouple reader.
type MAX31856Driver struct {
	*Driver
	thermocoupleType                       ThermocoupleType
	averageSample, noiseRejectionFrequency int
	mu                                     *sync.Mutex
	logger                                 commoninterface.Logger
}

type FaultState struct {
	Highthresh, Lowthresh, Refinlow, Refinhigh, Rtdinlow, Ovuv bool
}

func (d *MAX31856Driver) SetLogger(logger commoninterface.Logger) {
	d.logger = logger
}

const (
	MAX31856_CR0_REG         uint8 = 0x00
	MAX31856_CR0_AUTOCONVERT uint8 = 0x80
	MAX31856_CR0_1SHOT       uint8 = 0x40
	MAX31856_CR0_OCFAULT1    uint8 = 0x20
	MAX31856_CR0_OCFAULT0    uint8 = 0x10
	MAX31856_CR0_CJ          uint8 = 0x08
	MAX31856_CR0_FAULT       uint8 = 0x04
	MAX31856_CR0_FAULTCLR    uint8 = 0x02
	MAX31856_CR0_50HZ        uint8 = 0x01

	MAX31856_CR1_REG    uint8 = 0x01
	MAX31856_MASK_REG   uint8 = 0x02
	MAX31856_CJHF_REG   uint8 = 0x03
	MAX31856_CJLF_REG   uint8 = 0x04
	MAX31856_LTHFTH_REG uint8 = 0x05
	MAX31856_LTHFTL_REG uint8 = 0x06
	MAX31856_LTLFTH_REG uint8 = 0x07
	MAX31856_LTLFTL_REG uint8 = 0x08
	MAX31856_CJTO_REG   uint8 = 0x09
	MAX31856_CJTH_REG   uint8 = 0x0A
	MAX31856_CJTL_REG   uint8 = 0x0B
	MAX31856_LTCBH_REG  uint8 = 0x0C
	MAX31856_LTCBM_REG  uint8 = 0x0D
	MAX31856_LTCBL_REG  uint8 = 0x0E
	MAX31856_SR_REG     uint8 = 0x0F

	MAX31856_FAULT_CJRANGE uint8 = 0x80
	MAX31856_FAULT_TCRANGE uint8 = 0x40
	MAX31856_FAULT_CJHIGH  uint8 = 0x20
	MAX31856_FAULT_CJLOW   uint8 = 0x10
	MAX31856_FAULT_TCHIGH  uint8 = 0x08
	MAX31856_FAULT_TCLOW   uint8 = 0x04
	MAX31856_FAULT_OVUV    uint8 = 0x02
	MAX31856_FAULT_OPEN    uint8 = 0x01
)

var avgSelectionMap = map[int]uint8{1: 0x00, 2: 0x10, 4: 0x20, 8: 0x30, 16: 0x40}

type ThermocoupleType uint8

const (
	B   ThermocoupleType = 0b0000
	E   ThermocoupleType = 0b0001
	J   ThermocoupleType = 0b0010
	K   ThermocoupleType = 0b0011
	N   ThermocoupleType = 0b0100
	R   ThermocoupleType = 0b0101
	S   ThermocoupleType = 0b0110
	T   ThermocoupleType = 0b0111
	G8  ThermocoupleType = 0b1000
	G32 ThermocoupleType = 0b1100
)

// NewMAX31856Driver creates a new Gobot Driver for MAX31856 thermocouple reader
//
// Params:
//
//	a *Adaptor - the Adaptor to use with this Driver
//
// Optional params:
//
//	 spi.WithBusNumber(int):  bus to use with this driver
//		spi.WithChipNumber(int): chip to use with this driver
//	 spi.WithMode(int):    	 mode to use with this driver
//	 spi.WithBitCount(int):   number of bits to use with this driver
//	 spi.WithSpeed(int64):    speed in Hz to use with this driver
func NewMAX31856Driver(a spi.Connector, options ...func(*MAX31856Driver)) *MAX31856Driver {
	d := &MAX31856Driver{
		Driver:                  NewDriver(a, "MAX31856", spi.WithMode(1)),
		thermocoupleType:        K,
		noiseRejectionFrequency: 50,
		averageSample:           1,
		mu:                      &sync.Mutex{},
		logger:                  slog.Default(),
	}

	for _, option := range options {
		option(d)
	}
	d.afterStart = func() error {
		//# assert on any fault
		d.writeUint8(MAX31856_MASK_REG, 0x0)
		//  # configure open circuit faults
		d.writeUint8(MAX31856_CR0_REG, MAX31856_CR0_OCFAULT0)
		d.SetThermocoupleType(d.thermocoupleType)
		d.SetAverageSample(d.averageSample)
		d.SetNoiseRejection(50)
		return nil
	}

	return d
}
func WithLogger(logger commoninterface.Logger) func(*MAX31856Driver) {
	return func(driver *MAX31856Driver) {
		driver.SetLogger(logger)
	}
}
func WithThermocoupleType(t ThermocoupleType) func(*MAX31856Driver) {
	return func(driver *MAX31856Driver) {
		driver.thermocoupleType = t
	}
}
func WithAverageSample(nSamples int) func(*MAX31856Driver) {
	return func(driver *MAX31856Driver) {
		driver.averageSample = nSamples
	}
}
func WithNoiseRejection(frequency int) func(*MAX31856Driver) {
	return func(driver *MAX31856Driver) {
		driver.noiseRejectionFrequency = frequency
	}
}

// SetAverageSample sets the number of samples averaged per read
func (d *MAX31856Driver) SetAverageSample(nSamples int) error {
	avgValue, ok := avgSelectionMap[nSamples]
	if !ok {
		return fmt.Errorf("invalid nsamples")
	}
	if reg1, err := d.readUint8(MAX31856_CR1_REG); err != nil {
		return fmt.Errorf("read error")
	} else {
		reg1 &= 0b10001111
		reg1 |= avgValue
		err = d.writeUint8(MAX31856_CR1_REG, reg1)
		return err
	}
}

func (d *MAX31856Driver) SetThermocoupleType(thermocoupleType ThermocoupleType) error {
	//# get current value of CR1 Reg
	d.thermocoupleType = thermocoupleType
	confReg1, err := d.readUint8(MAX31856_CR1_REG)
	if err != nil {
		return err
	}
	confReg1 &= 0xF0
	//# add the new value for the TC type
	confReg1 |= (uint8(thermocoupleType) & 0x0F)
	return d.writeUint8(MAX31856_CR1_REG, confReg1)
}

// SetNoiseRejection sets the filter (50 or 60Hz)
func (d *MAX31856Driver) SetNoiseRejection(frequency int) error {
	d.logger.Debug("SetNoiseRejection", "frequency", frequency)
	if frequency != 50 && frequency != 60 {
		return fmt.Errorf("invalid frequency")
	}
	if reg0, err := d.readUint8(MAX31856_CR0_REG); err != nil {
		return fmt.Errorf("read error")
	} else {
		if frequency == 50 {
			reg0 |= MAX31856_CR0_50HZ
		} else {
			reg0 &= ^MAX31856_CR0_50HZ
		}

		err = d.writeUint8(MAX31856_CR0_REG, reg0)
		return err
	}

}

// Starts a one-shot measurement and returns immediately.
// A measurement takes approximately 160ms.
// Check the status of the measurement with `oneshot_pending`; when it is false,
func (d *MAX31856Driver) InitOneShotMeasurement() error {
	//read the current value of the first config register
	confReg0, err := d.readUint8(MAX31856_CR0_REG)
	if err != nil {
		return err
	}
	// and the complement to guarantee the autoconvert bit is unset
	confReg0 |= (MAX31856_CR0_AUTOCONVERT)
	// or the oneshot bit to ensure it is set
	confReg0 |= MAX31856_CR0_1SHOT

	// write it back with the new values, prompting the sensor to perform a measurement
	err = d.writeUint8(MAX31856_CR0_REG, confReg0)
	return err
}

// A boolean indicating the status of the one-shot flag.
// A True value means the measurement is still ongoing.
// A False value means measurement is complete.
func (d *MAX31856Driver) OneShotPending() bool {
	confReg0, err := d.readUint8(MAX31856_CR0_REG)
	if err != nil {
		d.logger.Error("OneShotPending Error", "err", err)
	}
	return (confReg0 & MAX31856_CR0_1SHOT) != 0
}
func (d *MAX31856Driver) waitOneshot() {
	for d.OneShotPending() {
		time.Sleep(10 * time.Millisecond)
	}
}
func (d *MAX31856Driver) performOneShotMeasurement() error {
	if err := d.InitOneShotMeasurement(); err != nil {
		return err
	}
	d.waitOneshot()

	return nil
}

// GetTemperature: Measure the temperature of the sensor and wait for the result.
//
//	Return value is in degrees Celsius.
func (d *MAX31856Driver) GetTemperature() (float64, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.performOneShotMeasurement(); err != nil {
		return 0, err
	}
	return d.UnpackTemperature()
}

// Reads the probe temperature from the register
func (d *MAX31856Driver) UnpackTemperature() (float64, error) {
	rawTemp := make([]byte, 3)
	if d.connection == nil {
		return 0, fmt.Errorf("cannot get data from device")
	}
	d.connection.ReadBlockData(MAX31856_LTCBH_REG, rawTemp)

	// effectively shift raw_read >> 12 to convert pseudo-float
	tempInt := (uint32(rawTemp[0]) << 16) | (uint32(rawTemp[1]) << 8) | uint32(rawTemp[2])
	return float64(tempInt) / 4096.0, nil
}

func (d *MAX31856Driver) writeUint8(address, val byte) error {
	//NEEDED: Address with 8th bit to 1 are write
	address |= 0x80
	if d.connection == nil {
		return fmt.Errorf("cannot get data from device")
	}
	if err := d.connection.WriteByteData(address, val); err != nil {
		return err
	}
	return nil
}

func (d *MAX31856Driver) readUint8(address byte) (uint8, error) {
	address &= 0x7F
	if d.connection == nil {
		return 0, fmt.Errorf("cannot get data from device")
	}
	if val, err := d.connection.ReadByteData(address); err != nil {
		return 0, err
	} else {
		return val, nil
	}

}
