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
    Host string
    ConfigFile *ssh_config.Config
    ClientConfig *ssh.ClientConfig
}

func  NewSshConfigFile(filePath string) *ssh_config.Config {
    if filePath == "" {
        filePath = filepath.Join(os.Getenv("HOME"), ".ssh", "config")
    }
    file, _ := os.Open(filePath) 
    configFile, _ := ssh_config.Decode(file)

    return configFile
}

func NewSshClient(host string, configFile *ssh_config.Config) *SshClient {
    sc := SshClient{
        Host: host,
        ConfigFile: configFile,
    }
    if _, err := (sc.ConfigFile.Get(host, "HostName")); err != nil {
        log.Error("Host not found in ssh-config: " + host)
    }

    identityFile, _ := sc.ConfigFile.Get(host, "IdentityFile")
    key, err := ioutil.ReadFile(identityFile)
    if err != nil {
        log.Fatal(err)
    }

    signer, err := ssh.ParsePrivateKey(key)
    if err != nil {
        log.Fatal(err)
    }

    user, _ := sc.ConfigFile.Get(host, "User")
    clientConfig := &ssh.ClientConfig{
    	User:            user,
    	Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
    }
    sc.ClientConfig = clientConfig
    return &sc
}

func (sc* SshClient) SendCmd(cmd string) ([]byte, error) {
    hostName, _ := sc.ConfigFile.Get(sc.Host, "HostName")
    port, _ := sc.ConfigFile.Get(sc.Host, "Port")
    network := fmt.Sprintf("%s:%s", hostName, port)
    conn, err := ssh.Dial(
        "tcp",
        network,
        sc.ClientConfig,
    )
    if err != nil {
        return nil, err
    }
    defer conn.Close()

    session, err := conn.NewSession()
    if err != nil {
        return nil, err
    }
    defer session.Close()

    output, err := session.CombinedOutput(cmd)
    return output, err
}

