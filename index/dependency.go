package index

type Dependency struct {
	Owner   string `json:"owner"`
	Name    string `json:"name"`
	Version string `json:"version"`
}
