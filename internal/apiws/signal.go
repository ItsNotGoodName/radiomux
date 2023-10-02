package apiws

type signal struct {
	C chan struct{}
}

func newSignal() signal {
	c := make(chan struct{}, 1)
	c <- struct{}{}
	return signal{
		C: c,
	}
}

func (s signal) Queue() {
	select {
	case s.C <- struct{}{}:
	default:
	}
}
