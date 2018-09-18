package oss

type windowsW struct {
}

func (w *windowsW) Info(m string) error {
	return nil
}

func (w *windowsW) Close() error {
	return nil
}

func New(flag Priority, tag string) (writer, error) {
	return new(windowsW), nil
}

func Dial(net, addr string, flag Priority, tag string) (writer, error) {
	return new(windowsW), nil
}