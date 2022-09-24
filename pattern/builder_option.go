package main

import (
	"fmt"
)

type Commander struct {
	name	string
}

type Battalion struct {
	name	string
	size	int
	power	[]CombatPower
}

type CombatPower struct {
	eq          Equipment
	annotations []string
	serialNo    SerialNo
}

/* object */
type Equipment interface {
	Attack()
}

type GundamMKII struct {}
func (mk2 *GundamMKII) Attack() {
	fmt.Println("Gudam MKII launches attack!")
}

type ZetaGundam struct {}
func (z *ZetaGundam) Attack() {
	fmt.Println("Zeta Gudam launches attack!")
}

type HyakuShiki struct {}
func (h *HyakuShiki) Attack() {
	fmt.Println("HyakuShiki launches attack!")
}

/* option func */
type CombatPowerOption interface {
	Annotate(*CombatPower)
}

/* option 1 */
type Remarks struct {
	annotations []string
}

func (r Remarks) Annotate(c *CombatPower) {
	c.annotations = r.annotations
}

/* option 2 */
type SerialNo int
func (s SerialNo) Annotate(c *CombatPower) {
	c.serialNo = s
}

/* implementation check */
var (
	_ CombatPowerOption = Remarks{}
	_ CombatPowerOption = SerialNo(0)
)

/* Builder */
type Builder struct {
	cmd	Commander
	cps	[]CombatPower
}

func (b *Builder) EquippedWith(e Equipment, opts ...CombatPowerOption) *Builder {
	cp := CombatPower{eq: e}
	for _, opt := range opts {
		opt.Annotate(&cp)
	}

	b.cps = append(b.cps, cp)
	return b
}

func (b *Builder) Ready() *Battalion {
	return &Battalion{
		name:  b.cmd.name,
		size:  len(b.cps),
		power: b.cps,
	}
}

func NewBattalionCommandedBy(c Commander) *Builder {
	return &Builder{cmd: c}
}

/* entrypoint */
func main() {
	cmd := Commander{name: "huo",}
	bt := NewBattalionCommandedBy(cmd).
	  EquippedWith(&GundamMKII{},
	    Remarks{annotations: []string{"Flying Armor"},}).
	  EquippedWith(&ZetaGundam{},
	    Remarks{annotations: []string{"Bio-Sensor", "Hyper Mega Launcher"},}, SerialNo(16)).
	  EquippedWith(&HyakuShiki{},
	    Remarks{annotations: []string{"IDE System", "Beam-Resistant Coating"},}, SerialNo(100)).
	  Ready()
	fmt.Printf("%v\n", bt)
}
