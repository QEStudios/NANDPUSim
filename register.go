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

func (r *Reg8) Get() byte {
	if !r.CanRead {
		Logger.Panic("attempted to read from Reg8 without read capability")
	}
	return r.val
}
func (r *Reg8) Set(v byte) {
	if !r.CanWrite {
		Logger.Panic("attempted to write to Reg8 without write capability")
	}
	r.val = v
}

type Reg16 struct {
	val uint16
	AccessFlags
}

func (r *Reg16) Get() uint16 {
	if !r.CanRead {
		Logger.Panic("attempted to read from Reg16 without read capability")
	}
	return r.val
}
func (r *Reg16) Set(v uint16) {
	if !r.CanWrite {
		Logger.Panic("attempted to write to Reg16 without write capability")
	}
	r.val = v
}

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

func (r *SplitReg16) Get() uint16 {
	if !r.CanRead {
		Logger.Panic("attempted to read from SplitReg16 without read capability")
	}
	return r.val
}
func (r *SplitReg16) Set(v uint16) {
	if !r.CanWrite {
		Logger.Panic("attempted to write to SplitReg16 without write capability")
	}
	r.val = v
}
func (h *splitHi) Get() byte {
	if !h.CanRead {
		Logger.Panic("attempted to read from splitHi without read capability")
	}
	return byte(h.parent.val >> 8)
}
func (h *splitHi) Set(v byte) {
	if !h.CanWrite {
		Logger.Panic("attempted to write to splitHi without write capability")
	}
	h.parent.val = (h.parent.val & 0x00FF) | (uint16(v) << 8)
}
func (l *splitLo) Get() byte {
	if !l.CanRead {
		Logger.Panic("attempted to read from splitLo without read capability")
	}
	return byte(l.parent.val & 0x00FF)
}
func (l *splitLo) Set(v byte) {
	if !l.CanWrite {
		Logger.Panic("attempted to write to splitLo without write capability")
	}
	l.parent.val = (l.parent.val & 0xFF00) | uint16(v)
}

const (
	REG_A  byte = 0x00
	REG_B  byte = 0x01
	REG_C  byte = 0x02
	REG_D  byte = 0x03
	REG_M1 byte = 0x04
	REG_M2 byte = 0x05
	REG_X  byte = 0x06
	REG_Y  byte = 0x07
	REG_J1 byte = 0x08
	REG_J2 byte = 0x09

	REG_M   byte = 0x00
	REG_XY  byte = 0x01
	REG_J   byte = 0x02
	REG_PC  byte = 0x03
	REG_INC byte = 0x04
)

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
