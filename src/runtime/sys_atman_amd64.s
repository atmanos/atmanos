#include "atman_asm.h"
#include "textflag.h"

TEXT runtime·exit(SB),NOSPLIT,$8-4
retry:
	MOVQ	$2, DI	// SCHEDOP_shutdown
	MOVQ	$0, (SP)
	MOVQ	SP, SI	// *reason (0 SHUTDOWN_poweroff)
	HYPERCALL($29)
	JMP	retry
	RET

// func crash()
TEXT runtime·crash(SB),NOSPLIT,$8-4
retry:
	MOVQ	$2, DI	// SCHEDOP_shutdown
	MOVQ	$3, (SP)
	MOVQ	SP, SI	// *reason (3 SHUTDOWN_crash)
	HYPERCALL($29)
	JMP	retry
	RET

// func usleep(ms uint32)
TEXT runtime·usleep(SB),NOSPLIT,$0-4
	MOVQ	$runtime·tasksleepus(SB), AX
	JMP	AX

TEXT runtime·open(SB),NOSPLIT,$0
	RET

TEXT runtime·read(SB),NOSPLIT,$0
	RET

TEXT runtime·closefd(SB),NOSPLIT,$0
	RET

TEXT runtime·write(SB),NOSPLIT,$0-28
	RET

// func nanotime() int64
TEXT runtime·nanotime(SB),NOSPLIT,$0-8
	MOVQ	$runtime·_nanotime(SB), AX
	JMP	AX

// func now() (sec int64, nsec int32)
TEXT time·now(SB),NOSPLIT,$0-16
	MOVQ	$runtime·_time_now(SB), AX
	JMP	AX

// set tls base to DI
TEXT runtime·settls(SB),NOSPLIT,$0
	MOVQ	DI, CX
	MOVQ	$0, DI	// SEGBASE_FS
	MOVQ	CX, SI	// TLS address
	MOVQ	$0, DX	// unused
	HYPERCALL($25)
	RET

TEXT runtime·hypercall(SB),NOSPLIT,$0
	MOVQ	a1+8(FP), DI
	MOVQ	a2+16(FP), SI
	MOVQ	a3+24(FP), DX
	MOVQ	a4+32(FP), R10
	MOVQ	a5+40(FP), R8
	MOVQ	a6+48(FP), R9
	HYPERCALL(trap+0(FP))
	MOVQ	AX, ret+56(FP)
	RET
