package tree

import (
	"errors"
	"testing"

	pb "github.com/NetAuth/Protocol"
)

func resetGroupConstructorMap() {
	gHookConstructors = make(map[string]GroupHookConstructor)
}

func TestGCRegisterAndInitialize(t *testing.T) {
	*debugChains = true
	resetGroupConstructorMap()
	defer resetGroupConstructorMap()

	RegisterGroupHookConstructor("null-hook", goodGroupConstructor)
	RegisterGroupHookConstructor("null-hook", goodGroupConstructor)

	if len(gHookConstructors) != 1 {
		t.Error("Duplicate hook registered")
	}

	RegisterGroupHookConstructor("bad-hook", badGroupConstructor)

	if len(gHookConstructors) != 2 {
		t.Error("bad-hook wasn't registered")
	}

	em := Manager{
		groupHooks: make(map[string]GroupHook),
	}

	em.InitializeGroupHooks()
	if len(em.groupHooks) != 1 {
		t.Error("bad-hook was initialized")
	}
}

func TestGCInitializeChainsOK(t *testing.T) {
	resetGroupConstructorMap()
	defer resetGroupConstructorMap()

	RegisterGroupHookConstructor("null-hook", goodGroupConstructor)
	RegisterGroupHookConstructor("null-hook2", goodGroupConstructor2)
	em := Manager{
		groupHooks:     make(map[string]GroupHook),
		groupProcesses: make(map[string][]GroupHook),
	}
	em.InitializeGroupHooks()

	c := map[string][]string{
		"TEST": []string{"null-hook", "null-hook2"},
	}

	if err := em.InitializeGroupChains(c); err != nil {
		t.Error(err)
	}
}

func TestGCInitializeBadHook(t *testing.T) {
	resetGroupConstructorMap()
	defer resetGroupConstructorMap()

	em := Manager{
		groupHooks:     make(map[string]GroupHook),
		groupProcesses: make(map[string][]GroupHook),
	}
	em.InitializeGroupHooks()

	c := map[string][]string{
		"TEST": []string{"unknown-hook"},
	}

	if err := em.InitializeGroupChains(c); err != ErrUnknownHook {
		t.Error(err)
	}
}

func TestGCCheckRequiredMissing(t *testing.T) {
	resetGroupConstructorMap()
	defer resetGroupConstructorMap()

	em := Manager{
		groupHooks:     make(map[string]GroupHook),
		groupProcesses: make(map[string][]GroupHook),
	}

	if err := em.CheckRequiredGroupChains(); err != ErrUnknownHookChain {
		t.Error("Passed with a required chain missing")
	}
}

func TestGCCheckRequiredEmpty(t *testing.T) {
	resetGroupConstructorMap()
	defer resetGroupConstructorMap()

	em := Manager{
		groupHooks:     make(map[string]GroupHook),
		groupProcesses: make(map[string][]GroupHook),
	}

	// This lets us do this without having hooks loaded, we just
	// register something into all the chains, and then kill one
	// of them at the end.
	for k := range defaultGroupChains {
		em.groupProcesses[k] = []GroupHook{
			&nullGroupHook{},
		}
	}

	em.groupProcesses["CREATE"] = nil

	if err := em.CheckRequiredGroupChains(); err != ErrEmptyHookChain {
		t.Error("Passed with an empty required chain")
	}
}

type nullGroupHook struct{}

func (*nullGroupHook) Name() string             { return "null-hook" }
func (*nullGroupHook) Priority() int            { return 50 }
func (*nullGroupHook) Run(_, _ *pb.Group) error { return nil }
func goodGroupConstructor(_ RefContext) (GroupHook, error) {
	return &nullGroupHook{}, nil
}

type nullGroupHook2 struct{}

func (*nullGroupHook2) Name() string             { return "null-hook2" }
func (*nullGroupHook2) Priority() int            { return 40 }
func (*nullGroupHook2) Run(_, _ *pb.Group) error { return nil }

func goodGroupConstructor2(_ RefContext) (GroupHook, error) {
	return &nullGroupHook2{}, nil
}

func badGroupConstructor(_ RefContext) (GroupHook, error) {
	return nil, errors.New("initialization error")
}

type errorGroupHook struct{}

func (*errorGroupHook) Name() string             { return "error-hook" }
func (*errorGroupHook) Priority() int            { return 50 }
func (*errorGroupHook) Run(_, _ *pb.Group) error { return errors.New("an Error") }
