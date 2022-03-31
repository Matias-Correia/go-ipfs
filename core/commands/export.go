package commands

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/ipfs/go-ipfs/core/commands/cmdenv"
	bitswap "github.com/ipfs/go-bitswap"
	"github.com/ipfs/go-ipfs-cmds"
	"github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/path"
)

var ExportCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "Export Logs to Database",
		ShortDescription: "Exports the logs gathered during the execution to the database",
	},

	Arguments: []cmds.Argument{},
	Options: []cmds.Option{},
	Run: func(env cmds.Environment) {
		node, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		bs := node.Exchange.(*bitswap.Bitswap)
		
		bs.Export()

	},
}