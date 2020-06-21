package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"github.com/open-dam/open-dam-worker/pkg/opendam"
	"github.com/sirupsen/logrus"
	"gocloud.dev/blob"
	"gocloud.dev/docstore"
	"gocloud.dev/gcerrors"
)

// ApiServicer defines the api actions for the API
type ApiServicer interface {
	GetAssets() (interface{}, error)
	PostAsset(AssetCreate) (interface{}, error)
	GetAsset(string) (interface{}, error)
	PutAsset(string, AssetUpdate) (interface{}, error)
	DeleteAsset(string) (interface{}, error)
	GetJob(string) (interface{}, error)
	PostJob(JobCreate) (interface{}, error)
}

// ApiService is a service that implents the logic for the ApiServicer
// This service should implement the business logic for every endpoint for the API.
// Include any external packages or services that will be required by this service.
type ApiService struct {
	docs      *docstore.Collection
	bucket    *blob.Bucket
	machinery *machinery.Server
	client    *http.Client
	logger    *logrus.Entry
}

// NewApiService creates an api service
func NewApiService() ApiServicer {
	logger := opendam.Logger()
	log.Set(logger)

	cnf, err := config.NewFromEnvironment(true)
	if err != nil {
		logger.WithError(err).Fatal("failed to build machinery config from environment")
	}

	logger.WithField("config", cnf).Debug("got config")

	server, err := machinery.NewServer(cnf)
	if err != nil {
		logger.WithError(err).Fatal("failed to start machinery server")
	}

	c, err := opendam.DocStoreFactory(os.Getenv("CONNECTION"))
	if err != nil {
		logger.WithError(err).Fatal("failed to open document store connection")
	}

	bucket, err := blob.OpenBucket(context.Background(), os.Getenv("BLOB_CONNECTION"))
	if err != nil {
		logger.WithError(err).Fatal("failed to open blob storage connection")
	}

	return &ApiService{
		docs:      c,
		bucket:    bucket,
		machinery: server,
		client:    &http.Client{},
		logger:    logger,
	}
}

