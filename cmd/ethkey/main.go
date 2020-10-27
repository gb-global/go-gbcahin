package main

import (
	"fmt"
	"os"

	"gbchain-org/go-gbchain/cmd/utils"
	"gopkg.in/urfave/cli.v1"
)

const (
	defaultKeyfileName = "keyfile.json"
)

// Git SHA1 commit hash of the release (set via linker flags)
var gitCommit = ""
var gitDate = ""

var app *cli.App

func init() {
	app = utils.NewApp(gitCommit, gitDate, "an Ethereum key manager")
	app.Commands = []cli.Command{
		commandGenerate,
		commandInspect,
		commandChangePassphrase,
		commandSignMessage,
		commandVerifyMessage,
	}
	cli.CommandHelpTemplate = utils.OriginCommandHelpTemplate
}

// Commonly used command line flags.
var (
	passphraseFlag = cli.StringFlag{
		Name:  "passwordfile",
		Usage: "the file that contains the password for the keyfile",
	}
	jsonFlag = cli.BoolFlag{
		Name:  "json",
		Usage: "output JSON instead of human-readable format",
	}
)

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
