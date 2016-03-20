#include "textflag.h"

// func clearbit(addr *uint64, n uint32)
TEXT ·clearbit(SB),NOSPLIT,$0-12
	MOVQ	addr+0(FP), AX
	MOVL	n+8(FP), BX
	LOCK
	BTRQ	BX, (AX)
	RET

TEXT ·hypervisorEventCallback(SB),NOSPLIT,$0
	MOVQ	0(SP), CX
	MOVQ	8(SP), R11
	ADDQ	$16, SP	// ignore CX, R11

	// Store $0 in code slot
	SUBQ	$8, SP
	MOVQ	$0x0, 0(SP)	// $0 in code

	// Save registers to stack
	SUBQ	$15*8, SP
	MOVQ	DI, 14*8(SP)
	MOVQ	SI, 13*8(SP)
	MOVQ	DX, 12*8(SP)
	MOVQ	CX, 11*8(SP)
	MOVQ	AX, 10*8(SP)
	MOVQ	R8, 9*8(SP)
	MOVQ	R9, 8*8(SP)
	MOVQ	R10, 7*8(SP)
	MOVQ	R11, 6*8(SP)
	MOVQ	BX, 5*8(SP)
	MOVQ	BP, 4*8(SP)
	MOVQ	R12, 3*8(SP)
	MOVQ	R13, 2*8(SP)
	MOVQ	R14, 1*8(SP)
	MOVQ	R15, 0*8(SP)

	// Call handleHypervisorCallback with pointer
	// to registers
	MOVQ	SP, AX
	SUBQ	$8, SP
	MOVQ	AX, 0(SP)
	CALL	·handleHypervisorCallback(SB)
	ADDQ	$8, SP

	// restore registers
	MOVQ	0*8(SP), R15
	MOVQ	1*8(SP), R14
	MOVQ	2*8(SP), R13
	MOVQ	3*8(SP), R12
	MOVQ	4*8(SP), BP
	MOVQ	5*8(SP), BX
	MOVQ	6*8(SP), R11
	MOVQ	7*8(SP), R10
	MOVQ	8*8(SP), R9
	MOVQ	9*8(SP), R8
	MOVQ	10*8(SP), AX
	MOVQ	11*8(SP), CX
	MOVQ	12*8(SP), DX
	MOVQ	13*8(SP), SI
	MOVQ	14*8(SP), DI
	ADDQ	$15*8, SP	// registers

	ADDQ	$0x8, SP	// ignore code

	PUSHQ	AX
	MOVQ	·_atman_shared_info(SB), AX
	MOVB	$0x0, 0x1(AX)
	POPQ	AX

	// Correct CS and SS for return to kernel space
	ORB	$3, 1*8(SP)
	ORB	$3, 4*8(SP)

	IRETQ
	RET

TEXT ·hypervisorFailsafeCallback(SB),NOSPLIT,$0
	IRETQ
	RET
