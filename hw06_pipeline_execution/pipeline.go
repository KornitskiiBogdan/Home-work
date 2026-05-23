package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		if done == nil {
			return in
		}
		return doneStage(in, done)
	}
	out := in
	for _, stage := range stages {
		if done != nil {
			out = doneStage(out, done)
		}
		out = stage(out)
	}
	if done != nil {
		out = doneStage(out, done)
	}
	return out
}

func doneStage(in In, done In) Out {
	out := make(Bi)
	go func() {
		defer close(out)
		for {
			if isDone(done) {
				drainIn(in)
				return
			}
			v, ok := getOrStop(in, done)
			if !ok {
				return
			}
			if !sendOrStop(out, v, in, done) {
				return
			}
		}
	}()
	return out
}

func getOrStop(in In, done In) (v interface{}, ok bool) {
	select {
	case <-done:
		drainIn(in)
		return nil, false
	case v, ok = <-in:
		return v, ok
	}
}

func sendOrStop(out Bi, v interface{}, in In, done In) bool {
	select {
	case <-done:
		drainIn(in)
		return false
	case out <- v:
		return true
	}
}

func drainIn(in In) {
	go func() {
		for {
			if _, open := <-in; !open {
				return
			}
		}
	}()
}

func isDone(done In) bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}
