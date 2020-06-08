package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/open-dam/open-dam-worker/pkg/opendam"
	"gocloud.dev/docstore"
	"gocloud.dev/gcerrors"
)

// ApiServicer defines the api actions for the API
type ApiServicer interface {
	GetAssets() (interface{}, error)
	PostAsset() (interface{}, error)
	GetAsset(string) (interface{}, error)
	PutAsset(string, AssetUpdate) (interface{}, error)
	DeleteAsset(string) (interface{}, error)
	GetJob(string) (interface{}, error)
}

// ApiService is a service that implents the logic for the ApiServicer
// This service should implement the business logic for every endpoint for the API.
// Include any external packages or services that will be required by this service.
type ApiService struct {
	c *docstore.Collection
}

// NewApiService creates an api service
func NewApiService() ApiServicer {
	c, err := opendam.DocStoreFactory(os.Getenv("CONNECTION"))
	if err != nil {
		fmt.Println("WTFFF")
	}
	return &ApiService{
		c: c,
	}
}

// GetAssets -
func (s *ApiService) GetAssets() (interface{}, error) {
	var assets []opendam.Asset
	iter := s.c.Query().Get(context.Background())
	defer iter.Stop()

	// Query.Get returns an iterator. Call Next on it until io.EOF.
	for {
		var a opendam.Asset
		err := iter.Next(context.Background(), &a)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		} else {
			assets = append(assets, a)
		}
	}
	return &opendam.Assets{Assets: assets}, nil
}

// PostAsset -
func (s *ApiService) PostAsset() (interface{}, error) {
	// TODO - update PostAsset with the required logic for this service method.
	// Add api_default_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.
	return nil, errors.New("service method 'PostAsset' not implemented")
}

// GetAsset -
func (s *ApiService) GetAsset(assetId string) (interface{}, error) {
	asset := opendam.Asset{AssetID: assetId}
	err := s.c.Get(context.Background(), &asset)
	return &asset, err
}

// PutAsset -
func (s *ApiService) PutAsset(assetId string, assetUpdate AssetUpdate) (interface{}, error) {
	//TODO use distributed lock on asset
	asset := opendam.Asset{AssetID: assetId}
	notFound := false
	if err := s.c.Get(context.Background(), &asset); err != nil {
		if gcerrors.Code(err) != gcerrors.NotFound {
			return nil, err
		}
		notFound = true
	}

	// merge requested asset changes with any existing values
	if assetUpdate.Kind != "" {
		asset.Kind = assetUpdate.Kind
	}

	//TODO handle duplicate
	asset.Formats = append(asset.Formats, assetUpdate.Formats...)

	//TODO handle duplicate
	asset.Tags = append(asset.Tags, assetUpdate.Tags...)

	for k, v := range assetUpdate.Metadata {
		asset.Metadata[k] = v
	}

	// automatically update version
	asset.Version.Timestamp = time.Now().Unix()

	if notFound {
		if err := s.c.Create(context.Background(), &asset); err != nil {
			return nil, err
		}
	}
	if err := s.c.Replace(context.Background(), &asset); err != nil {
		return nil, err
	}

	return asset, nil
}

// DeleteAsset -
func (s *ApiService) DeleteAsset(assetId string) (interface{}, error) {
	err := s.c.Delete(context.Background(), &opendam.Asset{AssetID: assetId})
	return nil, err
}

// GetJob -
func (s *ApiService) GetJob(jobId string) (interface{}, error) {
	// TODO - update GetJob with the required logic for this service method.
	// Add api_default_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.
	return nil, errors.New("service method 'GetJob' not implemented")
}
