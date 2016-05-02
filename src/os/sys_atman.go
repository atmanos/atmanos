package os

func init() {
	Args = []string{"atmanos"}
}

func hostname() (name string, err error) {
	return "atman", nil
}
