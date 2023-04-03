package commands

import (
	"crypto/md5"
	"fmt"
	"strings"
)

// h no neader
const PsCmd = "ps h -eo pid,user,lstart,tty,comm,args"

type Proc struct {
    Pid string
    Comm string
    User string
    Lstart string
    Tty string
    Args string
    Hash string
}

// parsing depends on the output of psCmd
func parseProc(s string) *Proc {
    spl := strings.Fields(s)
    return &Proc {
        Hash: fmt.Sprintf("%x", md5.Sum([]byte(s))),
        Pid: spl[0],
        User: spl[1],
        Lstart : strings.Join(spl[2:7], " "),
        Tty: spl[7],
        Comm: spl[8],
        Args: strings.Join(spl[9:], " "),
    }
}

type ProcMap map[string]*Proc

func (pm *ProcMap) ToString() string {
    out := ""
    for key := range(*pm) {
        out += fmt.Sprintf("%+v\n", (*pm)[key])
    }
    return out
}

func NewProcMap(s string) ProcMap {
    procs := make(ProcMap)
    spl := strings.Split(s, "\n")
    for _, line := range(spl) {
        if len(line) == 0 {
            continue
        }
        proc := parseProc(line)
        procs[proc.Hash] = proc
    }
    return procs
}

