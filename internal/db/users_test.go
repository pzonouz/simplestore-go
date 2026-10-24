package sqlc

import (
	"simplestore_go/utils"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) User {
	ctx := t.Context()
	createUserParams := CreateUserParams{Email: pgtype.Text{String: utils.CreateRandomString(10), Valid: true}, PhoneNumber: pgtype.Text{String: utils.CreateRandomPhoneNumberAsString(), Valid: true}}
	user, err := queries.CreateUser(ctx, createUserParams)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, user.Email, createUserParams.Email)
	require.Equal(t, user.PhoneNumber, createUserParams.PhoneNumber)
	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)
	return user
}
func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestUpdateUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	user2Params := UpdateUserParams{ID: user1.ID, Email: pgtype.Text{String: utils.CreateRandomString(9), Valid: true}, PhoneNumber: pgtype.Text{String: utils.CreateRandomPhoneNumberAsString(), Valid: true}}
	user2, err := queries.UpdateUser(t.Context(), user2Params)
	require.Equal(t, user1.ID, user2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user2.Email, user2Params.Email)
	require.Equal(t, user2.PhoneNumber, user2Params.PhoneNumber)
	require.NotZero(t, user2.ID)
	require.NotZero(t, user2.CreatedAt)
}

func TestGetUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	user2, err := queries.GetUser(t.Context(), user1.ID)
	require.NoError(t, err)
	require.Equal(t, user2.ID, user1.ID)
	require.Equal(t, user2.Email, user1.Email)
	require.Equal(t, user2.PhoneNumber, user1.PhoneNumber)
}

func TestDeleteUser(t *testing.T) {
	user := CreateRandomUser(t)
	err := queries.DeleteUser(t.Context(), user.ID)
	require.NoError(t, err)
	user2, err := queries.GetUser(t.Context(), user.ID)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, user2)
}
func TestGetUsers(t *testing.T) {
	users := []User{}
	for i := 0; i < 10; i++ {
		user := CreateRandomUser(t)
		users = append(users, user)
	}
	fetchedUser, err := queries.ListUsers(t.Context(), ListUsersParams{Limit: 5})
	require.NoError(t, err)
	require.Len(t, fetchedUser, 5)
}
