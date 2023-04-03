package main

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/voje/whack/internal/hosts"
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
        Flags: []cli.Flag {
            &cli.StringFlag {
                Name: "hosts",
                Value: "localhost,localhost",
                Usage: "Comma-separated list of target hosts",
            },
        },
    }

    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}

func run(ctx *cli.Context) error {
    log.Info("Starting Whack")

    sshConfigFile := sshclient.NewSshConfigFile("/home/kristjan/.ssh/config_vagrant")

    for _, host := range(strings.Split(ctx.String("hosts"), ",")) {
        h := hosts.NewHost(host, sshConfigFile)

        psOutput, err := h.Ps()
        if err != nil {
            log.Error(err)
        }
        log.Infof("[%s] >> %+v\n", h.Host, psOutput.ToString())
    }

    return nil
}

