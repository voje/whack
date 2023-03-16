package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/voje/whack/internal/whack"

	"github.com/urfave/cli/v2"
)

func main() {
    app := &cli.App{
        Name:  "whack",
        Usage: "Whack some illegal shells",
        Action: run,
    }

    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}

func run(*cli.Context) error {
    log.Info("Starting Whack")

    whack.SshConn()

    return nil
}

