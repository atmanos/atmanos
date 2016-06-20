package hypercall

import "unsafe"

func HypercallRaw(trap, a1, a2, a3, a4, a5, a6 uintptr) uintptr

const (
	GNTTABOP_setup_table = 2
	GTF_permit_access    = 1
)

const (
	GTF_readonly = 1 << 2
	GTF_reading  = 1 << 3
	GTF_writing  = 1 << 4
)

func GrantTableOp(cmd uintptr, op unsafe.Pointer) uintptr {
	const _grant_table_op = 20

	return HypercallRaw(
		_grant_table_op,
		cmd,
		uintptr(op),
		1,
		0,
		0,
		0,
	)
}

const (
	EVTCHNOP_alloc_unbound = 6
)

func EventChannelOp(cmd uintptr, op unsafe.Pointer) uintptr {
	const _event_channel_op = 32

	return HypercallRaw(
		_event_channel_op,
		cmd,
		uintptr(op),
		0,
		0,
		0,
		0,
	)
}
