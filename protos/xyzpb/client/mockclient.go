package client

import (
	"github.com/pedidopago/ms-template/protos/xyzpb"
	"github.com/pedidopago/ms-template/protos/xyzpb/mocks"
	"sync"
)

type MockClient struct {
	l      sync.Mutex
	mockcl *mocks.XYZServiceClient
}

func (cl *MockClient) Service() (xyzpb.XYZServiceClient, error) {
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
