package main

import (
	"log"
	"os"
)

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

type NANDPU struct {
	PC   Reg16 // Program Counter
	INST Reg8  // Instruction Register
	INC  Reg16 // Increment Register
	SP   Reg16 // Stack Pointer

	// Flags
	Zero     bool
	Carry    bool
	Sign     bool
	LessThan bool

	// 8 Bit Registers
	RegA Reg8
	RegB Reg8
	RegC Reg8
	RegD Reg8

	// 16 Bit Registers
	RegM  *SplitReg16
	RegXY *SplitReg16
	RegJ  *SplitReg16

	Reg8List  []Reg8Like
	Reg16List []Reg16Like

	Mem MemMap
}

var Logger *log.Logger

func NewNANDPU(romData []byte) *NANDPU {
	Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	c := NANDPU{}

	rom := NewROM(0x0000, 0x8000)
	rom.Init(romData)
	c.Mem.AddRegion(0x0000, 0x7FFF, rom)                    // 32K ROM (AT28C256)
	c.Mem.AddRegion(0x8000, 0xFFFF, NewRAM(0x8000, 0x8000)) // 32K RAM (CY62256N)

	c.PC.AccessFlags = AccessFlags{CanRead: true, CanWrite: true}
	c.INST.AccessFlags = AccessFlags{CanRead: true, CanWrite: true}
	c.INC.AccessFlags = AccessFlags{CanRead: true, CanWrite: false}
	c.SP.val = 0xFFFF
	c.SP.AccessFlags = AccessFlags{CanRead: true, CanWrite: true}

	c.RegA.AccessFlags = AccessFlags{CanRead: true, CanWrite: true}
	c.RegB.AccessFlags = AccessFlags{CanRead: true, CanWrite: true}
	c.RegC.AccessFlags = AccessFlags{CanRead: true, CanWrite: true}
	c.RegD.AccessFlags = AccessFlags{CanRead: true, CanWrite: true}

	c.RegM = NewSplitReg16(
		AccessFlags{CanRead: true, CanWrite: false}, // full M
		AccessFlags{CanRead: true, CanWrite: true},  // M.Hi
		AccessFlags{CanRead: true, CanWrite: true},  // M.Lo
	)
	c.RegXY = NewSplitReg16(
		AccessFlags{CanRead: true, CanWrite: true}, // full M
		AccessFlags{CanRead: true, CanWrite: true}, // M.Hi
		AccessFlags{CanRead: true, CanWrite: true}, // M.Lo
	)
	c.RegJ = NewSplitReg16(
		AccessFlags{CanRead: true, CanWrite: false},  // full M
		AccessFlags{CanRead: false, CanWrite: false}, // M.Hi
		AccessFlags{CanRead: false, CanWrite: false}, // M.Lo
	)

	c.RegM.Hi.AccessFlags = AccessFlags{CanRead: true, CanWrite: true}
	c.RegM.Lo.AccessFlags = AccessFlags{CanRead: true, CanWrite: true}

	c.RegXY.Hi.AccessFlags = AccessFlags{CanRead: true, CanWrite: true}
	c.RegXY.Lo.AccessFlags = AccessFlags{CanRead: true, CanWrite: true}

	c.RegJ.Hi.AccessFlags = AccessFlags{CanRead: false, CanWrite: true}
	c.RegJ.Lo.AccessFlags = AccessFlags{CanRead: false, CanWrite: true}

	c.Reg8List = []Reg8Like{
		&c.RegA,
		&c.RegB,
		&c.RegC,
		&c.RegD,
		c.RegM.Hi, c.RegM.Lo,
		c.RegXY.Hi, c.RegM.Lo,
		c.RegJ.Hi, c.RegJ.Lo,
	}

	c.Reg16List = []Reg16Like{
		c.RegM,
		c.RegXY,
		c.RegJ,
		&c.PC,
		&c.INC,
		&c.SP,
	}

	Logger.Println("Initialised NANDPU")
	return &c
}

func (c *NANDPU) getMemVal() byte {
	return c.Mem.Read(uint16(c.PC.Get()))
}

func (c *NANDPU) getInst() {
	c.INST.Set(c.getMemVal())
}

func (c *NANDPU) getReg8FromMem() (byte, Reg8Like) {
	targetIndex := c.getMemVal()
	target := c.Reg8List[targetIndex]
	return targetIndex, target
}
func (c *NANDPU) getReg16FromMem() (byte, Reg16Like) {
	targetIndex := c.getMemVal()
	target := c.Reg16List[targetIndex]
	return targetIndex, target
}

// Updates the Zero and Sign flags based off of the value argument, and Less Than based off of the A and B registers
func (c *NANDPU) updateFlags(value byte) {
	c.Zero = value == 0
	c.Sign = (value >> 7) == 1
	c.LessThan = c.RegB.Get() < c.RegC.Get()
}

func (c *NANDPU) increment16(value uint16) {
	c.INC.val = value + 1
	// We don't use the Set method here, because the INC register is configured to be read-only.
	// This special logic is the only thing that writes to the INC register.
}

func (c *NANDPU) pcInc() {
	c.increment16(c.PC.Get())
	c.PC.Set(c.INC.Get())
}

func (c *NANDPU) push(val byte) {
	c.Mem.Write(c.SP.Get(), val)
	c.SP.val -= 1
}

func (c *NANDPU) pop() byte {
	c.SP.val += 1
	return c.Mem.Read(c.SP.Get())
}

func (c *NANDPU) printFlags() {
	Logger.Printf("Flag values: Z=%d C=%d S=%d L=%d",
		boolToInt(c.Zero),
		boolToInt(c.Carry),
		boolToInt(c.Sign),
		boolToInt(c.LessThan),
	)
}

func (c *NANDPU) Step() bool { // TODO
	c.getInst()
	instName := OpcodeNames[c.INST.Get()]

	Logger.Printf("ADDR 0x%04X: Instruction byte: %s (0x%02X)\n", c.PC.Get(), instName, c.INST.Get())

	switch c.INST.Get() {
	case OP_NOP:
		Logger.Println("No operation.")
		c.pcInc()

	case OP_CMP:
		c.updateFlags(c.RegB.Get())
		c.Carry = (c.RegB.Get() & 0x01) == 1
		c.printFlags()
		c.pcInc()

	case OP_ADD:
		result := uint16(c.RegB.Get()) + uint16(c.RegC.Get())
		resultByte := byte(result)
		c.updateFlags(resultByte)
		c.Carry = result > 0xFF
		c.pcInc()
		targetIndex, target := c.getReg8FromMem()
		target.Set(resultByte)
		Logger.Printf("ADD regB (value %d) + regC (value %d) -> %s (new value %d)", c.RegB.Get(), c.RegC.Get(), Reg8Names[targetIndex], target.Get())
		c.printFlags()
		c.pcInc()

	case OP_SPECIAL_HALT:
		return false
	}

	return true
}
