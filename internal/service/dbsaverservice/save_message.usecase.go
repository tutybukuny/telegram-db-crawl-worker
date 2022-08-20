package dbsaverservice

import (
	"context"
	"fmt"

	"github.com/zelenin/go-tdlib/client"

	"crawl-worker/pkg/l"
)

func (s *serviceImpl) SaveMessage(ctx context.Context, message *client.Message) error {
	config, ok := s.configMap[fmt.Sprintf("%d", message.ChatId)]
	if !ok {
		s.ll.Debug("not configured channel, just ignored it", l.Int64("channel_id", message.ChatId))
		return nil
	}

	s.ll.Info("received message for saving", l.Object("message", message))

	s.dbMessageHelper.Save(ctx, config, message)
	return nil
}
