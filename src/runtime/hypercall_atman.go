package runtime

import "unsafe"

func hypercall(trap, a1, a2, a3, a4, a5, a6 uintptr) uintptr

func HYPERVISOR_console_io(op uint64, size uint64, data uintptr) uintptr {
	const _HYPERVISOR_console_io = 18

	return hypercall(
		_HYPERVISOR_console_io,
		uintptr(op),
		uintptr(size),
		data,
		0,
		0,
		0,
	)
}

func HYPERVISOR_event_channel_op(cmd uintptr, op unsafe.Pointer) uintptr {
	const _HYPERVISOR_event_channel_op = 32

	return hypercall(
		_HYPERVISOR_event_channel_op,
		cmd,
		uintptr(op),
		0,
		0,
		0,
		0,
	)
}

func HYPERVISOR_update_va_mapping(vaddr uintptr, val uintptr, flags uint64) uintptr {
	const _HYPERVISOR_update_va_mapping = 14

	return hypercall(
		_HYPERVISOR_update_va_mapping,
		vaddr,
		val,
		uintptr(flags),
		0,
		0,
		0,
	)
}

type mmuUpdate struct {
	ptr uintptr // machine address of PTE
	val uintptr // contents of new PTE
}

const DOMID_SELF = 0x7FF0

func HYPERVISOR_mmu_update(updates []mmuUpdate, domid uint16) uintptr {
	const _HYPERVISOR_mmu_update = 1

	return hypercall(
		_HYPERVISOR_mmu_update,
		uintptr(unsafe.Pointer(&updates[0])),
		uintptr(len(updates)),
		0, // done_out (unused)
		uintptr(domid),
		0,
		0,
	)
}

type mmuExtOp struct {
	cmd  uint32
	_    [4]byte
	arg1 uint64
	arg2 uint64
}

func HYPERVISOR_mmuext_op(ops []mmuExtOp, domid uint16) uintptr {
	const _HYPERVISOR_mmuext_op = 26

	return hypercall(
		_HYPERVISOR_mmuext_op,
		uintptr(unsafe.Pointer(&ops[0])),
		uintptr(len(ops)),
		0, // done_out (unused)
		uintptr(domid),
		0,
		0,
	)
}

func HYPERVISOR_sched_op(op uintptr, arg unsafe.Pointer) uintptr {
	const _HYPERVISOR_sched_op = 29

	return hypercall(
		_HYPERVISOR_sched_op,
		op,
		uintptr(arg),
		0,
		0,
		0,
		0,
	)
}
