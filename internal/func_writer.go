package internal

type funcWriter func([]byte) (int, error)

func (f funcWriter) Write(p []byte) (n int, err error) {
	return f(p)
}
