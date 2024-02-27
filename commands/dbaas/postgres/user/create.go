package user

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func CreateCmd() *core.Command {
	cmd := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{},
	)

	return cmd
}
