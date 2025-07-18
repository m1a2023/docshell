package handlers

import (
	"context"
	"docshell/internal/v1/docs/models"
	"docshell/internal/v1/docs/service"
	"docshell/internal/v1/utils"
	"net/http"
	"strconv"
)

func GetAllDocuments(w http.ResponseWriter, r *http.Request) {
	// Set context for chain
	ctx := context.Background()
	// Call next function and pass context
	service.GetAllDocuments(ctx, w, r)
}

func GetDocumentById(w http.ResponseWriter, r *http.Request) {
	// Read path value
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || 0 >= id { // Send 400
		utils.SendJSONResponse(w, models.ResponseCode{
			StatusCode: http.StatusBadRequest,
		})
		return
	}
	// Set context for chain
	ctx := context.Background()
	// Call next function and pass context
	service.GetDocumentById(ctx, w, r, id)
}
