package forwardingmessagehelper

import (
	"context"
	"fmt"
	"time"

	"github.com/zelenin/go-tdlib/client"

	mediamessagetype "crawl-worker/internal/constant/media-message-type"
	"crawl-worker/internal/pkg/config/dbsaverconfig"
	"crawl-worker/pkg/container"
	"crawl-worker/pkg/gpooling"
	"crawl-worker/pkg/l"
)

type ForwardingMessageHelper struct {
	ll               l.Logger                `container:"name"`
	gpooling         gpooling.IPool          `container:"name"`
	tdClient         *client.Client          `container:"name"`
	dbsaverConfigMap dbsaverconfig.ConfigMap `container:"name"`

	mediaAlbumMap map[int64]chan *client.Message
	dbChannelID   int64
}

func New(dbChannelID int64) *ForwardingMessageHelper {
	h := &ForwardingMessageHelper{
		mediaAlbumMap: make(map[int64]chan *client.Message),
		dbChannelID:   dbChannelID,
	}
	container.Fill(h)

	return h
}

func (h *ForwardingMessageHelper) Forward(ctx context.Context, config *dbsaverconfig.Config, message *client.Message) error {
	//if message.MediaAlbumId != 0 {
	//	mediaAlbumID := int64(message.MediaAlbumId)
	//	messages, ok := h.mediaAlbumMap[mediaAlbumID]
	//	if !ok {
	//		messages = make(chan *client.Message)
	//		h.mediaAlbumMap[mediaAlbumID] = messages
	//		h.gpooling.Submit(func() {
	//			h.sendGroupMessages(ctx, config, mediaAlbumID, messages)
	//		})
	//	}
	//	select {
	//	case messages <- message:
	//	}
	//	return nil
	//} else {
	files := h.buildFile(nil, message, config.ChannelType)
	return h.sendMessages(ctx, files)
	//}
}

func (h *ForwardingMessageHelper) buildFile(files []client.InputMessageContent, message *client.Message, messageType int) []client.InputMessageContent {
	if files == nil {
		files = make([]client.InputMessageContent, 0, 1)
	}

	caption := "nsfw"
	switch messageType {
	case mediamessagetype.SFW:
		caption = "sfw"
	case mediamessagetype.Others:
		caption = "others"
	}

	if message.MediaAlbumId != 0 {
		caption += fmt.Sprintf(";%d", message.MediaAlbumId)
	}

	var file client.InputMessageContent
	switch message.Content.MessageContentType() {
	case client.TypeMessageAnimation:
		content := message.Content.(*client.MessageAnimation)
		file = &client.InputMessageAnimation{
			Animation: &client.InputFileRemote{Id: content.Animation.Animation.Remote.Id},
			Caption:   &client.FormattedText{Text: caption},
		}
	case client.TypeMessageAudio:
		content := message.Content.(*client.MessageAudio)
		file = &client.InputMessageAudio{
			Audio:   &client.InputFileRemote{Id: content.Audio.Audio.Remote.Id},
			Caption: &client.FormattedText{Text: caption},
		}
	case client.TypeMessagePhoto:
		content := message.Content.(*client.MessagePhoto)
		file = &client.InputMessagePhoto{
			Photo:   &client.InputFileRemote{Id: content.Photo.Sizes[len(content.Photo.Sizes)-1].Photo.Remote.Id},
			Caption: &client.FormattedText{Text: caption},
		}
	case client.TypeMessageVideo:
		content := message.Content.(*client.MessageVideo)
		file = &client.InputMessageVideo{
			Video:   &client.InputFileRemote{Id: content.Video.Video.Remote.Id},
			Caption: &client.FormattedText{Text: caption},
		}
	default:
		h.ll.Error("unhandled message content type", l.String("message_content_type", message.Content.MessageContentType()))
	}

	return append(files, file)
}

func (h *ForwardingMessageHelper) sendMessages(
	ctx context.Context,
	messages []client.InputMessageContent,
) error {
	if len(messages) == 0 {
		return nil
	}

	msg := &client.SendMessageAlbumRequest{
		ChatId:               h.dbChannelID,
		InputMessageContents: messages,
	}

	sentMsg, err := h.tdClient.SendMessageAlbum(msg)
	if err != nil {
		h.ll.Error("cannot send message", l.Object("msg", msg), l.Error(err))
		return err
	} else {
		h.ll.Debug("sent message", l.Object("msg", sentMsg))
	}
	return nil
}

func (h *ForwardingMessageHelper) sendGroupMessages(ctx context.Context, config *dbsaverconfig.Config, mediaAlbumID int64, messages chan *client.Message) error {
	defer delete(h.mediaAlbumMap, mediaAlbumID)
	defer close(messages)

	var files []client.InputMessageContent
	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

loop:
	for {
		select {
		case message := <-messages:
			fs := h.buildFile(files, message, config.ChannelType)
			files = fs
		case <-timer.C:
			break loop
		}
	}

	return h.sendMessages(ctx, files)
}

func (h *ForwardingMessageHelper) getMessageType(message *client.Message) int {
	messageType := -1
	if message.ForwardInfo != nil {
		channelID := int64(0)
		switch message.ForwardInfo.Origin.MessageForwardOriginType() {
		case client.TypeMessageForwardOriginChannel:
			channelID = message.ForwardInfo.Origin.(*client.MessageForwardOriginChannel).ChatId
		case client.TypeMessageForwardOriginChat:
			channelID = message.ForwardInfo.Origin.(*client.MessageForwardOriginChat).SenderChatId
		}
		if forwardConfig, ok := h.dbsaverConfigMap[fmt.Sprintf("%d", channelID)]; ok {
			messageType = forwardConfig.ChannelType
		}
	}
	return messageType
}
