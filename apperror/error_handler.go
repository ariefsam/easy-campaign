package apperror

import (
	"encoding/json"
	"errors"
	"net/http"
)

func HandleError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"error": true,
	}

	var customError *CustomError
	if errors.As(err, &customError) {
		response["message"] = customError.Error()
		w.WriteHeader(customError.Code)
	} else {
		response["message"] = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}

	//nolint:errcheck
	json.NewEncoder(w).Encode(response) //nolint:errchkjson
}
