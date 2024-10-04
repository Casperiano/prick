package interfaces

import (
	prick "prick/internal/prick"
	"prick/internal/prick/common"
)

type Prickable interface {
	Poke(*prick.Api) error
	Patch(*prick.Api) error
	GetName() string
	GetLocation() string
	GetType() common.ResourceType
	ListPokes(*prick.Api) ([]*common.Poke, error)
}

type ListOptions struct{ ResourceGroup string }
