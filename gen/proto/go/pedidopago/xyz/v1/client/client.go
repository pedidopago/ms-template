package client

import (
	xyzv1 "github.com/pedidopago/ms-template/gen/proto/go/pedidopago/xyz/v1"
)

type Client interface {
	Service() (xyzv1.XYZServiceClient, error)
}
