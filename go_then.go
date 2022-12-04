package go_then

import (
	"errors"
	"sync"
	"time"
)

type Promise struct {
	resolveChannel chan any
	rejectChannel  chan error
	wg             sync.WaitGroup
}

func New() *Promise {

	p := new(Promise)
	p.resolveChannel = make(chan any, 1)
	p.rejectChannel = make(chan error, 1)

	return p
}

func (p *Promise) Resolve(i any) {
	p.resolveChannel <- i
}

func (p *Promise) Reject(i error) {
	p.rejectChannel <- i
}

func (p *Promise) Execute(i func()) *Promise {
	p.wg.Add(1)
	go i()

	return p
}

func (p *Promise) Wait() {
	p.wg.Wait()
}

func (p *Promise) Then(callback func(i any)) *Promise {
	go p.thenExecutor(callback)
	return p
}

func (p *Promise) thenExecutor(callback func(i any)) {
	for i := 0; i < 60; i++ {
		select {
		case o := <-p.resolveChannel:
			callback(o)
			p.wg.Done()
			return
		default:
		}
		time.Sleep(time.Second * 1)
	}

	p.rejectChannel <- errors.New("timeout")
}

func (p *Promise) Catch(errorHandler func(err error)) {
	go p.catchExecutor(errorHandler)
}

func (p *Promise) catchExecutor(errorHandler func(err error)) {
	for i := 0; i < 60; i++ {
		select {
		case o := <-p.rejectChannel:
			errorHandler(o)
			p.wg.Done()
			return
		default:
		}
		time.Sleep(time.Second * 1)
	}
}
