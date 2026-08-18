package profile

import (
	"buf.build/gen/go/kevin-labs/riotdata/connectrpc/go/kevin/riotdata/v1/riotdatav1connect"
)

type ProfileServiceHandler struct {
	riotdatav1connect.UnimplementedProfileServiceHandler

	service *ProfileService
}
