package cpu

import "errors"

// Current Program Status Register

func cpsrCanExecute(cpsr uint32, cond Condition) bool {
	switch cond {
	case GE:
		return cpsrSigned(cpsr) == cpsrOverflow(cpsr)
	case AL:
		return true
	default:
		panic(errors.New("todo"))
	}
}

func cpsrSigned(cpsr uint32) bool {
	return cpsr&0x8000 != 0
}

func cpsrOverflow(cpsr uint32) bool {
	return cpsr&0x1000 != 0
}
