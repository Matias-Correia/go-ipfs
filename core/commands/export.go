package commands

import (
	bitswap "github.com/ipfs/go-bitswap"
	cmds "github.com/ipfs/go-ipfs-cmds"
	"github.com/ipfs/go-ipfs/core/commands/cmdenv"
)

var ExportCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "Export Logs to Database",
		ShortDescription: "Exports the logs gathered during the execution to the database",
	},

	Arguments: []cmds.Argument{},
	Options:   []cmds.Option{},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		node, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		bs := node.Exchange.(*bitswap.Bitswap)

		bs.Export()
		return nil
	},
}
