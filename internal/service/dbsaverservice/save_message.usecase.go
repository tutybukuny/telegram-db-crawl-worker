package dbsaverservice

import (
	"context"
	"fmt"
	"strings"

	"github.com/zelenin/go-tdlib/client"

	"crawl-worker/internal/model/entity"
	"crawl-worker/pkg/l"
)

var channels = make(map[int64]*entity.Channel)

func (s *serviceImpl) SaveMessage(ctx context.Context, message *client.Message) error {
	config, ok := s.dbsaverConfigMap[fmt.Sprintf("%d", message.ChatId)]
	if !ok {
		s.ll.Debug("not configured channel, just ignored it", l.Int64("channel_id", message.ChatId))
		return nil
	}

	if s.isFilteredMessage(ctx, message) {
		s.ll.Debug("filtered message")
		return nil
	}

	var channel *entity.Channel
	if channel, ok = channels[message.ChatId]; !ok {
		chat, err := s.tdClient.GetChat(&client.GetChatRequest{ChatId: message.ChatId})
		if err == nil {
			if channel, err = s.channelRepo.GetOrCreate(ctx, chat.Id, chat.Title); err != nil {
				s.ll.Error("cannot get or create channel", l.Object("chat", chat), l.Error(err))
			} else {
				channels[channel.ID] = channel
			}
		} else {
			s.ll.Error("cannot get channel info", l.Int64("chat_id", message.ChatId), l.Error(err))
		}
	}

	s.ll.Info("received message for saving", l.Object("message", message))
	s.dbMessageHelper.Save(ctx, channel, config, message)
	return nil
}

func (s *serviceImpl) isFilteredMessage(ctx context.Context, message *client.Message) bool {
	switch message.Content.MessageContentType() {
	case client.TypeMessageAnimation:
		content := message.Content.(*client.MessageAnimation)
		return s.isFilteredContent(content.Caption.Text)
	case client.TypeMessageAudio:
		content := message.Content.(*client.MessageAudio)
		return s.isFilteredContent(content.Caption.Text)
	case client.TypeMessagePhoto:
		content := message.Content.(*client.MessagePhoto)
		return s.isFilteredContent(content.Caption.Text)
	case client.TypeMessageVideo:
		content := message.Content.(*client.MessageVideo)
		return s.isFilteredContent(content.Caption.Text)
	default:
		return true
	}
}

func (s *serviceImpl) isFilteredContent(content string) bool {
	content = strings.ToLower(content)
	for _, filteredContent := range s.filteredContents {
		if strings.Contains(content, filteredContent) {
			return true
		}
	}
	return false
}
