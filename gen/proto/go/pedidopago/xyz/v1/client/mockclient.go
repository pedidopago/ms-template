package client

import (
	"sync"

	xyzv1 "github.com/pedidopago/ms-template/gen/proto/go/pedidopago/xyz/v1"
	"github.com/pedidopago/ms-template/gen/proto/go/pedidopago/xyz/v1/mocks"
)

type MockClient struct {
	l      sync.Mutex
	mockcl *mocks.XYZServiceClient
}

func (cl *MockClient) Service() (xyzv1.XYZServiceClient, error) {
	cl.l.Lock()
	defer cl.l.Unlock()

	if cl.mockcl == nil {
		cl.mockcl = &mocks.XYZServiceClient{}
	}

	return cl.mockcl, nil
}

func (cl *MockClient) Mock() (*mocks.XYZServiceClient, error) {
	cl.l.Lock()
	defer cl.l.Unlock()

	if cl.mockcl == nil {
		cl.mockcl = &mocks.XYZServiceClient{}
	}

	return cl.mockcl, nil
}
