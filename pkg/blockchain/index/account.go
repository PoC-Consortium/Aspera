package index

import (
	"fmt"

	"github.com/ac0v/aspera/pkg/account"

	"github.com/blevesearch/bleve"
)

type AccountIndex interface {
	Index(a *account.Account) error
	ById(id uint64) []byte
	ByPublicKey(publicKey []byte) ([]byte, error)
	ByAddress(address string) ([]byte, error)
}

type accountIndex struct {
	index bleve.Index
}

type AccountIndexData struct {
	PublicKey string
	Address   string
}

func (_ *AccountIndexData) Type() string {
	return "account"
}

func NewAccountIndex() AccountIndex {
	return &accountIndex{
		index: openIndex("account"),
	}
}

func (ai *accountIndex) Index(a *account.Account) error {
	return ai.index.Index(fmt.Sprintf("%020d", a.Id), &AccountIndexData{
		Address:   a.Address,
		PublicKey: string(a.PublicKey),
	})
}

func (ai *accountIndex) ById(id uint64) []byte {
	return []byte(fmt.Sprintf("%020d", id))
}

func (ai *accountIndex) ByPublicKey(publicKey []byte) ([]byte, error) {
	query := bleve.NewMatchQuery(string(publicKey))
	query.SetField("PublicKey")
	searchRequest := bleve.NewSearchRequest(query)
	if res, err := ai.index.Search(searchRequest); err == nil {
		return primaryKey(res)
	} else {
		return nil, err
	}
}

func (ai *accountIndex) ByAddress(address string) ([]byte, error) {
	query := bleve.NewMatchQuery(address)
	query.SetField("Address")
	searchRequest := bleve.NewSearchRequest(query)
	if res, err := ai.index.Search(searchRequest); err == nil {
		return primaryKey(res)
	} else {
		return nil, err
	}
}
