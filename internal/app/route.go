package app

import (
	"net/http"

	"buf.build/gen/go/kevin-labs/lol-service/connectrpc/go/kevin/lolservice/v1/lolservicev1connect"
	"connectrpc.com/connect"
	"connectrpc.com/validate"
)

func routeProfileService(mux *http.ServeMux, handler lolservicev1connect.ProfileServiceHandler) {
	path, httpHandler := lolservicev1connect.NewProfileServiceHandler(handler,
		connect.WithInterceptors(validate.NewInterceptor()))

	mux.Handle(path, httpHandler)
}
