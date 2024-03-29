package arangostore

import (
	"context"
	"encoding/json"

	arangoDriver "github.com/arangodb/go-driver"
	"github.com/stretchr/testify/mock"
)

func convertToAny(i any, o any) {
	if o == nil {
		return
	}

	var err error

	b, err := json.Marshal(&i)
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(b, o); err != nil {
		panic(err)
	}
}

type MockArangoDB struct {
	mock.Mock
}

func (m *MockArangoDB) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockArangoDB) Info(ctx context.Context) (arangoDriver.DatabaseInfo, error) {
	args := m.Called(ctx)
	return args.Get(0).(arangoDriver.DatabaseInfo), args.Error(1)
}

func (m *MockArangoDB) EngineInfo(ctx context.Context) (arangoDriver.EngineInfo, error) {
	args := m.Called(ctx)
	return args.Get(0).(arangoDriver.EngineInfo), args.Error(1)
}

func (m *MockArangoDB) Remove(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockArangoDB) Collection(ctx context.Context, name string) (arangoDriver.Collection, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(arangoDriver.Collection), args.Error(1)
}

func (m *MockArangoDB) CollectionExists(ctx context.Context, name string) (bool, error) {
	args := m.Called(ctx, name)
	return args.Bool(0), args.Error(1)
}

func (m *MockArangoDB) Collections(ctx context.Context) ([]arangoDriver.Collection, error) {
	args := m.Called(ctx)
	return args.Get(0).([]arangoDriver.Collection), args.Error(1)
}

func (m *MockArangoDB) CreateCollection(ctx context.Context, name string, options *arangoDriver.CreateCollectionOptions) (arangoDriver.Collection, error) {
	args := m.Called(ctx, name, options)
	return args.Get(0).(arangoDriver.Collection), args.Error(1)
}

func (m *MockArangoDB) View(ctx context.Context, name string) (arangoDriver.View, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(arangoDriver.View), args.Error(1)
}

func (m *MockArangoDB) ViewExists(ctx context.Context, name string) (bool, error) {
	args := m.Called(ctx, name)
	return args.Bool(0), args.Error(1)
}

func (m *MockArangoDB) Views(ctx context.Context) ([]arangoDriver.View, error) {
	args := m.Called(ctx)
	return args.Get(0).([]arangoDriver.View), args.Error(1)
}

func (m *MockArangoDB) CreateArangoSearchView(ctx context.Context, name string, options *arangoDriver.ArangoSearchViewProperties) (arangoDriver.ArangoSearchView, error) {
	args := m.Called(ctx, name, options)
	return args.Get(0).(arangoDriver.ArangoSearchView), args.Error(1)
}

func (m *MockArangoDB) CreateArangoSearchAliasView(ctx context.Context, name string, options *arangoDriver.ArangoSearchAliasViewProperties) (arangoDriver.ArangoSearchViewAlias, error) {
	args := m.Called(ctx, name, options)
	return args.Get(0).(arangoDriver.ArangoSearchViewAlias), args.Error(1)
}

func (m *MockArangoDB) Graph(ctx context.Context, name string) (arangoDriver.Graph, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(arangoDriver.Graph), args.Error(1)
}

func (m *MockArangoDB) GraphExists(ctx context.Context, name string) (bool, error) {
	args := m.Called(ctx, name)
	return args.Bool(0), args.Error(1)
}

func (m *MockArangoDB) Graphs(ctx context.Context) ([]arangoDriver.Graph, error) {
	args := m.Called(ctx)
	return args.Get(0).([]arangoDriver.Graph), args.Error(1)
}

func (m *MockArangoDB) CreateGraph(ctx context.Context, name string, options *arangoDriver.CreateGraphOptions) (arangoDriver.Graph, error) {
	args := m.Called(ctx, name, options)
	return args.Get(0).(arangoDriver.Graph), args.Error(1)
}

func (m *MockArangoDB) CreateGraphV2(ctx context.Context, name string, options *arangoDriver.CreateGraphOptions) (arangoDriver.Graph, error) {
	args := m.Called(ctx, name, options)
	return args.Get(0).(arangoDriver.Graph), args.Error(1)
}

func (m *MockArangoDB) StartJob(ctx context.Context, options arangoDriver.PregelJobOptions) (string, error) {
	args := m.Called(ctx, options)
	return args.String(0), args.Error(1)
}

