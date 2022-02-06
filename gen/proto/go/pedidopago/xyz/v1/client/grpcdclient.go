package client

import (
	"sync"

	"github.com/pedidopago/ms-grpcd/pkg/grpcd"
	xyzv1 "github.com/pedidopago/ms-template/gen/proto/go/pedidopago/xyz/v1"
	"github.com/pedidopago/ms-template/internal/meta"
	"google.golang.org/grpc"
)

type LiveClient struct {
	GRPCDAddr string

	client grpcd.Client
	clconn *grpc.ClientConn
	l      sync.Mutex
}

func (cl *LiveClient) Service() (xyzv1.XYZServiceClient, error) {
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

	return xyzv1.NewXYZServiceClient(cl.clconn), nil
}
