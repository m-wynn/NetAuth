package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/bgentry/speakeasy"
	"github.com/google/subcommands"
)

// ChangeSecretCmd services the ChangeSecret RPC
type ChangeSecretCmd struct {
	entityID string
	secret   string
}

// Name of this cmdlet is 'change-secret'
func (*ChangeSecretCmd) Name() string { return "change-secret" }

// Synopsis for the cmdlet.
func (*ChangeSecretCmd) Synopsis() string { return "Change the secret for a given entity" }

// Usage in long form for the cmdlet.
func (*ChangeSecretCmd) Usage() string {
	return `change-secret --entity <ID>  --secret <secret>
Change the secret for the listed entity.  If no entity is provided the
entity specified by the top level flag will be used instead.`
}

// SetFlags is the interface function to set flags specific to this cmdlet.
func (p *ChangeSecretCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.entityID, "entity", getEntity(), "ID to change secret")
	f.StringVar(&p.secret, "secret", "", "New secret (omit for prompt)")
}

// Execute is the interface function to run the cmdlet.
func (p *ChangeSecretCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
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

	// Get the secret if it wasn't specified on the line
	if p.secret == "" {
		var err error
		p.secret, err = speakeasy.Ask("New Secret: ")
		if err != nil {
			fmt.Println(err)
			return subcommands.ExitFailure
		}
	}

	// Change the secret
	result, err := c.ChangeSecret(getEntity(), getSecret(), p.entityID, p.secret, t)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	fmt.Println(result.GetMsg())
	return subcommands.ExitSuccess
}