func (m *MockArangoDB) GetJob(ctx context.Context, id string) (*arangoDriver.PregelJob, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*arangoDriver.PregelJob), args.Error(1)
}

func (m *MockArangoDB) GetJobs(ctx context.Context) ([]*arangoDriver.PregelJob, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*arangoDriver.PregelJob), args.Error(1)
}

func (m *MockArangoDB) CancelJob(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockArangoDB) BeginTransaction(ctx context.Context, cols arangoDriver.TransactionCollections, opts *arangoDriver.BeginTransactionOptions) (arangoDriver.TransactionID, error) {
	args := m.Called(ctx, cols, opts)
	return args.Get(0).(arangoDriver.TransactionID), args.Error(1)
}

func (m *MockArangoDB) CommitTransaction(ctx context.Context, tid arangoDriver.TransactionID, opts *arangoDriver.CommitTransactionOptions) error {
	args := m.Called(ctx, tid, opts)
	return args.Error(0)
}

func (m *MockArangoDB) AbortTransaction(ctx context.Context, tid arangoDriver.TransactionID, opts *arangoDriver.AbortTransactionOptions) error {
	args := m.Called(ctx, tid, opts)
	return args.Error(0)
}

func (m *MockArangoDB) TransactionStatus(ctx context.Context, tid arangoDriver.TransactionID) (arangoDriver.TransactionStatusRecord, error) {
	args := m.Called(ctx, tid)
	return args.Get(0).(arangoDriver.TransactionStatusRecord), args.Error(1)
}

func (m *MockArangoDB) EnsureAnalyzer(ctx context.Context, analyzer arangoDriver.ArangoSearchAnalyzerDefinition) (bool, arangoDriver.ArangoSearchAnalyzer, error) {
	args := m.Called(ctx, analyzer)
	return args.Bool(0), args.Get(1).(arangoDriver.ArangoSearchAnalyzer), args.Error(2)
}

func (m *MockArangoDB) Analyzer(ctx context.Context, name string) (arangoDriver.ArangoSearchAnalyzer, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(arangoDriver.ArangoSearchAnalyzer), args.Error(1)
}

func (m *MockArangoDB) Analyzers(ctx context.Context) ([]arangoDriver.ArangoSearchAnalyzer, error) {
	args := m.Called(ctx)
	return args.Get(0).([]arangoDriver.ArangoSearchAnalyzer), args.Error(1)
}

func (m *MockArangoDB) Query(ctx context.Context, query string, bindVars map[string]any) (arangoDriver.Cursor, error) {
	args := m.Called(ctx, query, bindVars)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(arangoDriver.Cursor), args.Error(1)
}

func (m *MockArangoDB) ValidateQuery(ctx context.Context, query string) error {
	args := m.Called(ctx, query)
	return args.Error(0)
}

func (m *MockArangoDB) OptimizerRulesForQueries(ctx context.Context) ([]arangoDriver.QueryRule, error) {
	args := m.Called(ctx)
	return args.Get(0).([]arangoDriver.QueryRule), args.Error(1)
}

func (m *MockArangoDB) Transaction(ctx context.Context, action string, options *arangoDriver.TransactionOptions) (any, error) {
	args := m.Called(ctx, action, options)
	return args.Get(0), args.Error(1)
}

type MockArangoCursor struct {
	mock.Mock
}

func (m *MockArangoCursor) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockArangoCursor) HasMore() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockArangoCursor) ReadDocument(ctx context.Context, result any) (arangoDriver.DocumentMeta, error) {
	args := m.Called(ctx, result)
	convertToAny(args.Get(0), result)
	return args.Get(1).(arangoDriver.DocumentMeta), args.Error(2)
}

func (m *MockArangoCursor) Count() int64 {
	args := m.Called()
	return args.Get(0).(int64)
}

func (m *MockArangoCursor) Statistics() arangoDriver.QueryStatistics {
	args := m.Called()
	return args.Get(0).(arangoDriver.QueryStatistics)
}

func (m *MockArangoCursor) Extra() arangoDriver.QueryExtra {
	args := m.Called()
	return args.Get(0).(arangoDriver.QueryExtra)
}

type MockArangoCollection struct {
	mock.Mock
}

