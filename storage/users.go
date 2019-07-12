package storage

import (
	"context"

	"github.com/datalinkE/yet-another-chat/rpc"
)

type Users struct {
}

func (u *Users) Add(context context.Context, req *rpc.UsersAddRequest) (*rpc.UsersAddResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}
	return &rpc.UsersAddResponse{}, nil
}
