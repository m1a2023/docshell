package service

import (
	"context"
	utils "docshell/internal/v1/docs"
	"docshell/internal/v1/docs/models"
	"docshell/internal/v1/docs/repository"

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
