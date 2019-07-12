package storage

import (
	"context"

	"github.com/datalinkE/yet-another-chat/rpc"
)

type Chats struct{}

func (m *Chats) Add(ctx context.Context, req *rpc.ChatsAddRequest) (*rpc.ChatsAddResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}
	return &rpc.ChatsAddResponse{}, nil
}

func (m *Chats) Get(ctx context.Context, req *rpc.ChatsGetRequest) (*rpc.ChatsGetResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}
	return &rpc.ChatsGetResponse{}, nil
}
