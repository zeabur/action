package zbaction

// CleanupStack is a stack of cleanup functions,
// which acts like a `defer` block.
//
// The zero value of CleanupStack is safe to use.
type CleanupStack struct {
	fn []CleanupFn `exhaustruct:"optional"`
}

func (cs *CleanupStack) Push(fn CleanupFn) {
	if cs == nil {
		return
	}
	cs.fn = append(cs.fn, fn)
}

func (cs *CleanupStack) Pop() CleanupFn {
	if cs == nil {
		return nil
	}
	if len(cs.fn) == 0 {
		return nil
	}

	fn := cs.fn[len(cs.fn)-1]
	cs.fn = cs.fn[:len(cs.fn)-1]
	return fn
}

func (cs *CleanupStack) Run() {
	if cs == nil {
		return
	}

	for {
		fn := cs.Pop()
		if fn == nil {
			return
		}
		fn()
	}
}

func (cs *CleanupStack) WrapRun() CleanupFn {
	return func() {
		cs.Run()
	}
}
