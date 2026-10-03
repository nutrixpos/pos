package middlewares

import (
	"testing"
	"time"

	"github.com/nutrixpos/pos/modules/auth/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestJWT_GenerateAndValidate(t *testing.T) {
	j := NewJWTUtil("test-secret", 24)
	user := models.User{
		ID:       primitive.NewObjectID(),
		Username: "tester",
		Email:    "tester@example.com",
		Roles:    []string{"admin"},
	}

	token, err := j.GenerateToken(user)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := j.ValidateToken(token)
	require.NoError(t, err)
	assert.Equal(t, user.ID.Hex(), claims.UserID)
	assert.Equal(t, "tester", claims.Username)
	assert.Equal(t, "tester@example.com", claims.Email)
	assert.Equal(t, []string{"admin"}, claims.Roles)
	assert.Equal(t, "nutrix-pos", claims.Issuer)
	assert.WithinDuration(t, time.Now().Add(24*time.Hour), claims.ExpiresAt.Time, time.Minute)
}

func TestJWT_ValidateInvalid(t *testing.T) {
	j := NewJWTUtil("test-secret", 24)

	_, err := j.ValidateToken("not-a-token")
	assert.ErrorIs(t, err, ErrInvalidToken)

	other := NewJWTUtil("other-secret", 24)
	user := models.User{ID: primitive.NewObjectID(), Username: "tester"}
	token, err := other.GenerateToken(user)
	require.NoError(t, err)

	_, err = j.ValidateToken(token)
	assert.ErrorIs(t, err, ErrInvalidToken)
}

func TestJWT_RefreshToken(t *testing.T) {
	j := NewJWTUtil("test-secret", 24)
	user := models.User{
		ID:       primitive.NewObjectID(),
		Username: "tester",
		Email:    "tester@example.com",
	}

	token, err := j.GenerateToken(user)
	require.NoError(t, err)

	refreshed, err := j.RefreshToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, refreshed)

	claims, err := j.ValidateToken(refreshed)
	require.NoError(t, err)
	assert.Equal(t, user.ID.Hex(), claims.UserID)
}
