package imageToText

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	// "os"
	"strings"
	"testing"
)

func TestFailsIfRequestNotMultipartForm(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(Post))
	defer server.Close()

	resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer([]byte("")))
	if err != nil {
		t.Error(err)
	}

	// Fail test if status code is not 400
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected 400 but got %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	expected := "Content type not multipart/form-data"
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	// Fail test if error message is not correct
	trimmedResp := strings.Trim(string(body), "\n")
	if trimmedResp != expected {
		t.Errorf("Expected \"%v\" but got \"%v\"", expected, trimmedResp)
	}
}

func TestFailsIfNoFileProvided(t *testing.T) {
	// testFile, err := os.Open("file.txt")
	// if err != nil {
	// 	t.Error(err)
	// }
	// defer testFile.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// fileWriter, err := writer.CreateFormFile("file", "file.txt")
	// if err != nil {
	// 	t.Error(err)
	// }
	//
	// _, err = io.Copy(fileWriter, file)
	// if err != nil {
	// t.Error(err)
	// }

	writer.Close()

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "", &requestBody)
	if err != nil {
		t.Error(err)
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())

	Post(rr, req)

	resp := rr.Result()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected 400 but got %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	expected := "No file provided"
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	// Fail test if error message is not correct
	trimmedResp := strings.Trim(string(body), "\n")
	if trimmedResp != expected {
		t.Errorf("Expected \"%s\" but got \"%s\"", expected, trimmedResp)
	}
}

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
