package util

import (
	"github.com/gomodule/redigo/redis"
  "fmt"
  . "gshake/types"
)


func ContactScan(contact string, db redis.Conn) (keys []string, err error) {

 fmt.Printf("Notifs for %v \n", contact)

 iter := 0
 keys = []string{}
 matchstring := fmt.Sprintf("*%v*", contact)
 for {
   arr, arr_err := redis.Values(db.Do("SCAN", iter, "MATCH", matchstring))
   if arr_err != nil {
     err = arr_err
     return 
   } else {
     iter, _ = redis.Int(arr[0], nil)
     k, _ := redis.Strings(arr[1], nil)

     keys = append(keys, k...)
   }

   if(iter == 0) { break }
 }

 return keys, nil
}


func GetHashes(keys []string, db redis.Conn) ([]NameNotif, error) {
  notifs := make([]NameNotif, len(keys))

  for i, hash_key := range keys {

    fmt.Printf("GETTING HASH FOR %v\n", hash_key)

    values, err := redis.Values(db.Do("HGETALL", hash_key))

    if err != nil {
      fmt.Errorf("hashing error %v", err)
      return  notifs, err
    }

    n := new(NameNotif)
    err = redis.ScanStruct(values, n)

    if err != nil { 
      fmt.Errorf("ScanStruct Err: %v", err)
      return notifs, err
    }
    notifs[i] = *n

  }
  return notifs, nil
}


