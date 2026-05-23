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
			select {
			case <-done:
				go func() {
					for range in {
					}
				}()
				return
			default:

			}

			select {
			case <-done:
				go func() {
					for range in {
					}
				}()
				return

			case v, ok := <-in:
				if !ok {
					return
				}
				select {
				case <-done:
					go func() {
						for range in {
						}
					}()
					return
				case out <- v:
				}
			}
		}
	}()
	return out
}
