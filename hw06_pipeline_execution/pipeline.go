package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		out := startStage(stage, in, done)
		in = out
	}
	return in
}

func startStage(stage Stage, in In, done In) Out {
	out := make(Bi, 1)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				stageIn := make(Bi, 1)
				stageIn <- v
				close(stageIn)
				stageOut := stage(stageIn)
				out <- <-stageOut
			}
		}
	}()
	return out
}
