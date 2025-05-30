package main

type Reg8Like interface {
	Get() byte
	Set(byte)
}

type Reg16Like interface {
	Get() uint16
	Set(uint16)
}

type AccessFlags struct {
	CanRead  bool
	CanWrite bool
}

type Reg8 struct {
	val byte
	AccessFlags
}

func (r *Reg8) Get() byte  { return r.val }
func (r *Reg8) Set(v byte) { r.val = v }

type Reg16 struct {
	val uint16
	AccessFlags
}

func (r *Reg16) Get() uint16  { return r.val }
func (r *Reg16) Set(v uint16) { r.val = v }

type SplitReg16 struct {
	val uint16
	AccessFlags
	Hi *splitHi
	Lo *splitLo
}

type splitHi struct {
	parent *SplitReg16
	AccessFlags
}
type splitLo struct {
	parent *SplitReg16
	AccessFlags
}

func NewSplitReg16(flags16, flagsHi, flagsLo AccessFlags) *SplitReg16 {
	r := &SplitReg16{AccessFlags: flags16}
	r.Hi = &splitHi{parent: r, AccessFlags: flagsHi}
	r.Lo = &splitLo{parent: r, AccessFlags: flagsLo}
	return r
}

func (r *SplitReg16) Get() uint16  { return r.val }
func (r *SplitReg16) Set(v uint16) { r.val = v }
func (h *splitHi) Get() byte       { return byte(h.parent.val >> 8) }
func (h *splitHi) Set(v byte)      { h.parent.val = (h.parent.val & 0x00FF) | (uint16(v) << 8) }
func (l *splitLo) Get() byte       { return byte(l.parent.val & 0x00FF) }
func (l *splitLo) Set(v byte)      { l.parent.val = (l.parent.val & 0xFF00) | uint16(v) }

var Reg8Names = map[byte]string{
	0x00: "RegA",
	0x01: "RegB",
	0x02: "RegC",
	0x03: "RegD",
	0x04: "RegM.Hi",
	0x05: "RegM.LO",
	0x06: "RegXY.Hi",
	0x07: "RegXY.Lo",
	0x08: "RegJ.Hi",
	0x09: "RegJ.Lo",
}
var Reg16Names = map[byte]string{
	0x00: "RegM",
	0x01: "RegXY",
	0x02: "RegJ",
	0x03: "PC",
	0x04: "INC",
}
