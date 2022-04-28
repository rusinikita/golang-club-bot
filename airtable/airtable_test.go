package airtable_test

import (
	"context"
	"log"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"

	"github.com/rusinikita/gogoClub/airtable"
	"github.com/rusinikita/gogoClub/request"
)

func Test(t *testing.T) {
	ctx := context.Background()

	userID := strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Int())

	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	db := airtable.New()

	testUser := request.Request{
		UserID: userID,
		Name:   "Bla bla",
	}
	require.NoError(t, db.Create(ctx, testUser))

	var users []request.Request

	require.NoError(t, db.List(ctx, &users, airtable.Filter(request.Request{UserID: userID})))
	require.Len(t, users, 1)

	found := users[0]
	require.NotEmpty(t, found.RecordID)

	testUser.RecordID = found.RecordID
	require.Equal(t, found, testUser)

	require.NoError(t, db.Patch(ctx, request.Request{RecordID: testUser.RecordID, Name: "Boba"}))
	require.NoError(t, db.List(ctx, &users, airtable.Filter(request.Request{UserID: userID})))
	require.Len(t, users, 1)
	require.Equal(t, users[0].Name, "Boba")

	require.NoError(t, db.Delete(ctx, testUser))
	users = nil
	require.NoError(t, db.List(ctx, &users, airtable.Filter(request.Request{UserID: userID})))
	require.Empty(t, users)
}
