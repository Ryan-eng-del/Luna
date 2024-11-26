package service

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/config"
	"github.com/lunarianss/Luna/internal/api-server/core/app/apps/entities"
	accountDomain "github.com/lunarianss/Luna/internal/api-server/domain/account"
	appDomain "github.com/lunarianss/Luna/internal/api-server/domain/app"
	domain "github.com/lunarianss/Luna/internal/api-server/domain/app_running"
	chatDomain "github.com/lunarianss/Luna/internal/api-server/domain/chat"
	providerDomain "github.com/lunarianss/Luna/internal/api-server/domain/provider"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/web_app"
)

type WebMessageService struct {
	appRunningDomain *domain.AppRunningDomain
	accountDomain    *accountDomain.AccountDomain
	appDomain        *appDomain.AppDomain
	chatDomain       *chatDomain.ChatDomain
	providerDomain   *providerDomain.ModelProviderDomain
	config           *config.Config
}

func NewWebMessageService(appRunningDomain *domain.AppRunningDomain, accountDomain *accountDomain.AccountDomain, appDomain *appDomain.AppDomain, config *config.Config, providerDomain *providerDomain.ModelProviderDomain, chatDomain *chatDomain.ChatDomain) *WebMessageService {
	return &WebMessageService{
		appRunningDomain: appRunningDomain,
		accountDomain:    accountDomain,
		appDomain:        appDomain,
		config:           config,
		providerDomain:   providerDomain,
		chatDomain:       chatDomain,
	}
}

func (s *WebMessageService) ListConversations(ctx context.Context, appID, endUserID string, args *dto.ListConversationQuery, invokeFrom entities.InvokeForm) (*dto.ListConversationResponse, error) {

	endUser, err := s.appRunningDomain.AppRunningRepo.GetEndUserByID(ctx, endUserID)

	if err != nil {
		return nil, err
	}

	conversations, count, err := s.chatDomain.MessageRepo.FindEndUserConversationsOrderByUpdated(ctx, appID, string(invokeFrom), endUser, args.Limit, nil, nil, args.LastID, args.SortBy)

	if err != nil {
		return nil, err
	}

	conversationList := make([]*dto.WebConversationDetail, 0)

	for _, conversation := range conversations {
		conversationList = append(conversationList, dto.ConversationRecordToDetail(conversation))
	}

	hasMore := 0

	if len(conversations) == args.Limit {
		hasMore = 1
	}

	return &dto.ListConversationResponse{
		Data:    conversationList,
		Limit:   args.Limit,
		HasMore: hasMore,
		Count:   count,
	}, nil

}
