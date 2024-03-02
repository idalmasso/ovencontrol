package hwinterface

func (c piController) GetTemperature() (float64, error) {
	value, err := c.temperatureReader.GetTemperature()
	if err != nil {
		c.logger.Error("Error: %v", err)
		return 0, err
	}

	return value, nil
}
