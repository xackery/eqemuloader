package script

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

func (s *Script) githubReleaseDownload() (err error) {
	if s.Bin.ReleaseType != "github" {
		err = fmt.Errorf("invalid release type for github: %s", s.Bin.ReleaseType)
		return
	}
	client := &http.Client{}
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", s.Bin.ReleaseUser, s.Bin.ReleaseRepo)
	if s.Bin.ReleaseAccessToken != "" {
		url += "?access_token=" + s.Bin.ReleaseAccessToken
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

	release := &github.RepositoryRelease{}
	err = json.NewDecoder(resp.Body).Decode(&release)
	if err != nil {
		err = errors.Wrap(err, "failed to decode release")
		return
	}
	/*
		var latest *github.RepositoryRelease
		for _, release := range releases {
			if latest == nil {
				latest = release
				continue
			}
			if latest.GetCreatedAt().After(release.GetCreatedAt().Time) {
				latest = release
			}
		}*/
	if release == nil {
		err = fmt.Errorf("failed to get releases, bad response")
		return
	}
	for _, asset := range release.Assets {
		fmt.Println(asset.GetBrowserDownloadURL(), asset.GetName())
	}
	return
}
