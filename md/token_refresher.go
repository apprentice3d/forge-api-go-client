package md

import "github.com/outer-labs/forge-api-go-client/oauth"

type TokenRefresher interface {
	Bearer() *oauth.Bearer
	RefreshTokenIfRequired(auth oauth.ThreeLeggedAuth) error
}
