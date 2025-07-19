package handlers

import (
	"context"
	"docshell/internal/v1/docs/models"
	"docshell/internal/v1/docs/service"
	"docshell/internal/v1/utils"
	"encoding/json"
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

func CreateDocument(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form, specifies a maximum upload size.
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		utils.SendJSONResponse(w, models.ResponseCode{
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	// Struct to put in it parsed body
	var dc models.DocumentCreation

	// Read body
	body := r.FormValue("meta")
	if body == "" { // if empty
		utils.SendJSONResponse(w, models.ResponseCode{
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	// FormFile returns the first file for the given key `file`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, header, err := r.FormFile("file")
	if err != nil {
		utils.SendJSONResponse(w, models.ResponseCode{
			StatusCode: http.StatusBadRequest,
		})
		return
	}
	defer file.Close()

	// Try to decode body into the struct
	if err := json.Unmarshal([]byte(body), &dc); err != nil {
		utils.SendJSONResponse(w, models.ResponseCode{
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	// Set context for chain
	ctx := context.Background()
	// Call next function
	service.CreateDocument(ctx, w, r, file, header, dc)
}
