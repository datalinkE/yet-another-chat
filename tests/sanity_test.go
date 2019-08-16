package tests

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/datalinkE/yet-another-chat/internal/config"
	"github.com/datalinkE/yet-another-chat/internal/service"
	"github.com/datalinkE/yet-another-chat/rpc"
)

var (
	testUsers = map[string]*rpc.User{}
	testChats = map[string]*rpc.Chat{}
)

func TestMain(m *testing.M) {

	file, err := ioutil.TempFile("", "bolt.*.db")
	if err != nil {
		panic(err)
	}
	defer os.Remove(file.Name())

	go service.Run(&config.Config{
		ListenAddr: ":9000",
		BoltFile:   file.Name(),
	})

	time.Sleep(time.Second)

	os.Exit(m.Run())
}

func TestSanity01_Users(t *testing.T) {
	users := rpc.NewUsersJSONClient("http://:9000", &http.Client{Timeout: time.Second})
	ctx := context.Background()

	resp, err := users.Add(ctx, &rpc.UsersAddRequest{Username: "u1"})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, resp.Id > 0)

	_, err = users.Add(ctx, &rpc.UsersAddRequest{Username: "u1"})
	require.Error(t, err)

	resp, err = users.Add(ctx, &rpc.UsersAddRequest{Username: "u2"})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, resp.Id > 0)

	resp, err = users.Add(ctx, &rpc.UsersAddRequest{Username: "u3"})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, resp.Id > 0)

	usersResp, err := users.Get(ctx, &rpc.UsersGetRequest{})
	require.NoError(t, err)
	require.NotNil(t, usersResp)
	require.NotEmpty(t, usersResp.Users)

	for _, el := range usersResp.Users {
		testUsers[el.Username] = el
	}

	require.Equal(t, len(testUsers), 3)
}

func TestSanity02_Chats(t *testing.T) {
	chats := rpc.NewChatsJSONClient("http://:9000", &http.Client{Timeout: time.Second})
	ctx := context.Background()

	resp, err := chats.Add(ctx, &rpc.ChatsAddRequest{Name: "c1", Users: []int64{testUsers["u1"].Id, testUsers["u2"].Id}})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, resp.Id > 0)

	time.Sleep(time.Second)

	resp, err = chats.Add(ctx, &rpc.ChatsAddRequest{Name: "c2", Users: []int64{testUsers["u2"].Id, testUsers["u3"].Id}})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, resp.Id > 0)

	chatsResp, err := chats.Get(ctx, &rpc.ChatsGetRequest{User: testUsers["u1"].Id})
	require.NoError(t, err)
	require.NotNil(t, chatsResp)
	require.Equal(t, 1, len(chatsResp.Chats))
	require.Equal(t, "c1", chatsResp.Chats[0].Name)

	chatsResp, err = chats.Get(ctx, &rpc.ChatsGetRequest{User: testUsers["u2"].Id})
	require.NoError(t, err)
	require.NotNil(t, chatsResp)
	require.Equal(t, 2, len(chatsResp.Chats))

	for _, el := range chatsResp.Chats {
		testChats[el.Name] = el
	}

	require.Equal(t, "c2", chatsResp.Chats[0].Name)
	require.Equal(t, "c1", chatsResp.Chats[1].Name)
}

func TestSanity03_Messages(t *testing.T) {
	// TODO
	messages := rpc.NewMessagesJSONClient("http://:9000", &http.Client{Timeout: time.Second})
	ctx := context.Background()

	_, err := messages.Add(ctx, &rpc.MessagesAddRequest{
		Chat:   testChats["c1"].GetId(),
		Author: testUsers["u1"].GetId(),
		Text:   "u1 c1 text1",
	})
	require.NoError(t, err)

	_, err = messages.Add(ctx, &rpc.MessagesAddRequest{
		Chat:   testChats["c1"].GetId(),
		Author: testUsers["u1"].GetId(),
		Text:   "u1 c1 text2",
	})
	require.NoError(t, err)

	_, err = messages.Add(ctx, &rpc.MessagesAddRequest{
		Chat:   testChats["c1"].GetId(),
		Author: testUsers["u2"].GetId(),
		Text:   "u2 c1 text3",
	})
	require.NoError(t, err)

	resp, err := messages.Get(ctx, &rpc.MessagesGetRequest{Chat: testChats["c1"].GetId()})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, 3, len(resp.Messages))

	require.Contains(t, resp.Messages[0].Text, "text3")
	require.Contains(t, resp.Messages[1].Text, "text2")
	require.Contains(t, resp.Messages[2].Text, "text1")
}
