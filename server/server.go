package main

import (
	// "deal/middlewares"
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/amenabe22/chachata_backend/graph"
	"github.com/amenabe22/chachata_backend/graph/chans"
	"github.com/amenabe22/chachata_backend/graph/generated"
	"github.com/amenabe22/chachata_backend/graph/setup"
	"github.com/amenabe22/chachata_backend/middlewares"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	// _ "github.com/jinzhu/gorm/dialects/postgres" // using postgres sql
)

const defaultPort = "8081"

// TODO: move to conf package
type ckey string

func New() generated.Config {
	return generated.Config{
		Resolvers: &graph.Resolver{
			// Rooms:         map[string]*Chatroom{},
			AdminChans: map[string]*chans.CoreAdminChannel{},
		},
		Directives: generated.DirectiveRoot{
			User: func(ctx context.Context, obj interface{}, next graphql.Resolver, username string) (res interface{}, err error) {
				return next(context.WithValue(ctx, ckey("username"), username))
			},
		},
	}
}

func main() {
	mux := chi.NewRouter()
	c := cors.New(cors.Options{
		// allow everyone in
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	srv := handler.New(generated.NewExecutableSchema(New()))
	srv.AddTransport(transport.POST{})
	// added support for file upload transport
	srv.AddTransport(transport.MultipartForm{})
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})

	srv.Use(extension.Introspection{})
	mux.Use(middlewares.Middleware())
	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	mux.Handle("/query", c.Handler(srv))
	// fs := http.FileServer(http.Dir("./static"))
	// TODO: move hardcoded staic path to path finder variable
	fileServer := http.FileServer(http.Dir("/home/anonny/projects/fun/chachata/chachata_backend/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", neuter(fileServer)))
	// http.Handle("/", fileServer)
	setup.SetupModels()
	log.Printf("connect to http://0.0.0.0:%s/ for GraphQL playground", defaultPort)
	log.Fatal(http.ListenAndServe(":"+defaultPort, mux))
}

func neuter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}