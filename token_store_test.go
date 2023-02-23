package arangoStore

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/arangodb/go-driver"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/stretchr/testify/mock"
)

func TestNewTokenStore(t *testing.T) {
	type args struct {
		opts []TokenStoreOption
	}
	tests := []struct {
		name    string
		args    args
		want    *TokenStore
		wantErr bool
	}{
		{
			name: "new client store",
			args: args{
				opts: []TokenStoreOption{
					WithTokenStoreDatabase(new(MockArangoDB)),
					WithTokenStoreCollection("collection"),
				},
			},
			want: &TokenStore{
				db:         new(MockArangoDB),
				collection: "collection",
			},
		},
		{
			name: "new client store with default collection",
			args: args{
				opts: []TokenStoreOption{
					WithTokenStoreDatabase(new(MockArangoDB)),
				},
			},
			want: &TokenStore{
				db:         new(MockArangoDB),
				collection: DefaultTokenStoreCollection,
			},
		},
		{
			name: "new client store with no database",
			args: args{
				opts: []TokenStoreOption{
					WithTokenStoreCollection("collection"),
				},
			},
			wantErr: true,
		},
		{
			name: "new client store with invalid database",
			args: args{
				opts: []TokenStoreOption{
					WithTokenStoreDatabase(nil),
					WithTokenStoreCollection("collection"),
				},
			},
			wantErr: true,
		},
		{
			name: "new client store with invalid collection",
			args: args{
				opts: []TokenStoreOption{
					WithTokenStoreDatabase(new(MockArangoDB)),
					WithTokenStoreCollection(""),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := NewTokenStore(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTokenStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTokenStore() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenStore_Create(t *testing.T) {
	type fields struct {
		db         func(ctx context.Context, info oauth2.TokenInfo) driver.Database
		collection string
	}
	type args struct {
		ctx  context.Context
		info oauth2.TokenInfo
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "create token with code",
			fields: fields{
				db: func(ctx context.Context, info oauth2.TokenInfo) driver.Database {
					coll := new(MockArangoCollection)
					coll.On("CreateDocument", context.Background(), mock.Anything).Return(driver.DocumentMeta{}, nil)

					db := new(MockArangoDB)
					db.On("Collection", context.Background(), DefaultTokenStoreCollection).Return(coll, nil)

					return db
				},
				collection: DefaultTokenStoreCollection,
			},
			args: args{
				ctx: context.Background(),
				info: &models.Token{
					ClientID:      "client-id",
					UserID:        "user-id",
					Code:          "test-code",
					CodeCreateAt:  time.Now(),
					CodeExpiresIn: 10 * time.Second,
				},
			},
		},
		{
			name: "create token with access token",
			fields: fields{
				db: func(ctx context.Context, info oauth2.TokenInfo) driver.Database {
					coll := new(MockArangoCollection)
					coll.On("CreateDocument", context.Background(), mock.Anything).Return(driver.DocumentMeta{}, nil)

					db := new(MockArangoDB)
					db.On("Collection", context.Background(), DefaultTokenStoreCollection).Return(coll, nil)

					return db
				},
				collection: DefaultTokenStoreCollection,
			},
			args: args{
				ctx: context.Background(),
				info: &models.Token{
					ClientID:        "client-id",
					UserID:          "user-id",
					Access:          "test-access-token",
					AccessCreateAt:  time.Now(),
					AccessExpiresIn: 10 * time.Second,
				},
			},
		},
		{
			name: "create token with refresh token",
			fields: fields{
				db: func(ctx context.Context, info oauth2.TokenInfo) driver.Database {
					coll := new(MockArangoCollection)
					coll.On("CreateDocument", context.Background(), mock.Anything).Return(driver.DocumentMeta{}, nil)

					db := new(MockArangoDB)
					db.On("Collection", context.Background(), DefaultTokenStoreCollection).Return(coll, nil)

					return db
				},
				collection: DefaultTokenStoreCollection,
			},
			args: args{
				ctx: context.Background(),
				info: &models.Token{
					ClientID:         "client-id",
					UserID:           "user-id",
					Refresh:          "test-refresh-token",
					RefreshCreateAt:  time.Now(),
					RefreshExpiresIn: 10 * time.Second,
				},
			},
		},
		{
			name: "create token with collection error",
			fields: fields{
				db: func(ctx context.Context, info oauth2.TokenInfo) driver.Database {
					db := new(MockArangoDB)
					db.On("Collection", context.Background(), DefaultTokenStoreCollection).Return(nil, fmt.Errorf("error"))

					return db
				},
				collection: DefaultTokenStoreCollection,
			},
			args: args{
				ctx: context.Background(),
				info: &models.Token{
					ClientID:      "client-id",
					UserID:        "user-id",
					Code:          "test-code",
					CodeCreateAt:  time.Now(),
					CodeExpiresIn: 10 * time.Second,
				},
			},
			wantErr: true,
		},
		{
			name: "create token with error",
			fields: fields{
				db: func(ctx context.Context, info oauth2.TokenInfo) driver.Database {
					coll := new(MockArangoCollection)
					coll.On("CreateDocument", context.Background(), mock.Anything).Return(driver.DocumentMeta{}, fmt.Errorf("error"))

					db := new(MockArangoDB)
					db.On("Collection", context.Background(), DefaultTokenStoreCollection).Return(coll, nil)

					return db
				},
				collection: DefaultTokenStoreCollection,
			},
			args: args{
				ctx: context.Background(),
				info: &models.Token{
					ClientID:      "client-id",
					UserID:        "user-id",
					Code:          "test-code",
					CodeCreateAt:  time.Now(),
					CodeExpiresIn: 10 * time.Second,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := &TokenStore{
				db:         tt.fields.db(tt.args.ctx, tt.args.info),
				collection: tt.fields.collection,
			}
			if err := s.Create(tt.args.ctx, tt.args.info); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTokenStore_GetByAccess(t *testing.T) {
	type fields struct {
		db         func(ctx context.Context, access string, info oauth2.TokenInfo) driver.Database
		collection string
	}
	type args struct {
		ctx    context.Context
		access string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    oauth2.TokenInfo
		wantErr bool
	}{
		{
			name: "get token by access token",
			fields: fields{
				db: func(ctx context.Context, access string, info oauth2.TokenInfo) driver.Database {
					query := "FOR doc IN @@collection FILTER doc.access_token == @access_token RETURN doc"
					bindVars := map[string]interface{}{
						"@collection":  DefaultTokenStoreCollection,
						"access_token": access,
					}

					data, err := json.Marshal(info)
					if err != nil {
						t.Fatal(err)
					}

					var doc TokenStoreItem

					cursor := new(MockArangoCursor)
					cursor.On("Close").Return(nil)
					cursor.On("HasMore").Return(true, nil).Once()
					cursor.On("HasMore").Return(false, nil).Once()
					cursor.On("ReadDocument", ctx, &doc).Return(&TokenStoreItem{
						Key:     "test-key",
						Code:    info.GetCode(),
						Access:  info.GetAccess(),
						Refresh: info.GetRefresh(),
						Data:    data,
					}, driver.DocumentMeta{}, nil)

					db := new(MockArangoDB)
					db.On("Query", ctx, query, bindVars).Return(cursor, nil)

					return db
				},
				collection: DefaultTokenStoreCollection,
			},
			args: args{
				ctx:    context.Background(),
				access: "test-access-token",
			},
			want: &models.Token{
				ClientID:        "client-id",
				UserID:          "user-id",
				Access:          "test-access-token",
				AccessCreateAt:  time.Time{},
				AccessExpiresIn: 10 * time.Second,
			},
		},
		{
			name: "get token by access token with query error",
			fields: fields{
				db: func(ctx context.Context, access string, info oauth2.TokenInfo) driver.Database {
					query := "FOR doc IN @@collection FILTER doc.access_token == @access_token RETURN doc"
					bindVars := map[string]interface{}{
						"@collection":  DefaultTokenStoreCollection,
						"access_token": access,
					}

					db := new(MockArangoDB)
					db.On("Query", ctx, query, bindVars).Return(nil, fmt.Errorf("error"))

					return db
				},
				collection: DefaultTokenStoreCollection,
			},
			args: args{
				ctx:    context.Background(),
				access: "test-access-token",
			},
			wantErr: true,
		},
		{
			name: "get token by access token with read document error",
			fields: fields{
				db: func(ctx context.Context, access string, info oauth2.TokenInfo) driver.Database {
					query := "FOR doc IN @@collection FILTER doc.access_token == @access_token RETURN doc"
					bindVars := map[string]interface{}{
						"@collection":  DefaultTokenStoreCollection,
						"access_token": access,
					}

					cursor := new(MockArangoCursor)
					cursor.On("Close").Return(nil)
					cursor.On("HasMore").Return(true, nil).Once()
					cursor.On("HasMore").Return(false, nil).Once()
					cursor.On("ReadDocument", ctx, mock.Anything).Return(nil, driver.DocumentMeta{}, fmt.Errorf("error"))

					db := new(MockArangoDB)
					db.On("Query", ctx, query, bindVars).Return(cursor, nil)

					return db
				},
				collection: DefaultTokenStoreCollection,
			},
			args: args{
				ctx:    context.Background(),
				access: "test-access-token",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := &TokenStore{
				db:         tt.fields.db(tt.args.ctx, tt.args.access, tt.want),
				collection: tt.fields.collection,
			}
			got, err := s.GetByAccess(tt.args.ctx, tt.args.access)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByAccess() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByAccess() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenStore_GetByCode(t *testing.T) {
	type fields struct {
		db         func(ctx context.Context, code string, info oauth2.TokenInfo) driver.Database
		collection string
	}
	type args struct {
		ctx  context.Context
		code string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    oauth2.TokenInfo
		wantErr bool
	}{
		{
			name: "get token by code",
			fields: fields{
				db: func(ctx context.Context, code string, info oauth2.TokenInfo) driver.Database {
					query := "FOR doc IN @@collection FILTER doc.code == @code RETURN doc"
					bindVars := map[string]interface{}{
						"@collection": DefaultTokenStoreCollection,
						"code":        code,
					}

					data, err := json.Marshal(info)
					if err != nil {
						t.Fatal(err)
					}

					var doc TokenStoreItem

					cursor := new(MockArangoCursor)
					cursor.On("Close").Return(nil)
					cursor.On("HasMore").Return(true, nil).Once()
					cursor.On("HasMore").Return(false, nil).Once()
					cursor.On("ReadDocument", ctx, &doc).Return(&TokenStoreItem{
						Key:     "test-key",
						Code:    info.GetCode(),
						Access:  info.GetAccess(),
						Refresh: info.GetRefresh(),
						Data:    data,
					}, driver.DocumentMeta{}, nil)

					db := new(MockArangoDB)
					db.On("Query", ctx, query, bindVars).Return(cursor, nil)

					return db
				},
				collection: DefaultTokenStoreCollection,
			},
			args: args{
				ctx:  context.Background(),
				code: "test-code",
			},
			want: &models.Token{
				Code:          "test-code",
				CodeCreateAt:  time.Time{},
				CodeExpiresIn: 10 * time.Second,
			},
		},
		{
			name: "get token by code with query error",
			fields: fields{
				db: func(ctx context.Context, code string, info oauth2.TokenInfo) driver.Database {
					query := "FOR doc IN @@collection FILTER doc.code == @code RETURN doc"
					bindVars := map[string]interface{}{
						"@collection": DefaultTokenStoreCollection,
						"code":        code,
					}

					db := new(MockArangoDB)
					db.On("Query", ctx, query, bindVars).Return(nil, fmt.Errorf("error"))

					return db
				},
				collection: DefaultTokenStoreCollection,
			},
			args: args{
				ctx:  context.Background(),
				code: "test-code",
			},
			wantErr: true,
		},
		{
			name: "get token by code with cursor error",
			fields: fields{
				db: func(ctx context.Context, code string, info oauth2.TokenInfo) driver.Database {
					query := "FOR doc IN @@collection FILTER doc.code == @code RETURN doc"
					bindVars := map[string]interface{}{
						"@collection": DefaultTokenStoreCollection,
						"code":        code,
					}

					cursor := new(MockArangoCursor)
					cursor.On("Close").Return(nil)
					cursor.On("HasMore").Return(true, nil).Once()
					cursor.On("HasMore").Return(false, nil).Once()
					cursor.On("ReadDocument", ctx, mock.Anything).Return(nil, driver.DocumentMeta{}, fmt.Errorf("error"))

					db := new(MockArangoDB)
					db.On("Query", ctx, query, bindVars).Return(cursor, nil)

					return db
				},
				collection: DefaultTokenStoreCollection,
			},
			args: args{
				ctx:  context.Background(),
				code: "test-code",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := &TokenStore{
				db:         tt.fields.db(tt.args.ctx, tt.args.code, tt.want),
				collection: tt.fields.collection,
			}
			got, err := s.GetByCode(tt.args.ctx, tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByCode() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenStore_GetByRefresh(t *testing.T) {
	type fields struct {
		db         func(ctx context.Context, refresh string, want oauth2.TokenInfo) driver.Database
		collection string
	}
	type args struct {
		ctx     context.Context
		refresh string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    oauth2.TokenInfo
		wantErr bool
	}{
		{
			name: "get token by refresh token",
			fields: fields{
				db: func(ctx context.Context, refresh string, info oauth2.TokenInfo) driver.Database {
					query := "FOR doc IN @@collection FILTER doc.refresh_token == @refresh_token RETURN doc"
					bindVars := map[string]interface{}{
						"@collection":   DefaultTokenStoreCollection,
						"refresh_token": refresh,
					}

					data, err := json.Marshal(info)
					if err != nil {
						t.Fatal(err)
					}

					var doc TokenStoreItem

					cursor := new(MockArangoCursor)
					cursor.On("Close").Return(nil)
					cursor.On("HasMore").Return(true, nil).Once()
					cursor.On("HasMore").Return(false, nil).Once()
					cursor.On("ReadDocument", ctx, &doc).Return(&TokenStoreItem{
						Key:     "test-key",
						Code:    info.GetCode(),
						Refresh: info.GetRefresh(),
						Data:    data,
					}, driver.DocumentMeta{}, nil)

					db := new(MockArangoDB)
					db.On("Query", ctx, query, bindVars).Return(cursor, nil)

					return db
				},
				collection: DefaultTokenStoreCollection,
			},
			args: args{
				ctx:     context.Background(),
				refresh: "test-refresh-token",
			},
			want: &models.Token{
				ClientID:         "client-id",
				UserID:           "user-id",
				Refresh:          "test-refresh-token",
				RefreshCreateAt:  time.Time{},
				RefreshExpiresIn: 10 * time.Second,
			},
		},
		{
			name: "get token by refresh token with query error",
			fields: fields{
				db: func(ctx context.Context, refresh string, info oauth2.TokenInfo) driver.Database {
					query := "FOR doc IN @@collection FILTER doc.refresh_token == @refresh_token RETURN doc"
					bindVars := map[string]interface{}{
						"@collection":   DefaultTokenStoreCollection,
						"refresh_token": refresh,
					}

					db := new(MockArangoDB)
					db.On("Query", ctx, query, bindVars).Return(nil, fmt.Errorf("error"))

					return db
				},
				collection: DefaultTokenStoreCollection,
			},
			args: args{
				ctx:     context.Background(),
				refresh: "test-refresh-token",
			},
			wantErr: true,
		},
		{
			name: "get token by refresh token with read document error",
			fields: fields{
				db: func(ctx context.Context, refresh string, info oauth2.TokenInfo) driver.Database {
					query := "FOR doc IN @@collection FILTER doc.refresh_token == @refresh_token RETURN doc"
					bindVars := map[string]interface{}{
						"@collection":   DefaultTokenStoreCollection,
						"refresh_token": refresh,
					}

					cursor := new(MockArangoCursor)
					cursor.On("Close").Return(nil)
					cursor.On("HasMore").Return(true, nil).Once()
					cursor.On("HasMore").Return(false, nil).Once()
					cursor.On("ReadDocument", ctx, mock.Anything).Return(nil, driver.DocumentMeta{}, fmt.Errorf("error"))

					db := new(MockArangoDB)
					db.On("Query", ctx, query, bindVars).Return(cursor, nil)

					return db
				},
				collection: DefaultTokenStoreCollection,
			},
			args: args{
				ctx:     context.Background(),
				refresh: "test-refresh-token",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := &TokenStore{
				db:         tt.fields.db(context.Background(), tt.args.refresh, tt.want),
				collection: tt.fields.collection,
			}
			got, err := s.GetByRefresh(tt.args.ctx, tt.args.refresh)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByRefresh() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByRefresh() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenStore_RemoveByAccess(t *testing.T) {
	type fields struct {
		db         func(ctx context.Context, access string) driver.Database
		collection string
	}
	type args struct {
		ctx    context.Context
		access string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "remove token by access token",
			fields: fields{
				db: func(ctx context.Context, access string) driver.Database {
					query := "FOR doc IN @@collection FILTER doc.access_token == @access_token REMOVE doc IN @@collection"
					bindVars := map[string]interface{}{
						"@collection":  DefaultTokenStoreCollection,
						"access_token": access,
					}

					cursor := new(MockArangoCursor)
					cursor.On("Close").Return(nil)

					db := new(MockArangoDB)
					db.On("Query", ctx, query, bindVars).Return(cursor, nil)

					return db
				},
				collection: DefaultTokenStoreCollection,
			},
			args: args{
				ctx:    context.Background(),
				access: "test-access-token",
			},
		},
		{
			name: "remove token by access token with query error",
			fields: fields{
				db: func(ctx context.Context, access string) driver.Database {
					query := "FOR doc IN @@collection FILTER doc.access_token == @access_token REMOVE doc IN @@collection"
					bindVars := map[string]interface{}{
						"@collection":  DefaultTokenStoreCollection,
						"access_token": access,
					}

					db := new(MockArangoDB)
					db.On("Query", ctx, query, bindVars).Return(nil, fmt.Errorf("error"))

					return db
				},
				collection: DefaultTokenStoreCollection,
			},
			args: args{
				ctx:    context.Background(),
				access: DefaultTokenStoreCollection,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := &TokenStore{
				db:         tt.fields.db(context.Background(), tt.args.access),
				collection: tt.fields.collection,
			}
			if err := s.RemoveByAccess(tt.args.ctx, tt.args.access); (err != nil) != tt.wantErr {
				t.Errorf("RemoveByAccess() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTokenStore_RemoveByCode(t *testing.T) {
	type fields struct {
		db         func(ctx context.Context, code string) driver.Database
		collection string
	}
	type args struct {
		ctx  context.Context
		code string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "remove token by code",
			fields: fields{
				db: func(ctx context.Context, code string) driver.Database {
					query := "FOR doc IN @@collection FILTER doc.code == @code REMOVE doc IN @@collection"
					bindVars := map[string]interface{}{
						"@collection": DefaultTokenStoreCollection,
						"code":        code,
					}

					cursor := new(MockArangoCursor)
					cursor.On("Close").Return(nil)

					db := new(MockArangoDB)
					db.On("Query", ctx, query, bindVars).Return(cursor, nil)

					return db
				},
				collection: DefaultTokenStoreCollection,
			},
			args: args{
				ctx:  context.Background(),
				code: "test-code",
			},
		},
		{
			name: "remove token by code with query error",
			fields: fields{
				db: func(ctx context.Context, code string) driver.Database {
					query := "FOR doc IN @@collection FILTER doc.code == @code REMOVE doc IN @@collection"
					bindVars := map[string]interface{}{
						"@collection": DefaultTokenStoreCollection,
						"code":        code,
					}

					db := new(MockArangoDB)
					db.On("Query", ctx, query, bindVars).Return(nil, fmt.Errorf("error"))

					return db
				},
				collection: DefaultTokenStoreCollection,
			},
			args: args{
				ctx:  context.Background(),
				code: DefaultTokenStoreCollection,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := &TokenStore{
				db:         tt.fields.db(context.Background(), tt.args.code),
				collection: tt.fields.collection,
			}
			if err := s.RemoveByCode(tt.args.ctx, tt.args.code); (err != nil) != tt.wantErr {
				t.Errorf("RemoveByCode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTokenStore_RemoveByRefresh(t *testing.T) {
	type fields struct {
		db         func(ctx context.Context, refresh string) driver.Database
		collection string
	}
	type args struct {
		ctx     context.Context
		refresh string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "remove token by refresh token",
			fields: fields{
				db: func(ctx context.Context, refresh string) driver.Database {
					query := "FOR doc IN @@collection FILTER doc.refresh_token == @refresh_token REMOVE doc IN @@collection"
					bindVars := map[string]interface{}{
						"@collection":   DefaultTokenStoreCollection,
						"refresh_token": refresh,
					}

					cursor := new(MockArangoCursor)
					cursor.On("Close").Return(nil)

					db := new(MockArangoDB)
					db.On("Query", ctx, query, bindVars).Return(cursor, nil)

					return db
				},
				collection: DefaultTokenStoreCollection,
			},
			args: args{
				ctx:     context.Background(),
				refresh: "test-refresh-token",
			},
		},
		{
			name: "remove token by refresh token with query error",
			fields: fields{
				db: func(ctx context.Context, refresh string) driver.Database {
					query := "FOR doc IN @@collection FILTER doc.refresh_token == @refresh_token REMOVE doc IN @@collection"
					bindVars := map[string]interface{}{
						"@collection":   DefaultTokenStoreCollection,
						"refresh_token": refresh,
					}

					db := new(MockArangoDB)
					db.On("Query", ctx, query, bindVars).Return(nil, fmt.Errorf("error"))

					return db
				},
				collection: DefaultTokenStoreCollection,
			},
			args: args{
				ctx:     context.Background(),
				refresh: DefaultTokenStoreCollection,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := &TokenStore{
				db:         tt.fields.db(context.Background(), tt.args.refresh),
				collection: tt.fields.collection,
			}
			if err := s.RemoveByRefresh(tt.args.ctx, tt.args.refresh); (err != nil) != tt.wantErr {
				t.Errorf("RemoveByRefresh() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
