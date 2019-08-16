package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/asdine/storm"

	"github.com/datalinkE/yet-another-chat/rpc"
)

type Messages struct {
	DB *storm.DB
}

func NewMessages(db *storm.DB) *Messages {
	return &Messages{DB: db}
}

func Contains(container []int64, value int64) bool {
	for _, el := range container {
		if el == value {
			return true
		}
	}
	return false
}

func (m *Messages) Add(ctx context.Context, req *rpc.MessagesAddRequest) (*rpc.MessagesAddResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	chat := rpc.Chat{}

	err = m.DB.One("Id", req.GetChat(), &chat)
	if err != nil {
		return nil, err
	}

	if chat.Id == 0 {
		return nil, fmt.Errorf("no chat with id=%d", chat.Id)
	}

	if !Contains(chat.GetUsers(), req.GetAuthor()) {
		return nil, fmt.Errorf("no user with id=%d in chat=", req.GetAuthor(), req.GetChat())
	}

	message := rpc.Message{
		Chat:      req.GetChat(),
		Author:    req.GetAuthor(),
		Text:      req.GetText(),
		CreatedAt: time.Now().Format(time.RFC3339Nano),
	}

	err = m.DB.Save(&message) // TODO: transaction
	if err != nil {
		return nil, err
	}

	err = m.DB.UpdateField(&chat, "UpdatedAt", message.GetCreatedAt()) // TODO: transaction
	if err != nil {
		return nil, err
	}

	return &rpc.MessagesAddResponse{Id: message.GetId()}, nil
}

func (m *Messages) Get(ctx context.Context, req *rpc.MessagesGetRequest) (*rpc.MessagesGetResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	messages := []*rpc.Message{}
	err = m.DB.Find("Chat", req.GetChat(), &messages) // TODO: not working now - can't get chat messages
	if err != nil {
		return nil, err
	}

	return &rpc.MessagesGetResponse{Messages: messages}, nil
}
