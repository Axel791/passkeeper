package grpc

import (
	"context"
	"fmt"

	"github.com/Axel791/appkit"
	"github.com/Axel791/passkeeper_grpc/pb"
	groupdomain "github.com/Axel791/vault/internal/domains/group"
	userdomain "github.com/Axel791/vault/internal/domains/user"
)

type AuthRepo struct {
	authClient pb.AuthServiceClient
}

func NewAuthRepo(authClient pb.AuthServiceClient) *AuthRepo {
	return &AuthRepo{authClient: authClient}
}

func (r *AuthRepo) ValidateToken(ctx context.Context, token string) (userdomain.UserID, error) {
	user, err := r.authClient.ValidateToken(ctx, &pb.ValidateTokenRequest{Token: token})
	if err != nil {
		return userdomain.UserID{}, appkit.ForbiddenError("invalid token")
	}

	userID, err := userdomain.NewUserID(user.User.Id)
	if err != nil {
		return userdomain.UserID{}, appkit.ForbiddenError("invalid user id")
	}

	return userID, nil
}

func (r *AuthRepo) GetUserGroups(ctx context.Context, id userdomain.UserID) ([]groupdomain.Group, error) {
	groups, err := r.authClient.GetUserGroups(ctx, &pb.GetUserGroupsRequest{UserId: id.ToInt64()})
	if err != nil {
		return nil, appkit.Wrap(appkit.Unknown, "err fetch user group", err)
	}
	return toGroupDomain(groups)
}

func toGroupDomain(groups *pb.GetUserGroupsResponse) ([]groupdomain.Group, error) {
	var result []groupdomain.Group
	for _, group := range groups.Groups {
		groupID, err := groupdomain.NewGroupID(group.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to create group id: %w", err)
		}

		groupDomain := groupdomain.NewGroup(
			groupID,
			group.Name,
			group.Description,
		)

		result = append(result, groupDomain)
	}
	return result, nil
}
