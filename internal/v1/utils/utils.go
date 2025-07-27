package utils

import (
	"context"
	"crypto/sha512"
	"docshell/internal/v1/docs/models"
	"docshell/internal/v1/volume"
	"encoding/hex"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func UploadFile(ctx context.Context, path string, file io.Reader, handler *multipart.FileHeader) error {
	// Create path if not exists
	volDir := filepath.Join(volume.GetPath(), path)
	if err := CreateDir(volDir); err != nil {
		return err
	}

	// Create a temporary file in the volume's directory
	tempFile, err := os.CreateTemp(volDir, "upload_*")
	if err != nil {
		return err
	}
	// Ensure the temporary file is removed if an error occurs
	defer os.Remove(tempFile.Name())
	// Close the file when done
	defer tempFile.Close()

	// Read the uploaded file into a byte array
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	// Write the byte array to the temporary file
	if _, err := tempFile.Write(fileBytes); err != nil {
		return err
	}

	// Ensure all data is flushed to disk
	if err := tempFile.Sync(); err != nil {
		return err
	}

	// Close the file explicitly before renaming (required for Windows compatibility)
	if err := tempFile.Close(); err != nil {
		return err
	}

	// Rename the temporary file to the original filename
	newPath := filepath.Join(volDir, handler.Filename)
	if err := os.Rename(tempFile.Name(), newPath); err != nil {
		return err
	}

	// Check if the context was canceled
	if err := ctx.Err(); err != nil {
		// Optionally remove the renamed file if the context was canceled
		os.Remove(newPath)
		return err
	}

	return nil
}

func SendJSONResponse(w http.ResponseWriter, res any) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func SendJSONErrorResponse(w http.ResponseWriter, code int, msg string) {
	// Send http error
	http.Error(w, msg, code)
	// Send json error response

	SendJSONResponse(w,
		models.ErrorResponse{
			StatusCode: code,
			StatusText: http.StatusText(code),
			Message:    msg,
		},
	)
}

func GenerateHash(b []byte) (hash string, err error) {
	sha := sha512.New()
	_, err = sha.Write(b)
	if err != nil {
		return "", err
	}
	sum := sha.Sum(nil)
	return hex.EncodeToString(sum), nil
}

func CreateDir(path string) error {
	return os.MkdirAll(path, 0666)
}
