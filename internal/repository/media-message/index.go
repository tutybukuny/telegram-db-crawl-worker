package mediamessagerepo

import (
	"github.com/thnthien/impa/repository"

	"crawl-worker/internal/model/entity"
)

type IRepo interface {
	repository.IInsert[entity.MediaMessage, int64]
	repository.IFindByID[entity.MediaMessage, int64]
}
