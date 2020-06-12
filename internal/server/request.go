package server

import (
	"github.com/open-dam/open-dam-worker/pkg/opendam"
)

// AssetUpdate - A limited view of an asset with only editable fields. Formats, tags, and metadata are merged with any existing values
type AssetUpdate struct {
	Kind string `json:"kind,omitempty"`

	File opendam.File `json:"file,omitempty"`

	// additional assets/files associated with the asset
	Formats []opendam.Asset `json:"formats,omitempty"`

	// A list of metadata tags associated with the asset
	Tags []string `json:"tags,omitempty"`

	// Any user supplied metadata for the asset
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type AssetCreate struct {
	URL     string `json:"url"`
	AssetID string `json:"asset_id"`
}

type JobCreate struct {
	Tasks []Task `json:"tasks"`
}

type Task struct {
	Name string   `json:"name"`
	Args []string `json:"args"`
}
