package main

import (
	"log"
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

func (c *NANDPU) decrement16(value uint16) {
	c.INC.val = value - 1
}

func (c *NANDPU) pcInc() {
	c.increment16(c.PC.Get())
	c.PC.Set(c.INC.Get())
}

func (c *NANDPU) push(val byte) {
	c.Mem.Write(c.SP.Get(), val)
	c.decrement16(c.SP.Get())
	c.SP.Set(c.INC.Get())
}

func (c *NANDPU) pop() byte {
	c.increment16(c.SP.Get())
	c.SP.Set(c.INC.Get())
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

func (c *NANDPU) branchLogicImm(condition bool, opcode byte) {
	name := OpcodeNames[opcode]
	c.pcInc()
	addrLo := c.getMemVal()
	c.RegJ.Lo.Set(addrLo)
	c.pcInc()
	addrHi := c.getMemVal()
	c.RegJ.Hi.Set(addrHi)
	if condition {
		c.PC.Set(c.RegJ.Get())
		Logger.Printf("%s (condition met) -> jump to addr 0x%04X", name, c.RegJ.Get())
	} else {
		Logger.Printf("%s (condition not met) -> do not jump to addr 0x%04X", name, c.RegJ.Get())
		c.pcInc()
	}
}

func (c *NANDPU) branchLogicJ(condition bool, opcode byte) {
	name := OpcodeNames[opcode]

	if condition {
		c.PC.Set(c.RegJ.Get())
		Logger.Printf("%s (condition met) -> jump to addr 0x%04X", name, c.RegJ.Get())
	} else {
		Logger.Printf("%s (condition not met) -> do not jump to addr 0x%04X", name, c.RegJ.Get())
		c.pcInc()
	}
}

func (c *NANDPU) Step() bool {
	c.getInst()

	Logger.Printf("ADDR 0x%04X", c.PC.Get())

	switch c.INST.Get() {
	case OP_NOP:
		Logger.Println("No operation.")

	case OP_CMP:
		c.updateFlags(c.RegB.Get())
		c.Carry = (c.RegB.Get() & 0x01) == 1
		c.printFlags()

	case OP_ADD:
		result := uint16(c.RegB.Get()) + uint16(c.RegC.Get())
		resultByte := byte(result)
		c.updateFlags(resultByte)
		c.Carry = result > 0xFF
		c.pcInc()
		targetIndex, target := c.getReg8FromMem()
		prevRegBVal := c.RegB.Get()
		target.Set(resultByte)
		Logger.Printf("ADD regB (value %d) + regC (value %d) -> %s (new value %d)", prevRegBVal, c.RegC.Get(), Reg8Names[targetIndex], target.Get())
		c.printFlags()

	case OP_SUB:
		result := c.RegB.Get() - c.RegC.Get()
		c.updateFlags(result)
		c.Carry = c.RegC.Get() > c.RegB.Get()
		c.pcInc()
		targetIndex, target := c.getReg8FromMem()
		prevRegBVal := c.RegB.Get()
		target.Set(result)
		Logger.Printf("SUB regB (value %d) - regC (value %d) -> %s (new value %d)", prevRegBVal, c.RegC.Get(), Reg8Names[targetIndex], target.Get())
		c.printFlags()

	case OP_INC:
		result := c.RegB.Get() + 1
		c.updateFlags(result)
		c.Carry = c.RegB.Get() > result
		c.pcInc()
		targetIndex, target := c.getReg8FromMem()
		prevRegBVal := c.RegB.Get()
		target.Set(result)
		Logger.Printf("INC regB (value %d) + 1 -> %s (new value %d)", prevRegBVal, Reg8Names[targetIndex], target.Get())
		c.printFlags()

	case OP_DEC:
		result := c.RegB.Get() - 1
		c.updateFlags(result)
		c.Carry = result > c.RegB.Get()
		c.pcInc()
		targetIndex, target := c.getReg8FromMem()
		prevRegBVal := c.RegB.Get()
		target.Set(result)
		Logger.Printf("DEC regB (value %d) - 1 -> %s (new value %d)", prevRegBVal, Reg8Names[targetIndex], target.Get())
		c.printFlags()

	case OP_NAND:
		result := ^(c.RegB.Get() & c.RegC.Get())
		c.updateFlags(result)
		c.Carry = (result & 0x01) == 1
		c.pcInc()
		targetIndex, target := c.getReg8FromMem()
		prevRegBVal := c.RegB.Get()
		target.Set(result)
		Logger.Printf("NAND ~(regB (value %d) & regC (value %d)) -> %s (new value %d)", prevRegBVal, c.RegC.Get(), Reg8Names[targetIndex], target.Get())
		c.printFlags()

	case OP_SHR:
		oldCarry := c.Carry
		result := (c.RegB.Get() >> 1) | (byte(boolToInt(c.Carry)) << 7)
		c.updateFlags(result)
		c.Carry = (c.RegB.Get() >> 7) == 1
		c.pcInc()
		targetIndex, target := c.getReg8FromMem()
		prevRegBVal := c.RegB.Get()
		target.Set(result)
		Logger.Printf("SHR (regB (value %d) >> 1) | (carry (value %d) << 7) -> %s (new value %d)", prevRegBVal, boolToInt(oldCarry), Reg8Names[targetIndex], target.Get())
		c.printFlags()

	case OP_SHL:
		oldCarry := c.Carry
		result := (c.RegB.Get() << 1) | (byte(boolToInt(c.Carry)))
		c.updateFlags(result)
		c.Carry = (c.RegB.Get() >> 7) == 1
		c.pcInc()
		targetIndex, target := c.getReg8FromMem()
		prevRegBVal := c.RegB.Get()
		target.Set(result)
		Logger.Printf("SHR (regB (value %d) << 1) | carry (value %d) -> %s (new value %d)", prevRegBVal, boolToInt(oldCarry), Reg8Names[targetIndex], target.Get())
		c.printFlags()

	case OP_LDI:
		c.pcInc()
		val := c.getMemVal()
		c.pcInc()
		targetIndex, target := c.getReg8FromMem()
		prevTargetVal := target.Get()
		target.Set(val)
		Logger.Printf("LDI %d into %s (value %d)", val, Reg8Names[targetIndex], prevTargetVal)

	case OP_LDMI:
		c.pcInc()
		addrLo := c.getMemVal()
		c.RegM.Lo.Set(addrLo)
		c.pcInc()
		addrHi := c.getMemVal()
		c.RegM.Hi.Set(addrHi)
		val := c.Mem.Read(c.RegM.Get())
		c.pcInc()
		targetIndex, target := c.getReg8FromMem()
		prevTargetVal := target.Get()
		target.Set(val)
		Logger.Printf("LDMI addr 0x%04X (value %d) into %s (value %d)", c.RegM.Get(), val, Reg8Names[targetIndex], prevTargetVal)

	case OP_LDM:
		addr := c.RegM.Get()
		val := c.Mem.Read(addr)
		c.pcInc()
		targetIndex, target := c.getReg8FromMem()
		prevTargetVal := target.Get()
		target.Set(val)
		Logger.Printf("LDM from M register (addr 0x%04X) (value %d) into %s (value %d)", addr, val, Reg8Names[targetIndex], prevTargetVal)

	case OP_STOI:
		c.pcInc()
		sourceIndex, source := c.getReg8FromMem()
		c.pcInc()
		addrLo := c.getMemVal()
		c.RegM.Lo.Set(addrLo)
		c.pcInc()
		addrHi := c.getMemVal()
		c.RegM.Hi.Set(addrHi)
		prevMemVal := c.Mem.Read(c.RegM.Get())
		c.Mem.Write(c.RegM.Get(), source.Get())
		Logger.Printf("STOI from %s (value %d) into addr 0x%04X (value %d)", Reg8Names[sourceIndex], source.Get(), c.RegM.Get(), prevMemVal)

	case OP_STO:
		c.pcInc()
		sourceIndex, source := c.getReg8FromMem()
		addr := c.RegM.Get()
		prevMemVal := c.Mem.Read(addr)
		c.Mem.Write(addr, source.Get())
		Logger.Printf("STO from %s (value %d) into mem at M register (addr 0x%04X) (value %d)", Reg8Names[sourceIndex], source.Get(), addr, prevMemVal)

	case OP_PUSH:
		c.pcInc()
		sourceIndex, source := c.getReg8FromMem()
		c.push(source.Get())
		Logger.Printf("PUSH %s (value %d) onto stack", Reg8Names[sourceIndex], source.Get())

	case OP_POP:
		c.pcInc()
		targetIndex, target := c.getReg8FromMem()
		target.Set(c.pop())
		Logger.Printf("POP stack into %s (value %d)", Reg8Names[targetIndex], target.Get())

	case OP_MOV8:
		c.pcInc()
		sourceIndex, source := c.getReg8FromMem()
		c.pcInc()
		targetIndex, target := c.getReg8FromMem()
		oldTargetVal := target.Get()
		target.Set(source.Get())
		Logger.Printf("MOV8 from %s (value %d) into %s (value %d)", Reg8Names[sourceIndex], source.Get(), Reg8Names[targetIndex], oldTargetVal)

	case OP_MOV16:
		c.pcInc()
		sourceIndex, source := c.getReg16FromMem()
		c.pcInc()
		targetIndex, target := c.getReg16FromMem()
		oldTargetVal := target.Get()
		target.Set(source.Get())
		Logger.Printf("MOV16 from %s (value %d) into %s (value %d)", Reg16Names[sourceIndex], source.Get(), Reg16Names[targetIndex], oldTargetVal)

	case OP_JMPI:
		c.pcInc()
		addrLo := c.getMemVal()
		c.RegJ.Lo.Set(addrLo)
		c.pcInc()
		addrHi := c.getMemVal()
		c.RegJ.Hi.Set(addrHi)
		c.PC.Set(c.RegJ.Get())
		Logger.Printf("JMPI to addr 0x%04X", c.RegJ.Get())
		return true // Avoid incrementing the PC after the instruction has finished

	case OP_CALI:
		PCLo := byte(c.PC.Get() & 0x00FF)
		PCHi := byte((c.PC.Get() & 0xFF00) >> 8)
		c.push(PCLo)
		c.push(PCHi)

		c.pcInc()
		addrLo := c.getMemVal()
		c.RegJ.Lo.Set(addrLo)
		c.pcInc()
		addrHi := c.getMemVal()
		c.RegJ.Hi.Set(addrHi)
		c.PC.Set(c.RegJ.Get())
		Logger.Printf("CALI addr 0x%04X (SP now 0x%04X)", c.RegJ.Get(), c.SP.Get())
		return true // Avoid incrementing the PC after the instruction has finished

	case OP_JMP:
		c.PC.Set(c.RegJ.Get())
		Logger.Printf("JMP to addr 0x%04X", c.RegJ.Get())
		return true // Avoid incrementing the PC after the instruction has finished

	case OP_CALL:
		c.RegXY.Set(c.PC.Get())
		c.push(c.RegXY.Lo.Get())
		c.push(c.RegXY.Hi.Get())

		c.PC.Set(c.RegJ.Get())
		Logger.Printf("CALL addr 0x%04X (SP now 0x%04X)", c.RegJ.Get(), c.SP.Get())
		return true // Avoid incrementing the PC after the instruction has finished

	case OP_RET:
		c.RegJ.Hi.Set(c.pop())
		c.RegJ.Lo.Set(c.pop())

		c.PC.Set(c.RegJ.Get())
		Logger.Printf("RET to addr 0x%04X (SP now 0x%04X)", c.RegJ.Get(), c.SP.Get())

	case OP_BZSI:
		condition := c.Zero
		c.branchLogicImm(condition, OP_BZSI)
		return true // Avoid incrementing the PC after the instruction has finished
	case OP_BZCI:
		condition := !c.Zero
		c.branchLogicImm(condition, OP_BZCI)
		return true // Avoid incrementing the PC after the instruction has finished
	case OP_BCSI:
		condition := c.Carry
		c.branchLogicImm(condition, OP_BCSI)
		return true // Avoid incrementing the PC after the instruction has finished
	case OP_BCCI:
		condition := !c.Carry
		c.branchLogicImm(condition, OP_BCCI)
		return true // Avoid incrementing the PC after the instruction has finished
	case OP_BSSI:
		condition := c.Sign
		c.branchLogicImm(condition, OP_BSSI)
		return true // Avoid incrementing the PC after the instruction has finished
	case OP_BSCI:
		condition := !c.Sign
		c.branchLogicImm(condition, OP_BSCI)
		return true // Avoid incrementing the PC after the instruction has finished
	case OP_BLSI:
		condition := c.LessThan
		c.branchLogicImm(condition, OP_BLSI)
		return true // Avoid incrementing the PC after the instruction has finished
	case OP_BLCI:
		condition := !c.LessThan
		c.branchLogicImm(condition, OP_BLCI)
		return true // Avoid incrementing the PC after the instruction has finished

	case OP_BZS:
		condition := c.Zero
		c.branchLogicJ(condition, OP_BZS)
		return true // Avoid incrementing the PC after the instruction has finished
	case OP_BZC:
		condition := !c.Zero
		c.branchLogicJ(condition, OP_BZC)
		return true // Avoid incrementing the PC after the instruction has finished
	case OP_BCS:
		condition := c.Carry
		c.branchLogicJ(condition, OP_BCS)
		return true // Avoid incrementing the PC after the instruction has finished
	case OP_BCC:
		condition := !c.Carry
		c.branchLogicJ(condition, OP_BCC)
		return true // Avoid incrementing the PC after the instruction has finished
	case OP_BSS:
		condition := c.Sign
		c.branchLogicJ(condition, OP_BSS)
		return true // Avoid incrementing the PC after the instruction has finished
	case OP_BSC:
		condition := !c.Sign
		c.branchLogicJ(condition, OP_BSC)
		return true // Avoid incrementing the PC after the instruction has finished
	case OP_BLS:
		condition := c.LessThan
		c.branchLogicJ(condition, OP_BLS)
		return true // Avoid incrementing the PC after the instruction has finished
	case OP_BLC:
		condition := !c.LessThan
		c.branchLogicJ(condition, OP_BLC)
		return true // Avoid incrementing the PC after the instruction has finished

	case OP_SPECIAL_HALT:
		c.pcInc()
		return false
	}

	c.pcInc()

	Logger.Printf("STATE: PC=0x%04X A=0x%02X B=0x%02X C=0x%02X D=0x%02X M=0x%04X XY=0x%04X J=0x%04X SP=0x%04X INC=0x%04X | FLAGS Z=%t C=%t S=%t LT=%t",
		c.PC.Get(),
		c.RegA.Get(),
		c.RegB.Get(),
		c.RegC.Get(),
		c.RegD.Get(),
		c.RegM.Get(),
		c.RegXY.Get(),
		c.RegJ.Get(),
		c.SP.Get(),
		c.INC.Get(),
		c.Zero,
		c.Carry,
		c.Sign,
		c.LessThan,
	)

	return true
}
