package proxy

import (
	"bytes"
	"context"

	"github.com/golang/protobuf/jsonpb"
	"github.com/ldelossa/vim-grpc.vim/channel"
	"github.com/ldelossa/vim-grpc.vim/proto/env"
	pb "github.com/ldelossa/vim-grpc.vim/proto/env"
)

// EnvironmentService provides RPCs describing Vim's current
// editor environment.
type EnvironmentService struct {
	*Proxy
	pb.UnimplementedEnvServer
}

func NewEnvService(ctx context.Context, proxy *Proxy) *EnvironmentService {
	return &EnvironmentService{
		Proxy: proxy,
	}
}

func (env *EnvironmentService) GetEnv(ctx context.Context, req *env.GetEnvRequest) (*env.GetEnvResponse, error) {
	const (
		RPC = "GetEnv"
	)

	ch := env.Channel()
	if !ch.ChannelOpen() {
		return nil, channel.ErrChanClosed
	}

	m := jsonpb.Marshaler{
		EmitDefaults: false,
	}

	var b bytes.Buffer
	err := m.Marshal(&b, req)
	if err != nil {
		return nil, err
	}

	e := channel.Envelope{
		RPC:  RPC,
		Body: b.Bytes(),
	}

	e, err = ch.Send(ctx, &e).Wait(ctx)
	if err != nil {
		return nil, err
	}

	resp := &pb.GetEnvResponse{}
	err = jsonpb.Unmarshal(bytes.NewReader(e.Body), resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
