package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
)

// ChangeSecretCmd services the ChangeSecret RPC
type ChangeSecretCmd struct {
	ID     string
	secret string
}

// Name of this cmdlet is 'change-secret'
func (*ChangeSecretCmd) Name() string { return "change-secret" }

// Synopsis for the cmdlet.
func (*ChangeSecretCmd) Synopsis() string { return "Change the secret for a given entity" }

// Usage in long form for the cmdlet.
func (*ChangeSecretCmd) Usage() string {
	return `change-secret --ID <ID>  --secret <secret>
Change the secret for the listed entity.  If no entity is provided the
entity specified by the top level flag will be used instead.`
}

// SetFlags is the interface function to set flags specific to this cmdlet.
func (p *ChangeSecretCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.ID, "ID", "", "ID for the new entity")
	f.StringVar(&p.secret, "secret", "", "secret for the new entity")
}

// Execute is the interface function to run the cmdlet.
func (p *ChangeSecretCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// If the entity wasn't provided, use the one that was set
	// earlier.
	if p.ID == "" {
		p.ID = getEntity()
	}

	// Grab a client
	c, err := getClient()
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	// Get the authorization token
	t, err := c.GetToken(getEntity(), getSecret())
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	// Change the secret
	result, err := c.ChangeSecret(getEntity(), getSecret(), p.ID, p.secret, t)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	fmt.Println(result.GetMsg())
	return subcommands.ExitSuccess
}
