package api

import (
	//"log"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"gshake/hsd"
	. "gshake/types"
	"gshake/util"
	"net/http"
	//"reflect"
	"strconv"
)

func Verify(db redis.Conn) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		contact := vars["contact"]

		fmt.Printf("Verified %T", contact)
		http.Redirect(w, r, "/", 301)
	})
}

func Unsubscribe(db redis.Conn) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
    fmt.Printf("VARS: %v", vars)

    unsub_key := "unsubs:"
		contact := vars["contact"]

		// save unsub
		_, unsub_err := db.Do("SADD", unsub_key, contact)
		if unsub_err != nil {
			http.Error(w, unsub_err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("successfully unsubscribed")
	})
}

func NameInfo() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		vars := mux.Vars(r)
		data, hsd_err := hsd.NameInfo(vars["name"])

		if hsd_err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(hsd_err.Error())
			return
		}

		if json_err := json.NewEncoder(w).Encode(data); json_err != nil {
			http.Error(w, json_err.Error(), 500)
		}

	})
}

func WeekHandler(db redis.Conn) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		week, week_err := strconv.Atoi(vars["week"])
		if week_err != nil {
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
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")

    //Name := "YOYOYOYOY"
    Name := r.FormValue("name")
    fmt.Printf("$$$$$$$$$$$$$$$$$ %v \n", Name)

		data, hsd_err := hsd.NameInfo(Name)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if hsd_err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(hsd_err.Error())
			return
		}

		info := data.Result.Start

		fmt.Printf("error: %v reserved: %v, week: %v, start: %v\n", data.Error.Message, info.Reserved, info.Week, info.Start)

		if info.Reserved {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errors.New("Name Reserved"))
		}

		notif := NameNotif{
			Name:     Name,
			Contact:  r.FormValue("contact"),
			Notified: false,
			Week:     info.Week,
		}

		// 1 hash for every notif
		hash_key := notif.Name + ":" + notif.Contact // matches with name_key value
		// 2 sets for every week
		notif_key := fmt.Sprintf("notifs:%v", info.Week) // lookup all notifs for given week
		name_key := fmt.Sprintf("names:%v", info.Week)   // lookup all names for given week

		// save hash
		_, hash_err := db.Do("HMSET", redis.Args{hash_key}.AddFlat(notif)...)
		if hash_err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(hash_err)
			return
		}
		// save indexes
		_, notif_err := db.Do("SADD", notif_key, hash_key)
		if notif_err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(notif_err)
			return
		}
		_, name_err := db.Do("SADD", name_key, notif.Name)
		if name_err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(name_err)
			return
		}

		w.WriteHeader(http.StatusOK)
		if json_err := json.NewEncoder(w).Encode(notif); json_err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(json_err)
			return
		}
	})
}
