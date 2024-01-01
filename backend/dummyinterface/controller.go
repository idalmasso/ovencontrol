package dummyinterface

import "math/rand"

type DummyController struct {
}

func (d DummyController) GetTemperature() float32 {
	return rand.Float32()
}
