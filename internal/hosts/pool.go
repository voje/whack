package hosts

import "time"

type Pool struct {
    Hosts   map[string]*Host

    // Scan frequency in seconds
    Freq    int
}

func NewPool(freq int) *Pool {
    return &Pool {
        Freq: freq,
        Hosts: make(map[string]*Host),
    }
}

func (p *Pool) AddHost(h *Host) {
    p.Hosts[h.Host] = h
}

func (p *Pool) Scan() {

    for range time.Tick(time.Second * time.Duration(p.Freq)) {
        for h := range p.Hosts {
            p.Hosts[h].Scan() 
        }
    }

}

