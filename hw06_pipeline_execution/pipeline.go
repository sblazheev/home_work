package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	ch := in
	for _, s := range stages {
		ch = wrapperDone(s, ch, done)
	}

	out := make(Bi)
	go func() {
		defer close(out)
		defer func() {
			//nolint
			for range in {
			}
		}()
		for {
			select {
			case <-done:
				return
			case v, ok := <-ch:
				if !ok {
					return
				}
				out <- v
			}
		}
	}()
	return out
}

func wrapperDone(stage Stage, in In, done In) Out {
	out := make(Bi, 1)
	stageOut := stage(in)
	go func() {
		defer close(out)
		defer func() {
			//nolint
			for range stageOut {
			}
		}()
		for {
			select {
			case <-done:
				return
			default:
				select {
				case <-done:
					return
				case v, ok := <-stageOut:
					if !ok {
						return
					}
					out <- v
				}
			}
		}
	}()
	return out
}
