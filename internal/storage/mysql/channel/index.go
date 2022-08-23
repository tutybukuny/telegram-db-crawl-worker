package channelstore

import (
	"context"

	gormrepository "github.com/thnthien/impa/repository/gorm"
	"gorm.io/gorm"

	"crawl-worker/internal/model/entity"
)

type repoImpl struct {
	*gormrepository.BaseRepo
	*gormrepository.InsertBaseRepo[entity.Channel, int64]
	*gormrepository.UpdateBaseRepo[entity.Channel, int64]
	*gormrepository.FindByIDBaseRepo[entity.Channel, int64]
}

func New(db *gorm.DB) *repoImpl {
	base := gormrepository.NewBaseRepo[entity.Channel](db)
	insert := gormrepository.NewInsertBaseRepo[entity.Channel, int64](base)
	update := gormrepository.NewUpdateBaseRepo[entity.Channel, int64](base)
	find := gormrepository.NewFindByIDBaseRepo[entity.Channel, int64](base)
	return &repoImpl{base, insert, update, find}
}

func (r *repoImpl) GetOrCreate(ctx context.Context, channelID int64, channelName string) (*entity.Channel, error) {
	obj, err := r.FindByID(ctx, channelID)
	if err != nil {
		return nil, err
	}
	if obj == nil {
		obj = &entity.Channel{
			ID:   channelID,
			Name: channelName,
		}
		err = r.Insert(ctx, obj)
		if err != nil {
			return nil, err
		}
	}
	return obj, nil
}
