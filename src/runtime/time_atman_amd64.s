TEXT ·lfence(SB),NOSPLIT,$0
	BYTE	$0x0f; BYTE $0xae; BYTE $0xe8 // LFENCE
	RET
