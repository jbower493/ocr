package imageToText

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
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

func TestFailsIfFileIsUnsupportedMimeType(t *testing.T) {
	testFile, err := os.Open("../../../assets/testImages/test_webp_img.webp")
	if err != nil {
		t.Error(err)
	}
	defer testFile.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Set up a custom header for the file part to specify "Content-Type: image/png"
	partHeader := make(textproto.MIMEHeader)
	partHeader.Set("Content-Disposition", `form-data; name="image_to_text"; filename="image_to_text.webp"`)
	partHeader.Set("Content-Type", "image/webp") // Set the content type here

	// Create a new part in the writer using the custom headers
	part, err := writer.CreatePart(partHeader)
	if err != nil {
		t.Error(err)
	}

	_, err = io.Copy(part, testFile)
	if err != nil {
		t.Error(err)
	}

	writer.Close()

	server := httptest.NewServer(http.HandlerFunc(Post))
	defer server.Close()

	resp, err := http.Post(server.URL, writer.FormDataContentType(), &requestBody)
	if err != nil {
		t.Error(err)
	}

	// Fail test if status code is not 400
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected 400 but got %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	expected := "webp extention not supported"
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
