package hwinterface

import "log/slog"

func (c piController) GetTemperature() float64 {
	value, err := c.analogInput.GetTemperature()
	if err != nil {
		slog.Error("Error: %v", err)
		return 0
	}

	return value
}
