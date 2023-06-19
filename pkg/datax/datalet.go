package datax

import (
	"errors"
)

var (
	drivers = map[string]DataDriver{}
)

func RegisterDriver(name string, drv DataDriver) error {
	if drv == nil {
		return nil
	}
	if drv, ok := drivers[name]; ok {
		return errors.New("driver (" + name + ") already existed")
	} else {
		drivers[name] = drv
	}
	return nil
}
