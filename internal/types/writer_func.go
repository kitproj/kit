package types

type writeFunc func(p []byte) (n int, err error)

func (w writeFunc) Write(p []byte) (n int, err error) {
	return w(p)
}
