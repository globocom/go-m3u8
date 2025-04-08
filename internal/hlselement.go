package internal

// The HLSElement data type holds the following attributes:
//   - Name: The name of the Element (e.g. tag name).
//   - URI: The Uniform Resource Identifier of the Element (if applicable).
//   - Attrs: In-manifest Element attributes, in key-value format.
//   - Details: Not-in-manifest Element attributes, in key-value format.
type HLSElement struct {
	Name    string
	URI     string
	Attrs   map[string]string
	Details map[string]string
}
