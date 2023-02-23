package arangoStore

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/arangodb/go-driver"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
)

func TestNewClientStore(t *testing.T) {
	type args struct {
		opts []ClientStoreOption
	}
	tests := []struct {
		name    string
		args    args
		want    *ClientStore
		wantErr bool
	}{
		{
			name: "new client store",
			args: args{
				opts: []ClientStoreOption{
					WithClientStoreDatabase(new(MockArangoDB)),
					WithClientStoreCollection("collection"),
				},
			},
			want: &ClientStore{
				db:         new(MockArangoDB),
				collection: "collection",
			},
		},
		{
			name: "new client store with default collection",
			args: args{
				opts: []ClientStoreOption{
					WithClientStoreDatabase(new(MockArangoDB)),
				},
			},
			want: &ClientStore{
				db:         new(MockArangoDB),
				collection: DefaultClientStoreCollection,
			},
		},
		{
			name: "new client store with no database",
			args: args{
				opts: []ClientStoreOption{
					WithClientStoreCollection("collection"),
				},
			},
			wantErr: true,
		},
		{
			name: "new client store with invalid database",
			args: args{
				opts: []ClientStoreOption{
					WithClientStoreDatabase(nil),
					WithClientStoreCollection("collection"),
				},
			},
			wantErr: true,
		},
		{
			name: "new client store with invalid collection",
			args: args{
				opts: []ClientStoreOption{
					WithClientStoreDatabase(new(MockArangoDB)),
					WithClientStoreCollection(""),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := NewClientStore(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClientStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClientStore() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClientStore_Create(t *testing.T) {
	type fields struct {
		db         func(info oauth2.ClientInfo) driver.Database
		collection string
	}
	type args struct {
		info oauth2.ClientInfo
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "create client",
			fields: fields{
				db: func(info oauth2.ClientInfo) driver.Database {
					data, err := json.Marshal(info)
					if err != nil {
						t.Fatal(err)
					}

					coll := new(MockArangoCollection)
					coll.On("CreateDocument", context.Background(), &ClientStoreItem{
						Key:    info.GetID(),
						Secret: info.GetSecret(),
						Domain: info.GetDomain(),
						Data:   data,
					}).Return(driver.DocumentMeta{}, nil)

					db := new(MockArangoDB)
					db.On("Collection", context.Background(), DefaultClientStoreCollection).Return(coll, nil)

					return db
				},
				collection: DefaultClientStoreCollection,
			},
			args: args{
				info: &models.Client{
					ID:     "client-id",
					Secret: "client-secret",
					Domain: "example.com",
					Public: false,
					UserID: "user-id",
				},
			},
		},
		{
			name: "create client with collection error",
			fields: fields{
				db: func(info oauth2.ClientInfo) driver.Database {
					db := new(MockArangoDB)
					db.On("Collection", context.Background(), DefaultClientStoreCollection).Return(nil, fmt.Errorf("error"))

					return db
				},
				collection: DefaultClientStoreCollection,
			},
			args: args{
				info: &models.Client{
					ID:     "client-id",
					Secret: "client-secret",
					Domain: "example.com",
					Public: false,
					UserID: "user-id",
				},
			},
			wantErr: true,
		},
		{
			name: "create client with create document error",
			fields: fields{
				db: func(info oauth2.ClientInfo) driver.Database {
					data, err := json.Marshal(info)
					if err != nil {
						t.Fatal(err)
					}

					coll := new(MockArangoCollection)
					coll.On("CreateDocument", context.Background(), &ClientStoreItem{
						Key:    info.GetID(),
						Secret: info.GetSecret(),
						Domain: info.GetDomain(),
						Data:   data,
					}).Return(driver.DocumentMeta{}, fmt.Errorf("error"))

					db := new(MockArangoDB)
					db.On("Collection", context.Background(), DefaultClientStoreCollection).Return(coll, nil)

					return db
				},
				collection: DefaultClientStoreCollection,
			},
			args: args{
				info: &models.Client{
					ID:     "client-id",
					Secret: "client-secret",
					Domain: "example.com",
					Public: false,
					UserID: "user-id",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := &ClientStore{
				db:         tt.fields.db(tt.args.info),
				collection: tt.fields.collection,
			}
			if err := s.Create(tt.args.info); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClientStore_GetByID(t *testing.T) {
	type fields struct {
		db         func(ctx context.Context, key string, doc *ClientStoreItem) driver.Database
		collection string
	}
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    oauth2.ClientInfo
		wantErr bool
	}{
		{
			name: "get client by id",
			fields: fields{
				db: func(ctx context.Context, key string, doc *ClientStoreItem) driver.Database {
					var retDoc ClientStoreItem

					coll := new(MockArangoCollection)
					coll.On("ReadDocument", ctx, key, &retDoc).Return(doc, driver.DocumentMeta{}, nil)

					db := new(MockArangoDB)
					db.On("Collection", ctx, DefaultClientStoreCollection).Return(coll, nil)

					return db
				},
				collection: DefaultClientStoreCollection,
			},
			args: args{
				ctx: context.Background(),
				key: "client-id",
			},
			want: &models.Client{
				ID:     "client-id",
				Secret: "client-secret",
				Domain: "example.com",
				Public: false,
				UserID: "user-id",
			},
		},
		{
			name: "get client by id with collection error",
			fields: fields{
				db: func(ctx context.Context, key string, doc *ClientStoreItem) driver.Database {
					db := new(MockArangoDB)
					db.On("Collection", ctx, DefaultClientStoreCollection).Return(nil, fmt.Errorf("error"))

					return db
				},
				collection: DefaultClientStoreCollection,
			},
			args: args{
				ctx: context.Background(),
				key: "client-id",
			},
			wantErr: true,
		},
		{
			name: "get client by id with read document error",
			fields: fields{
				db: func(ctx context.Context, key string, doc *ClientStoreItem) driver.Database {
					coll := new(MockArangoCollection)
					coll.On("ReadDocument", ctx, key, doc).Return(nil, driver.DocumentMeta{}, fmt.Errorf("error"))

					db := new(MockArangoDB)
					db.On("Collection", ctx, DefaultClientStoreCollection).Return(coll, nil)

					return db
				},
				collection: DefaultClientStoreCollection,
			},
			args: args{
				ctx: context.Background(),
				key: "client-id",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			data, err := json.Marshal(tt.want)
			if err != nil {
				t.Fatal(err)
			}

			var storeItem ClientStoreItem
			if tt.want != nil {
				storeItem = ClientStoreItem{
					Key:    tt.want.GetID(),
					Secret: tt.want.GetSecret(),
					Domain: tt.want.GetDomain(),
					Data:   data,
				}
			}

			s := &ClientStore{
				db:         tt.fields.db(tt.args.ctx, tt.args.key, &storeItem),
				collection: tt.fields.collection,
			}
			got, err := s.GetByID(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}
