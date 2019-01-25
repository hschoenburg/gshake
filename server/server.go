package server

import(
  "net/http"
  "fmt"
  //"path"
)



func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
    // You can use the serve file helper to respond to 404 with
    // your request file.
    fmt.Printf("NOT FOUND URI: %s\n", r.RequestURI)
    //http.ServeFile(w, r, "../ui/build/index.html")
}


func IndexHandler() http.HandlerFunc {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "./ui/build/index.html")
  })
}
