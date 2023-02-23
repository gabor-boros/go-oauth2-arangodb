package arangoStore

import (
	"context"
	"encoding/json"

	arangoDriver "github.com/arangodb/go-driver"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
)

const (
	// DefaultClientStoreCollection is the default collection for storing clients.
	DefaultClientStoreCollection = "oauth2_clients"
)

// ClientStoreOption is a function that configures the ClientStore.
type ClientStoreOption func(*ClientStore) error

// WithClientStoreCollection configures the collection for the ClientStore.
func WithClientStoreCollection(collection string) ClientStoreOption {
	return func(s *ClientStore) error {
		if collection == "" {
			return ErrNoCollection
		}

		s.collection = collection

		return nil
	}
}

// WithClientStoreDatabase configures the database for the ClientStore.
func WithClientStoreDatabase(db arangoDriver.Database) ClientStoreOption {
	return func(s *ClientStore) error {
		if db == nil {
			return ErrNoDatabase
		}

		s.db = db

		return nil
	}
}

// ClientStoreItem data item
type ClientStoreItem struct {
	Key    string `json:"_key"`
	Secret string `json:"secret"`
	Domain string `json:"domain"`
	Data   []byte `json:"data"`
}

// ClientStore is a data struct that stores oauth2 client information.
type ClientStore struct {
	db         arangoDriver.Database
	collection string
}

// Create creates a new client in the store.
func (s *ClientStore) Create(info oauth2.ClientInfo) error {
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}

	coll, err := s.db.Collection(context.Background(), s.collection)
	if err != nil {
		return err
	}

	doc := &ClientStoreItem{
		Key:    info.GetID(),
		Secret: info.GetSecret(),
		Domain: info.GetDomain(),
		Data:   data,
	}

	_, err = coll.CreateDocument(context.Background(), doc)
	if err != nil {
		return err
	}

	return nil
}

// GetByID returns the client information by key from the store.
func (s *ClientStore) GetByID(ctx context.Context, key string) (oauth2.ClientInfo, error) {
	coll, err := s.db.Collection(ctx, s.collection)
	if err != nil {
		return nil, err
	}

	var client ClientStoreItem
	_, err = coll.ReadDocument(ctx, key, &client)
	if err != nil {
		return nil, err
	}

	var info models.Client
	err = json.Unmarshal(client.Data, &info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

// NewClientStore creates a new ClientStore.
func NewClientStore(opts ...ClientStoreOption) (*ClientStore, error) {
	s := &ClientStore{
		collection: DefaultClientStoreCollection,
	}

	for _, o := range opts {
		if err := o(s); err != nil {
			return nil, err
		}
	}

	if s.db == nil {
		return nil, ErrNoDatabase
	}

	return s, nil
}
