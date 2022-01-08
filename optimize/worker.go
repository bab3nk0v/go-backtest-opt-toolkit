package optimize

import (
	"fmt"
	"github.com/c-bata/goptuna"
	"runtime"
	"trade-optimizer/util"
)

type Worker struct {
	workerId int
	ctx      *OptimizeContext
	study    *goptuna.Study
}

func (w *Worker) run() error {
	defer func() {
		problem := recover()
		if e, ok := problem.(runtime.Error); ok {
			if e.Error() == "send on closed channel" {
				println(fmt.Sprintf("The system has closed the channel, terminating the worker %d", w.workerId))
			} else {
				panic(e)
			}
		} else {
			panic(problem)
		}
	}()
	return w.study.Optimize(
		w.ctx.computeObjective,
		util.MaxInt,
	)
}

func newWorker(workerId int, ctx *OptimizeContext, study *goptuna.Study) *Worker {
	return &Worker{
		workerId: workerId,
		ctx:      ctx,
		study:    study,
	}
}
