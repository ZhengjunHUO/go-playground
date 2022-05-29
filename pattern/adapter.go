package main

import (
	"fmt"
)

// 使用者，行为:用usb type C充电器给兼容该充电接口的设备充电，并未配备其他充电器
type user struct {}
func (u *user) useTypeCToChargeDevice(d device) {
	fmt.Println("Charging with USB TypeC charger...")
	d.chargeWithTypeC()
}

// 充电接口为TypeC的一类笔记本
type device interface {
	chargeWithTypeC()
}

// 充电接口为TypeC的设备，device的一个实例
type thinkpad struct {}
func (t *thinkpad) chargeWithTypeC() {
	fmt.Println("Thinkpad is charging with USB TypeC charger !")
}

// 充电接口为旧式DC的设备，并非device的一个实例，即user不能给它充电
type fujitsu struct {}
func (f *fujitsu) chargeWithDC() {
	fmt.Println("Fujitsu is charging with traditional DC charger !")
}

// 为旧式DC的设备提供的转接头，可使其适配type C
type fujitsuAdapter struct {
	fujitsuDevice	*fujitsu
}
func (fa *fujitsuAdapter) chargeWithTypeC() {
	fmt.Println("Use a USB Type C to DC adapter...")
	fa.fujitsuDevice.chargeWithDC()
}

func main() {
	you := &user{}

	// 可直接充电
	thkpad := &thinkpad{}
	you.useTypeCToChargeDevice(thkpad) 

	fmt.Println()
	
	// 需要转接头才能充电
	fjts := &fujitsuAdapter{&fujitsu{}}
	you.useTypeCToChargeDevice(fjts)
}
