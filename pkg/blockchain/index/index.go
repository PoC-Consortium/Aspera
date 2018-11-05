package index

import (
	"os"

	"github.com/blevesearch/bleve"
	"github.com/dgraph-io/badger"
)

func OpenAccountIndex() bleve.Index {
	return openIndex("account")
}

func openIndex(name string) bleve.Index {
	dir := "var/" + name + ".bleve"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		mapping := bleve.NewIndexMapping()
		if index, err := bleve.New(dir, mapping); err == nil {
			return index
		} else {
			panic(err)
		}
	} else {
		if index, err := bleve.Open(dir); err == nil {
			return index
		} else {
			panic(err)
		}
	}
}

func primaryKey(res *bleve.SearchResult) ([]byte, error) {
	switch len(res.Hits) {
	case 1:
		return []byte(res.Hits[0].ID), nil
	case 0:
		return nil, badger.ErrKeyNotFound
	default:
		panic("expected only one hit")
	}
}
