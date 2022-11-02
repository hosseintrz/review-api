package test

import (
	"context"
	"database/sql"
	db2 "github.com/hosseintrz/suggestion_api/internal/db/sqlc"
	"github.com/hosseintrz/suggestion_api/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func CreateRandomUser(t *testing.T) db2.User {
	params := db2.CreateUserParams{
		Username: util.RandomUsername(),
		Password: util.RandomPassword(),
	}
	user, err := testQueries.CreateUser(context.Background(), params)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, params.Username, user.Username)
	require.Equal(t, params.Password, user.Password)
	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)

	return user
}

func DeleteUser(t *testing.T, id int64) {
	err := testQueries.DeleteUser(context.Background(), id)
	require.NoError(t, err)

	user, err := testQueries.GetUser(context.Background(), id)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, user)
}

func TestCreateDeleteUser(t *testing.T) {
	user := CreateRandomUser(t)
	DeleteUser(t, user.ID)
}

func TestGetUser(t *testing.T) {
	testUser := CreateRandomUser(t)
	defer DeleteUser(t, testUser.ID)
	user, err := testQueries.GetUser(context.Background(), testUser.ID)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, testUser.Username, user.Username)
	require.WithinDuration(t, testUser.CreatedAt, user.CreatedAt, time.Second)
	require.Equal(t, testUser.Password, user.Password)
}

func TestListUsers(t *testing.T) {
	testUsers := make([]db2.User, 0)
	k := 10
	for i := 0; i < k; i++ {
		user := CreateRandomUser(t)
		testUsers = append(testUsers, user)
	}
	defer func() {
		for _, u := range testUsers {
			err := testQueries.DeleteUser(context.Background(), u.ID)
			require.NoError(t, err)
		}
	}()
	users, err := testQueries.ListUsers(context.Background())
	require.NoError(t, err)
	for _, user := range users {
		require.NotEmpty(t, user)
	}
}
