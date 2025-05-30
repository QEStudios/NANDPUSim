package main

import (
	"log"
	"os"
)

type NANDPU struct {
	PC   uint16 // Program Counter
	INST byte   // Instruction Register
	INC  uint16 // Increment Register

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
	nandpu := NANDPU{}
	nandpu.Mem.AddRegion(0x0000, 0x7FFF, NewRAM(0x0000, 0x8000)) // 32K RAM (CY62256N)
	nandpu.Mem.AddRegion(0x8000, 0xFFFF, NewROM(0x8000, 0x8000)) // 32K ROM (AT28C256)
	Logger.Println("Initialised NANDPU")
	return &nandpu
}

func (c *NANDPU) Step() { // TODO
	Logger.Println("Start of step")
}
