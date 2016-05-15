#include "atman_asm.h"
#include "go_asm.h"
#include "go_tls.h"
#include "textflag.h"

// func clearbit(addr *uint64, n uint32)
TEXT ·clearbit(SB),NOSPLIT,$0-12
	MOVQ	addr+0(FP), AX
	MOVL	n+8(FP), BX
	LOCK
	BTRQ	BX, (AX)
	RET

TEXT ·eventCallbackASM(SB),NOSPLIT,$0
	MOVQ	0(SP), CX
	MOVQ	8(SP), R11
	ADDQ	$16, SP	// ignore CX, R11

	// Save registers to stack
	SUBQ	$(cpuRegisters_code+8), SP
	MOVQ	$0x0, (cpuRegisters_code)(SP)
	SAVE_ALL

	// Save stack pointer
	MOVQ	SP, CX

	// Get current g
	get_tls(AX)
	MOVQ	g(AX), BX

	// Switch to gsignal stack
	MOVQ	g_m(BX), BX
	MOVQ	m_gsignal(BX), R10
	MOVQ	(g_stack+stack_hi)(R10), BP
	MOVQ	BP, SP

	// Call handleHypervisorCallback with pointer
	// to registers
	SUBQ	$24, SP
	MOVQ	CX, 0(SP)
	MOVQ	SP, 8(SP)
	CALL	·eventCallback(SB)

	// Restore original stack pointer
	MOVQ	0(SP), SP

	RESTORE_ALL
	ADDQ	$(cpuRegisters_code+8), SP

	PUSHQ	AX
	MOVQ	·_atman_shared_info(SB), AX
	MOVB	$0, (xenSharedInfo_VCPUInfo+vcpuInfo_UpcallMask)(AX)
	POPQ	AX

	// Correct CS and SS for return to kernel space
	ORB	$3, 1*8(SP)
	ORB	$3, 4*8(SP)

	IRETQ
	RET

TEXT ·eventFailsafe(SB),NOSPLIT,$0
	IRETQ
	RET

TEXT ·atmansettls(SB),NOSPLIT,$-8
	MOVQ	tls+0(FP), DI
	CALL runtime·settls(SB)
	RET
