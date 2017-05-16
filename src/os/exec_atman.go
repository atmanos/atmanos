package os

import "time"

var Kill Signal

func (p *Process) signal(sig Signal) error {
	return nil
}

func (*Process) wait() (*ProcessState, error) {
	return nil, ErrInvalid
}

func (*Process) release() error {
	return nil
}

func (*Process) kill() error {
	return nil
}

func findProcess(pid int) (p *Process, err error) {
	return newProcess(pid, 0), nil
}

func startProcess(name string, argv []string, attr *ProcAttr) (p *Process, err error) {
	return nil, ErrNotExist
}

type ProcessState struct{}

func (p *ProcessState) Pid() int {
	return -1
}

func (p *ProcessState) exited() bool {
	return true
}

func (p *ProcessState) success() bool {
	return false
}

func (p *ProcessState) sys() interface{} {
	return nil
}

func (p *ProcessState) sysUsage() interface{} {
	return nil
}

func (p *ProcessState) userTime() time.Duration {
	return 0
}

func (p *ProcessState) systemTime() time.Duration {
	return 0
}

func (p *ProcessState) String() string {
	return "exit status: exited"
}
