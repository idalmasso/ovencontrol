package hwinterface

import "time"

func (d *piController) OpenAir() error {
	//first set the open relay as 1
	if err := d.airCompressorOpen.On(); err != nil {
		return err
	}
	if err := d.airCompressorPower.On(); err != nil {
		return err
	}
	go func() {
		time.Sleep(10 * time.Second)
		d.airCompressorPower.Off()
	}()
	return nil
}

func (d *piController) CloseAir() error {
	//first set the open relay as 0
	if err := d.airCompressorOpen.Off(); err != nil {
		return err
	}
	if err := d.airCompressorPower.On(); err != nil {
		return err
	}
	go func() {
		time.Sleep(10 * time.Second)
		d.airCompressorPower.Off()
	}()
	return nil
}
