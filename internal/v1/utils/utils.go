package utils

import (
	"context"
	"docshell/internal/v1/volume"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func UploadFile(ctx context.Context, file multipart.File, handler *multipart.FileHeader) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Create a temporary file within our temp-images directory that follows
			// a particular naming pattern
			tempFile, err := os.CreateTemp(volume.GetPath(), "*")
			if err != nil {
				return err
			}
			defer tempFile.Close()

			// read all of the contents of our uploaded file into a
			// byte array
			fileBytes, err := io.ReadAll(file)
			if err != nil {
				return err
			}
			// write this byte array to our temporary file
			if _, err := tempFile.Write(fileBytes); err != nil {
				return err
			}

			go func() { // TODO Fix bug with renaming
				// rename temporary file to origin
				tmp := tempFile.Name()
				origin := filepath.Join(volume.GetPath(), handler.Filename)
				if err := os.Rename(tmp, origin); err != nil {
					// return err
				}
				// return that we have successfully uploaded file
			}()
			return nil
		}
	}
}

func SendJSONResponse(w http.ResponseWriter, res any) {
	json.NewEncoder(w).Encode(res)
}
