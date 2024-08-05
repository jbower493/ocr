package textRecognition

import (
	"github.com/otiai10/gosseract/v2"
)

// The goserract library is showing type errors because the types are not being recognized correctly in WSL for some reason, even though the types in the repo in windows are fine, and it executes fine in docker.
func Base64BytesToText(base64ByteSlice []byte) (string, error) {
	client := gosseract.NewClient()
	defer client.Close()

	setImageFromBytesErr := client.SetImageFromBytes(base64ByteSlice)
	if setImageFromBytesErr != nil {
		return "", setImageFromBytesErr
	}

	text, getTextErr := client.Text()
	if getTextErr != nil {
		return "", getTextErr
	}

	return text, nil
}