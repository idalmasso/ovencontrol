package spi

import (
	"log/slog"
	"sync"
	"time"

	"gobot.io/x/gobot/v2/drivers/spi"
)

// MAX31856Driver is a driver for the MAX31856 thermocouple reader.
type MAX31856Driver struct {
	*Driver
	ThermocoupleType ThermocoupleType
	mu               *sync.Mutex
}

type FaultState struct {
	Highthresh, Lowthresh, Refinlow, Refinhigh, Rtdinlow, Ovuv bool
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

var AVGSEL_CONSTS = map[int]uint16{1: 0x00, 2: 0x10, 4: 0x20, 8: 0x30, 16: 0x40}

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
func NewMAX31856Driver(a spi.Connector, options ...func(spi.Config)) *MAX31856Driver {
	d := &MAX31856Driver{
		Driver:           NewDriver(a, "MAX31856", spi.WithMode(1)),
		ThermocoupleType: S,
		mu:               &sync.Mutex{},
	}

	for _, option := range options {
		option(d)
	}
	d.afterStart = func() error {
		slog.Info("after start")
		//# assert on any fault
		d.WriteUint8(MAX31856_MASK_REG, 0x0)
		//  # configure open circuit faults
		d.WriteUint8(MAX31856_CR0_REG, MAX31856_CR0_OCFAULT0)
		// # set thermocouple type
		d.SetThermocoupleType(d.ThermocoupleType)

		v, err := d.ReadUint8(MAX31856_CR1_REG)
		slog.Info("value reg1", "value", v, "err", err)
		return nil
	}

	return d
}

// Starts a one-shot measurement and returns immediately.
// A measurement takes approximately 160ms.
// Check the status of the measurement with `oneshot_pending`; when it is false,
func (d *MAX31856Driver) InitOneShotMeasurement() error {
	//read the current value of the first config register
	slog.Info("InitOneShotMeasurement")
	slog.Info("InitOneShotMeasurement Before readByte")
	confReg0, err := d.ReadUint8(MAX31856_CR0_REG)
	if err != nil {
		slog.Error("Error", "err", err)
		return err
	}
	slog.Info("InitOneShotMeasurement After Readbyte", "readValue", confReg0)
	// and the complement to guarantee the autoconvert bit is unset
	confReg0 |= (MAX31856_CR0_AUTOCONVERT)
	// or the oneshot bit to ensure it is set
	confReg0 |= MAX31856_CR0_1SHOT
	slog.Info("InitOneShotMeasurement BeforeWrite", "writeValue", confReg0)
	// write it back with the new values, prompting the sensor to perform a measurement
	err = d.WriteUint8(MAX31856_CR0_REG, confReg0)
	if err != nil {
		slog.Error("InitOneShotMeasurement write Error", "err", err)
	}
	confReg0, err = d.ReadUint8(MAX31856_CR0_REG)
	if err != nil {
		slog.Error("InitOneShotMeasurement Error", "err", err)
	}
	slog.Info("InitOneShotMeasurement AfterWrite", "Reread value", confReg0)
	return err
}

// A boolean indicating the status of the one-shot flag.
// A True value means the measurement is still ongoing.
// A False value means measurement is complete.
func (d *MAX31856Driver) OneShotPending() bool {
	slog.Info("OneShotPending")
	confReg0, err := d.ReadUint8(MAX31856_CR0_REG)
	if err != nil {
		slog.Error("OneShotPending Error", "err", err)
	}
	slog.Info("OneShotPending return", "confreg", confReg0, "valRet", (confReg0&MAX31856_CR0_1SHOT) != 0)
	return (confReg0 & MAX31856_CR0_1SHOT) != 0
}
func (d *MAX31856Driver) waitOneshot() {
	slog.Info("waitOneShot")
	for d.OneShotPending() {
		time.Sleep(10 * time.Millisecond)
	}
}
func (d *MAX31856Driver) performOneShotMeasurement() error {
	slog.Info("performOneShotMeasurement")
	if err := d.InitOneShotMeasurement(); err != nil {
		slog.Info("performOneShotMeasurement err", "err", err)
		return err
	}
	slog.Info("performOneShotMeasurement Before wait")
	d.waitOneshot()

	return nil
}

// GetTemperature: Measure the temperature of the sensor and wait for the result.
//
//	Return value is in degrees Celsius.
func (d *MAX31856Driver) GetTemperature() (float64, error) {
	slog.Info("GetTemperature")
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.performOneShotMeasurement(); err != nil {
		slog.Info("GetTemperature Err", "err", err)
		return 0, err
	}
	return d.UnpackTemperature(), nil
}

// Reads the probe temperature from the register
func (d *MAX31856Driver) UnpackTemperature() float64 {
	slog.Info("UnpackTemperature")
	rawTemp := make([]byte, 3)
	d.connection.ReadBlockData(MAX31856_LTCBH_REG, rawTemp)
	slog.Info("UnpackTemperature read", "data", rawTemp)
	slog.Info("UnbpackTemp", "data0", rawTemp[0], "data1", rawTemp[1], "data2", rawTemp[2])
	// shift to remove extra byte from unpack needing 4 bytes
	//raw_temp >>= 8d

	// effectively shift raw_read >> 12 to convert pseudo-float
	tempInt := (uint32(rawTemp[0]) << 16) | (uint32(rawTemp[1]) << 8) | uint32(rawTemp[2])
	slog.Info("UnpcakTemperature", "tempInt", tempInt)
	return float64(tempInt) / 4096.0
}

func (d *MAX31856Driver) SetThermocoupleType(thermocoupleType ThermocoupleType) error {
	//# get current value of CR1 Reg

	confReg1, err := d.ReadUint8(MAX31856_CR1_REG)
	if err != nil {
		return err
	}
	confReg1 &= 0xF0
	//# add the new value for the TC type
	confReg1 |= uint8(thermocoupleType) & 0x0F
	return d.WriteUint8(MAX31856_CR1_REG, confReg1)
}

func (d *MAX31856Driver) WriteUint8(address, val byte) error {
	//NEEDED: Address with 7th bit to 1 are write
	address = (address | 0x80) & 0xFF
	if err := d.connection.WriteByteData(address, val); err != nil {
		return err
	}
	return nil
}

func (d *MAX31856Driver) ReadUint8(address byte) (uint8, error) {
	readVal := make([]byte, 1)
	address &= 0x7F
	if err := d.connection.ReadCommandData([]byte{address}, readVal); err != nil {
		return 0, err
	}
	return uint8(readVal[0]), nil
}

func (d *MAX31856Driver) ReadUint16(address byte) (uint16, error) {
	readVal := make([]byte, 2)
	address &= 0x7F
	if err := d.connection.ReadCommandData([]byte{address}, readVal); err != nil {
		return 0, err
	}
	return (uint16(readVal[0]) << 8) | uint16(readVal[1]), nil
}
