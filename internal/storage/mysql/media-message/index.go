package mediamessagestore

import (
	"github.com/thnthien/impa/repository/gorm"
	"gorm.io/gorm"

	"crawl-worker/internal/model/entity"
)

type repoImpl struct {
	*gormrepository.BaseRepo
	*gormrepository.InsertBaseRepo[entity.MediaMessage, int64]
	*gormrepository.FindByIDBaseRepo[entity.MediaMessage, int64]
}

func New(db *gorm.DB) *repoImpl {
	base := gormrepository.NewBaseRepo[entity.MediaMessage](db)
	insert := gormrepository.NewInsertBaseRepo[entity.MediaMessage, int64](base)
	find := gormrepository.NewFindByIDBaseRepo[entity.MediaMessage, int64](base)
	return &repoImpl{base, insert, find}
}
