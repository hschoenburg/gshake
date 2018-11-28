package main

import (
  "gshake/handlers"
  //"html"
  "context"
  "flag"
  "fmt"
  "github.com/gorilla/mux"
  "net/http"
  "log"
  "os"
	"github.com/gomodule/redigo/redis"
  "os/signal"
  "time"
)

// ToDO
//Utilities to build
// weekly notifier
// dictionary INDEXER
// web GUI with "Names this week"
// Connect to remote HSD node


func main() {
	var wait time.Duration
  flag.DurationVar(&wait, "graceful-timeout", time.Second * 15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
  flag.Parse()

  // Create Pool of Redis Connections

  pool := newPool()
  conn := pool.Get()
  defer conn.Close()


	r := mux.NewRouter()
	r.HandleFunc("/", handlers.Hello)
  r.HandleFunc("/info/{name}", handlers.NameInfo).Methods("GET")

  r.HandleFunc("/notify", handlers.NotifyHandler(conn)).Methods("POST")

  r.HandleFunc("/notifs/{contact}", handlers.NotifsHandler(conn)).Methods("GET")

  r.HandleFunc("/names/{week}", handlers.WeekHandler(conn)).Methods("GET")



	srv := &http.Server{
			Addr:         "0.0.0.0:8080",
			// Good practice to set timeouts to avoid Slowloris attacks.
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
			Handler: r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
    go func() {
	      fmt.Println("Server Up. Listening on 8080")
        if err := srv.ListenAndServe(); err != nil {
            log.Println(err)
        }
    }()

    c := make(chan os.Signal, 1)
    // We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
    // SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
    signal.Notify(c, os.Interrupt)

    // Block until we receive our signal.
    <-c

    // Create a deadline to wait for.
    ctx, cancel := context.WithTimeout(context.Background(), wait)
    defer cancel()
    // Doesn't block if no connections, but will otherwise wait
    // until the timeout deadline.
    srv.Shutdown(ctx)
    // Optionally, you could run srv.Shutdown in a goroutine and block on
    // <-ctx.Done() if your application should wait for other services
    // to finalize based on context cancellation.
    log.Println("shutting down")
    os.Exit(0)
}


func newPool() *redis.Pool {
	return &redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle: 80,
		// max number of connections
		MaxActive: 12000,
		// Dial is an application supplied function for creating and
		// configuring a connection.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}



