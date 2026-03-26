package handlers

import (
	"net/http"

	"uniswap-campus-marketplace/apiresponse"
)

func writeSuccess(w http.ResponseWriter, status int, data interface{}) {
	apiresponse.WriteSuccess(w, status, data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	apiresponse.WriteError(w, status, message)
}
