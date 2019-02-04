package script

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/codeclysm/extract"
	"github.com/pkg/errors"
)

// DownloadRemote fetches a remote url and gets a file, to dl locally
func (s *Script) fileDownload(url string, dst string) (err error) {
	out, err := os.Create(dst)
	if err != nil {
		err = errors.Wrapf(err, "failed to create %s", dst)
		return
	}
	defer out.Close()
	resp, err := http.Get(url)
	if err != nil {
		err = errors.Wrapf(err, "failed to get %s", url)
		return
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		err = errors.Wrapf(err, "failed to download %s", url)
		return
	}
	return
}

func (s *Script) fileExtract(src string, dst string) (ok bool, err error) {
	ok = true
	ext := strings.ToLower(filepath.Ext(src))
	_, dstfile := filepath.Split(src)
	dstfile = strings.Replace(dstfile, ext, "", -1)
	if s.IsVerbose {
		fmt.Println("extension of", src, "is", ext)
	}

	switch ext {
	case ".zip":
		err = fileExtract(src, dst)
		if err != nil {
			err = errors.Wrapf(err, "failed to unzip %s", src)
			return
		}
	case ".bz2":
		err = fileExtract(src, dst)
		if err != nil {
			err = errors.Wrapf(err, "failed to unbz2 %s", src)
			return
		}
	case ".tar":
		err = fileExtract(src, dst)
		if err != nil {
			err = errors.Wrapf(err, "failed to untar %s", src)
			return
		}
	case ".gzip":
		dst += dstfile
		err = fileExtract(src, dst)
		if err != nil {
			err = errors.Wrapf(err, "failed to ungzip %s", src)
			return
		}
	case ".gz":
		dst += dstfile
		err = fileExtract(src, dst)
		if err != nil {
			err = errors.Wrapf(err, "failed to ungzip %s", src)
			return
		}
	default:
		ok = false
	}
	return
}

func fileGZip(source string, target string) (err error) {
	reader, err := os.Open(source)
	if err != nil {
		return err
	}

	filename := filepath.Base(source)
	//target = filepath.Join(target, fmt.Sprintf("%s.gz", filename))
	writer, err := os.Create(target)
	if err != nil {
		return err
	}
	defer writer.Close()

	archiver := gzip.NewWriter(writer)
	archiver.Name = filename
	defer archiver.Close()

	_, err = io.Copy(archiver, reader)
	return err
}

func fileExtract(src string, dst string) (err error) {
	fmt.Println("extracting to", dst)
	f, err := os.Open(src)
	if err != nil {
		return
	}
	defer f.Close()
	ctx := context.Background()
	err = extract.Archive(ctx, f, dst, nil)
	if err != nil {
		return
	}
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func fileCopy(src, dst string) (err error) {
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
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
