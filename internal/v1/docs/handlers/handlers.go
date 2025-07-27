package handlers

import (
	"context"
	"docshell/internal/v1/docs/models"
	"docshell/internal/v1/docs/service"
	"docshell/internal/v1/utils"
	"encoding/json"
	"fmt"
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
	if err != nil || 0 >= id {
		msg := fmt.Sprintf("Path value 'id=%v' incorrect", id)
		utils.SendJSONErrorResponse(w, http.StatusBadRequest, msg)
		return
	}
	// Set context for chain
	ctx := context.Background()
	// Call next function and pass context
	service.GetDocumentById(ctx, w, r, id)
}

func CreateDocument(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form, specifies a maximum upload size.
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		msg := "During parsing multupart form"
		utils.SendJSONErrorResponse(w, http.StatusBadRequest, msg)
		return
	}

	// Struct to put in it parsed body
	var dc models.DocumentCreation

	// Read body
	body := r.FormValue("meta")
	if body == "" { // if empty
		msg := "Body is empty"
		utils.SendJSONErrorResponse(w, http.StatusBadRequest, msg)
		return
	}

	// FormFile returns the first file for the given key `file`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, header, err := r.FormFile("file")
	if err != nil {
		msg := "Form file incorrect"
		utils.SendJSONErrorResponse(w, http.StatusBadRequest, msg)
		return
	}
	defer file.Close()

	// Try to decode body into the struct
	if err := json.Unmarshal([]byte(body), &dc); err != nil {
		msg := "JSON is incorrect"
		utils.SendJSONErrorResponse(w, http.StatusBadRequest, msg)
		return
	}

	// Set context for chain
	ctx := context.Background()
	// Call next function
	service.CreateDocument(ctx, w, r, file, header, dc)
}
