package storage

import (
	"context"
	"time"

	"github.com/asdine/storm"

	"github.com/datalinkE/yet-another-chat/rpc"
)

type Users struct {
	DB *storm.DB
}

func (u *Users) Add(ctx context.Context, req *rpc.UsersAddRequest) (*rpc.UsersAddResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	user := rpc.User{
		Username:  req.Username,
		CreatedAt: time.Now().Format(time.RFC3339Nano),
	}

	err = u.DB.Save(&user)
	if err != nil {
		return nil, err
	}

	return &rpc.UsersAddResponse{Id: user.GetId()}, nil
}

func (u *Users) Get(ctx context.Context, req *rpc.UsersGetRequest) (*rpc.UsersGetResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	users := []*rpc.User{}
	err = u.DB.All(&users)
	if err != nil {
		return nil, err
	}

	return &rpc.UsersGetResponse{Users: users}, nil
}
