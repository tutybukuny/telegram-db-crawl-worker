package forwarderservice

import (
	"context"
	"time"

	"github.com/zelenin/go-tdlib/client"
	"go.uber.org/ratelimit"

	"crawl-worker/internal/pkg/config/dbsaverconfig"
	forwardingmessagehelper "crawl-worker/internal/pkg/helper/forwarding-message"
	channelrepo "crawl-worker/internal/repository/channel"
	"crawl-worker/internal/repository/media-message"
	"crawl-worker/pkg/container"
	"crawl-worker/pkg/gpooling"
	"crawl-worker/pkg/l"
)

type IService interface {
	ForwardMessage(ctx context.Context, message *client.Message) error
}

type serviceImpl struct {
	ll               l.Logger                `container:"name"`
	gpooling         gpooling.IPool          `container:"name"`
	tdClient         *client.Client          `container:"name"`
	mediaMessageRepo mediamessagerepo.IRepo  `container:"name"`
	channelRepo      channelrepo.IRepo       `container:"name"`
	dbsaverConfigMap dbsaverconfig.ConfigMap `container:"name"`
	filteredContents []string                `container:"name"`

	forwarder *forwardingmessagehelper.ForwardingMessageHelper
	limiter   ratelimit.Limiter
}

func New(dbChannelID int64) *serviceImpl {
	service := &serviceImpl{
		forwarder: forwardingmessagehelper.New(dbChannelID),
		limiter:   ratelimit.New(300, ratelimit.Per(time.Minute)),
	}
	container.Fill(service)

	return service
}
