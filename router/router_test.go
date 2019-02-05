package router

import (
   "fmt"
   "gshake/mocks"
   "github.com/golang/mock/gomock" 
   "net/http"
   "net/http/httptest"
   "testing"
)


func TestNotifyHandler(t *testing.T) {

  contact := "hans@where.com"
  name := "hans"

  mockCtrl := gomock.NewController(t)
  defer mockCtrl.Finish()
  mockConn := mocks.NewMockConn(mockCtrl)

  //mockConn.EXPECT().Do("HMSET", "unsubs:", contact).Return(nil, nil).Times(1)

  path := fmt.Sprintf("/unsubscribe/%s", contact)
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := BuildRouter(mockConn)

  handler.ServeHTTP(rr, req)

  if status := rr.Code; status != http.StatusOK {
    t.Errorf("Notify handler returned wrong status %v", status)
  }
}


func TestUnsubscribeHandler(t *testing.T) {

  contact := "hans@where.com"

  mockCtrl := gomock.NewController(t)
  defer mockCtrl.Finish()
  mockConn := mocks.NewMockConn(mockCtrl)

  mockConn.EXPECT().Do("SADD", "unsubs:", contact).Return(nil, nil).Times(1)

  path := fmt.Sprintf("/unsubscribe/%s", contact)
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := BuildRouter(mockConn)

  handler.ServeHTTP(rr, req)

  if status := rr.Code; status != http.StatusOK {
    t.Errorf("Notify handler returned wrong status %v", status)
  }
}


