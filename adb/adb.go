package adb

import goadb "github.com/zach-klippenstein/goadb"

var Adb *goadb.Adb

func init() {
	Adb, _ = goadb.New()
}

func DeviceBySerial(serial string) *goadb.Device {
	return Adb.Device(goadb.DeviceWithSerial(serial))
}
