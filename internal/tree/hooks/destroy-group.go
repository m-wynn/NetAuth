package hooks

import (
	"github.com/NetAuth/NetAuth/internal/db"
	"github.com/NetAuth/NetAuth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

// DestroyGroup removes an entity from the system.
type DestroyGroup struct {
	tree.BaseHook
	db.DB
}

// Run will request the underlying datastore to remove the group,
// returning any status provided.  If the group Name is not specified
// in g, it will be obtained from dg.
func (d *DestroyGroup) Run(g, dg *pb.Group) error {
	// This hook is somewhat special since it might be called
	// after a processing pipeline, or just to remove a group.
	if g.GetName() == "" {
		g.Name = dg.Name
	}
	return d.DeleteGroup(g.GetName())
}

func init() {
	tree.RegisterGroupHookConstructor("destroy-group", NewDestroyGroup)
}

// NewDestroyGroup returns an initialized DestroyGroup hook for use.
func NewDestroyGroup(c tree.RefContext) (tree.GroupHook, error) {
	return &DestroyGroup{tree.NewBaseHook("destroy-group", 99), c.DB}, nil
}
