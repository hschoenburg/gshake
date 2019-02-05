package router

import (
	"github.com/gorilla/mux"
  "github.com/gomodule/redigo/redis"
	"gshake/api"
	"gshake/server"
	"net/http"
)

func BuildRouter(conn redis.Conn) *mux.Router {

	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(server.NotFoundHandler)

	r.HandleFunc("/", server.IndexHandler())

	r.HandleFunc("/info/{name}", api.NameInfo()).Methods("GET")

	r.HandleFunc("/notify", api.NotifyHandler(conn)).Methods("POST")

	r.HandleFunc("/notifs/{contact}", api.NotifsHandler(conn)).Methods("GET")

	r.HandleFunc("/names/{week}", api.WeekHandler(conn)).Methods("GET")

	r.HandleFunc("/verify/{contact}", api.Verify(conn)).Methods("GET")

	r.HandleFunc("/unsubscribe/{contact}", api.Unsubscribe(conn)).Methods("GET")

	const STATIC_DIR = "/ui/build/"

	r.PathPrefix(STATIC_DIR).Handler(http.StripPrefix(STATIC_DIR, http.FileServer(http.Dir("."+STATIC_DIR))))

	return r
}
