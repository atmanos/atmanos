package runtime

//go:nosplit
func lfence()

//go:nosplit
func _nanotime() (ns int64) {
	systemstack(func() {
		ns = shadowTimeInfo.nanotime()
	})

	return ns
}

//go:nosplit
func _time_now() (int64, int32) {
	return shadowTimeInfo.timeNow()
}

var shadowTimeInfo timeInfo

// timeInfo shadows time-related values stored in xenSharedInfo
// and vcpuTimeInfo structures.
type timeInfo struct {
	BootVersion uint32
	BootSec     int64
	BootNsec    int64

	SystemVersion uint32
	SystemNsec    uint64
	TSC           uint64
	TSCMul        uint32
	TSCShift      int8
}

// nanotime returns a monotonically increasing nanosecond time value.
func (t *timeInfo) nanotime() int64 {
	t.checkSystemTime()

	return int64(t.SystemNsec) + t.nsSinceSystem()
}

// nsSinceSystem returns the nanoseconds that have elapsed since
// t.SystemNsec was stored.
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

// checkSystemTime ensures the system clock values are up-to-date.
func (t *timeInfo) checkSystemTime() {
	src := &_atman_shared_info.VCPUInfo[0].Time

	for t.needsUpdate(t.SystemVersion, &src.Version) {
		t.SystemVersion = atomicload(&src.Version)

		lfence()
		t.SystemNsec = src.SystemNsec
		t.TSC = src.TSC
		t.TSCMul = src.TSCMul
		t.TSCShift = src.TSCShift
		lfence()
	}
}

func (t *timeInfo) needsUpdate(shadow uint32, src *uint32) bool {
	latest := atomicload(src)

	return shadow != latest || latest&1 == 1
}

func (t *timeInfo) timeNow() (int64, int32) {
	t.checkBootTime()

	var (
		sec  = t.BootSec
		nsec = t.BootNsec + t.nanotime()
	)

	// move whole seconds to second counter
	sec += nsec / 1e9
	nsec %= 1e9

	return sec, int32(nsec)
}

// checkBootTime ensures the boot (wall) clock values are up-to-date.
func (t *timeInfo) checkBootTime() {
	src := _atman_shared_info

	for t.needsUpdate(t.BootVersion, &src.WcVersion) {
		t.BootVersion = atomicload(&src.WcVersion)

		lfence()
		t.BootSec = int64(src.WcSec)
		t.BootNsec = int64(src.WcNsec)
		lfence()
	}
}
