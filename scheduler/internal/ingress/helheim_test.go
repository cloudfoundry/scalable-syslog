// This file was generated by github.com/nelsam/hel.  Do not
// edit this code by hand unless you *really* know what you're
// doing.  Expect any changes made manually to be overwritten
// the next time hel regenerates this file.

package ingress_test

import (
	"net/http"

	v1 "github.com/cloudfoundry-incubator/scalable-syslog/internal/api/v1"
	"github.com/cloudfoundry-incubator/scalable-syslog/scheduler/internal/ingress"
)

type mockCounter struct {
	CountCalled chan bool
	CountOutput struct {
		Ret0 chan int
	}
}

func newMockCounter() *mockCounter {
	m := &mockCounter{}
	m.CountCalled = make(chan bool, 100)
	m.CountOutput.Ret0 = make(chan int, 100)
	return m
}
func (m *mockCounter) Count() int {
	m.CountCalled <- true
	return <-m.CountOutput.Ret0
}

type mockBindingReader struct {
	FetchBindingsCalled chan bool
	FetchBindingsOutput struct {
		Bindings chan ingress.Bindings
		Err      chan error
	}
}

func newMockBindingReader() *mockBindingReader {
	m := &mockBindingReader{}
	m.FetchBindingsCalled = make(chan bool, 100)
	m.FetchBindingsOutput.Bindings = make(chan ingress.Bindings, 100)
	m.FetchBindingsOutput.Err = make(chan error, 100)
	return m
}
func (m *mockBindingReader) FetchBindings() (appBindings ingress.Bindings, err error) {
	m.FetchBindingsCalled <- true
	return <-m.FetchBindingsOutput.Bindings, <-m.FetchBindingsOutput.Err
}

type mockAdapterPool struct {
	ListCalled chan bool
	ListOutput struct {
		Bindings chan [][]*v1.Binding
		Err      chan error
	}
	CreateCalled chan bool
	CreateInput  struct {
		Binding chan *v1.Binding
	}
	CreateOutput struct {
		Err chan error
	}
	DeleteCalled chan bool
	DeleteInput  struct {
		Binding chan *v1.Binding
	}
	DeleteOutput struct {
		Err chan error
	}
}

func newMockAdapterPool() *mockAdapterPool {
	m := &mockAdapterPool{}
	m.ListCalled = make(chan bool, 100)
	m.ListOutput.Bindings = make(chan [][]*v1.Binding, 100)
	m.ListOutput.Err = make(chan error, 100)
	m.CreateCalled = make(chan bool, 100)
	m.CreateInput.Binding = make(chan *v1.Binding, 100)
	m.CreateOutput.Err = make(chan error, 100)
	m.DeleteCalled = make(chan bool, 100)
	m.DeleteInput.Binding = make(chan *v1.Binding, 100)
	m.DeleteOutput.Err = make(chan error, 100)
	return m
}
func (m *mockAdapterPool) List() (bindings [][]*v1.Binding, err error) {
	m.ListCalled <- true
	return <-m.ListOutput.Bindings, <-m.ListOutput.Err
}
func (m *mockAdapterPool) Create(binding *v1.Binding) (err error) {
	m.CreateCalled <- true
	m.CreateInput.Binding <- binding
	return <-m.CreateOutput.Err
}
func (m *mockAdapterPool) Delete(binding *v1.Binding) (err error) {
	m.DeleteCalled <- true
	m.DeleteInput.Binding <- binding
	return <-m.DeleteOutput.Err
}

type mockGetter struct {
	GetCalled chan bool
	GetInput  struct {
		NextID chan int
	}
	GetOutput struct {
		Resp chan *http.Response
		Err  chan error
	}
}

func newMockGetter() *mockGetter {
	m := &mockGetter{}
	m.GetCalled = make(chan bool, 100)
	m.GetInput.NextID = make(chan int, 100)
	m.GetOutput.Resp = make(chan *http.Response, 100)
	m.GetOutput.Err = make(chan error, 100)
	return m
}
func (m *mockGetter) Get(nextID int) (resp *http.Response, err error) {
	m.GetCalled <- true
	m.GetInput.NextID <- nextID
	return <-m.GetOutput.Resp, <-m.GetOutput.Err
}