// GetAssets -
func (s *ApiService) GetAssets() (interface{}, error) {
	var assets []opendam.Asset
	iter := s.docs.Query().Get(context.Background())
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
func (s *ApiService) PostAsset(assetCreate AssetCreate) (interface{}, error) {
	assetID := assetCreate.AssetID
	if assetID != "" {
		if _, err := s.GetAsset(assetID); err != nil {
			return nil, err
		}
	}

	if assetID == "" {
		assetID = uuid.New().String()
	}

	kind, contentType, err := s.assetKind(assetCreate.URL)
	if err != nil {
		return nil, err
	}

	f, err := s.upload(assetCreate.URL, assetID)
	if err != nil {
		return nil, err
	}
	f.ContentType = contentType

	b, err := s.bucket.ReadAll(context.Background(), assetID)
	if err != nil {
		return nil, err
	}
	s.logger.Info(len(b))

	extract, _ := tasks.NewSignature("extract", []tasks.Arg{
		{Type: "string", Value: assetID},
		{Type: "string", Value: assetID},
	})
	sigs := []*tasks.Signature{extract}
	switch kind {
	case "image":
		imageAnalysis, _ := tasks.NewSignature("imageanalysis", []tasks.Arg{
			{Type: "string", Value: assetID},
			{Type: "string", Value: assetID},
		})
		imageCreation, _ := tasks.NewSignature("imagecreation", []tasks.Arg{
			{Type: "string", Value: assetID},
			{Type: "string", Value: assetID},
		})
		sigs = append(sigs, imageAnalysis, imageCreation)
	case "audio":
		soundwave, _ := tasks.NewSignature("soundwave", []tasks.Arg{
			{Type: "string", Value: assetID},
			{Type: "string", Value: assetID},
		})
		sigs = append(sigs, soundwave)

	}

	chain, _ := tasks.NewChain(sigs...)
	_, err = s.machinery.SendChain(chain)
	if err != nil {
		return nil, fmt.Errorf("send chain err %s", err.Error())
	}

	asset, err := s.PutAsset(assetID, AssetUpdate{
		Kind: kind,
		File: f,
	})
	if err != nil {
		// TODO cancel job somehow
		return nil, err
	}

	// TODO handle "job id"
	return asset, nil
}

func (s *ApiService) assetKind(URL string) (string, string, error) {
	resp, err := s.client.Get(URL)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	mime, err := mimetype.DetectReader(resp.Body)
	if err != nil {
		return "", "", err
	}
	m := strings.ToLower(mime.String())
	s.logger.WithFields(logrus.Fields{
		"mime":         m,
		"ext":          mime.Extension(),
		"content_type": resp.Header.Get("Content-Type"),
	}).Debug("mime type detected")

	if strings.HasPrefix(m, "image") {
		return "image", m, nil
	} else if strings.HasPrefix(m, "audio") {
		return "audio", m, nil
	} else if strings.HasPrefix(m, "video") {
		return "video", m, nil
	} else if strings.HasPrefix(m, "text") {
		return "text", m, nil
	}
	return "unknown", m, nil

}

func (s *ApiService) upload(url, assetID string) (opendam.File, error) {
	// check that url is already on blob storage and return file with key
	var file opendam.File
	wr, err := s.bucket.NewWriter(context.Background(), assetID, nil)
	if err != nil {
		return file, err
	}
	defer wr.Close()

	resp, err := s.client.Get(url)
	if err != nil {
		return file, err
	}
	defer resp.Body.Close()

	_, err = wr.ReadFrom(resp.Body)
	if err != nil {
		return file, err
	}

	return opendam.File{
		Name:   "", //TODO how do we get filename if frontend client is uploading files to storage
		Source: assetID,
	}, nil
}

// GetAsset -
func (s *ApiService) GetAsset(assetId string) (interface{}, error) {
	asset := opendam.Asset{AssetID: assetId}
	err := s.docs.Get(context.Background(), &asset)
	return &asset, err
}

// PutAsset -
func (s *ApiService) PutAsset(assetId string, assetUpdate AssetUpdate) (interface{}, error) {
	//TODO use distributed lock on asset
	asset := opendam.Asset{AssetID: assetId}
	notFound := false
	if err := s.docs.Get(context.Background(), &asset); err != nil {
		if gcerrors.Code(err) != gcerrors.NotFound {
			return nil, err
		}
		notFound = true
	}

	// merge requested asset changes with any existing values
	if assetUpdate.Kind != "" {
		asset.Kind = assetUpdate.Kind
	}

	//TODO merge file data
	asset.File = &assetUpdate.File

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
		if err := s.docs.Create(context.Background(), &asset); err != nil {
			return nil, err
		}
	}
	if err := s.docs.Replace(context.Background(), &asset); err != nil {
		return nil, err
	}

	return asset, nil
}

// DeleteAsset -
func (s *ApiService) DeleteAsset(assetId string) (interface{}, error) {
	// TODO this needs to trigger a delete workflow to cleanup all assets
	err := s.docs.Delete(context.Background(), &opendam.Asset{AssetID: assetId})
	return nil, err
}

// GetJob -
func (s *ApiService) GetJob(jobId string) (interface{}, error) {
	// TODO - update GetJob with the required logic for this service method.
	// Add api_default_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.
	return nil, errors.New("service method 'GetJob' not implemented")
}

// PostJob -
func (s *ApiService) PostJob(jobCreate JobCreate) (interface{}, error) {
	var sigs []*tasks.Signature
	for _, t := range jobCreate.Tasks {
		var args []tasks.Arg
		for _, a := range t.Args {
			args = append(args, tasks.Arg{Type: "string", Value: a})
		}
		task, _ := tasks.NewSignature(t.Name, args)
		sigs = append(sigs, task)
	}
	chain, _ := tasks.NewChain(sigs...)
	_, err := s.machinery.SendChain(chain)
	if err != nil {
		return nil, fmt.Errorf("send chain err %s", err.Error())
	}

	return nil, nil
}
