package repo

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/model/v1"
)

type MessageRepo interface {
	CreateMessage(ctx context.Context, message *model.Message) (*model.Message, error)
	CreateConversation(ctx context.Context, conversation *model.Conversation) (*model.Conversation, error)

	UpdateMessage(ctx context.Context, message *model.Message) error
	UpdateConversationUpdateAt(ctx context.Context, appID string, conversation *model.Conversation) error

	GetMessageByID(ctx context.Context, messageID string) (*model.Message, error)
	GetConversationByID(ctx context.Context, conversationID string) (*model.Conversation, error)
	GetConversationByUser(ctx context.Context, appId string, conversationID string, user model.BaseAccount) (*model.Conversation, error)
}
