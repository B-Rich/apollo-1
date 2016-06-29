// Apollo provides `net/context`-aware middleware chaining
package apollo

import (
	"fmt"
	"net/http"
	"strconv"

	"golang.org/x/net/context"
)

func handlerZero(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("h0\n"))
}

func handlerOne(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("h1\n"))
	return nil
}

func handlerContext(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if value, ok := FromContext(ctx); ok {
		contents := strconv.Itoa(value) + "\n"
		w.Write([]byte(contents))
		return nil
	}
	return fmt.Errorf("value not in context\n")
}

func handlerError(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return fmt.Errorf("err1\n")
}

func middleZero(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("m0\n"))
		h.ServeHTTP(w, r)
	})
}

func middleOne(h Handler) Handler {
	return HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		w.Write([]byte("m1\n"))
		return h.ServeHTTP(ctx, w, r)
	})
}

func middleTwo(h Handler) Handler {
	return HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		w.Write([]byte("m2\n"))
		return h.ServeHTTP(ctx, w, r)
	})
}

func middleAddError(h Handler) Handler {
	return HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		err := h.ServeHTTP(ctx, w, r)
		w.Write([]byte("error:" + err.Error()))

		return fmt.Errorf("found an error\n")
	})
}

func middleHandleError(h Handler) Handler {
	return HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		err := h.ServeHTTP(ctx, w, r)
		if err != nil {
			w.Write([]byte("error:" + err.Error()))
		}

		return nil
	})
}

// TestContext
type key int

const testKey key = 0

func NewTestContext(ctx context.Context, dummy int) context.Context {
	return context.WithValue(ctx, testKey, dummy)
}

func FromContext(ctx context.Context) (int, bool) {
	dummy, ok := ctx.Value(testKey).(int)
	return dummy, ok
}
