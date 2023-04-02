package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/voje/whack/internal/sshclient"

	"github.com/urfave/cli/v2"
)

func main() {
    app := &cli.App{
        Name:  "whack",
        Usage: `Whack some illegal processes
        Descripiton: "This application connects to remote machines via ssh 
        and list their running processes. It keeps track of existing processes 
        and alerts on newly spawned processes.`,
        Action: run,
    }

    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}

func run(*cli.Context) error {
    log.Info("Starting Whack")

    sc := sshclient.NewSshClient("/home/kristjan/.ssh/config_vagrant")

    sc.SshConn("wraith")

    return nil
}

