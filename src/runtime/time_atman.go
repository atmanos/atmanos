package runtime

//go:nosplit
func lfence()

var shadowTimeInfo timeInfo

//go:nosplit
func _nanotime() (ns int64) {
	systemstack(func() {
		shadowTimeInfo.load(_atman_shared_info)
		ns = shadowTimeInfo.nanotime()
	})

	return ns
}

//go:nosplit
func _time_now() (int64, int32) {
	shadowTimeInfo.load(_atman_shared_info)

	return shadowTimeInfo.timeNow()
}

// timeInfo shadows time-related values stored in xenSharedInfo
// and vcpuTimeInfo structures.
type timeInfo struct {
	BootSec  int64
	BootNsec int64

	SystemNsec uint64
	TSC        uint64
	TSCMul     uint32
	TSCShift   int8

	Version uint32
}

// load atomically populates t from info.
func (t *timeInfo) load(info *xenSharedInfo) {
	src := &info.VCPUInfo[0].Time

	if t.Version == atomicload(&src.Version) {
		return
	}

	t.BootSec = int64(info.WcSec)
	t.BootNsec = int64(info.WcNsec)

	for {
		t.Version = atomicload(&src.Version)

		lfence()
		t.SystemNsec = src.SystemNsec
		t.TSC = src.TSC
		t.TSCMul = src.TSCMul
		t.TSCShift = src.TSCShift
		lfence()

		newVersion := atomicload(&src.Version)

		if newVersion&1 == 1 {
			continue
		}

		if t.Version == newVersion {
			return
		}
	}
}

func (t *timeInfo) nsSinceSystem() int64 {
	diff := uint64(cputicks()) - t.TSC

	if t.TSCShift < 0 {
		diff >>= uint8(-t.TSCShift)
	} else {
		diff <<= uint8(t.TSCShift)
	}

	diff *= uint64(t.TSCMul)
	diff >>= 32

	return int64(diff)
}

func (t *timeInfo) nanotime() int64 {
	return int64(t.SystemNsec) + t.nsSinceSystem()
}

func (t *timeInfo) timeNow() (int64, int32) {
	var (
		sec  = t.BootSec
		nsec = t.BootNsec + t.nanotime()
	)

	// move whole seconds to second counter
	sec += nsec / 1e9
	nsec %= 1e9

	return sec, int32(nsec)
}
