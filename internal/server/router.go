package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// A Route defines the parameters for an api endpoint
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
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
