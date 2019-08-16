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

var testUsers []*rpc.User

func TestSanity01_Users(t *testing.T) {
	u := rpc.NewUsersJSONClient("http://:9000", &http.Client{Timeout: time.Second})
	ctx := context.Background()

	resp, err := u.Add(ctx, &rpc.UsersAddRequest{Username: "u1"})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, resp.Id > 0)

	_, err = u.Add(ctx, &rpc.UsersAddRequest{Username: "u1"})
	require.Error(t, err)

	resp, err = u.Add(ctx, &rpc.UsersAddRequest{Username: "u2"})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, resp.Id > 0)

	usersResp, err := u.Get(ctx, &rpc.UsersGetRequest{})
	require.NoError(t, err)
	require.NotNil(t, usersResp)
	require.NotEmpty(t, usersResp.Users)

	testUsers = usersResp.Users

	usernames := []string{}

	for _, el := range testUsers {
		usernames = append(usernames, el.Username)
	}

	require.Equal(t, len(usernames), 2)
	require.Contains(t, usernames, "u1")
	require.Contains(t, usernames, "u2")
}

func TestSanity02_Chats(t *testing.T) {
	// TODO
}

func TestSanity03_Messages(t *testing.T) {
	// TODO
}
