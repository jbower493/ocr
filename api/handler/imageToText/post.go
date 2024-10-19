package imageToText

import (
	"net/http"
	"ocr/internal/httpHelpers"
)

func Post(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusNotFound)
		return
	}

	// Bitwise operator, shifts the bits of a number to the left, effectively multiplying it by 2. So this is like doing 10 * 2^20, which is aparently equal to 10MB (arg unit is bytes)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file_to_optimize")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	mimeType := header.Header.Get("Content-Type")
	// param3 := r.FormValue("param3") // get extra form fields out of the request

	// Respond
	type ResponseData struct {
		Text    string `json:"text"`
		Success bool   `json:"success"`
	}
	response := ResponseData{
		Text:    mimeType,
		Success: true,
	}
	httpHelpers.HandleSuccessResponse(w, response)
}
