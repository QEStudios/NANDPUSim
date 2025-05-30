package main

func main() {
	nandpu := NewNANDPU([]byte{
		OP_CMP,
		OP_ADD,
	})
	nandpu.Step()
	nandpu.Step()
}
