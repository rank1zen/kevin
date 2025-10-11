package frontend

import "github.com/rank1zen/kevin/internal"

// Handler provides the API for server operations.
type Handler struct {
	// Datasource handles the business logic. A nil value indicates that
	// Handler will use the zero value [internal.Datasource].
	Datasource *internal.Datasource
}
