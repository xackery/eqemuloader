package script

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
)

// binCheck checks if a repo exists, and clones if not
// then it will scan bin releases for any new binaries
func (s *Script) binCheck() (err error) {
	_, err = git.PlainOpen(s.Bin.Directory)
	if err != nil {
		err = s.binClone()
		if err != nil {
			err = errors.Wrap(err, "cloning")
			return
		}
		return
	}
	return
}

func (s *Script) binClone() (err error) {
	if s.Bin.URL == "" {
		fmt.Println("skipping cloning, bin.url not set")
		return
	}
	if s.Bin.AuthType == "" {
		fmt.Println("skipping cloning, bin.auth_type not set")
		return
	}
	fmt.Println("Cloning", s.Bin.URL, "to", s.Bin.Directory)
	fmt.Println("(this may take a while!!)")

	auth, err := getAuth(s.Bin.AuthType, s.Bin.AuthUsername, s.Bin.AuthPassword, s.Bin.AuthKey)
	if err != nil {
		err = errors.Wrap(err, "auth")
		return
	}
	_, err = git.PlainClone(s.Bin.Directory, false, &git.CloneOptions{
		URL:               s.Bin.URL,
		Progress:          os.Stdout,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Auth:              auth,
	})
	if err != nil {
		err = errors.Wrap(err, "failed to clone git")
		os.Remove(s.Bin.Directory)
		return
	}
	return
}

func (s *Script) binPull() (err error) {
	repo, err := git.PlainOpen(s.Bin.Directory)
	if err != nil {
		err = s.binClone()
		if err != nil {
			err = errors.Wrap(err, "cloning")
			return
		}
		return
	}

	// Get the working directory for the repository
	w, err := repo.Worktree()
	if err != nil {
		err = errors.Wrap(err, "getting worktree")
		return
	}

	/*st, err := w.Status()
	if err != nil {
		err = errors.Wrap(err, "getting status")
		return
	}

	if st.IsClean() {
		err = fmt.Errorf("repository is not clean")
		return
	}*/

	auth, err := getAuth(s.Bin.AuthType, s.Bin.AuthUsername, s.Bin.AuthPassword, s.Bin.AuthKey)
	if err != nil {
		err = errors.Wrap(err, "auth failed")
		return
	}
	// Pull the latest changes from the origin remote and merge into the current branch
	err = w.Pull(&git.PullOptions{
		Auth:              auth,
		RemoteName:        "origin",
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			fmt.Println("already up to date.")
			err = nil
			return
		}
		err = errors.Wrap(err, "pulling")
		return
	}

	// Print the latest commit that was just pulled
	ref, err := repo.Head()
	if err != nil {
		err = errors.Wrap(err, "head failed")
		return
	}
	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		err = errors.Wrap(err, "ref hash failed")
		return
	}
	fmt.Println(commit)
	return
}

func (s *Script) binReleaseDownload() (err error) {
	switch strings.ToLower(s.Bin.ReleaseType) {
	case "github":
		err = s.githubReleaseDownload()
		if err != nil {
			return
		}
	case "gitea":
		err = s.giteaReleaseDownload("bin", s.Bin.ReleaseURL, s.Bin.ReleaseUser, s.Bin.ReleaseRepo, s.Bin.ReleaseAccessToken, s.Bin.Directory)
		if err != nil {
			return
		}
	default:
		err = fmt.Errorf("unknown bin release type: %s", s.Bin.ReleaseType)
	}
	return
}
