package service

import (
	"auth-app/internal/middleware"
	"auth-app/internal/repository"
	"auth-app/pkg/proto/auth"
	"context"
	"errors"
	"fmt"
	"strconv"
)

type UserService interface {
	Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error)
	GetAllUsers(ctx context.Context, req *auth.GetUsersRequest) (*auth.GetUsersResponse, error)
	CreateUser(context.Context, *auth.CreateUserRequest) (*auth.CreateUserResponse, error)
	UpdateUser(context.Context, *auth.UpdateUserRequest) (*auth.UpdateUserResponse, error)
	DeleteUser(context.Context, *auth.DeleteUserRequest) (*auth.DeleteUserResponse, error)
	mustEmbedUnimplementedAuthServiceServer()
}

func (s *UserServiceImpl) mustEmbedUnimplementedAuthServiceServer() {
	//TODO implement me
	panic("implement me")
}

type UserServiceImpl struct {
	Repo repository.UserRepository
	auth.UnimplementedAuthServiceServer
}

func NewUserService(repo repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{Repo: repo}
}

func (s *UserServiceImpl) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	// Fetch user from database
	user, err := s.Repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Validate password
	if !middleware.ComparePassword(user.Password, req.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token := middleware.GenerateToken(&user)

	// Store session in Redis with expiration (e.g., 24 hours)
	err = s.Repo.SetRedis(ctx, "session:"+token, req.Email)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	// Return response
	return &auth.LoginResponse{
		AccessToken: token,
	}, nil
}

func (s *UserServiceImpl) GetAllUsers(ctx context.Context, req *auth.GetUsersRequest) (*auth.GetUsersResponse, error) {
	users, err := s.Repo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	var pbUsers []*auth.User
	for _, u := range users {
		var pbRoles []*auth.Role
		for _, role := range u.Roles {
			pbRoles = append(pbRoles, &auth.Role{
				RoleId:   strconv.Itoa(role.ID),
				RoleName: role.Name,
			})
		}

		pbUsers = append(pbUsers, &auth.User{
			Roles:      pbRoles,
			Email:      u.Email,
			LastAccess: u.LastAccess.String(),
		})
	}

	return &auth.GetUsersResponse{Users: pbUsers}, nil
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, req *auth.CreateUserRequest) (*auth.CreateUserResponse, error) {
	hasPermission, err := s.Repo.CheckRoleRight(req.RoleId, "r_create")
	if err != nil || !hasPermission {
		return nil, errors.New("unauthorized: missing create permissions")
	}

	userID, err := s.Repo.CreateUser(req.Name, req.Email, req.Password, req.RoleId)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &auth.CreateUserResponse{UserId: userID}, nil
}

func (s *UserServiceImpl) UpdateUser(ctx context.Context, req *auth.UpdateUserRequest) (*auth.UpdateUserResponse, error) {
	hasPermission, err := s.Repo.CheckRoleRight(req.RoleId, "r_update")
	if err != nil || !hasPermission {
		return nil, errors.New("unauthorized: missing update permissions")
	}

	err = s.Repo.UpdateUser(req.UserId, req.Name, req.Email, req.Password, req.RoleId)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &auth.UpdateUserResponse{Message: "User updated successfully"}, nil
}

func (s *UserServiceImpl) DeleteUser(ctx context.Context, req *auth.DeleteUserRequest) (*auth.DeleteUserResponse, error) {
	hasPermission, err := s.Repo.CheckRoleRight(req.RoleId, "r_delete")
	if err != nil || !hasPermission {
		return nil, errors.New("unauthorized: missing delete permissions")
	}

	err = s.Repo.DeleteUser(req.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}

	return &auth.DeleteUserResponse{Message: "User deleted successfully"}, nil
}
