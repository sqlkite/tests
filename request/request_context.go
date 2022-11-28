package request

// Same as request, but the handler is:
//   handler(*fasthttp.RequestCtx, T)
//
// Instead of
//   handler(*fasthttp.RequestCtx)
//
// In other words, the handler is receiving additional
// information (like an app/request-specific "context", or
// an "env" as we call it in our systems)
//
// Surprises me that Func() and Func[T] as well as
// Struct and Struct[T] conflict with each other and we
// need to have unique names for each.

import (
	"testing"

	"github.com/valyala/fasthttp"
	"src.sqlkite.com/utils/http"
)

func ReqT[T any](t *testing.T, env T) RequestBuilderT[T] {
	return RequestBuilderT[T]{env, Req(t)}
}

type RequestBuilderT[T any] struct {
	env T
	rb  RequestBuilder
}

func (r RequestBuilderT[T]) Path(path string) RequestBuilderT[T] {
	r.rb = r.rb.Path(path)
	return r
}

func (r RequestBuilderT[T]) Method(method string) RequestBuilderT[T] {
	r.rb = r.rb.Method(method)
	return r
}

func (r RequestBuilderT[T]) Header(key string, value string) RequestBuilderT[T] {
	r.rb = r.rb.Header(key, value)
	return r
}

func (r RequestBuilderT[T]) ProjectId(id string) RequestBuilderT[T] {
	r.rb = r.rb.ProjectId(id)
	return r
}

func (r RequestBuilderT[T]) Query(query map[string]string) RequestBuilderT[T] {
	r.rb = r.rb.Query(query)
	return r
}

func (r RequestBuilderT[T]) Body(body any) RequestBuilderT[T] {
	r.rb = r.rb.Body(body)
	return r
}

func (r RequestBuilderT[T]) UserValue(key string, value any) RequestBuilderT[T] {
	r.rb = r.rb.UserValue(key, value)
	return r
}

func (r RequestBuilderT[T]) Host(host string) RequestBuilderT[T] {
	r.rb = r.rb.Host(host)
	return r
}

func (r RequestBuilderT[T]) Get(handler func(*fasthttp.RequestCtx, T) (http.Response, error)) response {
	return r.Method("GET").Request(handler)
}

func (r RequestBuilderT[T]) Post(handler func(*fasthttp.RequestCtx, T) (http.Response, error)) response {
	return r.Method("POST").Request(handler)
}

func (r RequestBuilderT[T]) Put(handler func(*fasthttp.RequestCtx, T) (http.Response, error)) response {
	return r.Method("PUT").Request(handler)
}

func (r RequestBuilderT[T]) Delete(handler func(*fasthttp.RequestCtx, T) (http.Response, error)) response {
	return r.Method("DELETE").Request(handler)
}

func (r RequestBuilderT[T]) Request(handler func(*fasthttp.RequestCtx, T) (http.Response, error)) response {
	conn := r.rb.Conn()
	res, err := handler(conn, r.env)
	if res != nil {
		res.Write(conn)
	} else {
		http.ServerError().Write(conn)
	}

	// r2? really? :dealwithit:
	r2 := Res(r.rb.t, conn)
	r2.Err = err
	return r2
}
