package oauth

import (
	"sync"
	"time"
)

type RefreshableToken struct {
	bearer          *Bearer
	TokenExpireTime time.Time
	readMutex       sync.Mutex
	writeMutex      sync.Mutex
}

func NewRefreshableToken(bearer *Bearer, expiryTime time.Time) *RefreshableToken {
	return &RefreshableToken{
		bearer:          bearer,
		TokenExpireTime: expiryTime,
	}
}

func (t *RefreshableToken) RefreshTokenIfRequired(auth ThreeLeggedAuth) error {
	t.writeMutex.Lock()
	defer t.writeMutex.Unlock()

	// Check if token has expired
	now := time.Now()
	expiryTime := t.TokenExpireTime
	if now.Before(expiryTime) {
		return nil
	}

	refreshedBearer, err := auth.RefreshToken(t.bearer.RefreshToken, "data:read")
	if err != nil {
		return err
	}

	// Refresh "now" and add new token expiration time to API struct along with new credentials
	now = time.Now()
	newExpiryTime := now.Add(time.Second * time.Duration(refreshedBearer.ExpiresIn))
	t.TokenExpireTime = newExpiryTime

	t.bearer.AccessToken = refreshedBearer.AccessToken
	t.bearer.ExpiresIn = refreshedBearer.ExpiresIn
	t.bearer.RefreshToken = refreshedBearer.RefreshToken
	t.bearer.TokenType = refreshedBearer.TokenType

	return nil
}

func (t *RefreshableToken) Bearer() *Bearer {
	t.readMutex.Lock()
	defer t.readMutex.Unlock()
	return t.bearer
}
