package hooks

import (
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/NetAuth/NetAuth/internal/db/memdb"
	"github.com/NetAuth/NetAuth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

func TestSaveGroup(t *testing.T) {
	mdb, err := memdb.New()
	if err != nil {
		t.Fatal(err)
	}

	hook, err := NewSaveGroup(tree.RefContext{DB: mdb})
	if err != nil {
		t.Fatal(err)
	}

	g := &pb.Group{Name: proto.String("fooGroup")}

	if err := hook.Run(g, &pb.Group{}); err != nil {
		t.Fatal(err)
	}

	if _, err := mdb.LoadGroup("fooGroup"); err != nil {
		t.Fatal("Group wasn't saved", err)
	}
}
