package dbmessagehelper

import (
	"context"
	"time"

	"github.com/zelenin/go-tdlib/client"

	"crawl-worker/internal/model/entity"
	"crawl-worker/internal/pkg/config/dbsaverconfig"
	mediamessagerepo "crawl-worker/internal/repository/media-message"
	"crawl-worker/pkg/container"
	"crawl-worker/pkg/gpooling"
	"crawl-worker/pkg/l"
)

type DBMessageHelper struct {
	ll               l.Logger               `container:"name"`
	gpooling         gpooling.IPool         `container:"name"`
	mediaMessageRepo mediamessagerepo.IRepo `container:"name"`

	ChannelID     int64
	ForwardToIDs  []int64
	mediaAlbumMap map[int64]chan *client.Message
}

func New() *DBMessageHelper {
	h := &DBMessageHelper{
		mediaAlbumMap: make(map[int64]chan *client.Message),
	}
	container.Fill(h)

	return h
}

func (h *DBMessageHelper) Save(ctx context.Context, config *dbsaverconfig.Config, message *client.Message) error {
	if message.MediaAlbumId != 0 {
		mediaAlbumID := int64(message.MediaAlbumId)
		messages, ok := h.mediaAlbumMap[mediaAlbumID]
		if !ok {
			messages = make(chan *client.Message)
			h.mediaAlbumMap[mediaAlbumID] = messages
			h.gpooling.Submit(func() {
				h.saveGroupMessages(ctx, config, mediaAlbumID, messages)
			})
		}
		select {
		case messages <- message:
		}
		return nil
	} else {
		files := h.buildFile(nil, message)
		return h.saveMessages(ctx, config, files)
	}
}

func (h *DBMessageHelper) buildFile(files []entity.Message, message *client.Message) []entity.Message {
	if files == nil {
		files = make([]entity.Message, 0, 1)
	}

	switch message.Content.MessageContentType() {
	case client.TypeMessageAnimation:
		content := message.Content.(*client.MessageAnimation)
		files = append(files, entity.Message{
			Type:   entity.MediaTypeAnimation,
			FileID: content.Animation.Animation.Remote.Id,
		})
	case client.TypeMessageAudio:
		content := message.Content.(*client.MessageAudio)
		files = append(files, entity.Message{
			Type:   entity.MediaTypeAudio,
			FileID: content.Audio.Audio.Remote.Id,
		})
	case client.TypeMessagePhoto:
		content := message.Content.(*client.MessagePhoto)
		files = append(files, entity.Message{
			Type:   entity.MediaTypePhoto,
			FileID: content.Photo.Sizes[len(content.Photo.Sizes)-1].Photo.Remote.Id,
		})
	case client.TypeMessageVideo:
		content := message.Content.(*client.MessageVideo)
		files = append(files, entity.Message{
			Type:   entity.MediaTypeVideo,
			FileID: content.Video.Video.Remote.Id,
		})
	default:
		h.ll.Error("unhandled message content type", l.String("message_content_type", message.Content.MessageContentType()))
	}

	return files
}

func (h *DBMessageHelper) saveMessages(ctx context.Context, config *dbsaverconfig.Config, messages []entity.Message) error {
	if len(messages) == 0 {
		return nil
	}

	msg := &entity.MediaMessage{
		Messages: messages,
		Type:     config.ChannelType,
	}

	err := h.mediaMessageRepo.Insert(ctx, msg)
	if err != nil {
		h.ll.Error("cannot create media message", l.Object("msg", msg), l.Error(err))
		return err
	} else {
		h.ll.Debug("saved message", l.Object("msg", msg))
	}
	return nil
}

func (h *DBMessageHelper) saveGroupMessages(ctx context.Context, config *dbsaverconfig.Config, mediaAlbumID int64, messages chan *client.Message) error {
	defer delete(h.mediaAlbumMap, mediaAlbumID)
	defer close(messages)

	var files []entity.Message
	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

loop:
	for {
		select {
		case message := <-messages:
			files = h.buildFile(files, message)
		case <-timer.C:
			break loop
		}
	}

	return h.saveMessages(ctx, config, files)
}
