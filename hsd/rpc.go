package hsd

import (
  "bytes"
  "errors"
  "encoding/json"
  "fmt"
  "net/http"
  "io/ioutil"
)

const HSD_URL="http://127.0.0.1:13037"
  
type NameInfoData struct {
	Result struct {
	  Info struct {} 		`json:"info"`
		Start StartData   `json:"start"`
  }
  Error struct {
    Message string  `json:"message"`
    Code int        `json:"code"`
  }                 `json:"error"`
}

type StartData struct {
		Reserved bool `json:"reserved"`
		Week     int  `json:"week"`
		Start    int  `json:"start"`
}
	


func NameInfo(name string) (NameInfoData, error) {

  req := map[string]interface{}{
    "method": "getnameinfo",
    "params": []string{name},
  }

  args, err := json.Marshal(req)

  resp, err := http.Post(HSD_URL, "application/json", bytes.NewBuffer(args))

  defer resp.Body.Close()

  body, err := ioutil.ReadAll(resp.Body)

  fmt.Printf("HSD: %v\n", string(body))

	data := NameInfoData{}
	err = json.Unmarshal(body, &data)

  if(data.Error.Message != "") {
    return NameInfoData{}, errors.New(data.Error.Message)
  }

  if err != nil {
    return NameInfoData{}, fmt.Errorf("NameInfo Error: %v", err)
  }

  return data, nil
}
