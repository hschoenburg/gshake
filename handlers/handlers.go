package handlers

import (
  //"log"
  "net/http"
	"errors"
  "github.com/gorilla/mux"
  "gshake/hsd"
  "gshake/util"
  . "gshake/types"
  "fmt"
  "encoding/json"
  "github.com/gomodule/redigo/redis"
  //"reflect"
  "strconv"
)

func Index() http.HandlerFunc {

  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

    http.ServeFile(w, r, "./index.html")
  })
}

func NameInfo () http.HandlerFunc {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

    vars := mux.Vars(r)
    data, hsd_err := hsd.NameInfo(vars["name"])

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)

    json_err := json.NewEncoder(w).Encode(data); 
    
    // idk do all the errs at once?
    if json_err != nil {
      http.Error(w, json_err.Error(), 500)
      return
    }
    if hsd_err != nil {
      http.Error(w, hsd_err.Error(), 500)
    return
    }
  })
}

func WeekHandler(db redis.Conn) http.HandlerFunc {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    week, week_err := strconv.Atoi(vars["week"])
    if(week_err != nil) {
      http.Error(w, week_err.Error(), 500)
    }

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
    vars := mux.Vars(r)

    keys, scan_err := util.ContactScan(vars["contact"], db)
    fmt.Printf("KEYS: %v\n", keys)

    if scan_err != nil {
      http.Error(w, scan_err.Error(), 500)
    }

    notifs, _ := util.GetHashes(keys, db)
    fmt.Printf("notifs %v", notifs)

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    if json_err := json.NewEncoder(w).Encode(notifs); json_err != nil {
      http.Error(w, json_err.Error(), 500)
    }
  })
}


func NotifyHandler(db redis.Conn) http.HandlerFunc {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		Name := r.FormValue("name")
		data, hsd_err :=  hsd.NameInfo(Name)
    info := data.Result.Start

    fmt.Printf("error: %v reserved: %v, week: %v, start: %v\n", data.Error.Message, info.Reserved, info.Week, info.Start)

    if(hsd_err != nil) {
      http.Error(w, hsd_err.Error(), 500)
      return
    }

		if(info.Reserved) {
			panic(errors.New("NAME RESERVED"))
		}

    notif := NameNotif{
      Name: Name,
      Contact: r.FormValue("contact"),
      Notified: false,
      Verified: false,
      Week: info.Week,
    }

    // 1 hash for every notif
    hash_key := notif.Name + ":" + notif.Contact // matches with name_key value
    // 2 sets for every week
    notif_key := fmt.Sprintf("notifs:%v", info.Week) // lookup all notifs for given week
    name_key := fmt.Sprintf("names:%v", info.Week) // lookup all names for given week

    // save hash
    _, hash_err := db.Do("HMSET", redis.Args{hash_key}.AddFlat(notif)...)
    if hash_err != nil {
      http.Error(w, hash_err.Error(), 500)
      return
    }
    // save indexes
    _, notif_err := db.Do("SADD", notif_key, hash_key)
    if notif_err != nil {
      http.Error(w, notif_err.Error(), 500)
      return
    }
    _, name_err := db.Do("SADD", name_key, notif.Name)
    if name_err != nil {
      http.Error(w, name_err.Error(), 500)
      return
    }

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    if json_err := json.NewEncoder(w).Encode(notif); json_err != nil {
      http.Error(w, json_err.Error(), 500)
      return
    }
  })
}




