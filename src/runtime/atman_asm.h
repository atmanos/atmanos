#define _PAGE_ROUND_UP(REGISTER) \
	ADDQ	$0x0000000000000fff, REGISTER	\
	ANDQ	$0xfffffffffffff000, REGISTER

#define HYPERCALL(TRAP) \
	MOVQ	TRAP, CX				\
	IMULQ	$32, CX					\
	MOVQ	$runtime·_atman_hypercall_page(SB), BX	\
	_PAGE_ROUND_UP(BX)				\
	ADDQ	CX, BX					\
        CALL    BX

#define READ_FS_BASE(r) \
	MOVQ	$0xc0000100, CX	\ // MSR_FS_BASE
	RDMSR			\
	SHLQ	$32, DX		\ // DX <<= 32
	ADDQ	DX, AX		\ // AX = DX + AX
	MOVQ	AX, r

#define RESTORE_ALL \
	MOVQ	(cpuRegisters_rdi)(SP), DI	\
	MOVQ	(cpuRegisters_rsi)(SP), SI	\
	MOVQ	(cpuRegisters_rdx)(SP), DX	\
	MOVQ	(cpuRegisters_rcx)(SP), CX	\
	MOVQ	(cpuRegisters_rax)(SP), AX	\
	MOVQ	(cpuRegisters_r8)(SP), R8	\
	MOVQ	(cpuRegisters_r9)(SP), R9	\
	MOVQ	(cpuRegisters_r10)(SP), R10	\
	MOVQ	(cpuRegisters_r11)(SP), R11	\
	MOVQ	(cpuRegisters_rbx)(SP), BX	\
	MOVQ	(cpuRegisters_rbp)(SP), BP	\
	MOVQ	(cpuRegisters_r12)(SP), R12	\
	MOVQ	(cpuRegisters_r13)(SP), R13	\
	MOVQ	(cpuRegisters_r14)(SP), R14	\
	MOVQ	(cpuRegisters_r15)(SP), R15	\

#define SAVE_ALL \
	MOVQ	DI, (cpuRegisters_rdi)(SP)	\
	MOVQ	SI, (cpuRegisters_rsi)(SP)	\
	MOVQ	DX, (cpuRegisters_rdx)(SP)	\
	MOVQ	CX, (cpuRegisters_rcx)(SP)	\
	MOVQ	AX, (cpuRegisters_rax)(SP)	\
	MOVQ	R8, (cpuRegisters_r8)(SP)	\
	MOVQ	R9, (cpuRegisters_r9)(SP)	\
	MOVQ	R10, (cpuRegisters_r10)(SP)	\
	MOVQ	R11, (cpuRegisters_r11)(SP)	\
	MOVQ	BX, (cpuRegisters_rbx)(SP)	\
	MOVQ	BP, (cpuRegisters_rbp)(SP)	\
	MOVQ	R12, (cpuRegisters_r12)(SP)	\
	MOVQ	R13, (cpuRegisters_r13)(SP)	\
	MOVQ	R14, (cpuRegisters_r14)(SP)	\
	MOVQ	R15, (cpuRegisters_r15)(SP)
