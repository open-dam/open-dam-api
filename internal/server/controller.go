package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// A Controller binds http requests to an api service and writes the service results to the http response
type Controller struct {
	service ApiServicer
}

// NewController creates an api controller
func NewController(s ApiServicer) *Controller {
	return &Controller{service: s}
}

// GetAssets -
func (c *Controller) GetAssets(w http.ResponseWriter, r *http.Request) {
	result, err := c.service.GetAssets()
	if err != nil {
		EncodeErrorResponse(err, w)
		return
	}

	EncodeJSONResponse(result, nil, w)
}

// PostAsset -
func (c *Controller) PostAsset(w http.ResponseWriter, r *http.Request) {
	assetCreate := AssetCreate{}
	if err := json.NewDecoder(r.Body).Decode(&assetCreate); err != nil {
		EncodeErrorResponse(err, w)
		return
	}
	result, err := c.service.PostAsset(assetCreate)
	if err != nil {
		EncodeErrorResponse(err, w)
		return
	}

	EncodeJSONResponse(result, nil, w)
}

// GetAsset -
func (c *Controller) GetAsset(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	assetId := params["asset_id"]
	result, err := c.service.GetAsset(assetId)
	if err != nil {
		EncodeErrorResponse(err, w)
		return
	}

	EncodeJSONResponse(result, nil, w)
}

// PutAsset -
func (c *Controller) PutAsset(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	assetId := params["asset_id"]
	assetUpdate := AssetUpdate{}
	if err := json.NewDecoder(r.Body).Decode(&assetUpdate); err != nil {
		EncodeErrorResponse(err, w)
		return
	}
	result, err := c.service.PutAsset(assetId, assetUpdate)
	if err != nil {
		EncodeErrorResponse(err, w)
		return
	}

	EncodeJSONResponse(result, nil, w)
}

// DeleteAsset -
func (c *Controller) DeleteAsset(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	assetId := params["asset_id"]
	result, err := c.service.DeleteAsset(assetId)
	if err != nil {
		EncodeErrorResponse(err, w)
		return
	}

	EncodeJSONResponse(result, nil, w)
}

// GetJob -
func (c *Controller) GetJob(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	jobId := params["job_id"]
	result, err := c.service.GetJob(jobId)
	if err != nil {
		EncodeErrorResponse(err, w)
		return
	}

	EncodeJSONResponse(result, nil, w)
}

// PostJob -
func (c *Controller) PostJob(w http.ResponseWriter, r *http.Request) {
	jobCreate := JobCreate{}
	if err := json.NewDecoder(r.Body).Decode(&jobCreate); err != nil {
		EncodeErrorResponse(err, w)
		return
	}
	result, err := c.service.PostJob(jobCreate)
	if err != nil {
		EncodeErrorResponse(err, w)
		return
	}

	EncodeJSONResponse(result, nil, w)
}
