package main

import (
	"fmt"
	"time"

	//"github.com/go-chi/chi"

	//"github.com/gorilla/mux"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"log"
	"net/http"
)

/*type MuxHandler struct {
	GetHandler  http.Handler
	PostHandler http.Handler
}

func (h MuxHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetHandler.ServeHTTP(w, r)
	case http.MethodPost:
		h.PostHandler.ServeHTTP(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}*/

/*func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, "we've got panic here!")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("request: %s %s - %v\n",
			r.Method,
			r.URL.EscapedPath(),
			time.Since(start),
		)
	})
}*/

func RecoverMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if e := recover(); e != nil {
				ctx.SetStatusCode(http.StatusInternalServerError)
				fmt.Fprintln(ctx, "we've got panic here!")
			}
		}()
		next(ctx)
	})
}
func LoggingMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		start := time.Now()
		next(ctx)
		log.Printf("request: %s %s - %v\n",
			ctx.Method(),
			ctx.Path(),
			time.Since(start),
		)
	})
}

func main() {

	router := router.New()

	router.GET("/", func(ctx *fasthttp.RequestCtx) {
		ctx.WriteString("GET HANDLER\n")
	})

	router.GET("/{id}", func(ctx *fasthttp.RequestCtx) {
		id, ok := ctx.UserValue("id").(string)
		if !ok {
			ctx.WriteString("INVALID RESOURCE ID\n")
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}
		fmt.Fprintf(ctx, "GET BY ID HANDLER. RESOURCE ID IS %s\n", id)
	})

	log.Fatal(fasthttp.ListenAndServe(":8020", LoggingMiddleware(RecoverMiddleware(router.Handler))))

	/*router := chi.NewRouter()

	router.Use(LoggingMiddleware)
	router.Use(RecoverMiddleware)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		//panic("panic!")
		fmt.Fprintln(w, "GET HANDLER")
	})

	router.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		fmt.Fprintln(w, "GET BY ID HANDLER. RESOURCE ID IS", id)
	})

	router.Get("/{id}/name/{name}", func(w http.ResponseWriter, r *http.Request) {
		id, name := chi.URLParam(r, "id"), chi.URLParam(r, "name")
		fmt.Fprintf(w, "GET BY ID HANDLER WITH NAME. RESOURCE ID IS %s AND NAME IS %s\n",
			id, name)
	})

	router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "POST HANDLER")
	})

	log.Fatal(http.ListenAndServe(":8020", router))*/

	/*router := mux.NewRouter()

	router.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		fmt.Fprintln(w, "GET BY ID HANDLER. RESOURCE ID IS", vars["id"])
	}).Methods(http.MethodGet)

	router.HandleFunc("/{id}/name/{name}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		fmt.Fprintf(w, "GET BY ID HANDLER WITH NAME. RESOURCE ID IS %s AND NAME IS %s\n",
			vars["id"], vars["name"])
	}).Methods(http.MethodGet)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "POST HANDLER")
	}).Methods(http.MethodPost)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "GET HANDLER")
	}).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":8020", router))*/
}
