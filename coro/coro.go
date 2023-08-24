package coro

import (
	"slices"
)

const routineCancelled = "coroutine cancelled"

type Yield func()

func New(resume func(yield Yield)) *Routine[struct{}] {
	return WithReturn(func(y YieldReturn[struct{}]) {
		resume(func() {
			y(struct{}{})
		})
	})
}

type YieldReturn[V any] func(V)

func WithReturn[V any](resume func(YieldReturn[V])) *Routine[V] {
	r := &Routine[V]{ // 1 alloc
		resumed: make(chan struct{}), // 1 alloc
		done:    make(chan V),        // 1 alloc
		status:  Suspended,
	}
	go r.start(resume) // 3 allocs

	return r
}

type Routine[V any] struct {
	done    chan V
	resumed chan struct{}
	status  Status
}

func (r *Routine[V]) start(f func(YieldReturn[V])) { // 1 alloc
	defer r.recoverAndDestroy()

	_, ok := <-r.resumed // 2 allocs
	if !ok {
		panic(routineCancelled)
	}

	r.status = Running
	f(r.yield)
}

func (r *Routine[V]) yield(v V) {
	r.done <- v
	r.status = Suspended
	if _, ok := <-r.resumed; !ok {
		panic(routineCancelled)
	}
}

func (r *Routine[V]) recoverAndDestroy() {
	p := recover()
	if p != nil && p != routineCancelled {
		panic("coroutine panicked")
	}
	r.status = Dead
	close(r.done)
}

func (r *Routine[V]) Resume() (value V, hasMore bool) {
	if r.status == Dead {
		return
	}

	r.resumed <- struct{}{}
	value, hasMore = <-r.done
	return
}

func (r *Routine[V]) Status() Status {
	return r.status
}

func (r *Routine[V]) Cancel() {
	if r.status == Dead {
		return
	}

	close(r.resumed)
	<-r.done
}

type Status string

const (
	// Normal    Status = "normal"    // This coroutine is currently waiting in coresume for another coroutine. (Either for the running coroutine, or for another normal coroutine)
	Running   Status = "running"   // This is the coroutine that's currently running - aka the one that just called costatus.
	Suspended Status = "suspended" // This coroutine is not running - either it has yielded or has never been resumed yet.
	Dead      Status = "dead"      // This coroutine has either returned or died due to an error.
)

type Routines []*Routine[struct{}]

func (r Routines) ResumeAll() Routines {
	for _, rout := range r {
		rout.Resume()
	}
	return slices.DeleteFunc(r, func(r *Routine[struct{}]) bool {
		return r.Status() == Dead
	})
}
