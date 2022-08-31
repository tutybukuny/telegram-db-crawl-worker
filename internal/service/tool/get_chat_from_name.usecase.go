package toolservice

import (
	"context"

	"github.com/zelenin/go-tdlib/client"

	"crawl-worker/pkg/l"
)

func (s *serviceImpl) GetChatFromName(ctx context.Context, chatName string) error {
	chats, err := s.tdClient.SearchChatsOnServer(&client.SearchChatsOnServerRequest{Query: chatName, Limit: 100})
	if err != nil {
		s.ll.Error("cannot search chats", l.String("chat_name", chatName), l.Error(err))
		return err
	}

	for _, chatID := range chats.ChatIds {
		chat, err := s.tdClient.GetChat(&client.GetChatRequest{ChatId: chatID})
		if err != nil {
			s.ll.Error("cannot get chat", l.Int64("chat_id", chatID), l.Error(err))
			continue
		}
		s.ll.Info("get chat info", l.Int64("chat_id", chatID), l.String("chat", chat.Title))
	}
	return nil
}
