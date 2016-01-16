// func clearbit(addr *uint64, n uint32)
TEXT ·clearbit(SB),NOSPLIT,$0-12
	MOVQ	addr+0(FP), BX
	MOVL	n+8(FP), AX
	LOCK
	BTRQ	AX, 0(BX)
