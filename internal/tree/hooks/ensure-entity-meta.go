package hooks

import (
	"github.com/NetAuth/NetAuth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

// EnsureEntityMeta has one function: to ensure that the metadata
// struct on an entity is not nil.
type EnsureEntityMeta struct {
	tree.BaseHook
}

// Run will apply an empty metadata struct if one is not already
// present.
func (*EnsureEntityMeta) Run(e, de *pb.Entity) error {
	if e.Meta == nil {
		e.Meta = &pb.EntityMeta{}
	}
	return nil
}

func init() {
	tree.RegisterEntityHookConstructor("ensure-entity-meta", NewEnsureEntityMeta)
}

// NewEnsureEntityMeta returns an initialized hook to the caller.
func NewEnsureEntityMeta(c tree.RefContext) (tree.EntityHook, error) {
	return &EnsureEntityMeta{tree.NewBaseHook("ensure-entity-meta", 20)}, nil
}
