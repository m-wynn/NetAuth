package hooks

import (
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/NetAuth/NetAuth/internal/crypto/nocrypto"
	"github.com/NetAuth/NetAuth/internal/db/memdb"
	"github.com/NetAuth/NetAuth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

func TestLoadExisting(t *testing.T) {
	mdb, err := memdb.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := tree.RefContext{
		DB: mdb,
	}

	hook, err := NewCreateEntityIfMissing(ctx)
	if err != nil {
		t.Fatal(err)
	}

	err = mdb.SaveEntity(&pb.Entity{
		ID:     proto.String("foo"),
		Number: proto.Int32(42),
	})
	if err != nil {
		t.Fatal(err)
	}

	e := &pb.Entity{}
	de := &pb.Entity{
		ID: proto.String("foo"),
	}

	if err := hook.Run(e, de); err != nil {
		t.Fatal(err)
	}

	if e.GetNumber() != 42 {
		t.Fatal("Existing entity not retrieved")
	}

}

func TestCreateNew(t *testing.T) {
	mdb, err := memdb.New()
	if err != nil {
		t.Fatal(err)
	}
	crypto, err := nocrypto.New()
	if err != nil {
		t.Fatal(err)
	}

	ctx := tree.RefContext{
		DB:     mdb,
		Crypto: crypto,
	}

	hook, err := NewCreateEntityIfMissing(ctx)
	if err != nil {
		t.Fatal(err)
	}

	e := &pb.Entity{}
	de := &pb.Entity{
		ID:     proto.String("foo"),
		Secret: proto.String("foo"),
	}

	if err := hook.Run(e, de); err != nil {
		t.Fatal(err)
	}
	if e.GetID() != "foo" {
		t.Fatal("Bad Entity")
	}

	// This check only works because we're using nocrypto which
	// stores plaintext secrets.
	if e.GetSecret() != "foo" {
		t.Fatal("Secret not set")
	}
}
