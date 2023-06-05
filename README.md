# go-oauth2-arangodb

[![GoDoc](https://godoc.org/github.com/gabor-boros/go-oauth2-arangodb?status.svg)](https://godoc.org/github.com/gabor-boros/go-oauth2-arangodb)
[![Go Report Card](https://goreportcard.com/badge/github.com/gabor-boros/go-oauth2-arangodb)](https://goreportcard.com/report/github.com/gabor-boros/go-oauth2-arangodb)
[![Maintainability](https://api.codeclimate.com/v1/badges/fc29b0acda61b0ec6689/maintainability)](https://codeclimate.com/github/gabor-boros/go-oauth2-arangodb/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/fc29b0acda61b0ec6689/test_coverage)](https://codeclimate.com/github/gabor-boros/go-oauth2-arangodb/test_coverage)

This package is an [ArangoDB] storage implementation for [go-oauth2] using
ArangoDB's official [go-driver].

The package is following semantic versioning and is not tied to the versioning
of [go-oauth2].

[ArangoDB]: https://www.arangodb.com/
[go-oauth2]: https://github.com/go-oauth2/oauth2
[go-driver]: https://github.com/arangodb/go-driver

## Installation

```bash
go get github.com/gabor-boros/go-oauth2-arangodb
```

## Example usage

```go
package main

import (
	"context"
	"os"

	arangoDriver "github.com/arangodb/go-driver"
	arangoHTTP "github.com/arangodb/go-driver/http"

	"github.com/go-oauth2/oauth2/v4/manage"

	arangostore "github.com/gabor-boros/go-oauth2-arangodb"
)

func main() {
	conn, _ := arangoHTTP.NewConnection(arangoHTTP.ConnectionConfig{
		Endpoints: []string{os.Getenv("ARANGO_URL")},
	})

	client, _ := arangoDriver.NewClient(arangoDriver.ClientConfig{
		Connection:     conn,
		Authentication: arangoDriver.BasicAuthentication(os.Getenv("ARANGO_USER"), os.Getenv("ARANGO_PASSWORD")),
	})

	db, _ := client.Database(context.Background(), os.Getenv("ARANGO_DB"))

	clientStore, _ := arangostore.NewClientStore(
		arangostore.WithClientStoreDatabase(db),
		arangostore.WithClientStoreCollection("oauth2_clients"),
	)

	tokenStore, _ := arangostore.NewTokenStore(
		arangostore.WithTokenStoreDatabase(db),
		arangostore.WithTokenStoreCollection("oauth2_tokens"),
	)

	manager := manage.NewDefaultManager()
	manager.MapTokenStorage(tokenStore)
	manager.MapClientStorage(clientStore)

	// ...
}
```

## Contributing

Contributions are welcome! Please open an issue or a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file
for details.
