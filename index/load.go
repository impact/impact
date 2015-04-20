package index

import (
	"fmt"
	"log"

	"github.com/xogeny/impact/config"
)

func LoadIndex(verbose bool) (*Index, error) {
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
		if verbose {
			log.Printf("Loading index data from %s", index_url)
		}
		err = ind.ParseIndex(index_url)
		if err != nil {
			return nil, fmt.Errorf("Error parsing index at %s: %v", index_url, err)
		}
	}

	return ind, nil
}
