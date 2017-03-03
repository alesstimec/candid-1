// Copyright 2017 Canonical Ltd.

package admincmd

import (
	"time"

	"github.com/juju/cmd"
	"github.com/juju/gnuflag"
	"github.com/juju/idmclient/params"
	"golang.org/x/net/context"
	"gopkg.in/errgo.v1"
	"gopkg.in/macaroon-bakery.v2-unstable/bakery"
)

type showCommand struct {
	userCommand

	out cmd.Output
}

func newShowCommand() cmd.Command {
	return &showCommand{}
}

var showDoc = `
The show command shows the details for the specified user.

    user-admin show -e bob@example.com
`

func (c *showCommand) Info() *cmd.Info {
	return &cmd.Info{
		Name:    "show",
		Purpose: "show user details",
		Doc:     showDoc,
	}
}

func (c *showCommand) SetFlags(f *gnuflag.FlagSet) {
	c.userCommand.SetFlags(f)

	c.out.AddFlags(f, "smart", cmd.DefaultFormatters)
}

func (c *showCommand) Run(ctxt *cmd.Context) error {
	ctx := context.Background()
	username, err := c.lookupUser(ctxt)
	if err != nil {
		return errgo.Mask(err)
	}
	client, err := c.Client(ctxt)
	if err != nil {
		return errgo.Mask(err)
	}
	u, err := client.User(ctx, &params.UserRequest{
		Username: username,
	})
	if err != nil {
		return errgo.Mask(err)
	}
	user := user{
		Username:      string(u.Username),
		ExternalID:    u.ExternalID,
		Owner:         string(u.Owner),
		Groups:        []string{},
		SSHKeys:       []string{},
		LastLogin:     timeString(u.LastLogin),
		LastDischarge: timeString(u.LastDischarge),
	}
	if u.ExternalID != "" {
		user.Name = &u.FullName
		user.Email = &u.Email
	} else {
		user.PublicKeys = u.PublicKeys
	}
	if len(u.IDPGroups) > 0 {
		user.Groups = u.IDPGroups
	}
	if len(u.SSHKeys) > 0 {
		user.SSHKeys = u.SSHKeys
	}
	return c.out.Write(ctxt, user)
}

func timeString(t *time.Time) string {
	if t == nil || t.IsZero() {
		return "never"
	}
	return t.Format(time.RFC3339)
}

// user represents a user in the system.
type user struct {
	Username      string              `json:"username" yaml:"username"`
	ExternalID    string              `json:"external-id,omitempty" yaml:"external-id,omitempty"`
	Name          *string             `json:"name,omitempty" yaml:"name,omitempty"`
	Email         *string             `json:"email,omitempty" yaml:"email,omitempty"`
	Owner         string              `json:"owner,omitempty" yaml:"owner,omitempty"`
	PublicKeys    []*bakery.PublicKey `json:"public-keys,omitempty" yaml:"public-keys,omitempty"`
	Groups        []string            `json:"groups" yaml:"groups"`
	SSHKeys       []string            `json:"ssh-keys" yaml:"ssh-keys"`
	LastLogin     string              `json:"last-login" yaml:"last-login"`
	LastDischarge string              `json:"last-discharge" yaml:"last-discharge"`
}
