package scheduler

import (
	"time"
	_ "unsafe"

	"github.com/df-mc/dragonfly/server/world"
)

type Scheduler[P, R any] struct {
	*world.World
}

func (s *Scheduler[P, R]) Schedule(duration time.Duration, fn world.ExecFunc) {
	time.Sleep(duration)
	s.Exec(fn)
}

func (s *Scheduler[P, R]) ReturnExec(f ReturnExecFunc[R]) R {
	c := make(chan struct{})
	t := &returnTransaction[R]{c: c, f: f}
	pushToQueue(s.World, t)
	<-c
	return t.returnValue
}

type ReturnExecFunc[T any] func(tx *world.Tx) T
type returnTransaction[T any] struct {
	returnValue T
	c           chan struct{}
	f           ReturnExecFunc[T]
}

func (rtx *returnTransaction[T]) Run(w *world.World) {
	tx := &Tx{w: w}
	rtx.returnValue = rtx.f(tx.Get())
	tx.close()
	close(rtx.c)
}

func (s *Scheduler[P, R]) ParamExec(f ParamExecFunc[P], param P) chan struct{} {
	c := make(chan struct{})
	t := paramTransaction[P]{c: c, f: f, param: param}
	pushToQueue(s.World, t)
	return c
}

type ParamExecFunc[T any] func(tx *world.Tx, param T)

type paramTransaction[T any] struct {
	param T
	c     chan struct{}
	f     ParamExecFunc[T]
}

func (ptx paramTransaction[T]) Run(w *world.World) {
	tx := &Tx{w: w}
	ptx.f(tx.Get(), ptx.param)
	tx.close()
	close(ptx.c)
}

func (s *Scheduler[P, R]) PipeExec(param P, f PipeExecFunc[P, R]) R {
	c := make(chan struct{})
	t := &pipeTransaction[P, R]{c: c, param: param, f: f}
	pushToQueue(s.World, t)
	<-c
	return t.returnValue
}

type PipeExecFunc[P, R any] func(tx *world.Tx, param P) R

type pipeTransaction[P, R any] struct {
	param       P
	returnValue R
	c           chan struct{}
	f           PipeExecFunc[P, R]
}

func (ptx *pipeTransaction[P, R]) Run(w *world.World) {
	tx := &Tx{w: w}
	ptx.returnValue = ptx.f(tx.Get(), ptx.param)
	tx.close()
	close(ptx.c)
}
