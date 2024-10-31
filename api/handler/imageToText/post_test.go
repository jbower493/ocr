package imageToText

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPost(t *testing.T) {
	req, err := http.NewRequest("POST", "/image-to-text", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Post)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"message":"hello"}`
	respBody, _ := ioutil.ReadAll(rr.Body)
	if string(respBody) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", string(respBody), expected)
	}
}
