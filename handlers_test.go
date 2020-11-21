package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"server/handlers"
)

func TestHandlers(t *testing.T){
	testEquityEndpoint := "/equity/TSLA"
	testIndexEndpoint := "/index/FTSE"
	req1, err := http.NewRequest("GET", testEquityEndpoint, nil)
	req2, err := http.NewRequest("GET", testIndexEndpoint, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr1 := httptest.NewRecorder()
	rr2 := httptest.NewRecorder()
	EquityHandler := http.HandlerFunc(handlers.GetEquity)
	IndexHandler := http.HandlerFunc(handlers.GetIndex)
	EquityHandler.ServeHTTP(rr1, req1)
	IndexHandler.ServeHTTP(rr2, req2)

	if status1 := rr1.Code; status1 != http.StatusOK {
		t.Errorf("Error showing wrong status code of: %v", status1)
	}

	if status2 := rr2.Code; status2 != http.StatusOK {
		t.Errorf("Error showing wrong status code of: %v", status2)
	}

	if contentType1 := rr1.Header().Get("Content-Type"); contentType1 != "application/json" {
        t.Errorf("content type header does not match %v",contentType1)
    }
	if contentType2 := rr2.Header().Get("Content-Type"); contentType2 != "application/json" {
        t.Errorf("content type header does not match %v",contentType2)
    }
}