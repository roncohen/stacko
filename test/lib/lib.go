package lib

import (
	"compress/gzip"
)

func WillError() {
	var ptr *gzip.Reader
	ptr.Close()
}
