package generator

// FileConfig contains the transformer plugin settings to be applied ot the file
type FileConfig struct {
	// Version is a version of transformer rules to be applied to the file
	Version int32
	// Debug is a flag for displaying extra info
	Debug bool
}
