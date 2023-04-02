package sshclient

import (
	"strings"
)

// h no neader
const psCmd = "ps h -eo pid,user,lstart,tty,comm,args"

type Proc struct {
    Pid string
    User string
    Lstart string
    Tty string
    Comm string
    Args string
}

// parsing depends on the output of psCmd
func parseProc(s string) *Proc {
    spl := strings.Fields(s)
    return &Proc {
        Pid: spl[0],
        User: spl[1],
        Lstart : strings.Join(spl[2:7], " "),
        Tty: spl[7],
        Comm: spl[8],
        Args: strings.Join(spl[9:], " "),
    }
}

type PsOutput struct {
    Procs []*Proc
}

func parsePsOutput(s string) *PsOutput {
    procs := []*Proc{}
    spl := strings.Split(s, "\n")
    for _, line := range(spl) {
        if len(line) == 0 {
            continue
        }
        procs = append(procs, parseProc(line))
    }
    // TODO is this an array copy?
    return &PsOutput{
        Procs: procs,
    }
}

