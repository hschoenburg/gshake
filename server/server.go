package server

import(
  "net/http"
)



func IndexHandler(buildPath string) http.HandlerFunc {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "./index.html")
  })
}
