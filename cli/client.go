package cli

import (
	"github.com/urfave/cli"
)

func Request(ctx *cli.Context) {
	addr := ctx.GlobalString("listen-addr")
	_ = addr
}