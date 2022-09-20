package main

import (
	"fmt"
)

/* Builder interface */
type IFGunBuilder interface {
	setMeleeWeapon()
	setMainWeapon()
	setShield()
	setLongRangeWeapon()
	getGundam() Gundam
}

func getGunBuilder(builderType string) IFGunBuilder {
	if builderType == "classic" {
		return newStdBuilder()
	}

	if builderType == "nextgen" {
		return newNextgenBuilder()
	}

	return nil
}

/* Product */
type Gundam struct {
	meleeWeapon	string
	mainWeapon	string
	shield		string
	longrangeWeapon	string
}

/* Builder Implementation */
type StandardBuilder struct {
	Gundam
}

func newStdBuilder() *StandardBuilder {
	return &StandardBuilder{}
}

func (s *StandardBuilder) setMeleeWeapon() {
	s.meleeWeapon = "Beam Saber"
}

func (s *StandardBuilder) setMainWeapon() {
	s.mainWeapon = "BOWA*XBR-M-79-07G Beam Rifle"
}

func (s *StandardBuilder) setShield() {
	s.shield = "RX*M-Sh-008/S-01025 Shield"
}

func (s *StandardBuilder) setLongRangeWeapon() {
	s.longrangeWeapon = "BLASH*XHB-L-03/N-STD Hyper Bazooka"
}

func (s *StandardBuilder) getGundam() Gundam {
	return s.Gundam
}


type NextGenBuilder struct {
	Gundam
}

func newNextgenBuilder() *NextGenBuilder {
	return &NextGenBuilder{}
}

func (n *NextGenBuilder) setMeleeWeapon() {
	n.meleeWeapon = "Beam Tonfa"
}

func (n *NextGenBuilder) setMainWeapon() {
	n.mainWeapon = "Beam Magnum"
}

func (n *NextGenBuilder) setShield() {
	n.shield = "Armed Armor DE"
}

func (n *NextGenBuilder) setLongRangeWeapon() {
	n.longrangeWeapon = "Hyper Bazooka"
}

func (n *NextGenBuilder) getGundam() Gundam {
	return n.Gundam
}

/* Director: how to build */
type Gnaku struct {
	builder IFGunBuilder
}

func newGnaku(b IFGunBuilder) *Gnaku {
	return &Gnaku{
		builder: b,
	}
}

func (g *Gnaku) setGunBuilder(b IFGunBuilder) {
	g.builder = b
}

func (g *Gnaku) buildGundam() Gundam {
	if g.builder == nil {
		return Gundam{}
	}

	g.builder.setMeleeWeapon()
	g.builder.setMainWeapon()
	g.builder.setShield()
	g.builder.setLongRangeWeapon()
	return g.builder.getGundam()
}


func main() {
	fmt.Println(newGnaku(getGunBuilder("classic")).buildGundam())
	fmt.Println(newGnaku(getGunBuilder("nextgen")).buildGundam())
	fmt.Println(newGnaku(getGunBuilder("foo")).buildGundam())
}
