package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func inWrapper(in In, done In) Out {
	out := make(Bi)
	go func() {
		defer func() {
			close(out)
			//nolint
			for range in {
			}
		}()
		for {
			select {
			case v, ok := <-in:
				if !ok {
					return
				}
				out <- v
			case <-done:
				return
			}
		}
	}()
	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		in = inWrapper(in, done)
		in = stage(in)
	}

	return in
}
