package zipper

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"

	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type HttpHandler struct {
}

func (h HttpHandler) Zip(src *Source) (ZipReadCloser, error) {
	client := CtxHttpClient(src)
	path := src.Path
	cleanFunc := func() error {
		return nil
	}
	u, _ := url.Parse(src.Path)
	username := ""
	password := ""
	if u.User != nil && u.User.Username() != "" {
		username = u.User.Username()
		password, _ = u.User.Password()
	}
	u.User = nil
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	if username != "" {
		req.SetBasicAuth(username, password)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	err = h.checkRespHttpError(resp)
	if err != nil {
		return nil, err
	}
	if IsTarFile(path) {
		defer resp.Body.Close()
		return h.tar2Zip(resp.Body)
	}
	if IsTarGzFile(path) {
		defer resp.Body.Close()
		return h.targz2Zip(resp.Body)
	}
	return NewZipFile(resp.Body, resp.ContentLength, cleanFunc), nil
}
func (h HttpHandler) checkRespHttpError(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	content := ""
	if err == nil {
		content = string(b)
	}
	return fmt.Errorf(
		"Error occured when dowloading file: %d %s: \n%s",
		resp.StatusCode,
		http.StatusText(resp.StatusCode),
		content,
	)
}
func (h HttpHandler) Detect(src *Source) bool {
	path := src.Path
	return IsWebURL(path) && (IsZipFile(path) || IsTarFile(path) || IsTarGzFile(path))
}

func (h HttpHandler) targz2Zip(r io.ReadCloser) (*ZipFile, error) {
	gzf, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return h.tar2Zip(gzf)
}
func (h HttpHandler) tar2Zip(r io.ReadCloser) (*ZipFile, error) {
	zipFile, err := ioutil.TempFile("", "downloads-zipper")
	if err != nil {
		return nil, err
	}
	cleanFunc := func() error {
		return os.Remove(zipFile.Name())
	}
	err = h.writeTarToZip(r, zipFile)
	if err != nil {
		zipFile.Close()
		return nil, err
	}
	zipFile.Close()
	file, err := os.Open(zipFile.Name())
	if err != nil {
		return nil, err
	}
	fs, _ := file.Stat()
	return NewZipFile(file, fs.Size(), cleanFunc), nil
}
func (h HttpHandler) writeTarToZip(r io.Reader, zipFile *os.File) error {
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()
	tarReader := tar.NewReader(r)
	hasRootFolder := false
	i := 0
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fileInfo := header.FileInfo()
		if i == 0 && fileInfo.IsDir() {
			hasRootFolder = true
			continue
		}
		zipHeader, err := zip.FileInfoHeader(fileInfo)
		if err != nil {
			return err
		}
		if !hasRootFolder {
			zipHeader.Name = header.Name
		} else {
			splitFile := strings.Split(header.Name, "/")
			zipHeader.Name = strings.Join(splitFile[1:], "/")
		}
		if !fileInfo.IsDir() {
			zipHeader.Method = zip.Deflate
		}
		w, err := zipWriter.CreateHeader(zipHeader)
		if err != nil {
			return err
		}
		i++
		if fileInfo.IsDir() {
			continue
		}
		_, err = io.Copy(w, tarReader)
	}
	return nil
}
func (h HttpHandler) Sha1(src *Source) (string, error) {
	client := CtxHttpClient(src)
	path := src.Path
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	err = h.checkRespHttpError(resp)
	if err != nil {
		return "", err
	}
	return GetSha1FromReader(resp.Body)
}
func (h HttpHandler) Name() string {
	return "http"
}
