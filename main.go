package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

type Route struct {
	Logger  bool
	Tester  bool
	Handler http.Handler
}

type App struct {
	User *Route
}

type User struct{}

type key int

const (
	ctxTestKey key = 1
	ctxUserID      = 2
)

func main() {
	app := &App{
		User: &Route{
			Logger: true,
			Tester: true,
		},
	}

	http.ListenAndServe(":8080", app)
}

func shiftPath(p string) (head, tail string)

func (h *App) log(next http.Handler) http.Handler

func (h *App) test(next http.Handler) http.Handler

func (h *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var next *Route
	var head string

	head, r.URL.Path = shiftPath(r.URL.Path)
	if len(head) == 0 {
		next = &Route{
			Logger: true,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("home page"))
			}),
		}
	} else if head == "user" {
		var i interface{} = User{}
		next = &Route{
			Logger:  true,
			Tester:  true,
			Handler: i.(http.Handler),
		}
	} else {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	if next.Logger {
		next.Handler = h.log(next.Handler)
	}

	if next.Tester {
		next.Handler = h.test(next.Handler)
	}

	next.Handler.ServeHTTP(w, r)
}

func (u User) Detail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	v := ctx.Value(ctxTestKey)
	id := ctx.Value(ctxUserID)
	w.Write([]byte(fmt.Sprintf("value of context is %s for user id %d.", v, id)))
}

func (u User) Profile(w http.ResponseWriter, r *http.Request)

func (u User) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)
	if head == "profile" {
		u.Profile(w, r)
		return
	} else if head == "detail" {
		head, _ := shiftPath(r.URL.Path)
		i, err := strconv.Atoi(head)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, ctxUserID, i)
		u.Detail(w, r.WithContext(ctx))
		return
	}

	http.Error(w, "not found", http.StatusNotFound)
}
