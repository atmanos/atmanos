TEXT _rt0_amd64_atman(SB),NOSPLIT,$-8
	CLD
	MOVQ	SI, runtime·_atman_start_info(SB)
	MOVQ	$main(SB), AX
	JMP	AX

TEXT main(SB),NOSPLIT,$-8
	MOVQ	$runtime·rt0_go(SB), AX
	JMP	AX
