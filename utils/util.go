package utils

import (
	"encoding/json"
	"net/http"

)

func Message(status bool, message string, code int) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message, "code": code}
}

func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
