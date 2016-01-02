#include "go_asm.h"
#include "go_tls.h"
#include "textflag.h"

// func taskstart(fn, _, mp, gp unsafe.Pointer)
TEXT ·taskstart(SB),NOSPLIT,$0-24
	MOVQ	(SP), R12
	MOVQ	16(SP), R8
	MOVQ	24(SP), R9

	// set m->procid to current task ID
	MOVQ	$runtime·taskcurrent(SB), BX
	MOVQ	(BX), AX
	MOVQ	AX, m_procid(R8)
	
	// Set FS to point at m->tls.
	LEAQ	m_tls(R8), DI
	CALL	runtime·settls(SB)

	// Set up new stack
	get_tls(CX)
	MOVQ	R8, g_m(R9)
	MOVQ	R9, g(CX)
	CALL	runtime·stackcheck(SB)

	// Call fn
	CALL	R12

	// Exit if function returns
	CALL	runtime·taskexit(SB)

	RET // unreachable

// func contextsave(*Context) int
TEXT ·contextsave(SB),NOSPLIT,$0-16
	MOVQ	ctx+0(FP), DI
	MOVQ	(SP), CX
	MOVQ	CX, 128(DI)	// save ip to rip
	MOVQ	SP, 152(DI)	// save sp to rsp

	MOVQ	$0xc0000100, CX	// MSR_FS_BASE
	RDMSR
	SHLQ	$32, DX	// DX <<= 32
	ADDQ	DX, AX	// AX = DX + AX
	MOVQ	AX, 184(DI)	// save tls

	MOVQ	$0, ret+8(FP)
	RET

// func contextload(*Context)
TEXT ·contextload(SB),NOSPLIT,$0-8
	MOVQ	ctx+0(FP), DI
	MOVQ	152(DI), R8	// save sp
	MOVQ	128(DI), R9	// save ip
	MOVQ	184(DI), DI
	CALL	runtime·settls(SB) // restore tls
	MOVQ	R8, SP
	MOVQ	R9, (SP)	// set return address
	MOVQ	$1, 16(SP)	// set return value
	RET
