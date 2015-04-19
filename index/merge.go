package index

// This function merges two indices.  The libraries from i2 will be
// appended to the list from i1.
func (i1 *Index) Merge(i2 Index) error {
	// We simply add all Library entries from i2 to the
	// end of the list of libraries in i1.
	i1.Libraries = append(i1.Libraries, i2.Libraries...)

	return nil
}
