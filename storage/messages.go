package storage

import (
	"context"

	"github.com/datalinkE/yet-another-chat/rpc"
)

type Messages struct{}

func (m *Messages) Add(ctx context.Context, req *rpc.MessagesAddRequest) (*rpc.MessagesAddResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}
	return &rpc.MessagesAddResponse{}, nil
}

func (m *Messages) Get(ctx context.Context, req *rpc.MessagesGetRequest) (*rpc.MessagesGetResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}
	return &rpc.MessagesGetResponse{}, nil
}
