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
	"github.com/amenabe22/chachata_backend/graph/model"
	"github.com/amenabe22/chachata_backend/graph/setup"
	"github.com/go-chi/chi"
	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"

	// "github.com/tinrab/retry"
	"gorm.io/gorm"

	"github.com/go-chi/jwtauth/v5"
	// _ "github.com/jinzhu/gorm/dialects/postgres" // using postgres sql
)

const defaultPort = "8080"

// TODO: move to conf package
type ckey string

var db *gorm.DB

func New() generated.Config {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	// retry.ForeverSleep(2*time.Second, func(_ int) error {
	// 	_, err := client.Ping().Result()
	// 	return err
	// })
	return generated.Config{
		Resolvers: &graph.Resolver{
			RedisClient: client,
			Rooms:       map[string]*model.Chatroom{},
			AdminChans:  map[string]*chans.CoreAdminChannel{},
			Coredb:      setup.SetupModels(),
			// db,
		},
		Directives: generated.DirectiveRoot{
			User: func(ctx context.Context, obj interface{}, next graphql.Resolver, username string) (res interface{}, err error) {
				return next(context.WithValue(ctx, ckey("username"), username))
			},
		},
	}
}

var tokenAuth *jwtauth.JWTAuth

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
	// authWare, _ := middlewares.Middleware()
	srv.Use(extension.Introspection{})
	// mux.Use(authWare.MiddlewareFunc())
	// mux.Use(authWare.MiddlewareFunc())
	mux.Use(model.JwtMiddleware())

	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))

	// handler.GraphQL(go_orders_graphql_api.NewExecutableSchema(go_orders_graphql_api.Config{Resolvers: &go_orders_graphql_api.Resolver{
	//     DB: db,
	// }})))
	mux.Handle("/query", c.Handler(srv))
	// fs := http.FileServer(http.Dir("./static"))
	// TODO: move hardcoded staic path to path finder variable
	fileServer := http.FileServer(http.Dir("/home/anonny/projects/fun/chachata/chachata_backend/static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", neuter(fileServer)))
	// http.Handle("/", fileServer)

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
