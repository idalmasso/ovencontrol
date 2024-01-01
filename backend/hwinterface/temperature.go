package hwinterface

func (c piController) GetTemperature() (float64, error) {
	value, err := c.analogInput.AnalogRead("0")
	if err != nil {
		return 0, err
	}
	ovenTemperature := float64(value) / 255 * 1350
	return ovenTemperature, nil
}