func (m *MockArangoCollection) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockArangoCollection) Checksum(ctx context.Context, withRevisions bool, withData bool) (arangoDriver.CollectionChecksum, error) {
	args := m.Called(ctx, withRevisions, withData)
	return args.Get(0).(arangoDriver.CollectionChecksum), args.Error(1)
}

func (m *MockArangoCollection) Database() arangoDriver.Database {
	args := m.Called()
	return args.Get(0).(arangoDriver.Database)
}

func (m *MockArangoCollection) Status(ctx context.Context) (arangoDriver.CollectionStatus, error) {
	args := m.Called(ctx)
	return args.Get(0).(arangoDriver.CollectionStatus), args.Error(1)
}

func (m *MockArangoCollection) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockArangoCollection) Statistics(ctx context.Context) (arangoDriver.CollectionStatistics, error) {
	args := m.Called(ctx)
	return args.Get(0).(arangoDriver.CollectionStatistics), args.Error(1)
}

func (m *MockArangoCollection) Revision(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m *MockArangoCollection) Properties(ctx context.Context) (arangoDriver.CollectionProperties, error) {
	args := m.Called(ctx)
	return args.Get(0).(arangoDriver.CollectionProperties), args.Error(1)
}

func (m *MockArangoCollection) SetProperties(ctx context.Context, options arangoDriver.SetCollectionPropertiesOptions) error {
	args := m.Called(ctx, options)
	return args.Error(0)
}

func (m *MockArangoCollection) Shards(ctx context.Context, details bool) (arangoDriver.CollectionShards, error) {
	args := m.Called(ctx, details)
	return args.Get(0).(arangoDriver.CollectionShards), args.Error(1)
}

func (m *MockArangoCollection) Load(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockArangoCollection) Unload(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockArangoCollection) Remove(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockArangoCollection) Truncate(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockArangoCollection) Index(ctx context.Context, name string) (arangoDriver.Index, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(arangoDriver.Index), args.Error(1)
}

func (m *MockArangoCollection) IndexExists(ctx context.Context, name string) (bool, error) {
	args := m.Called(ctx, name)
	return args.Bool(0), args.Error(1)
}

func (m *MockArangoCollection) Indexes(ctx context.Context) ([]arangoDriver.Index, error) {
	args := m.Called(ctx)
	return args.Get(0).([]arangoDriver.Index), args.Error(1)
}

func (m *MockArangoCollection) EnsureFullTextIndex(ctx context.Context, fields []string, options *arangoDriver.EnsureFullTextIndexOptions) (arangoDriver.Index, bool, error) { //nolint:staticcheck
	args := m.Called(ctx, fields, options)
	return args.Get(0).(arangoDriver.Index), args.Bool(1), args.Error(2)
}

func (m *MockArangoCollection) EnsureGeoIndex(ctx context.Context, fields []string, options *arangoDriver.EnsureGeoIndexOptions) (arangoDriver.Index, bool, error) {
	args := m.Called(ctx, fields, options)
	return args.Get(0).(arangoDriver.Index), args.Bool(1), args.Error(2)
}

func (m *MockArangoCollection) EnsureHashIndex(ctx context.Context, fields []string, options *arangoDriver.EnsureHashIndexOptions) (arangoDriver.Index, bool, error) {
	args := m.Called(ctx, fields, options)
	return args.Get(0).(arangoDriver.Index), args.Bool(1), args.Error(2)
}

func (m *MockArangoCollection) EnsurePersistentIndex(ctx context.Context, fields []string, options *arangoDriver.EnsurePersistentIndexOptions) (arangoDriver.Index, bool, error) {
	args := m.Called(ctx, fields, options)
	return args.Get(0).(arangoDriver.Index), args.Bool(1), args.Error(2)
}

func (m *MockArangoCollection) EnsureSkipListIndex(ctx context.Context, fields []string, options *arangoDriver.EnsureSkipListIndexOptions) (arangoDriver.Index, bool, error) {
	args := m.Called(ctx, fields, options)
	return args.Get(0).(arangoDriver.Index), args.Bool(1), args.Error(2)
}

func (m *MockArangoCollection) EnsureTTLIndex(ctx context.Context, field string, expireAfter int, options *arangoDriver.EnsureTTLIndexOptions) (arangoDriver.Index, bool, error) {
	args := m.Called(ctx, field, expireAfter, options)
	return args.Get(0).(arangoDriver.Index), args.Bool(1), args.Error(2)
}

