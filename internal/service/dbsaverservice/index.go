package dbsaverservice

import (
	"context"

	"github.com/zelenin/go-tdlib/client"

	"crawl-worker/internal/pkg/config/dbsaverconfig"
	dbmessagehelper "crawl-worker/internal/pkg/helper/db-message"
	"crawl-worker/internal/repository/media-message"
	"crawl-worker/pkg/container"
	"crawl-worker/pkg/gpooling"
	"crawl-worker/pkg/l"
)

type IService interface {
	SaveMessage(ctx context.Context, message *client.Message) error
}

type serviceImpl struct {
	ll               l.Logger                `container:"name"`
	gpooling         gpooling.IPool          `container:"name"`
	tdClient         *client.Client          `container:"name"`
	mediaMessageRepo mediamessagerepo.IRepo  `container:"name"`
	dbsaverConfigMap dbsaverconfig.ConfigMap `container:"name"`
	filteredContents []string                `container:"name"`

	dbMessageHelper *dbmessagehelper.DBMessageHelper
}

func New() *serviceImpl {
	service := &serviceImpl{
		dbMessageHelper: dbmessagehelper.New(),
	}
	container.Fill(service)

	return service
}
