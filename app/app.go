package app

import (
	"fmt"
	"net/http"
)

type App struct {
	mux *http.ServeMux
}

func NewApp() *App {
	a := &App{
		mux: http.NewServeMux(),
	}

	a.bindHandlers()

	return a
}

func (a *App) Listen(host string) {
	http.ListenAndServe(host, a.mux)
}

func (a *App) bindHandlers() {
	a.mux.HandleFunc("/", a.rootHandler)
}

func (a *App) rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
}
