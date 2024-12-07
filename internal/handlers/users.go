package handlers

import (
	"context"

	genapi "github.com/h-okay/golang-api-example/gen/api"
	gendb "github.com/h-okay/golang-api-example/gen/db"
)

type UserHandlers struct {
	queries *gendb.Queries
}

func NewUserHandlers(queries *gendb.Queries) *UserHandlers {
	return &UserHandlers{queries: queries}
}

func (h *UserHandlers) ListUsers(ctx context.Context, request genapi.ListUsersRequestObject) (genapi.ListUsersResponseObject, error) {
	page := 1
	if request.Params.Page != nil {
		page = *request.Params.Page
	}

	limit := 10
	if request.Params.Limit != nil {
		limit = *request.Params.Limit
	}

	offset := (page - 1) * limit

	users, err := h.queries.ListUsers(ctx, gendb.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	// Convert db users to API users
	apiUsers := make([]genapi.User, len(users))
	for i, user := range users {
		apiUsers[i] = genapi.User{
			Id:        int(user.ID),
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: &user.CreatedAt.Time,
		}
	}

	return genapi.ListUsers200JSONResponse(apiUsers), nil
}

func (h *UserHandlers) CreateUser(ctx context.Context, request genapi.CreateUserRequestObject) (genapi.CreateUserResponseObject, error) {
	newUser, err := h.queries.CreateUser(ctx, gendb.CreateUserParams{
		Email: request.Body.Email,
		Name:  request.Body.Name,
	})
	if err != nil {
		return nil, err
	}

	return genapi.CreateUser201JSONResponse{
		Id:        int(newUser.ID),
		Email:     newUser.Email,
		Name:      newUser.Name,
		CreatedAt: &newUser.CreatedAt.Time,
	}, nil
}

func (h *UserHandlers) GetUserById(ctx context.Context, request genapi.GetUserByIdRequestObject) (genapi.GetUserByIdResponseObject, error) {
	user, err := h.queries.GetUserByID(ctx, int32(request.Id))
	if err != nil {
		return nil, err
	}

	return genapi.GetUserById200JSONResponse{
		Id:        int(user.ID),
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: &user.CreatedAt.Time,
	}, nil
}
