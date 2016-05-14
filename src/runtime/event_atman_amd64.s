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
	MOVQ	DI, (cpuRegisters_rdi)(SP)
	MOVQ	SI, (cpuRegisters_rsi)(SP)
	MOVQ	DX, (cpuRegisters_rdx)(SP)
	MOVQ	CX, (cpuRegisters_rcx)(SP)
	MOVQ	AX, (cpuRegisters_rax)(SP)
	MOVQ	R8, (cpuRegisters_r8)(SP)
	MOVQ	R9, (cpuRegisters_r9)(SP)
	MOVQ	R10, (cpuRegisters_r10)(SP)
	MOVQ	R11, (cpuRegisters_r11)(SP)
	MOVQ	BX, (cpuRegisters_rbx)(SP)
	MOVQ	BP, (cpuRegisters_rbp)(SP)
	MOVQ	R12, (cpuRegisters_r12)(SP)
	MOVQ	R13, (cpuRegisters_r13)(SP)
	MOVQ	R14, (cpuRegisters_r14)(SP)
	MOVQ	R15, (cpuRegisters_r15)(SP)

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

	// Restore registers from stack
	MOVQ	(cpuRegisters_rdi)(SP), DI
	MOVQ	(cpuRegisters_rsi)(SP), SI
	MOVQ	(cpuRegisters_rdx)(SP), DX
	MOVQ	(cpuRegisters_rcx)(SP), CX
	MOVQ	(cpuRegisters_rax)(SP), AX
	MOVQ	(cpuRegisters_r8)(SP), R8
	MOVQ	(cpuRegisters_r9)(SP), R9
	MOVQ	(cpuRegisters_r10)(SP), R10
	MOVQ	(cpuRegisters_r11)(SP), R11
	MOVQ	(cpuRegisters_rbx)(SP), BX
	MOVQ	(cpuRegisters_rbp)(SP), BP
	MOVQ	(cpuRegisters_r12)(SP), R12
	MOVQ	(cpuRegisters_r13)(SP), R13
	MOVQ	(cpuRegisters_r14)(SP), R14
	MOVQ	(cpuRegisters_r15)(SP), R15
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
