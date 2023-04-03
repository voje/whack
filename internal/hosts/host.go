package hosts

import (
	"github.com/kevinburke/ssh_config"
	"github.com/voje/whack/internal/commands"
	"github.com/voje/whack/internal/sshclient"
	log "github.com/sirupsen/logrus"
)

type Host struct {
    Host string
    SshClient *sshclient.SshClient
    ProcMap commands.ProcMap
}

func NewHost (host string, configFile *ssh_config.Config) *Host {
    log.Info("Init host: " + host) 
    return &Host {
        Host: host,
        SshClient: sshclient.NewSshClient(host, configFile),
        ProcMap: make(commands.ProcMap),
    }
}

// UpdateProcs updates the list of Host processes and returns 
// a list of newly spawned processes
func (h *Host) UpdateProcs(procs commands.ProcMap) (commands .ProcMap) {
    newProcs := make(commands.ProcMap)
    for hash := range procs {
        if _, ok := h.ProcMap[hash]; !ok {
            newProcs[hash] = procs[hash]
        }
    }
    h.ProcMap = procs
    return newProcs
}

func (h *Host) Ps() (*commands.ProcMap, error) {
    b, err := h.SshClient.SendCmd(commands.PsCmd)
    if err != nil {
        return nil, err
    }
    
    procMap := commands.NewProcMap(string(b))
    if err != nil {
        return nil, err
    }

    return &procMap, nil
}

