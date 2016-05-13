#include "textflag.h"

TEXT ·lfence(SB),NOSPLIT,$0
	LFENCE
	RET
