package scanner_int

type BucketObject struct {
	Key  string
	Size int64
}

type Config struct {
	Aws struct {
		Enabled    bool
		BucketName string
	}
	Local struct {
		Enabled bool
		Path    string
	}
}
