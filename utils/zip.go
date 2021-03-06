package utils

import (
	"compress/flate"

	"github.com/mholt/archiver"

	"github.com/baetyl/baetyl-go/v2/errors"
)

var defaultZip = &archiver.Zip{
	CompressionLevel:     flate.DefaultCompression,
	MkdirAll:             true,
	SelectiveCompression: true,
	OverwriteExisting:    true,
}

// Zip zip source files to destination file
func Zip(sources []string, destination string) error {
	return errors.Trace(defaultZip.Archive(sources, destination))
}

// Unzip unzip source file to destination
func Unzip(source, destination string) error {
	return errors.Trace(defaultZip.Unarchive(source, destination))
}
