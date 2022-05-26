package main

import (
	"fmt"
)

const POOLSIZE = 16

type Datacenter struct {
	ID		int
	Nodes		int
	IdleNodes	int
}

func NewDatacenter(id int) *Datacenter {
	return &Datacenter{
		ID:		id,
		Nodes:		POOLSIZE,
		IdleNodes:	POOLSIZE,
	}
}

func (d *Datacenter) DoGraphicCalc(nodesNeeded int) Command {
	return &GraphicCalcCmd {nodesNeeded, d}
}

func (d *Datacenter) DoStatsCalc(nodesNeeded int) Command {
	return &StatsCalcCmd {nodesNeeded, d}
}

func (d *Datacenter) DoFreeAllNodes() Command {
	return &FreeAllNodesCmd{d}
}

type Command interface {
	Exec()
}

type GraphicCalcCmd struct {
	nodesNeeded	int
	dc		*Datacenter
}

func (g *GraphicCalcCmd) Exec() {
	if g.nodesNeeded <= g.dc.IdleNodes {
		g.dc.IdleNodes -= g.nodesNeeded
		fmt.Printf("[Graphical Task Launching] Reserved %d nodes from datacenter %d ...\n", g.nodesNeeded, g.dc.ID)
	}else{
		fmt.Printf("[Graphical Task Launching] Failed to reserved %d nodes from datacenter %d, abort \n", g.nodesNeeded, g.dc.ID)
	}
}

type StatsCalcCmd struct {
	nodesNeeded	int
	dc		*Datacenter
}

func (s *StatsCalcCmd) Exec() {
	if s.nodesNeeded <= s.dc.IdleNodes {
		s.dc.IdleNodes -= s.nodesNeeded
		fmt.Printf("[Statistical Task Launching] Reserved %d nodes from datacenter %d ...\n", s.nodesNeeded, s.dc.ID)
	}else{
		fmt.Printf("[Statistical Task Launching] Failed to reserved %d nodes from datacenter %d, abort \n", s.nodesNeeded, s.dc.ID)
	}
}

type FreeAllNodesCmd struct {
	dc	*Datacenter
}

func (f *FreeAllNodesCmd) Exec() {
	fmt.Printf("[WARN] Stop all running task in datacenter %d\n", f.dc.ID)
	f.dc.IdleNodes = POOLSIZE
}

type Operator struct {
	id	 int
	commands []Command
}

func (o *Operator) ExecAll() {
	for _, c := range o.commands {
		c.Exec()
	}
}

func main() {
	dc := NewDatacenter(1)
	op := Operator{
		id: 0, 
		commands: []Command {
			dc.DoGraphicCalc(8),
			dc.DoStatsCalc(2),
			dc.DoStatsCalc(4),
			dc.DoGraphicCalc(4),
			dc.DoFreeAllNodes(),
			dc.DoGraphicCalc(4),
		},
	}

	op.ExecAll()
}
