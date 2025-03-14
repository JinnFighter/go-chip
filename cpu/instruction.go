package cpu

type IInstruction interface {
	SetupValues(instructionBytes uint16)
	Execute(cpu *CpuInstance)
}
