package service

import (
	"bytes"
	"context"
	"docshell/internal/v1/docs/models"
	"docshell/internal/v1/docs/repository"
	"docshell/internal/v1/utils"
	"docshell/internal/v1/volume"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"sync"

	"docshell/internal/v1/storage"
	"net/http"
	"time"
)

func GetAllDocuments(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	// Set timeout context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Get db connection
	con := storage.GetConnection()

	// Get all documents from repository
	docs, err := repository.GetAllDocuments(ctx, con)
	if err != nil {
		msg := "Database error: could not read docs"
		utils.SendJSONErrorResponse(w, http.StatusInternalServerError, msg)
		return
	}

	select {
	case <-ctx.Done(): // Context exceeded
		if ctx.Err() == context.DeadlineExceeded {
			msg := "Request timeout"
			utils.SendJSONErrorResponse(w, http.StatusRequestTimeout, msg)
		} else { // Other error
			msg := "Unknown error"
			utils.SendJSONErrorResponse(w, http.StatusInternalServerError, msg)
		}
	default: // Success
		// Build response
		res := models.ResponseMultipleDocuments{
			StatusCode: http.StatusOK,
			Documents:  make([]models.Document, len(docs)),
		}
		// Copy documents
		copy(res.Documents, docs)
		utils.SendJSONResponse(w, res)
	}
}

func GetDocumentById(ctx context.Context, w http.ResponseWriter, r *http.Request, id int) {
	// Set timeout context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Get db connection
	con := storage.GetConnection()

	// Get document
	doc, err := repository.GetDocumentById(ctx, con, id)
	if err != nil {
		msg := "Internal server error"
		utils.SendJSONErrorResponse(w, http.StatusInternalServerError, msg)
		return
	}
	// Send NoContent if empty
	if doc == (models.Document{}) {
		msg := "Requested document not found"
		utils.SendJSONErrorResponse(w, http.StatusNotFound, msg)
		return
	}

	select {
	case <-ctx.Done(): // Context exceed
		msg := "Internal context exceed"
		utils.SendJSONErrorResponse(w, http.StatusRequestTimeout, msg)
	default: // Success
		utils.SendJSONResponse(w, models.ResponseSingleDocument{
			StatusCode: http.StatusOK,
			Document:   doc,
		})
	}
}

func CreateDocument(ctx context.Context, w http.ResponseWriter, r *http.Request,
	file io.Reader, header *multipart.FileHeader, dc models.DocumentCreation) {
	// Read all flie to memory
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		msg := "Document can not be read"
		utils.SendJSONErrorResponse(w, http.StatusInternalServerError, msg)
		return
	}

	// Fill DocumentCreation fields
	dc.Title = header.Filename
	dc.Size = header.Size
	dc.Path = filepath.Clean(dc.Path)
	dc.Hash, err = utils.GenerateHash(fileBytes)
	if err != nil {
		msg := "Hash generation fault"
		utils.SendJSONErrorResponse(w, http.StatusInternalServerError, msg)
		return
	}

	// Set context to cancel if error occured
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)
	// Set error chan to save goroutines error
	errChan := make(chan error, 2)
	// Set wait group to synchronize goroutines
	var wg sync.WaitGroup
	wg.Add(2)

	// Get db connection
	con := storage.GetConnection()

	// Create document that will send user
	var doc models.Document

	// Save document record to database
	go func() {
		defer wg.Done()
		if doc, err = repository.CreateDocument(ctx, con, dc); err != nil {
			errChan <- err
			cancel(err)
		}

		// May be caused by a hash column integrity violation
		// http code 409 (Conflict) here
		if (doc == models.Document{}) {
			// TODO Change hash UNIQUENESS
			text := http.StatusText(http.StatusConflict)
			err := errors.New(text)
			errChan <- err
			cancel(err)
		}
	}()

	// Saves file to specified directory
	go func() {
		defer wg.Done()
		addPath := dc.Path
		reader := bytes.NewReader(fileBytes)
		if err := utils.UploadFile(ctx, addPath, reader, header); err != nil {
			errChan <- err
			cancel(err)
		}
	}()
	// Waits goroutines' end
	wg.Wait()
	// Close chan
	close(errChan)

	// If errors occured, write to logs
	for err := range errChan {
		if err != nil {
			log.Println(err)
		}
	}

	// Sends Response
	select {
	case <-ctx.Done():
		cause := context.Cause(ctx)
		// Choose correct http code for response
		if cause.Error() == http.StatusText(http.StatusConflict) {
			msg := "Document already exists"
			utils.SendJSONErrorResponse(w, http.StatusConflict, msg)
		} else {
			msg := "Internal context exceed"
			utils.SendJSONErrorResponse(w, http.StatusConflict, msg)
		}
	default:
		utils.SendJSONResponse(w, models.ResponseSingleDocument{
			StatusCode: http.StatusOK,
			Document:   doc,
		})
	}
}

func DownloadDocument(ctx context.Context, w http.ResponseWriter, r *http.Request, path string) {
	// Build path to file
	vol := volume.GetPath()
	path = filepath.Join(vol, path)
	// Open file
	file, err := os.Open(path)
	if err != nil {
		msg := fmt.Sprintf("Could not open file %v", path)
		utils.SendJSONErrorResponse(w, http.StatusInternalServerError, msg)
		return
	}
	defer file.Close()
	// Set headers
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", file.Name()))
	w.Header().Set("Content-Type", "application/octet-stream")
	// Save file to user
	_, err = io.Copy(w, file)
	if err != nil {
		msg := fmt.Sprintf("Could not copy file %v", path)
		utils.SendJSONErrorResponse(w, http.StatusInternalServerError, msg)
		return
	}
}
