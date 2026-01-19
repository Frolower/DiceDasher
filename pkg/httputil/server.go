//мам можно мне https://github.com/gorilla/mux
//у нас есть https://github.com/gorilla/mux дома
//https://github.com/gorilla/mux дома:

package httputil

import "net/http"

type Route struct {
	handlers map[string]http.HandlerFunc
}

type Router struct {
	routes map[string]*Route
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]*Route),
	}
}

func (r *Router) Handle(path string) *Route {
	route := &Route{
		handlers: make(map[string]http.HandlerFunc),
	}
	r.routes[path] = route
	return route
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	route, ok := r.routes[req.URL.Path]
	if !ok {
		http.NotFound(w, req)
		return
	}

	handler, ok := route.handlers[req.Method]
	if !ok {
		w.Header().Set("Allow", allowedMethods(route))
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	handler(w, req)
}

func allowedMethods(route *Route) string {
	allow := ""
	for m := range route.handlers {
		if allow != "" {
			allow += ", "
		}
		allow += m
	}
	return allow
}

func (r *Route) GET(h http.HandlerFunc) *Route {
	r.handlers[http.MethodGet] = h
	return r
}

func (r *Route) POST(h http.HandlerFunc) *Route {
	r.handlers[http.MethodPost] = h
	return r
}
