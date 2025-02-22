package helper

import (
	"campaign/apperror"
	"encoding/json"
	"errors"
	"net/http"
)

func HandleError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"error": true,
	}

	var customError *apperror.CustomError
	if errors.As(err, &customError) {
		response["message"] = customError.Error()
		w.WriteHeader(customError.Code)
	} else {
		response["message"] = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(response)
}
