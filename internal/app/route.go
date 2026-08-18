package app

import (
	"net/http"

	"buf.build/gen/go/kevin-labs/riotdata/connectrpc/go/kevin/riotdata/v1/riotdatav1connect"
	"connectrpc.com/connect"
	"connectrpc.com/validate"
)

func routeProfileService(mux *http.ServeMux, handler riotdatav1connect.ProfileServiceHandler) {
	path, httpHandler := riotdatav1connect.NewProfileServiceHandler(
		handler,
		connect.WithInterceptors(validate.NewInterceptor()),
	)

	mux.Handle(path, httpHandler)
}
