package v1

import (
	"context"

	userdomain "github.com/Axel791/auth/internal/domains/user"
	"github.com/Axel791/auth/internal/usecases/group"
	"github.com/Axel791/auth/internal/usecases/user"
	"github.com/Axel791/passkeeper_grpc/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer

	getUserGroups group.GetGroupByUserID
	validateToken user.Validate
}

func NewAuthServer(
	getUserGroups group.GetGroupByUserID,
	validateToken user.Validate,
) *AuthServer {
	return &AuthServer{
		getUserGroups: getUserGroups,
		validateToken: validateToken,
	}
}

func (s *AuthServer) GetUserGroups(
	ctx context.Context,
	req *pb.GetUserGroupsRequest,
) (*pb.GetUserGroupsResponse, error) {
	userID, err := userdomain.NewUserID(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	groups, err := s.getUserGroups.Execute(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get groups: %v", err)
	}

	out := &pb.GetUserGroupsResponse{}
	for _, g := range groups {
		out.Groups = append(out.Groups, &pb.Group{
			Id:          g.ID,
			Name:        g.Name,
			Description: g.Description,
			CreatedAt:   timestamppb.New(g.CreatedAt),
		})
	}
	return out, nil
}

func (s *AuthServer) ValidateToken(
	ctx context.Context,
	req *pb.ValidateTokenRequest,
) (*pb.ValidateTokenResponse, error) {

	userDTO, err := s.validateToken.Execute(ctx, req.GetToken())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "validate token: %v", err)
	}

	return &pb.ValidateTokenResponse{
		User: &pb.User{
			Id:    userDTO.ID,
			Email: userDTO.Email,
		},
	}, nil
}
