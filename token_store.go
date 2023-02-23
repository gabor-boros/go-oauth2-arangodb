package arangoStore

import (
	"context"
	"encoding/json"
	"time"

	arangoDriver "github.com/arangodb/go-driver"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
)

const (
	// DefaultTokenStoreCollection is the default collection for storing tokens.
	DefaultTokenStoreCollection = "oauth2_tokens" // nolint: gosec
)

// TokenStoreOption is a function that configures the TokenStore.
type TokenStoreOption func(*TokenStore) error

// WithTokenStoreCollection configures the collection for the TokenStore.
func WithTokenStoreCollection(collection string) TokenStoreOption {
	return func(s *TokenStore) error {
		if collection == "" {
			return ErrNoCollection
		}

		s.collection = collection

		return nil
	}
}

// WithTokenStoreDatabase configures the database for the TokenStore.
func WithTokenStoreDatabase(db arangoDriver.Database) TokenStoreOption {
	return func(s *TokenStore) error {
		if db == nil {
			return ErrNoDatabase
		}

		s.db = db

		return nil
	}
}

// TokenStoreItem data item
type TokenStoreItem struct {
	Key       string    `db:"_key"`
	Code      string    `db:"code"`
	Access    string    `db:"access_token"`
	Refresh   string    `db:"refresh_token"`
	Data      []byte    `db:"data"`
	CreatedAt time.Time `db:"created_at"`
	ExpiresAt time.Time `db:"expires_at"`
}

// TokenStore is a data struct that stores oauth2 client information.
type TokenStore struct {
	db         arangoDriver.Database
	collection string
}

func (s *TokenStore) getByQuery(ctx context.Context, query string, bindVars map[string]any) (oauth2.TokenInfo, error) {
	cursor, err := s.db.Query(ctx, query, bindVars)
	if err != nil {
		return nil, err
	}
	defer func(cursor arangoDriver.Cursor) {
		if err := cursor.Close(); err != nil {
			panic(err)
		}
	}(cursor)

	var doc TokenStoreItem
	for cursor.HasMore() {
		_, err := cursor.ReadDocument(ctx, &doc)
		if err != nil {
			return nil, err
		}
	}

	var info models.Token
	if err = json.Unmarshal(doc.Data, &info); err != nil {
		return nil, err
	}

	return &info, nil
}

func (s *TokenStore) removeByQuery(ctx context.Context, query string, bindVars map[string]any) error {
	cursor, err := s.db.Query(ctx, query, bindVars)
	if err != nil {
		return err
	}

	return cursor.Close()
}

// Create creates a new client in the store.
func (s *TokenStore) Create(ctx context.Context, info oauth2.TokenInfo) error {
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}

	coll, err := s.db.Collection(ctx, s.collection)
	if err != nil {
		return err
	}

	doc := TokenStoreItem{
		Data:      data,
		CreatedAt: time.Now(),
	}

	if code := info.GetCode(); code != "" {
		doc.Code = code
		doc.ExpiresAt = info.GetCodeCreateAt().Add(info.GetCodeExpiresIn())
	} else {
		if access := info.GetAccess(); access != "" {
			doc.Access = info.GetAccess()
			doc.ExpiresAt = info.GetAccessCreateAt().Add(info.GetAccessExpiresIn())
		}

		if refresh := info.GetRefresh(); refresh != "" {
			doc.Refresh = info.GetRefresh()
			doc.ExpiresAt = info.GetRefreshCreateAt().Add(info.GetRefreshExpiresIn())
		}
	}

	_, err = coll.CreateDocument(ctx, doc)
	if err != nil {
		return err
	}

	return nil
}

// GetByCode returns the token by its authorization code.
func (s *TokenStore) GetByCode(ctx context.Context, code string) (oauth2.TokenInfo, error) {
	query := "FOR doc IN @@collection FILTER doc.code == @code RETURN doc"
	bindVars := map[string]interface{}{
		"@collection": s.collection,
		"code":        code,
	}

	return s.getByQuery(ctx, query, bindVars)
}

// GetByAccess returns the token by its access token.
func (s *TokenStore) GetByAccess(ctx context.Context, access string) (oauth2.TokenInfo, error) {
	query := "FOR doc IN @@collection FILTER doc.access_token == @access_token RETURN doc"
	bindVars := map[string]interface{}{
		"@collection":  s.collection,
		"access_token": access,
	}

	return s.getByQuery(ctx, query, bindVars)
}

// GetByRefresh returns the token by its refresh token.
func (s *TokenStore) GetByRefresh(ctx context.Context, refresh string) (oauth2.TokenInfo, error) {
	query := "FOR doc IN @@collection FILTER doc.refresh_token == @refresh_token RETURN doc"
	bindVars := map[string]interface{}{
		"@collection":   s.collection,
		"refresh_token": refresh,
	}

	return s.getByQuery(ctx, query, bindVars)
}

// RemoveByCode deletes the token by its authorization code.
func (s *TokenStore) RemoveByCode(ctx context.Context, code string) error {
	query := "FOR doc IN @@collection FILTER doc.code == @code REMOVE doc IN @@collection"
	bindVars := map[string]interface{}{
		"@collection": s.collection,
		"code":        code,
	}

	return s.removeByQuery(ctx, query, bindVars)
}

func (s *TokenStore) RemoveByAccess(ctx context.Context, access string) error {
	query := "FOR doc IN @@collection FILTER doc.access_token == @access_token REMOVE doc IN @@collection"
	bindVars := map[string]interface{}{
		"@collection":  s.collection,
		"access_token": access,
	}

	return s.removeByQuery(ctx, query, bindVars)
}

func (s *TokenStore) RemoveByRefresh(ctx context.Context, refresh string) error {
	query := "FOR doc IN @@collection FILTER doc.refresh_token == @refresh_token REMOVE doc IN @@collection"
	bindVars := map[string]interface{}{
		"@collection":   s.collection,
		"refresh_token": refresh,
	}

	return s.removeByQuery(ctx, query, bindVars)
}

// NewTokenStore creates a new TokenStore.
func NewTokenStore(opts ...TokenStoreOption) (*TokenStore, error) {
	s := &TokenStore{
		collection: DefaultTokenStoreCollection,
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
