package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTestRecursionHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/test_recursion", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(TestRecursionHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// print the return
	// fmt.Print(rr.Body.String())

	expected := `{"alive": true}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
