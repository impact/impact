package index

import (
	"fmt"
	"github.com/xogeny/impact/config"
)

func LoadIndex() (*Index, error) {
	// Read settings
	settings, err := config.ReadSettings()
	if err != nil {
		return nil, fmt.Errorf("Error reading user settings: %v", err)
	}

	ind := NewIndex()

	// Parse the various indices (according to the order defined in
	// settings).

	// N.B. - The order of the indices in the settings indicates the
	// order they will be searched in.  So entries in earlier indices
	// will match before entries in later indices.
	for _, index_url := range settings.Indices {
		err = ind.ParseIndex(index_url)
		if err != nil {
			return nil, fmt.Errorf("Error parsing index at %s: %v", ind, err)
		}
	}

	return ind, nil
}
