package utils

import (
	"encoding/json"
	"net/http"
)

func SetHeaders(w http.ResponseWriter) {
	// w.Header().Add("application/json", )
}

func SendJSONResponse(w http.ResponseWriter, res any) {
	json.NewEncoder(w).Encode(res)
}
