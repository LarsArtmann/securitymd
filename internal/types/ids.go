package types

import id "github.com/larsartmann/go-branded-id"

type (
	IDID        = id.ID[IDBrand, string]
	IDBrand     struct{}
	PolicyID    = id.ID[PolicyBrand, string]
	PolicyBrand struct{}
)
