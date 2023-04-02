package sshclient

type Host struct {
    Procs ProcMap
}

// UpdateProcs updates the list of Host processes and returns 
// a list of newly spawned processes
func (h *Host) UpdateProcs(procs ProcMap) (newProcs ProcMap) {
    for hash, _ := range procs {
        if _, ok := h.Procs[hash]; !ok {
            newProcs[hash] = procs[hash]
        }
    }
    h.Procs = procs
    return newProcs
}

