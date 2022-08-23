package channelrepo

import (
	"context"

	"github.com/thnthien/impa/repository"

	"crawl-worker/internal/model/entity"
)

type IRepo interface {
	repository.IInsert[entity.Channel, int64]
	repository.IFindByID[entity.Channel, int64]
	repository.IUpdate[entity.Channel, int64]

	GetOrCreate(ctx context.Context, channelID int64, channelName string) (*entity.Channel, error)
}
