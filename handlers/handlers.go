package handlers

import (
  //"log"
  "net/http"
	"errors"
  "github.com/gorilla/mux"
  "gshake/hsd"
  "fmt"
  "encoding/json"
  "github.com/gomodule/redigo/redis"
  //"reflect"
  "strconv"
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

func WeekHandler(db redis.Conn) http.HandlerFunc {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    week, week_err := strconv.Atoi(vars["week"])
    if(week_err != nil) {
      http.Error(w, week_err.Error(), 500)
    }

    // SMEMBERS name:list:week "all the names available this week" 


    // 2 sets for every week
    //notif_key := fmt.Sprintf("notifs:%v", 2) // lookup all notifs for given week
    name_key := fmt.Sprintf("names:%v", week) // lookup all names for given week

   names, names_err := redis.Strings(db.Do("SMEMBERS", name_key))
    
    if names_err != nil {
      http.Error(w, names_err.Error(), 500)
    }

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    if json_err := json.NewEncoder(w).Encode(names); json_err != nil {
      http.Error(w, json_err.Error(), 500)
    }
  })
}


func NotifsHandler(db redis.Conn) http.HandlerFunc {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // get all notifs for a given email addr
    //vars := mux.Vars(r)

    /*
    email := mux.Vars(r)["email"]
    fmt.Println("Notifs for %v", email)

    iter := 0
    keys = []string{}
    for {
      arr, arr_err := redis.Values(db.Do("SCAN", iter, "MATCH", email))
      if arr_err != nil {
        http.Error(w, arr_err.Error(), 500)
      }

      iter, _ = redis.Int(arr[0], nil)

    //Query(email)
    // SCAN 0 MATCH "email" -- returns all notif hash keys for email
    // foreach
    // HGETALL key

    //Query(week)
    // SMEMBERS notifs:week "all the name:email pairs for that week" 
    // SMEMBERS names:week "all the names available this week" 
    // for each notif:list n
    // HGETALL n

    /*
    _, notifs_err := db.Do("HGETALL", redis.Args{hash_key}.AddFlat(notif)...)
    if notifs_err != nil {
      http.Error(w, notifs_err.Error(), 500)
    }
    */


  })
}

func NotifyHandler(db redis.Conn) http.HandlerFunc {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		Name := r.FormValue("name")
		data :=  hsd.NameInfo(Name)
		fmt.Printf("reserved: %v, week: %v, start: %v", data.Reserved, data.Week, data.Start)

		if(data.Reserved) {
			panic(errors.New("NAME RESERVED"))
		}

    notif := NameNotif{
      Name: Name,
      Contact: r.FormValue("contact"),
      Notified: false,
      Verified: false,
      Week: data.Start,
    }

    //Save(week, clevertld, email)
    // 
    // Weekly Indexes
    // SADD notif:list:45 "clevertld:email"
    // SADD name:list:45 "clevertld"
    //
    // Notif Hashes
    // HMSET clevertld:email "week 34, verified true, notified, false"


    //Utilities to build
    // weekly notifier
    // dictionary INDEXER
    // web GUI with "Names this week"
    // 1 hash for every notif
    hash_key := notif.Name + ":" + notif.Contact // matches with name_key value
    // 2 sets for every week
    notif_key := fmt.Sprintf("notifs:%v", data.Start) // lookup all notifs for given week
    name_key := fmt.Sprintf("names:%v", data.Start) // lookup all names for given week

    // save hash
    _, hash_err := db.Do("HMSET", redis.Args{hash_key}.AddFlat(notif)...)
    if hash_err != nil {
      http.Error(w, hash_err.Error(), 500)
    }
    // save indexes
    _, notif_err := db.Do("SADD", notif_key, hash_key)
    if notif_err != nil {
      http.Error(w, notif_err.Error(), 500)
    }
    _, name_err := db.Do("SADD", name_key, notif.Name)
    if name_err != nil {
      http.Error(w, name_err.Error(), 500)
    }


    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    if json_err := json.NewEncoder(w).Encode(notif); json_err != nil {
      http.Error(w, json_err.Error(), 500)
    }
  })
}




