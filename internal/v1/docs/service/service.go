package service

import (
	"bytes"
	"context"
	"docshell/internal/v1/docs/models"
	"docshell/internal/v1/docs/repository"
	"docshell/internal/v1/utils"
	"io"
	"log"
	"mime/multipart"
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
		utils.SendJSONResponse(w, models.ResponseCode{
			StatusCode: http.StatusInternalServerError,
		})
		return
	}
	// Build response
	res := models.ResponseMultipleDocuments{
		StatusCode: http.StatusOK,
		Documents:  make([]models.Document, len(docs)),
	}
	// Copy documents
	copy(res.Documents, docs)

	select {
	case <-ctx.Done(): // Context exceeded
		if ctx.Err() == context.DeadlineExceeded {
			utils.SendJSONResponse(w, models.ResponseCode{
				StatusCode: http.StatusRequestTimeout,
			})
		} else { // Other error
			utils.SendJSONResponse(w, models.ResponseCode{
				StatusCode: http.StatusInternalServerError,
			})
		}
	default: // Success
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
		utils.SendJSONResponse(w, models.ResponseCode{
			StatusCode: http.StatusInternalServerError,
		})
		log.Println(err)
		return
	}
	// Send NoContent if empty
	if doc == (models.Document{}) {
		utils.SendJSONResponse(w, models.ResponseCode{
			StatusCode: http.StatusNoContent,
		})
		return
	}
	// Build response
	res := models.ResponseSingleDocument{
		StatusCode: http.StatusOK,
		Document:   doc,
	}

	select {
	case <-ctx.Done(): // Context exceed
		utils.SendJSONResponse(w, models.ResponseCode{
			StatusCode: http.StatusRequestTimeout,
		})
	default: // Success
		utils.SendJSONResponse(w, res)
	}
}

func CreateDocument(ctx context.Context, w http.ResponseWriter, r *http.Request, file io.Reader, header *multipart.FileHeader, dc models.DocumentCreation) {
	// Read all flie to memory
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		utils.SendJSONResponse(w, models.ResponseCode{
			StatusCode: http.StatusInternalServerError,
		})
	}

	// Fill DocumentCreation
	dc.Size = int(header.Size)
	dc.Hash, err = utils.GenerateHash(fileBytes)
	if err != nil {
		utils.SendJSONResponse(w, models.ResponseCode{
			StatusCode: http.StatusBadRequest,
		})
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
	go func() { // Save to db
		defer wg.Done()
		var err error
		doc, err = repository.CreateDocument(ctx, con, dc)
		if err != nil {
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
	if ctx.Err() != nil {
		utils.SendJSONResponse(w, models.ResponseCode{
			StatusCode: http.StatusInternalServerError,
		})
	} else {
		utils.SendJSONResponse(w, models.ResponseSingleDocument{
			StatusCode: http.StatusOK,
			Document:   doc,
		})
	}
}
