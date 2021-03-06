package ctl

import (
	"context"
	"flag"
	"fmt"

	pb "github.com/NetAuth/Protocol"

	"github.com/google/subcommands"
)

// ModifyMetaCmd requests the server to modify the EntityMeta section of an entity.
type ModifyMetaCmd struct {
	entityID       string
	GECOS          string
	PrimaryGroup   string
	legalName      string
	displayName    string
	homedir        string
	shell          string
	graphicalShell string
	badgeNumber    string
}

// Name of this cmdlet is 'modify-meta'
func (*ModifyMetaCmd) Name() string { return "modify-meta" }

// Synopsis returns short-form usage information.
func (*ModifyMetaCmd) Synopsis() string { return "Modify meta-data on an entity" }

// Usage returns long-form usage information.
func (*ModifyMetaCmd) Usage() string {
	return `modify-meta --entity <ID> [fields-to-be-modified]
Modify an entity by updating the named fields to the provided values.
`
}

// SetFlags sets the cmdlet specific flags.
func (p *ModifyMetaCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.entityID, "entity", getEntity(), "ID for the entity to modify")
	f.StringVar(&p.PrimaryGroup, "primary-group", "NO_CHANGE", "Primary Group")
	f.StringVar(&p.GECOS, "GECOS", "NO_CHANGE", "Entity GECOS field")
	f.StringVar(&p.legalName, "legalName", "NO_CHANGE", "Legal name associated with the entity")
	f.StringVar(&p.displayName, "displayName", "NO_CHANGE", "Display name associated with the entity")
	f.StringVar(&p.homedir, "homedir", "NO_CHANGE", "Home directory for the entity")
	f.StringVar(&p.shell, "shell", "NO_CHANGE", "User command interpreter to be used by the entity")
	f.StringVar(&p.graphicalShell, "graphicalShell", "NO_CHANGE", "Graphical shell to be used by the entity")
	f.StringVar(&p.badgeNumber, "badgeNumber", "NO_CHANGE", "Badge number for the entity")
}

// Execute runs the cmdlet.
func (p *ModifyMetaCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Grab a client
	c, err := getClient()
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	// Get the authorization token
	t, err := getToken(c, getEntity())
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	meta := &pb.EntityMeta{}

	// This if block is kind of a hack, it is needed to ensure
	// that fields that weren't set to be modified in the command
	// line flags don't get set in the struct and so don't get
	// overwritten on the server.  If this isn't done here then it
	// means that the server only remembers the last thing to
	// change.  If of course you want to literally set a field to
	// "NO_CHANGE" this is somewhat impossible to do with the CLI.
	if p.PrimaryGroup != "NO_CHANGE" {
		meta.PrimaryGroup = &p.PrimaryGroup
	}
	if p.GECOS != "NO_CHANGE" {
		meta.GECOS = &p.GECOS
	}
	if p.legalName != "NO_CHANGE" {
		meta.LegalName = &p.legalName
	}
	if p.displayName != "NO_CHANGE" {
		meta.DisplayName = &p.displayName
	}
	if p.homedir != "NO_CHANGE" {
		meta.Home = &p.homedir
	}
	if p.shell != "NO_CHANGE" {
		meta.Shell = &p.shell
	}
	if p.graphicalShell != "NO_CHANGE" {
		meta.GraphicalShell = &p.graphicalShell
	}
	if p.badgeNumber != "NO_CHANGE" {
		meta.BadgeNumber = &p.badgeNumber
	}

	result, err := c.ModifyEntityMeta(p.entityID, t, meta)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}
	if result.GetMsg() != "" {
		fmt.Println(result.GetMsg())
	}

	return subcommands.ExitSuccess
}
