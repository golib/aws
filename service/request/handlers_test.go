package request_test

import (
	"testing"

	"github.com/golib/assert"

	"github.com/golib/aws/service"
	"github.com/golib/aws/service/request"
)

func TestHandlerList(t *testing.T) {
	s := ""
	r := &request.Request{}
	l := request.HandlerList{}
	l.PushBack(func(r *request.Request) {
		s += "a"
		r.Data = s
	})
	l.Run(r)
	assert.Equal(t, "a", s)
	assert.Equal(t, "a", r.Data)
}

func TestMultipleHandlers(t *testing.T) {
	r := &request.Request{}
	l := request.HandlerList{}
	l.PushBack(func(r *request.Request) { r.Data = nil })
	l.PushFront(func(r *request.Request) { r.Data = service.Bool(true) })
	l.Run(r)
	if r.Data != nil {
		t.Error("Expected handler to execute")
	}
}

func TestNamedHandlers(t *testing.T) {
	l := request.HandlerList{}
	named := request.NamedHandler{Name: "Name", Fn: func(r *request.Request) {}}
	named2 := request.NamedHandler{Name: "NotName", Fn: func(r *request.Request) {}}
	l.PushBackNamed(named)
	l.PushBackNamed(named)
	l.PushBackNamed(named2)
	l.PushBack(func(r *request.Request) {})
	assert.Equal(t, 4, l.Len())
	l.Remove(named)
	assert.Equal(t, 2, l.Len())
}

func TestLoggedHandlers(t *testing.T) {
	expectedHandlers := []string{"name1", "name2"}
	l := request.HandlerList{}
	loggedHandlers := []string{}
	l.AfterEachFn = request.HandlerListLogItem
	cfg := service.Config{Logger: service.LoggerFunc(func(args ...interface{}) {
		loggedHandlers = append(loggedHandlers, args[2].(string))
	})}

	named1 := request.NamedHandler{Name: "name1", Fn: func(r *request.Request) {}}
	named2 := request.NamedHandler{Name: "name2", Fn: func(r *request.Request) {}}
	l.PushBackNamed(named1)
	l.PushBackNamed(named2)
	l.Run(&request.Request{Config: cfg})

	assert.Equal(t, expectedHandlers, loggedHandlers)
}

func TestStopHandlers(t *testing.T) {
	l := request.HandlerList{}
	stopAt := 1
	l.AfterEachFn = func(item request.HandlerListRunItem) bool {
		return item.Index != stopAt
	}

	called := 0
	l.PushBackNamed(request.NamedHandler{Name: "name1", Fn: func(r *request.Request) {
		called++
	}})
	l.PushBackNamed(request.NamedHandler{Name: "name2", Fn: func(r *request.Request) {
		called++
	}})
	l.PushBackNamed(request.NamedHandler{Name: "name3", Fn: func(r *request.Request) {
		assert.Fail(t, "third handler should not be called")
	}})
	l.Run(&request.Request{})

	assert.Equal(t, 2, called, "Expect only two handlers to be called")
}
