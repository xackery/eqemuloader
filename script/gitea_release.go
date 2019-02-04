package script

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"code.gitea.io/sdk/gitea"

	"github.com/pkg/errors"
)

func (s *Script) giteaReleaseDownload(name string, releaseURL string, user string, repo string, accessToken string, directory string) (err error) {

	client := &http.Client{}
	url := fmt.Sprintf("http://%s/api/v1/repos/%s/%s/releases", releaseURL, user, repo)
	if accessToken != "" {
		url += "?access_token=" + accessToken
	}
	if s.IsVerbose {
		fmt.Println("url:", url)
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		err = errors.Wrap(err, "failed to prepare request")
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		err = errors.Wrap(err, "failed to get releases")
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = fmt.Errorf("invalid status code: %d", resp.StatusCode)
		return
	}

	releases := []*gitea.Release{}
	err = json.NewDecoder(resp.Body).Decode(&releases)
	if err != nil {
		err = errors.Wrap(err, "failed to decode release")
		return
	}
	var latest *gitea.Release
	for _, release := range releases {
		if latest == nil {
			latest = release
			continue
		}
		if latest.PublishedAt.Before(release.PublishedAt) {
			latest = release
		}
	}
	if latest == nil {
		err = fmt.Errorf("failed to get latest release")
		return
	}
	if s.cacheGet(fmt.Sprintf("%s_last_release", name)) == latest.TagName {
		if s.IsVerbose {
			fmt.Println("releases up to date")
		}
		return
	}
	err = s.cacheSet(fmt.Sprintf("%s_last_release", name), latest.TagName)
	if err != nil {
		err = errors.Wrap(err, "failed to save latest release cache")
		return
	}
	fmt.Println("new release detected", latest.PublishedAt.Format("2006-01-02"), latest.TagName)

	var ok bool
	for _, asset := range latest.Attachments {
		fmt.Println(asset.Name)
		src := fmt.Sprintf("%s%s", directory, asset.Name)
		dst := directory
		err = s.fileDownload(asset.DownloadURL, src)
		if err != nil {
			err = errors.Wrapf(err, "failed to download %s", asset.DownloadURL)
			return
		}
		ok, err = s.fileExtract(src, dst)
		if err != nil {
			return
		}
		if ok { //remove extracted contents
			if s.IsVerbose {
				fmt.Println("removing zipped file", src)
			}
			os.Remove(src)
		}
	}
	//change permissions of executable files
	files := []string{"zone", "world", "ucs", "queryserv", "shared_memory"}
	for _, file := range files {
		target := fmt.Sprintf("%s%s", directory, file)
		os.Chmod(target, 0775)
	}
	return
}
