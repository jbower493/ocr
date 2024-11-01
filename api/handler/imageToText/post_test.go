package imageToText

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFailsIfRequestNotMultipartForm(t *testing.T) {
	req, err := http.NewRequest("POST", "/image-to-text", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Post)

	handler.ServeHTTP(rr, req)

	status := rr.Code
	if status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expected := "Content type not multipart/form-data"
	resp, _ := io.ReadAll(rr.Body)
	trimmedResp := strings.Trim(string(resp), "\n")
	if trimmedResp != expected {
		t.Errorf("Expected \"%v\" but got \"%v\"", expected, trimmedResp)
	}
}

/*func TestFailsIfNoFileProvided(t *testing.T) {
	req, err := http.NewRequest("POST", "/image-to-text", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Post)

	handler.ServeHTTP(rr, req)

	status := rr.Code
	if status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}*/

/*func TestSucceedsIfFileProvided(t *testing.T) {
	req, err := http.NewRequest("POST", "/image-to-text", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Post)

	handler.ServeHTTP(rr, req)

	status := rr.Code
	if status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"message":"hello"}`
	respBody, _ := ioutil.ReadAll(rr.Body)
	if string(respBody) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", string(respBody), expected)
	}
}*/
