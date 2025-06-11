package scenarios

import (
	"context"

	"github.com/Axel791/appkit"
	groupdomains "github.com/Axel791/auth/internal/domains/group"
	userdomain "github.com/Axel791/auth/internal/domains/user"
	"github.com/Axel791/auth/internal/usecases/group/dto"
	"github.com/Axel791/auth/internal/usecases/group/providers"
	log "github.com/sirupsen/logrus"
)

type GetGroupByUserID struct {
	logger          *log.Logger
	groupRepository providers.GroupRepository
}

func NewGetGroupByUserID(
	logger *log.Logger,
	groupRepository providers.GroupRepository,
) *GetGroupByUserID {
	return &GetGroupByUserID{
		logger:          logger,
		groupRepository: groupRepository,
	}
}

func (s *GetGroupByUserID) Execute(ctx context.Context, userID userdomain.UserID) ([]dto.Group, error) {
	groups, err := s.groupRepository.GetUserGroups(ctx, userID)
	if err != nil {
		s.logger.WithError(err).Error("failed to fetch groups")

		return nil, appkit.Wrap(appkit.Unknown, "error getting user groups", err)
	}

	return toGroupsDTO(groups), nil
}

func toGroupsDTO(groups []groupdomains.Group) []dto.Group {
	var dtoGroups []dto.Group
	for _, group := range groups {
		dtoGroups = append(
			dtoGroups,
			dto.Group{
				ID:          group.ID().ToInt64(),
				Name:        group.Name(),
				Description: group.Description(),
				CreatedAt:   group.CreatedAt(),
			},
		)
	}
	return dtoGroups
}
