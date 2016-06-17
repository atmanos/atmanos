TEXT ·MemoryBarrier(SB), $0
	MFENCE
	RET

TEXT ·MemoryBarrierWrite(SB), $0
	SFENCE
	RET
