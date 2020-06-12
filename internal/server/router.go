package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type ApiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewRouter creates a new router for any number of api routers
func NewRouter(c *Controller) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(handlers.CORS(), Logger)

	router.Path("/assets").Methods(http.MethodGet).Name("GetAssets").HandlerFunc(c.GetAssets)
	router.Path("/assets").Methods(http.MethodPost).Name("PostAsset").HandlerFunc(c.PostAsset)
	router.Path("/assets/{asset_id}").Methods(http.MethodGet).Name("GetAsset").HandlerFunc(c.GetAsset)
	router.Path("/assets/{asset_id}").Methods(http.MethodPut).Name("PutAsset").HandlerFunc(c.PutAsset)
	router.Path("/assets/{asset_id}").Methods(http.MethodDelete).Name("DeleteAsset").HandlerFunc(c.DeleteAsset)
	router.Path("/jobs").Methods(http.MethodPost).Name("PostJob").HandlerFunc(c.PostJob)
	router.Path("/jobs/{job_id}").Methods(http.MethodGet).Name("GetJob").HandlerFunc(c.GetJob)

	return router
}

// EncodeJSONResponse uses the json encoder to write an interface to the http response with an optional status code
func EncodeJSONResponse(i interface{}, status *int, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if status != nil {
		w.WriteHeader(*status)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	return json.NewEncoder(w).Encode(i)
}

// EncodeErrorResponse writes the error response for any error that occurred during the request
func EncodeErrorResponse(err error, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	status := http.StatusInternalServerError

	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(&ApiError{Code: status, Message: err.Error()})
}
