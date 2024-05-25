package memcache

import (
	"errors"
	"io/ioutil"
	"strings"
)

type Memcache struct {
	path  string
	cache map[int64]string
}

func NewMemcache(path string) (*Memcache, error) {
	if path == "" {
		return nil, errors.New("invalid path")
	}

	cache, err := readCategoriesFromFile(path)
	if err != nil {
		return nil, err
	}

	return &Memcache{
		path:  path,
		cache: cache,
	}, nil
}

func readCategoriesFromFile(filename string) (map[int64]string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	categories := make(map[int64]string)
	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		categories[int64(i)] = line
	}

	return categories, nil
}

func (m *Memcache) Lookup(id int64) (string, error) {
	if res, ok := m.cache[id]; ok {
		return res, nil
	}
	return "No category", errors.New("cat no hit")
}
