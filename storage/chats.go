package storage

import (
	"context"
	"time"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"

	"github.com/datalinkE/yet-another-chat/rpc"
)

type Chats struct {
	DB *storm.DB
}

func (c *Chats) Add(ctx context.Context, req *rpc.ChatsAddRequest) (*rpc.ChatsAddResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	createdAt := time.Now().Format(time.RFC3339Nano)
	chat := rpc.Chat{
		Name:      req.Name,
		Users:     req.Users,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}

	err = c.DB.Save(&chat)
	if err != nil {
		return nil, err
	}

	return &rpc.ChatsAddResponse{Id: chat.GetId()}, nil
}

func (m *Chats) Get(ctx context.Context, req *rpc.ChatsGetRequest) (*rpc.ChatsGetResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	chats := []*rpc.Chat{}
	err = m.DB.Select(q.Eq("Id", req.User)).OrderBy("UpdatedAt").Find(&chats)
	if err != nil {
		return nil, err
	}

	return &rpc.ChatsGetResponse{Chats: chats}, nil
}
