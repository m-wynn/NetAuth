package hooks

import (
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/NetAuth/NetAuth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

func TestMergeEntityMeta(t *testing.T) {
	hook, err := NewMergeEntityMeta(tree.RefContext{})
	if err != nil {
		t.Fatal(err)
	}

	e := &pb.Entity{Meta: &pb.EntityMeta{}}
	de := &pb.Entity{
		Meta: &pb.EntityMeta{
			GECOS:  proto.String("PFY"),
			Groups: []string{"not-to-be-merged"},
		},
	}

	if err := hook.Run(e, de); err != nil {
		t.Fatal(err)
	}

	if len(e.GetMeta().GetGroups()) != 0 || e.GetMeta().GetGECOS() != "PFY" {
		t.Fatal("Spec error - please trace hook")
	}
}
