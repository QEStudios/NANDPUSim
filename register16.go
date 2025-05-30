package main

type Reg16 struct {
	val uint16
}

func (r *Reg16) Get() uint16  { return r.val }
func (r *Reg16) Set(v uint16) { r.val = v }
func (r *Reg16) Lo() byte     { return byte(r.val & 0x00FF) }
func (r *Reg16) Hi() byte     { return byte(r.val >> 8) }
func (r *Reg16) SetLo(v byte) { r.val = (r.val & 0xFF00) | uint16(v) }
func (r *Reg16) SetHi(v byte) { r.val = (r.val & 0x00FF) | (uint16(v) << 8) }
