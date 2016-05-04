package runtime

//go:nosplit
func lfence()

var shadowTimeInfo timeInfo

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

// checkBootTime ensures the boot (wall) clock values are up-to-date.
func (t *timeInfo) checkBootTime() {
	src := _atman_shared_info

	t.check(&t.BootVersion, &src.WcVersion, t.updateBootTime)
}

func (t *timeInfo) updateBootTime() {
	src := _atman_shared_info
	t.BootSec = int64(src.WcSec)
	t.BootNsec = int64(src.WcNsec)
}

// checkMonotonicTime ensures the system clock values are up-to-date.
func (t *timeInfo) checkSystemTime() {
	src := &_atman_shared_info.VCPUInfo[0].Time

	t.check(&t.SystemVersion, &src.Version, t.updateSystemTime)
}

func (t *timeInfo) updateSystemTime() {
	src := &_atman_shared_info.VCPUInfo[0].Time

	t.SystemNsec = src.SystemNsec
	t.TSC = src.TSC
	t.TSCMul = src.TSCMul
	t.TSCShift = src.TSCShift
}

// check atomically syncronizes the shadow and src versions
// calling update if the versions disagree.
func (t *timeInfo) check(shadow, src *uint32, update func()) {
	if *shadow == atomicload(src) {
		return
	}

	for {
		*shadow = atomicload(src)

		lfence()
		update()
		lfence()

		new := atomicload(src)

		if new&1 == 1 {
			continue
		}

		if *shadow == new {
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
	t.checkSystemTime()

	return int64(t.SystemNsec) + t.nsSinceSystem()
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
