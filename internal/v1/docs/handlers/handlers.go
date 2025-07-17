package handlers

import (
	"context"
	"docshell/internal/v1/docs/service"
	"net/http"
)

func GetAllDocuments(w http.ResponseWriter, r *http.Request) {
	// Set context for chain
	ctx := context.Background()
	// Call next function and pass context
	service.GetAllDocuments(ctx, w, r)
}
