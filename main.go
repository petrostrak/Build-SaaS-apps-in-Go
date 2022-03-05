package main

import "net/http"

type Route struct {
	Logger  bool
	Tester  bool
	Handler http.Handler
}

type App struct {
	User *Route
}

type User struct{}

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
