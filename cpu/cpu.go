package cpu

import (
	"encoding/binary"
	"errors"
	"log"
	"strconv"
)

var (
	ErrNotImplementedYet = errors.New("not implemented yet")
)

const (
	OPCODE_ARM_SIZE = 4
)

type Condition uint8

type OpCode uint32

type InstructionKind uint8

const (
	Branch InstructionKind = iota
	BranchLink
	Mov
)

func (me InstructionKind) String() string {
	return [...]string{"Branch", "BranchLink"}[me]
}

const (
	EQ Condition = iota
	NE
	CS
	CC
	MI
	PL
	VS
	VC
	HI
	LS
	GE
	LT
	GT
	LE
	AL
	NV
)

func armMode(opCode OpCode) InstructionKind {
	if 0x0F_00_00_00&uint32(opCode) == 0x0A_00_00_00 {
		return Branch
	}
	if 0x0F_00_00_00&uint32(opCode) == 0x0B_00_00_00 {
		return Branch
	}
	if 0x0D_EF_00_00&uint32(opCode) == 0x01_A0_00_00 {
		return Mov
	}
	panic(ErrNotImplementedYet)
}

func (me Condition) String() string {
	return [...]string{
		"EQ",
		"NE",
		"CS",
		"CC",
		"MI",
		"PL",
		"VS",
		"VC",
		"HI",
		"LS",
		"GE",
		"LT",
		"GT",
		"LE",
		"AL",
		"NV",
	}[me]
}

type Cpu struct {
	rom            []byte
	programCounter int
	cpsr           uint32
	registers      [16]uint32
}

func NewCpu(rom []byte) *Cpu {
	return &Cpu{
		rom:            rom,
		programCounter: 0,
		cpsr:           0,
	}
}

func (cpu *Cpu) branch(opCode OpCode) {
	offset := opCode & 0x00_FF_FF_FF
	log.Printf("offset: %d", offset)

	cpu.programCounter += 8 + int(offset)*4
	log.Printf("PC: %d", cpu.programCounter)
}

func (cpu *Cpu) branchLink(opCode OpCode) {
	panic(ErrNotImplementedYet)
}

func (cpu *Cpu) mov(opCode OpCode) {
	// bits [24-21] are the RD
	rd := (opCode & 0x00_00_F0_00) >> 12
	log.Printf("RD: %d", rd)

	// 25th bit is I = Immediate Flag
	immediate := (opCode&0x02_00_00_00)>>25 != 0
	log.Printf("Immediate: %v", immediate)

	// 20th bit is S = Condition Set
	if opCode&0x00_08_00_00 > 0 {
		panic(errors.New("condition set"))
	}

	if immediate {
		// bits [7-0] are the immediate value
		immediate_value := opCode & 0x00_00_00_FF
		log.Printf("Value: %d", immediate_value)

		// the instruction is MOV RD, immediate_value
		cpu.registers[rd] = uint32(immediate_value)
	} else {
		panic(ErrNotImplementedYet)
	}

	// N.B: I'm not sure where this has to be executed
	cpu.programCounter += int(OPCODE_ARM_SIZE)
}

func (cpu *Cpu) Execute(opCode OpCode, instruction InstructionKind) {
	switch instruction {
	case Branch:
		cpu.branch(opCode)
	case BranchLink:
		cpu.branchLink(opCode)
	case Mov:
		cpu.mov(opCode)
	default:
		panic(ErrNotImplementedYet)
	}
}

func (cpu *Cpu) Step() {
	opCode := cpu.Fetch()
	condition, instruction := cpu.Decode(opCode)

	if cpsrCanExecute(uint32(opCode), condition) {
		cpu.Execute(opCode, instruction)
	}
}

func (cpu *Cpu) Fetch() OpCode {
	instructionIndex := cpu.programCounter
	endInstruction := instructionIndex + OPCODE_ARM_SIZE
	op := cpu.rom[instructionIndex:endInstruction]
	opBE := binary.LittleEndian.Uint32(op)
	log.Printf("fetch op -> %s\n", strconv.FormatUint(uint64(opBE), 2))
	return OpCode(opBE)
}

func (cpu *Cpu) Decode(opCode OpCode) (Condition, InstructionKind) {
	condition := Condition(opCode >> 28) // latest 4 bit (28..32)
	log.Printf("condition -> %v", condition)
	instruction := armMode(opCode)
	log.Printf("instruction -> %v", instruction)
	return condition, instruction
}
