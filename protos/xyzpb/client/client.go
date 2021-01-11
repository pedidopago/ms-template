package client

import (
	"github.com/pedidopago/ms-template/protos/xyzpb"
)

type Client interface {
	Service() (xyzpb.XYZServiceClient, error)
}
