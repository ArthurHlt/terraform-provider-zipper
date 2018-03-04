package zipper

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"path/filepath"
	"strings"
)

const (
	chunkForSha1 = 5 * 1024
)

// check if path is a web url
func IsWebURL(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}

var ZIP_FILE_EXT []string = []string{
	".zip",
	".jar",
	".war",
}
var TAR_FILE_EXT []string = []string{
	".tar",
}
var GZIP_FILE_EXT []string = []string{
	".gz",
	".gzip",
}
var TARGZ_FILE_EXT []string = []string{
	".tgz",
}

// Create sha1 from a reader by loading in maximum 5kb
func GetSha1FromReader(reader io.Reader) (string, error) {
	buf := new(bytes.Buffer)
	_, err := io.CopyN(buf, reader, chunkForSha1)
	if err != nil && err != io.EOF {
		return "", err
	}
	// we don't want to retrieve everything
	// so we close if it closeable
	if closer, ok := reader.(io.Closer); ok {
		err := closer.Close()
		if err != nil {
			return "", err
		}
	}

	h := sha1.New()
	h.Write(buf.Bytes())
	return hex.EncodeToString(h.Sum(nil)), nil
}

// check if file has one if extensions given
func HasExtFile(path string, extensions ...string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	if ext == "" {
		return false
	}
	for _, extension := range extensions {
		if extension == ext {
			return true
		}
	}
	return false
}

// check if file is a zip file
func IsZipFile(path string) bool {
	return HasExtFile(path, ZIP_FILE_EXT...)
}

// check if file is a tar file
func IsTarFile(path string) bool {
	return HasExtFile(path, TAR_FILE_EXT...)
}

// check if file is a tgz file
func IsTarGzFile(path string) bool {
	isTgz := HasExtFile(path, TARGZ_FILE_EXT...)
	if isTgz {
		return true
	}
	isGz := HasExtFile(path, GZIP_FILE_EXT...)
	if !isGz {
		return false
	}
	return IsTarFile(filepath.Ext(strings.TrimSuffix(path, filepath.Ext(path))))
}
