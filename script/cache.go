package script

import (
	"fmt"
	"os"

	toml "github.com/pelletier/go-toml"
	"github.com/pkg/errors"
)

func (s *Script) cacheSet(key string, value string) (err error) {
	if s.cache.Entries == nil {
		s.cache.Entries = make(map[string]string)
	}
	s.cache.Entries[key] = value
	err = s.cacheSave()
	if err != nil {
		return
	}
	return
}

func (s *Script) cacheGet(key string) string {
	val, ok := s.cache.Entries[key]
	if !ok {
		return ""
	}
	return val
}

func (s *Script) cacheLoad() (err error) {
	fi, err := os.Stat("./.cache")
	if err != nil {
		err = os.Mkdir("./.cache", 0755)
		if err != nil {
			err = errors.Wrap(err, "failed to make .cache directory")
			return
		}
		fi, err = os.Stat("./.cache")
		if err != nil {
			err = errors.Wrap(err, "failed cache creation")
			return
		}
	}
	if !fi.IsDir() {
		err = fmt.Errorf(".cache is not a directory but needs to be (perhaps a file exists?")
		return
	}

	f, err := os.Open("./.cache/loader.cache")
	if err != nil {
		s.cache.Entries = make(map[string]string)
		err = s.cacheSave()
		if err != nil {
			err = errors.Wrap(err, "failed to create new cache")
			return
		}
		return
	}
	defer f.Close()

	tree, err := toml.LoadReader(f)
	if err != nil {
		err = errors.Wrap(err, "failed to load config")
		return
	}

	err = tree.Unmarshal(&s.cache)
	if err != nil {
		err = errors.Wrap(err, "failed to unmarshal")
		return
	}
	if s.cache.Entries == nil {
		s.cache.Entries = make(map[string]string)
	}
	return
}

func (s *Script) cacheSave() (err error) {
	f, err := os.Create("./.cache/loader.cache")
	if err != nil {
		f, err = os.Create("./.cache/loader.cache")
		if err != nil {
			err = errors.Wrap(err, "failed to create cache file")
			return
		}
	}
	defer f.Close()
	err = toml.NewEncoder(f).Encode(s.cache)
	if err != nil {
		return
	}
	return
}
