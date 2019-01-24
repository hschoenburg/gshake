package hsd

import (
  "bytes"
  "errors"
  "encoding/json"
  "fmt"
  "net/http"
  "io/ioutil"
  . "gshake/types"
)

//const HSD_URL="x:gobabygo@192.241.237.167:13037"
const HSD_URL="http://x:gobabygo@192.241.237.167:13037"


func NameInfo(name string) (NameInfoData, error) {

  req := map[string]interface{}{
    "method": "getnameinfo",
    "params": []string{name},
  }

  empty := NameInfoData{}

  args, err := json.Marshal(req)
  if(err != nil) { return empty, err }


  resp, err := http.Post(HSD_URL, "application/json", bytes.NewBuffer(args))
  if(err != nil) { return empty, err }

  defer resp.Body.Close()

  body, err := ioutil.ReadAll(resp.Body)
  if(err != nil) { return empty, err }

  fmt.Printf("HSD: %v\n", string(body))

	data := NameInfoData{}
	err = json.Unmarshal(body, &data)
  if(err != nil) { return empty, err }

  if(data.Error.Message != "") {
    return empty, errors.New(data.Error.Message)
  }


  return data, nil
}
