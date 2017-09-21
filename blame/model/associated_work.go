package model

import (
	"github.com/heysquirrel/tribe/apis"
	"github.com/heysquirrel/tribe/git"
)

type AssociatedWork struct {
	Context   Context
	WorkItems apis.WorkItems
}

type AssociatedContributors struct {
	Context      Context
	Contributors git.Contributors
}
