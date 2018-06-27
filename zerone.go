package zoutil

import (
	"sync"
	"time"

	"github.com/ironzhang/zerone"
	"github.com/ironzhang/zerone/zclient"
	"github.com/ironzhang/zerone/zserver"
)

type Zerone struct {
	Node          string
	Timeout       time.Duration
	ClientVerbose int
	ServerVerbose int

	zerone  zerone.Zerone
	clients sync.Map
	servers sync.Map
}

func Open(node string, opts Options) (*Zerone, error) {
	zopts, err := opts.zerone()
	if err != nil {
		return nil, err
	}
	z, err := zerone.NewZerone(zopts)
	if err != nil {
		return nil, err
	}
	return &Zerone{
		Node:          node,
		Timeout:       time.Duration(opts.Timeout),
		ClientVerbose: opts.ClientVerbose,
		ServerVerbose: opts.ServerVerbose,
		zerone:        z,
	}, nil
}

func (z *Zerone) Close() error {
	z.clients.Range(func(key, value interface{}) bool {
		c := value.(*zclient.Client)
		c.Close()
		return true
	})

	z.servers.Range(func(key, value interface{}) bool {
		s := value.(*zserver.Server)
		s.Close()
		return true
	})

	return z.zerone.Close()
}

func (z *Zerone) NewClient(service string) (*zclient.Client, error) {
	if v, ok := z.clients.Load(service); ok {
		return v.(*zclient.Client), nil
	}

	c, err := z.zerone.NewClient(z.Node, service)
	if err != nil {
		return nil, err
	}
	c.SetTraceVerbose(z.ClientVerbose)

	if v, loaded := z.clients.LoadOrStore(service, c); loaded {
		c.Close()
		return v.(*zclient.Client), nil
	}
	return c, nil
}

func (z *Zerone) NewServer(service string) (*zserver.Server, error) {
	if v, ok := z.servers.Load(service); ok {
		return v.(*zserver.Server), nil
	}

	s, err := z.zerone.NewServer(z.Node, service)
	if err != nil {
		return nil, err
	}
	s.SetTraceVerbose(z.ServerVerbose)

	if v, loaded := z.servers.LoadOrStore(service, s); loaded {
		s.Close()
		return v.(*zserver.Server), nil
	}
	return s, nil
}
