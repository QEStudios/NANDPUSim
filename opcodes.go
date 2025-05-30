package main

const (
	// Args: None
	//
	// Changes flags: None
	//
	// Description: No operation.
	OP_NOP byte = 0x00

	// Args: None
	//
	// Changes flags: Zero, Carry, Sign, Less Than
	//
	// Description: Updates the flags based on registers B and C and performs no other operations.
	// The Carry flag is set to the highest bit of 8-bit register B.
	OP_CMP byte = 0x10

	// Args: 1=[out 8-bit reg]
	//
	// Changes flags: Zero, Carry, Sign, Less Than
	//
	// Description: Stores the result of (B + C) into the output 8-bit register (arg 1).
	// All flags are updated based on result.
	// The Carry flag is set when the result overflows past 256.
	OP_ADD byte = 0x11

	// Args: 1=[out 8-bit reg]
	//
	// Changes flags: Zero, Carry, Sign, Less Than
	//
	// Description: Stores the result of (B - C) into the output 8-bit register (arg 1).
	// All flags are updated based on the result.
	// The Carry flag is set when the result underflows below 0.
	OP_SUB byte = 0x12

	// Args: 1=[out 8-bit reg]
	//
	// Changes flags: Zero, Carry, Sign, Less Than
	//
	// Description: Stores the result of (B + 1) into the output 8-bit register (arg 1).
	// All flags are updated based on the result.
	// The Carry flag is set when the result overflows past 256.
	OP_INC byte = 0x13

	// Args: 1=[out 8-bit reg]
	//
	// Changes flags: Zero, Carry, Sign, Less Than
	//
	// Description: Stores the result of (B - 1) into the output 8-bit register (arg 1).
	// All flags are updated based on the result.
	// The Carry flag is set when the result underflows below 0.
	OP_DEC byte = 0x14

	// Args: 1=[out 8-bit reg]
	//
	// Changes flags: Zero, Carry, Sign, Less Than
	//
	// Description: Stores the result of (B NAND 1) into the output 8-bit register (arg 1).
	// All flags are updated based on the result.
	// The Carry flag is set to the highest bit of the result.
	OP_NAND byte = 0x15

	// Args: 1=[out 8-bit reg]
	//
	// Changes flags: Zero, Carry, Sign, Less Than
	//
	// Description: Stores the result of (B >> 1) into the output 8-bit register (arg 1).
	// The Carry flag is used as the most significant bit in the result.
	// All flags are updated based on the result.
	// The Carry flag is set to the bit shifted out (the original least significant bit of B).
	OP_SHR byte = 0x16

	// Args: 1=[out 8-bit reg]
	//
	// Changes flags: Zero, Carry, Sign, Less Than
	//
	// Description: Stores the result of (B << 1) into the output 8-bit register (arg 1).
	// The Carry flag is used as the least significant bit in the result.
	// All flags are updated based on the result.
	// The Carry flag is set to the bit shifted out (the original most significant bit of B).
	OP_SHL byte = 0x17

	// Args: 1=[immediate], 2=[out 8-bit reg]
	//
	// Changes flags: None
	//
	// Description: Stores the immediate value (arg 1) into the output 8-bit register (arg 2).
	OP_LDI byte = 0x20

	// Args: 1=[memory addr], 2=[out 8-bit reg]
	//
	// Changes flags: None
	//
	// Description: Stores the value at the memory address (arg 1) into the output 8-bit register (arg 2).
	OP_LDMI byte = 0x21

	// Args: 1=[out 8-bit reg]
	//
	// Changes flags: None
	//
	// Description: Stores the value at the memory address stored in the M register into the output 8-bit register (arg 2).
	OP_LDM byte = 0x22

	// Args: 1=[in 8-bit reg], 2=[mem addr]
	//
	// Changes flags: None
	//
	// Description: Stores the value from the input 8-bit register (arg 1) into the memory address (arg 2).
	OP_STOI byte = 0x23

	// Args: 1=[out 8-bit reg]
	//
	// Changes flags: None
	//
	// Description: Stores the value from the input 8-bit register (arg 1) into the memory address stored in the M register.
	OP_STO byte = 0x24

	// Args: 1=[in 8-bit reg]
	//
	// Changes flags: None
	//
	// Description: Pushes the value from the input 8-bit register (arg 1) onto the stack,
	// and increments the Stack Pointer by 1.
	OP_PUSH byte = 0x25

	// Args: 1=[out 8-bit reg]
	//
	// Changes flags: None
	//
	// Description: Pops the value off the top of the stack into the output 8-bit register (arg 1),
	// and decrements the Stack Pointer by 1.
	OP_POP byte = 0x26

	// Args: 1=[source 8-bit reg], 2=[dest 8-bit reg]
	//
	// Changes flags: None
	//
	// Description: Loads the value in the source 8-bit register (arg 1) into the destination 8-bit register (arg 2).
	OP_MOV8 byte = 0x30

	// Args: 1=[source 16-bit reg], 2=[dest 16-bit reg]
	//
	// Changes flags: None
	//
	// Description: Loads the value in the source 16-bit register (arg 1) into the destination 16-bit register (arg 2).
	OP_MOV16 byte = 0x31

	// Args: 1=[low byte of addr], 2=[high byte of addr]
	//
	// Changes flags: None
	//
	// Description: Jumps execution to the given address.
	OP_JMPI byte = 0x40

	// Args: 1=[low byte of addr], 2=[high byte of addr]
	//
	// Changes flags: None
	//
	// Description: Pushes the current Program Counter value onto the stack (first the low byte, then the high byte),
	// and then jumps execution to the given address.
	OP_CALI byte = 0x41

	// Args: 1=[low byte of addr], 2=[high byte of addr]
	//
	// Changes flags: None
	//
	// Description: Jumps execution to the address stored in the J register.
	OP_JMP byte = 0x42

	// Args: 1=[low byte of addr], 2=[high byte of addr]
	//
	// Changes flags: None
	//
	// Description: Pushes the current Program Counter value onto the stack (first the low byte, then the high byte),
	// and then jumps execution to the address stored in the J register.
	OP_CALL byte = 0x43

	// Args: None
	//
	// Changes flags: None
	//
	// Description: Pops the top of the stack into the Program Counter (first the high byte, then the low byte), and then jumps execution to that address.
	OP_RET byte = 0x44

	// Args: 1=[low byte of addr], 2=[high byte of addr]
	//
	// Changes flags: None
	//
	// Description: If the Zero flag is set, jumps execution to the given address; otherwise performs no operation.
	OP_BZSI byte = 0x50

	// Args: 1=[low byte of addr], 2=[high byte of addr]
	//
	// Changes flags: None
	//
	// Description: If the Zero flag is clear, jumps execution to the given address; otherwise performs no operation.
	OP_BZCI byte = 0x51

	// Args: 1=[low byte of addr], 2=[high byte of addr]
	//
	// Changes flags: None
	//
	// Description: If the Carry flag is set, jumps execution to the given address; otherwise performs no operation.
	OP_BCSI byte = 0x52

	// Args: 1=[low byte of addr], 2=[high byte of addr]
	//
	// Changes flags: None
	//
	// Description: If the Carry flag is clear, jumps execution to the given address; otherwise performs no operation.
	OP_BCCI byte = 0x53

	// Args: 1=[low byte of addr], 2=[high byte of addr]
	//
	// Changes flags: None
	//
	// Description: If the Sign flag is set, jumps execution to the given address; otherwise performs no operation.
	OP_BSSI byte = 0x54

	// Args: 1=[low byte of addr], 2=[high byte of addr]
	//
	// Changes flags: None
	//
	// Description: If the Sign flag is clear, jumps execution to the given address; otherwise performs no operation.
	OP_BSCI byte = 0x55

	// Args: 1=[low byte of addr], 2=[high byte of addr]
	//
	// Changes flags: None
	//
	// Description: If the Less Than flag is set, jumps execution to the given address; otherwise performs no operation.
	OP_BLSI byte = 0x56

	// Args: 1=[low byte of addr], 2=[high byte of addr]
	//
	// Changes flags: None
	//
	// Description: If the Less Than flag is clear, jumps execution to the given address; otherwise performs no operation.
	OP_BLCI byte = 0x57

	// Args: 1=[low byte of addr], 2=[high byte of addr]
	//
	// Changes flags: None
	//
	// Description: If the Zero flag is set, jumps execution to the address stored in the J register; otherwise performs no operation.
	OP_BZS byte = 0x60

	// Args: 1=[low byte of addr], 2=[high byte of addr]
	//
	// Changes flags: None
	//
	// Description: If the Zero flag is clear, jumps execution to the address stored in the J register; otherwise performs no operation.
	OP_BZC byte = 0x61

	// Args: 1=[low byte of addr], 2=[high byte of addr]
	//
	// Changes flags: None
	//
	// Description: If the Carry flag is set, jumps execution to the address stored in the J register; otherwise performs no operation.
	OP_BCS byte = 0x62

	// Args: 1=[low byte of addr], 2=[high byte of addr]
	//
	// Changes flags: None
	//
	// Description: If the Carry flag is clear, jumps execution to the address stored in the J register; otherwise performs no operation.
	OP_BCC byte = 0x63

	// Args: 1=[low byte of addr], 2=[high byte of addr]
	//
	// Changes flags: None
	//
	// Description: If the Sign flag is set, jumps execution to the address stored in the J register; otherwise performs no operation.
	OP_BSS byte = 0x64

	// Args: 1=[low byte of addr], 2=[high byte of addr]
	//
	// Changes flags: None
	//
	// Description: If the Sign flag is clear, jumps execution to the address stored in the J register; otherwise performs no operation.
	OP_BSC byte = 0x65

	// Args: 1=[low byte of addr], 2=[high byte of addr]
	//
	// Changes flags: None
	//
	// Description: If the Less Than flag is set, jumps execution to the address stored in the J register; otherwise performs no operation.
	OP_BLS byte = 0x66

	// Args: 1=[low byte of addr], 2=[high byte of addr]
	//
	// Changes flags: None
	//
	// Description: If the Less Than flag is clear, jumps execution to the address stored in the J register; otherwise performs no operation.
	OP_BLC byte = 0x67
)

var OpcodeNames = map[byte]string{
	0x00: "NOP",
	0x10: "CMP",
	0x11: "ADD",
	0x12: "SUB",
	0x13: "INC",
	0x14: "DEC",
	0x15: "NAND",
	0x16: "SHR",
	0x17: "SHL",
	0x20: "LDI",
	0x21: "LDMI",
	0x22: "LDM",
	0x23: "STOI",
	0x24: "STO",
	0x25: "PUSH",
	0x26: "POP",
	0x30: "MOV8",
	0x31: "MOV16",
	0x40: "JMPI",
	0x41: "CALLI",
	0x42: "JMP",
	0x43: "CALL",
	0x44: "RET",
	0x50: "BZSI",
	0x51: "BZCI",
	0x52: "BCSI",
	0x53: "BCCI",
	0x54: "BSSI",
	0x55: "BSCI",
	0x56: "BLSI",
	0x57: "BLCI",
	0x60: "BZS",
	0x61: "BZC",
	0x62: "BCS",
	0x63: "BCC",
	0x64: "BSS",
	0x65: "BSC",
	0x66: "BLS",
	0x67: "BLC",
}
