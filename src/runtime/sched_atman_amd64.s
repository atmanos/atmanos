#include "atman_asm.h"
#include "go_asm.h"
#include "go_tls.h"
#include "textflag.h"

// func taskstart(fn, mp, gp unsafe.Pointer)
TEXT ·taskstart(SB),NOSPLIT,$0-24
	MOVQ	(SP), R12
	MOVQ	8(SP), R8
	MOVQ	16(SP), R9

	// set m->procid to current task ID
	MOVQ	$runtime·taskcurrent(SB), BX
	MOVQ	(Task_ID)(BX), AX
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

// func contextsave(ctx *Context, after uintptr)
TEXT ·contextsave(SB),NOSPLIT,$0-16
	MOVQ	ctx+0(FP), DI

	MOVQ	(SP), CX
	MOVQ	CX, (Context_r+cpuRegisters_rip)(DI)

	MOVW	CS, (Context_r+cpuRegisters_cs)(DI)
	MOVW	SS, (Context_r+cpuRegisters_ss)(DI)

	PUSHFQ
	MOVQ	(SP), CX
	MOVQ	CX, (Context_r+cpuRegisters_rflags)(DI)
	POPFQ

	MOVQ	SP, CX	// Save SP without return address
	ADDQ	$8, CX
	MOVQ	CX, (Context_r+cpuRegisters_rsp)(DI)

	READ_FS_BASE(BX)
	MOVQ	BX, (Context_tls)(DI)

	MOVQ	after+8(FP), CX
	CMPQ	CX, $0
	JZ	skipcallback
	CALL	CX

skipcallback:

	RET

// func contextload(*Context)
TEXT ·contextload(SB),NOSPLIT,$0-8
	MOVQ	ctx+0(FP), AX
	MOVQ	(Context_tls)(AX), DI
	CALL	runtime·settls(SB)
	MOVQ	ctx+0(FP), SP
	RESTORE_ALL
	ADDQ	$(cpuRegisters_code+8), SP
	IRETQ
