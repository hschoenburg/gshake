package hsd

import (
  "bytes"
  "encoding/json"
  "log"
  "net/http"
  "io/ioutil"
)

const HSD_URL="http://127.0.0.1:13037"
  
type NameInfoData struct {
	Result struct {
				Info struct {} 		`json:"info"`
				Start startData `json:"start"`
	} `json:"result"`
}

type startData struct {
		Reserved bool `json:"reserved"`
		Week     int  `json:"week"`
		Start    int  `json:"start"`
}
	


func NameInfo(name string) startData {

  req := map[string]interface{}{
    "method": "getnameinfo",
    "params": []string{name},
  }

  args, err := json.Marshal(req)

  if err != nil {
    log.Fatal("JSON marshall failure: %s", err)
  }

  resp, err := http.Post(HSD_URL, "application/json", bytes.NewBuffer(args))

  defer resp.Body.Close()

  if err != nil {
    log.Fatal("RPC client connection failure: %s", err)
  }

  body, err := ioutil.ReadAll(resp.Body)

  if err != nil {
  }

	data := NameInfoData{}
	err = json.Unmarshal(body, &data)

  if err != nil {
    log.Fatal("RPC client connection failure: %s", err)
  }

  return data.Result.Start
}
