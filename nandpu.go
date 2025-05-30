package main

import (
	"log"
	"os"
)

type NANDPU struct {
	PC   uint16 // Program Counter
	INST byte   // Instruction Register
	INC  uint16 // Increment Register

	// Flags
	Zero     bool
	Carry    bool
	Sign     bool
	LessThan bool

	// 8 Bit Registers
	RegA byte
	RegB byte
	RegC byte
	RegD byte

	// 16 Bit Registers
	RegM  Reg16
	RegXY Reg16
	RegJ  Reg16

	Mem MemMap
}

var Logger *log.Logger

func NewNANDPU() *NANDPU {
	Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	c := NANDPU{}
	c.Mem.AddRegion(0x0000, 0x7FFF, NewROM(0x0000, 0x8000)) // 32K ROM (AT28C256)
	c.Mem.AddRegion(0x8000, 0xFFFF, NewRAM(0x8000, 0x8000)) // 32K RAM (CY62256N)
	Logger.Println("Initialised NANDPU")
	return &c
}

func (c *NANDPU) GetInst() byte {
	return c.Mem.Read(uint16(c.INST))
}

// Updates the Zero and Sign flags based off of the value argument, and Less Than based off of the A and B registers
func (c *NANDPU) UpdateFlags(value byte) {
	c.Zero = value == 0
	c.Sign = (value >> 7) == 1
	c.LessThan = c.RegB < c.RegC
}

func (c *NANDPU) Step() { // TODO
	Logger.Println("Start of step")
}
