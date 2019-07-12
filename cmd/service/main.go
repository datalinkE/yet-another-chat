package main

import (
	"net/http"
	"strings"

	"github.com/datalinkE/yet-another-chat/rpc"
	"github.com/datalinkE/yet-another-chat/storage"
)

func CapitalizerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.Title(r.URL.Path)
		h.ServeHTTP(w, r)
	})
}

func main() {
	usersHandler := rpc.NewUsersServer(&storage.Users{}, nil)
	chatsHandler := rpc.NewChatsServer(&storage.Chats{}, nil)
	messagesHandler := rpc.NewMessagesServer(&storage.Messages{}, nil)

	mux := http.NewServeMux()
	mux.Handle(usersHandler.PathPrefix(), usersHandler)
	mux.Handle(chatsHandler.PathPrefix(), chatsHandler)
	mux.Handle(messagesHandler.PathPrefix(), messagesHandler)

	http.ListenAndServe(":9000", CapitalizerMiddleware(mux))
}
