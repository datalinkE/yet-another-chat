package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/protobuf"
	"github.com/gorilla/handlers"

	"github.com/datalinkE/yet-another-chat/rpc"
	"github.com/datalinkE/yet-another-chat/storage"
)

var DBName = "my.db" // TODO: config/env

func main() {
	db, err := storm.Open(DBName, storm.Codec(protobuf.Codec))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	usersHandler := rpc.NewUsersServer(storage.NewUsers(db), nil)
	chatsHandler := rpc.NewChatsServer(storage.NewChats(db), nil)
	messagesHandler := rpc.NewMessagesServer(storage.NewMessages(db), nil)

	mux := http.NewServeMux()
	mux.Handle(usersHandler.PathPrefix(), usersHandler)
	mux.Handle(chatsHandler.PathPrefix(), chatsHandler)
	mux.Handle(messagesHandler.PathPrefix(), messagesHandler)

	muxWithMiddlewares := CapitalizerMiddleware(mux) // NOTE: need this to fulfill http path requirements
	muxWithMiddlewares = handlers.LoggingHandler(os.Stdout, muxWithMiddlewares)

	http.ListenAndServe(":9000", muxWithMiddlewares)
}

func CapitalizerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.Title(r.URL.Path)
		h.ServeHTTP(w, r)
	})
}
