package client

import (
	"sync"

	"github.com/pedidopago/ms-grpcd/pkg/grpcd"
	"github.com/pedidopago/ms-template/internal/meta"
	"github.com/pedidopago/ms-template/protos/xyzpb"
	"google.golang.org/grpc"
)

type LiveClient struct {
	GRPCDAddr string

	client grpcd.Client
	clconn *grpc.ClientConn
	l      sync.Mutex
}

func (cl *LiveClient) Service() (xyzpb.XYZServiceClient, error) {
	cl.l.Lock()
	defer cl.l.Unlock()

	if cl.client == nil {
		var err error
		cl.client, err = grpcd.NewClient(cl.GRPCDAddr)
		if err != nil {
			return nil, err
		}
	}
	if cl.clconn == nil {
		var err error
		cl.clconn, err = cl.client.Conn(meta.ServiceName())
		if err != nil {
			return nil, err
		}
	}

	return xyzpb.NewXYZServiceClient(cl.clconn), nil
}
