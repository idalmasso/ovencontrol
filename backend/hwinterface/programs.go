package hwinterface

import "encoding/binary"

func (c *piController) GetPercentual() float64 {
	return c.actualPercentual
}
func (c *piController) GetMaxPower() float64 {
	return c.maxPower
}
func (c *piController) SetPercentual(f float64) error {
	c.actualPercentual = f
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], uint64(f*255))
	return c.ssrPowerController.SetPower(b[0])
}
func (c *piController) InitStartProgram() error {
	c.logger.Info("Init start program")
	var err error
	err = c.ledOvenWorking.On()
	if err != nil {
		c.ledOvenWorking.Off()
		return err
	}
	err = c.ovenRelayPower.On()
	if err != nil {
		c.ovenRelayPower.Off()
		return err
	}
	return nil
}

func (c *piController) EndProgram() error {
	c.logger.Info("End program")
	var err error
	err = c.ledOvenWorking.Off()
	if err != nil {
		c.ovenRelayPower.Off()
		return err
	}
	return c.ovenRelayPower.Off()

}
