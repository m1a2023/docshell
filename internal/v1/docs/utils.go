package utils

import (
	"encoding/json"
	"net/http"
)

func SendJSONResponse(w http.ResponseWriter, res any) {
	json.NewEncoder(w).Encode(res)
}
