package script

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

// Wget fetches a remote url and gets a file, to dl locally
func (s *Script) Wget(url string, dst string) (err error) {
	fmt.Println("downloading", url)
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
