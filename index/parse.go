package index

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
)

// The ParseIndex function reads index information from a given URL and
// then merges it into the associated index.
//
// TODO: This will eventually need to be able to handle different
// index versions (by parsing the Version field in the index first and
// then determining how, or if, to parse it).
func (index *Index) ParseIndex(index_url string) error {
	// Parse the URL to break it down
	u, err := url.Parse(index_url)
	if err != nil {
		return fmt.Errorf("Error parsing url: %v", err)
	}

	var bytes []byte

	// Depending on the scheme, extract the underlying bytes
	switch u.Scheme {
	case "file":
		// If it is a file, open the file...
		f, err := os.Open(u.Path)
		if err != nil {
			return fmt.Errorf("Unable to open file '%s': %v", u.Path, err)
		}

		// ...and read the contents
		data, err := ioutil.ReadAll(f)
		if err != nil {
			return fmt.Errorf("Unable to read file '%s': %v", u.Path, err)
		}

		bytes = data
	//case "http":
	//	fallthrough
	//case "https":
	default:
		// If we don't suppor the scheme, throw an error
		return fmt.Errorf("Unsupported URL scheme '%s', unable to download", u.Scheme)
	}

	// Now create an empty index to read the bytes into
	contents := Index{}

	// Unmarshal the bytes as JSON into the new empty index
	err = json.Unmarshal(bytes, &contents)
	if err != nil {
		return fmt.Errorf("Unable to parse JSON at %s: %v", index_url, err)
	}

	// Assuming everything worked, merge the next information into the existing
	// index.
	// N.B. - The entries in contents will be appended to the index.  So a search
	// for the first matching entry will always favor the original index over the
	// entries just parsed.  This is by design.
	err = index.Merge(contents)
	if err != nil {
		return fmt.Errorf("Error merging indices: %v", err)
	}

	return nil
}
