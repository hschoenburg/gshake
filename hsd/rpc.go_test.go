package hsd

import (
  "testing"
  "reflect"
)

func TestNameInfo(t *testing.T) {

  name := "hans"

  req, err := NameInfo(name)
  
  if(err != nil) {
    t.Fatal("Error requestion NameInfoData: %v", err)
  }

  r := reflect.TypeOf(req)
  if(r != NameInfoData) {
    t.Fatal("Wrong Response: %T", r)
  }

}


