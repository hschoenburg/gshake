package handlers

import (
  //"log"
  "net/http"
  "github.com/gorilla/mux"
  "fmt"
  "encoding/json"
  "github.com/gomodule/redigo/redis"
)

type NameNotif struct {
  Name string   `json:"name"`
  Contact string  `json:"contact"`
  Week int      `json:"week"`
  Notified bool `json:"notified"`
  Verified bool `json:"verified"`
}

type NameNotifs []NameNotif


func Hello(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "Welcome to GShake")
}

func NameStatus (w http.ResponseWriter, r *http.Request) {
  // next step query hsd-rpc
  vars := mux.Vars(r)
  fmt.Fprintln(w, "GETTING STATUS ", vars["name"])
}

func Notifs(w http.ResponseWriter, r *http.Request) {
  // query redis reply with all notifs for email
  // for dev purposes onlu
  vars := mux.Vars(r)
  fmt.Println("Notifs for %v", vars["email"])
  notifs := NameNotifs{}
  if err := json.NewEncoder(w).Encode(notifs); err != nil {
    panic(err)
  }
}

func NotifsHandler(db redis.Conn) http.HandlerFunc {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

    // get all notifs for a given email addr
    //vars := mux.Vars(r)

  })
}

func NotifyHandler(db redis.Conn) http.HandlerFunc {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // next step save to redis
    // save week in key as well TODO

    notif := NameNotif{
      Name: r.FormValue("name"),
      Contact: r.FormValue("contact"),
      Notified: false,
      Verified: false,
      Week: 1,
    }

    //Save(week, clevertld, email)
    // 
    // SADD notif:list:45 "clevertld:email"
    // SADD name:list:45 "clevertld"
    // HMSET clevertld:email "week 34, verified true, notified, false"


    //Query(week)
    // HGETALL notif:list:week "all the name:email pairs for that week" 
    // HGETALL name:list:week "all the names available this week" 


    //Utilities to build
    // weekly notifier
    // dictionary INDEXER
    // web GUI with "Names this week"
    
    key := notif.Name + ":" + notif.Contact


    //hash := fmt.Sprintf("name %v contact %v week %v verified %v notified %v", notif.Name, notif.Contact, 22, false, false)

    _, err := db.Do("HMSET", redis.Args{key}.AddFlat(notif)...)
    if err != nil {
      http.Error(w, err.Error(), 500)
    }

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusOK)
    if err = json.NewEncoder(w).Encode(notif); err != nil {
      http.Error(w, err.Error(), 500)
  })
}




