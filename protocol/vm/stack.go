package vm

import (
	"github.com/holiman/uint256"

	"coingod/math/checked"
)

func opToAltStack(vm *virtualMachine) error {
	if err := vm.applyCost(2); err != nil {
		return err
	}

	if len(vm.dataStack) == 0 {
		return ErrDataStackUnderflow
	}
	// no standard memory cost accounting here
	vm.altStack = append(vm.altStack, vm.dataStack[len(vm.dataStack)-1])
	vm.dataStack = vm.dataStack[:len(vm.dataStack)-1]
	return nil
}

func opFromAltStack(vm *virtualMachine) error {
	if err := vm.applyCost(2); err != nil {
		return err
	}

	if len(vm.altStack) == 0 {
		return ErrAltStackUnderflow
	}

	// no standard memory cost accounting here
	vm.dataStack = append(vm.dataStack, vm.altStack[len(vm.altStack)-1])
	vm.altStack = vm.altStack[:len(vm.altStack)-1]
	return nil
}

func op2Drop(vm *virtualMachine) error {
	if err := vm.applyCost(2); err != nil {
		return err
	}

	for i := 0; i < 2; i++ {
		if _, err := vm.pop(false); err != nil {
			return err
		}
	}
	return nil
}

func op2Dup(vm *virtualMachine) error {
	return nDup(vm, 2)
}

func op3Dup(vm *virtualMachine) error {
	return nDup(vm, 3)
}

func nDup(vm *virtualMachine, n int) error {
	if err := vm.applyCost(int64(n)); err != nil {
		return err
	}

	if len(vm.dataStack) < n {
		return ErrDataStackUnderflow
	}

	for i := 0; i < n; i++ {
		if err := vm.pushDataStack(vm.dataStack[len(vm.dataStack)-n], false); err != nil {
			return err
		}
	}
	return nil
}

func op2Over(vm *virtualMachine) error {
	if err := vm.applyCost(2); err != nil {
		return err
	}

	if len(vm.dataStack) < 4 {
		return ErrDataStackUnderflow
	}

	for i := 0; i < 2; i++ {
		if err := vm.pushDataStack(vm.dataStack[len(vm.dataStack)-4], false); err != nil {
			return err
		}
	}
	return nil
}

func op2Rot(vm *virtualMachine) error {
	if err := vm.applyCost(2); err != nil {
		return err
	}

	if len(vm.dataStack) < 6 {
		return ErrDataStackUnderflow
	}

	newStack := make([][]byte, 0, len(vm.dataStack))
	newStack = append(newStack, vm.dataStack[:len(vm.dataStack)-6]...)
	newStack = append(newStack, vm.dataStack[len(vm.dataStack)-4:]...)
	newStack = append(newStack, vm.dataStack[len(vm.dataStack)-6])
	newStack = append(newStack, vm.dataStack[len(vm.dataStack)-5])
	vm.dataStack = newStack
	return nil
}

func op2Swap(vm *virtualMachine) error {
	if err := vm.applyCost(2); err != nil {
		return err
	}

	if len(vm.dataStack) < 4 {
		return ErrDataStackUnderflow
	}

	newStack := make([][]byte, 0, len(vm.dataStack))
	newStack = append(newStack, vm.dataStack[:len(vm.dataStack)-4]...)
	newStack = append(newStack, vm.dataStack[len(vm.dataStack)-2:]...)
	newStack = append(newStack, vm.dataStack[len(vm.dataStack)-4])
	newStack = append(newStack, vm.dataStack[len(vm.dataStack)-3])
	vm.dataStack = newStack
	return nil
}

func opIfDup(vm *virtualMachine) error {
	if err := vm.applyCost(1); err != nil {
		return err
	}

	item, err := vm.top()
	if err != nil {
		return err
	}

	if AsBool(item) {
		return vm.pushDataStack(item, false)
	}
	return nil
}

func opDepth(vm *virtualMachine) error {
	if err := vm.applyCost(1); err != nil {
		return err
	}

	return vm.pushBigInt(uint256.NewInt(uint64(len(vm.dataStack))), false)
}

func opDrop(vm *virtualMachine) error {
	if err := vm.applyCost(1); err != nil {
		return err
	}

	_, err := vm.pop(false)
	return err
}

func opDup(vm *virtualMachine) error {
	return nDup(vm, 1)
}

func opNip(vm *virtualMachine) error {
	if err := vm.applyCost(1); err != nil {
		return err
	}

	top, err := vm.top()
	if err != nil {
		return err
	}

	// temporarily pop off the top value with no standard memory accounting
	vm.dataStack = vm.dataStack[:len(vm.dataStack)-1]
	if _, err = vm.pop(false); err != nil {
		return err
	}
	// now put the top item back
	vm.dataStack = append(vm.dataStack, top)
	return nil
}

func opOver(vm *virtualMachine) error {
	if err := vm.applyCost(1); err != nil {
		return err
	}

	if len(vm.dataStack) < 2 {
		return ErrDataStackUnderflow
	}

	return vm.pushDataStack(vm.dataStack[len(vm.dataStack)-2], false)
}

func opPick(vm *virtualMachine) error {
	if err := vm.applyCost(2); err != nil {
		return err
	}

	n, err := vm.popBigInt(false)
	if err != nil {
		return err
	}

	off, ok := checked.AddInt64(int64(n.Uint64()), 1)
	if !ok {
		return ErrBadValue
	}

	dataStackSize := int64(len(vm.dataStack))
	if dataStackSize < off {
		return ErrDataStackUnderflow
	}

	return vm.pushDataStack(vm.dataStack[dataStackSize-off], false)
}

func opRoll(vm *virtualMachine) error {
	if err := vm.applyCost(2); err != nil {
		return err
	}

	n, err := vm.popBigInt(false)
	if err != nil {
		return err
	}

	off, ok := checked.AddInt64(int64(n.Uint64()), 1)
	if !ok {
		return ErrBadValue
	}

	return rot(vm, off)
}

func opRot(vm *virtualMachine) error {
	if err := vm.applyCost(2); err != nil {
		return err
	}

	return rot(vm, 3)
}

func rot(vm *virtualMachine, n int64) error {
	if n < 1 {
		return ErrBadValue
	}

	if int64(len(vm.dataStack)) < n {
		return ErrDataStackUnderflow
	}

	index := int64(len(vm.dataStack)) - n
	newStack := make([][]byte, 0, len(vm.dataStack))
	newStack = append(newStack, vm.dataStack[:index]...)
	newStack = append(newStack, vm.dataStack[index+1:]...)
	newStack = append(newStack, vm.dataStack[index])
	vm.dataStack = newStack
	return nil
}

func opSwap(vm *virtualMachine) error {
	if err := vm.applyCost(1); err != nil {
		return err
	}

	l := len(vm.dataStack)
	if l < 2 {
		return ErrDataStackUnderflow
	}

	vm.dataStack[l-1], vm.dataStack[l-2] = vm.dataStack[l-2], vm.dataStack[l-1]
	return nil
}

func opTuck(vm *virtualMachine) error {
	if err := vm.applyCost(1); err != nil {
		return err
	}

	if len(vm.dataStack) < 2 {
		return ErrDataStackUnderflow
	}

	top2 := make([][]byte, 2)
	copy(top2, vm.dataStack[len(vm.dataStack)-2:])
	// temporarily remove the top two items without standard memory accounting
	vm.dataStack = vm.dataStack[:len(vm.dataStack)-2]
	if err := vm.pushDataStack(top2[1], false); err != nil {
		return err
	}

	vm.dataStack = append(vm.dataStack, top2...)
	return nil
}
