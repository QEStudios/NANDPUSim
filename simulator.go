package main

func main() {
	nandpu := NewNANDPU([]byte{
		OP_CMP,
		OP_INC, REG_B,
		OP_INC, REG_B,
		OP_INC, REG_B,
		// Register B should now hold the value 3
		OP_ADD, REG_C,
		OP_ADD, REG_D,
		OP_ADD, REG_C,
		OP_SUB, REG_A,
		OP_SPECIAL_HALT,
	})

	running := true
	for running {
		running = nandpu.Step()
	}
}
