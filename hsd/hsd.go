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
const HSD_URL="http://x:gobabygo@localhost:14037"


type RpcClient struct {
  NodeUrl string
}

func NewRpcClient() RpcClient {
  client := RpcClient {}  
  client.NodeUrl = HSD_URL
  return client
}


func (r RpcClient) NameInfo(name string) (NameInfoData, error) {

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


	data := NameInfoData{}
  data.Name = name
	err = json.Unmarshal(body, &data)

  fmt.Printf("HSD: %+v\n", data)
  if(err != nil) { return empty, err }

  if(data.Error.Message != "") {
    return empty, errors.New(data.Error.Message)
  }
  return data, nil
}
