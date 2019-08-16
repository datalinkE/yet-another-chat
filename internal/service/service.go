package service

import (
	"net/http"
	"os"
	"strings"

	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/protobuf"
	"github.com/gorilla/handlers"

	"github.com/datalinkE/yet-another-chat/internal/config"
	"github.com/datalinkE/yet-another-chat/internal/storage"
	"github.com/datalinkE/yet-another-chat/rpc"
)

func Run(cfg *config.Config) error {
	db, err := storm.Open(cfg.BoltFile, storm.Codec(protobuf.Codec))
	if err != nil {
		return err
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

	return http.ListenAndServe(cfg.ListenAddr, muxWithMiddlewares)
}

func CapitalizerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.Title(r.URL.Path)
		h.ServeHTTP(w, r)
	})
}
