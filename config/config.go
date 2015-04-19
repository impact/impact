package config

import (
	"os"
	"path"
	"path/filepath"
)

func ReadSettings() (Settings, error) {
	dir, _ := filepath.Abs(path.Join(os.Getenv("GOPATH"), "src", "github.com", "xogeny",
		"impact", "sample_index.json"))

	return Settings{
		//		Indices: []string{"https://impact.modelica.org/impact_data.json"},
		Indices: []string{"file://" + dir},
	}, nil
}
