package hw06pipelineexecution

import (
	"log"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	for _, s := range stages {
		in = wrapperDone(s, in, done)
	}

	out := make(Bi)
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
				out <- v
			}
		}
	}()
	return out
}

func wrapperDone(stage Stage, in In, done In) Out {
	out := make(Bi)
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
			case v, ok := <-stageOut:
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
