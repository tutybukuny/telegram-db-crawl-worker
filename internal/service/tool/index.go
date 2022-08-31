package toolservice

import (
	"context"

	"github.com/zelenin/go-tdlib/client"

	"crawl-worker/pkg/container"
	"crawl-worker/pkg/l"
)

type IService interface {
	GetChatFromName(ctx context.Context, chatName string) error
}

type serviceImpl struct {
	ll       l.Logger       `container:"name"`
	tdClient *client.Client `container:"name"`
}

func New() *serviceImpl {
	service := &serviceImpl{}
	container.Fill(service)

	return service
}
