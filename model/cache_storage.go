package model

import (
	"bytes"

	json "github.com/json-iterator/go"

	"github.com/coocood/freecache"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
)

type Cache struct {
	cache *freecache.Cache
}

func (c *Cache) Initialize() error {
	cacheSize := 100 * 1024 * 1024
	c.cache = freecache.NewCache(cacheSize)
	return nil
}

func (c *Cache) Get(s string) ([]byte, error) {
	return c.cache.Get([]byte(s))
}

func (c *Cache) GetAll() ([]byte, error) {
	a := [][]byte{}
	iter := c.cache.NewIterator()
	for {
		if i := iter.Next(); i == nil {
			break
		} else {
			a = append(a, i.Value)
		}

	}
	b := bytes.Join(a, []byte(`,`))

	// insert '[' to the front
	b = append(b, 0)
	copy(b[1:], b[0:])
	b[0] = byte('[')

	// append ']'
	b = append(b, ']')

	return b, nil
}

func (c *Cache) Create(i interface{}) (string, error) {
	if data, err := json.Marshal(i); err != nil {
		return "", err
	} else {
		id := uuid.NewV4().String()
		log.Debug().Msgf("saving %s to cache", id)
		return id, c.cache.Set([]byte(id), data, 300)
	}
}

func (c *Cache) Update(s string, i interface{}) (string, error) {
	if data, err := json.Marshal(i); err != nil {
		return "", err
	} else {
		log.Debug().Msgf("updating %s in cache", s)
		return s, c.cache.Set([]byte(s), data, 300)
	}
}

func (c *Cache) Delete(s string) error {
	panic("implement me")
}
