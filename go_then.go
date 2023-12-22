package go_then

import (
	"context"
	"sync"
)

type Promise struct {
	wg             *sync.WaitGroup
	ctx            context.Context
	resolveChannel chan any
	rejectChannel  chan error
	callback       func(i any)
	errorHandler   func(err error)
}

type Resolver func(i any)
type Rejector func(i error)

func New(ctx context.Context, executor func(resolve Resolver, reject Rejector)) *Promise {

	p := new(Promise)
	p.ctx = ctx
	p.resolveChannel = make(chan any, 1)
	p.rejectChannel = make(chan error, 1)
	p.wg = &sync.WaitGroup{}

	var resolver Resolver = func(i any) {
		p.resolveChannel <- i
		p.wg.Done()
		close(p.resolveChannel)
	}
	var rejector Rejector = func(i error) {
		p.rejectChannel <- i
		p.wg.Done()
		close(p.rejectChannel)
	}

	p.wg.Add(2)
	go p.checker()
	go executor(resolver, rejector)

	return p
}

func (p *Promise) Then(callback func(i any)) *Promise {
	p.callback = callback

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
			if p.errorHandler != nil {
				p.errorHandler(p.ctx.Err())
			}
			return
		case o := <-p.resolveChannel:
			if p.callback != nil {
				p.callback(o)
			}
			return
		case err := <-p.rejectChannel:
			if p.errorHandler != nil {
				p.errorHandler(err)
			}
			return
		}
	}
}
