package go_then

import (
	"context"
	"sync"
)

type Promiser interface {
	// TODO:
}

type Resolver func(i any)
type Rejector func(i error)

type Promise struct {
	wg             *sync.WaitGroup
	ctx            context.Context
	resolveChannel chan any
	rejectChannel  chan error
	callback       func(i any)
	errorHandler   func(err error)
	executor       func(Resolver, Rejector)
}

func New(ctx context.Context, executor func(resolve Resolver, reject Rejector)) *Promise {

	p := new(Promise)
	p.ctx = ctx
	p.executor = executor

	return p
}

func (p *Promise) Then(callback func(i any)) *Promise {
	p.callback = callback
	p.resolveChannel = make(chan any, 1)
	p.rejectChannel = make(chan error, 1)
	p.wg = &sync.WaitGroup{}

	var resolver Resolver = func(i any) {
		p.resolveChannel <- i
	}

	var rejector Rejector = func(i error) {
		p.rejectChannel <- i
	}

	p.wg.Add(1)
	go p.checker()
	go p.executor(resolver, rejector)

	return p
}

func (p *Promise) Catch(errorHandler func(i error)) *Promise {
	p.errorHandler = errorHandler

	return p
}

func (p *Promise) Wait() {
	if p.wg == nil {
		return
	}

	p.wg.Wait()
}

func (p *Promise) checker() {
	defer p.wg.Done()
	for {
		select {
		case <-p.ctx.Done():
			p.errorHandler(p.ctx.Err())
			close(p.rejectChannel)
			close(p.resolveChannel)
			return
		case o := <-p.resolveChannel:
			p.callback(o)
			close(p.rejectChannel)
			close(p.resolveChannel)
			return
		case err := <-p.rejectChannel:
			if p.errorHandler != nil {
				p.errorHandler(err)
			}
			close(p.rejectChannel)
			close(p.resolveChannel)
			return
		}
	}
}
