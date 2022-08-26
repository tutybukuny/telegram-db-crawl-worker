package listener

import (
	"context"

	"github.com/zelenin/go-tdlib/client"

	forwarderservice "crawl-worker/internal/service/forwarder"
	"crawl-worker/pkg/container"
	"crawl-worker/pkg/gpooling"
	"crawl-worker/pkg/l"
)

type TelegramListener struct {
	ll               l.Logger                  `container:"name"`
	gpooling         gpooling.IPool            `container:"name"`
	tdClient         *client.Client            `container:"name"`
	forwarderService forwarderservice.IService `container:"name"`
}

func New() *TelegramListener {
	listener := &TelegramListener{}
	container.Fill(listener)

	return listener
}

func (tl *TelegramListener) Listen() {
	listener := tl.tdClient.GetListener()
	defer listener.Close()
	ctx := context.Background()

	for update := range listener.Updates {
		envelop, ok := update.(*client.UpdateNewMessage)
		if !ok {
			continue
		}
		message := envelop.Message

		tl.forwarderService.ForwardMessage(ctx, message)
	}
}
