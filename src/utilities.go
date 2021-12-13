package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// CopyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()
	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}
	return
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
func CopyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)
	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}
	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}
	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}
	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}
			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}
	return
}

//func FileNameWithoutPathOrExtension(fileName string) string {
//	var fileNameWithoutExtension string
//	if pos := strings.LastIndexByte(fileName, '.'); pos != -1 {
//		fileNameWithoutExtension = fileName[:pos]
//	}
//	return filepath.Base(fileNameWithoutExtension)
//}

func GetIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

var markdownLinkRegex = regexp.MustCompile("\\[[\\w\\d\\s]+\\]\\(([\\w\\d./?=#:\\-_]+)\\)")
var htmlLinkRegex = regexp.MustCompile("<a href=\"([\\w\\d./?=#:\\-\\(\\)_]+)\">")

func ExtractLinkURLs(content, baseUrlForRelativeLinks string) []string {
	links := make([]string, 0)
	for _, match := range markdownLinkRegex.FindAllStringSubmatch(content, -1) {
		parsedLink, err := ReturnAbsoluteLink(match[1], baseUrlForRelativeLinks)
		if err == nil {
			links = append(links, parsedLink)
		}
	}
	for _, match := range htmlLinkRegex.FindAllStringSubmatch(content, -1) {
		parsedLink, err := ReturnAbsoluteLink(match[1], baseUrlForRelativeLinks)
		if err == nil {
			links = append(links, parsedLink)
		}
	}
	for _, link := range links {
		parsedLink, err := url.Parse(link)
		if err != nil {
			logger.Error(err.Error())
		} else {
			if !parsedLink.IsAbs() {

			}
		}
	}
	return links
}

func ReturnAbsoluteLink(link, baseUrlForRelativeLinks string) (string, error) {
	newLink := link
	parsedLink, err := url.Parse(link)
	if err != nil {
		return newLink, err
	} else {
		if !parsedLink.IsAbs() {
			switch {
			case strings.HasPrefix(link, "/"):
				newLink = baseUrlForRelativeLinks + link
			case strings.HasPrefix(link, "./"):
				newLink = baseUrlForRelativeLinks + link[1:(len(link)-1)]
			default:
				newLink = baseUrlForRelativeLinks + "/" + link
			}

		}
	}
	return newLink, err
}
