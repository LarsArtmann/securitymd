package types

import "github.com/larsartmann/go-composable-business-types/id"

type (
	IDID     = id.ID[IDBrand, string]
	IDBrand  struct{}
	PolicyID = id.ID[PolicyBrand, string]
	PolicyBrand struct{}
)
