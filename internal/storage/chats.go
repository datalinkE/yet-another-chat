package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/asdine/storm"

	"github.com/datalinkE/yet-another-chat/rpc"
)

type Chats struct {
	DB *storm.DB
}

func NewChats(db *storm.DB) *Chats {
	return &Chats{DB: db}
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

	// TODO: "bad" select, try to use indexes, maybe with more buckets
	err = m.DB.Select().
		OrderBy("CreatedAt").
		Reverse().
		Each(
			&rpc.Chat{},
			func(v interface{}) error {
				el, ok := v.(*rpc.Chat)
				if !ok {
					return fmt.Errorf("can't convert object to chat")
				}
				if Contains(el.Users, req.User) {
					chats = append(chats, el)
				}
				return nil
			},
		)
	if err != nil {
		return nil, err
	}

	return &rpc.ChatsGetResponse{Chats: chats}, nil
}
