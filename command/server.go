package command

import (
	"context"
	"fmt"
	"strings"

	"github.com/ironzhang/x-pearls/log"
	"github.com/ironzhang/zoutil"
)

type Server struct {
	zerone *zoutil.Zerone
}

func NewServer() *Server {
	return &Server{}
}

func (p *Server) GetLogLevel(ctx context.Context, args interface{}, level *string) error {
	*level = log.GetLevel()
	return nil
}

func (p *Server) SetLogLevel(ctx context.Context, level string, reply interface{}) error {
	return log.SetLevel(level)
}

func (p *Server) GetTraceVerbose(ctx context.Context, object TraceObject, verbose *int) error {
	switch strings.ToUpper(object.Type) {
	case "C":
		if object.Service == "" {
			*verbose = p.zerone.GetClientVerbose()
		} else {
			c, ok := p.zerone.GetClient(object.Service)
			if !ok {
				return fmt.Errorf("service %q client not found", object.Service)
			}
			*verbose = c.GetTraceVerbose()
		}
	case "S":
		if object.Service == "" {
			*verbose = p.zerone.GetServerVerbose()
		} else {
			s, ok := p.zerone.GetServer(object.Service)
			if !ok {
				return fmt.Errorf("service %q server not found", object.Service)
			}
			*verbose = s.GetTraceVerbose()
		}
	default:
		return fmt.Errorf("%q object type is unknown", object.Type)
	}
	return nil
}

func (p *Server) SetTraceVerbose(ctx context.Context, args SetTraceVerboseArgs, reply interface{}) error {
	switch strings.ToUpper(args.Type) {
	case "C":
		if args.Service == "" {
			p.zerone.SetClientVerbose(args.Verbose)
		} else {
			c, ok := p.zerone.GetClient(args.Service)
			if !ok {
				return fmt.Errorf("service %q client not found", args.Service)
			}
			c.SetTraceVerbose(args.Verbose)
		}
	case "S":
		if args.Service == "" {
			p.zerone.SetServerVerbose(args.Verbose)
		} else {
			s, ok := p.zerone.GetServer(args.Service)
			if !ok {
				return fmt.Errorf("service %q client not found", args.Service)
			}
			s.SetTraceVerbose(args.Verbose)
		}
	case "A":
		p.zerone.SetClientVerbose(args.Verbose)
		p.zerone.SetServerVerbose(args.Verbose)
	default:
		return fmt.Errorf("%q object type is unknown", args.Type)
	}
	return nil
}