func (m *MockArangoCollection) EnsureZKDIndex(ctx context.Context, fields []string, options *arangoDriver.EnsureZKDIndexOptions) (arangoDriver.Index, bool, error) {
	args := m.Called(ctx, fields, options)
	return args.Get(0).(arangoDriver.Index), args.Bool(1), args.Error(2)
}

func (m *MockArangoCollection) EnsureInvertedIndex(ctx context.Context, options *arangoDriver.InvertedIndexOptions) (arangoDriver.Index, bool, error) {
	args := m.Called(ctx, options)
	return args.Get(0).(arangoDriver.Index), args.Bool(1), args.Error(2)
}

func (m *MockArangoCollection) DocumentExists(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

func (m *MockArangoCollection) ReadDocument(ctx context.Context, key string, result any) (arangoDriver.DocumentMeta, error) {
	args := m.Called(ctx, key, result)
	convertToAny(args.Get(0), &result)
	return args.Get(1).(arangoDriver.DocumentMeta), args.Error(2)
}

func (m *MockArangoCollection) ReadDocuments(ctx context.Context, keys []string, results any) (arangoDriver.DocumentMetaSlice, arangoDriver.ErrorSlice, error) {
	args := m.Called(ctx, keys, results)
	convertToAny(args.Get(0), &results)
	return args.Get(1).(arangoDriver.DocumentMetaSlice), args.Get(2).(arangoDriver.ErrorSlice), args.Error(3)
}

func (m *MockArangoCollection) CreateDocument(ctx context.Context, document any) (arangoDriver.DocumentMeta, error) {
	args := m.Called(ctx, document)
	return args.Get(0).(arangoDriver.DocumentMeta), args.Error(1)
}

func (m *MockArangoCollection) CreateDocuments(ctx context.Context, documents any) (arangoDriver.DocumentMetaSlice, arangoDriver.ErrorSlice, error) {
	args := m.Called(ctx, documents)
	return args.Get(0).(arangoDriver.DocumentMetaSlice), args.Get(1).(arangoDriver.ErrorSlice), args.Error(2)
}

func (m *MockArangoCollection) UpdateDocument(ctx context.Context, key string, update any) (arangoDriver.DocumentMeta, error) {
	args := m.Called(ctx, key, update)
	return args.Get(0).(arangoDriver.DocumentMeta), args.Error(1)
}

func (m *MockArangoCollection) UpdateDocuments(ctx context.Context, keys []string, updates any) (arangoDriver.DocumentMetaSlice, arangoDriver.ErrorSlice, error) {
	args := m.Called(ctx, keys, updates)
	return args.Get(0).(arangoDriver.DocumentMetaSlice), args.Get(1).(arangoDriver.ErrorSlice), args.Error(2)
}

func (m *MockArangoCollection) ReplaceDocument(ctx context.Context, key string, document any) (arangoDriver.DocumentMeta, error) {
	args := m.Called(ctx, key, document)
	return args.Get(0).(arangoDriver.DocumentMeta), args.Error(1)
}

func (m *MockArangoCollection) ReplaceDocuments(ctx context.Context, keys []string, documents any) (arangoDriver.DocumentMetaSlice, arangoDriver.ErrorSlice, error) {
	args := m.Called(ctx, keys, documents)
	return args.Get(0).(arangoDriver.DocumentMetaSlice), args.Get(1).(arangoDriver.ErrorSlice), args.Error(2)
}

func (m *MockArangoCollection) RemoveDocument(ctx context.Context, key string) (arangoDriver.DocumentMeta, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(arangoDriver.DocumentMeta), args.Error(1)
}

func (m *MockArangoCollection) RemoveDocuments(ctx context.Context, keys []string) (arangoDriver.DocumentMetaSlice, arangoDriver.ErrorSlice, error) {
	args := m.Called(ctx, keys)
	return args.Get(0).(arangoDriver.DocumentMetaSlice), args.Get(1).(arangoDriver.ErrorSlice), args.Error(2)
}

func (m *MockArangoCollection) ImportDocuments(ctx context.Context, documents any, options *arangoDriver.ImportDocumentOptions) (arangoDriver.ImportDocumentStatistics, error) {
	args := m.Called(ctx, documents, options)
	return args.Get(0).(arangoDriver.ImportDocumentStatistics), args.Error(1)
}
