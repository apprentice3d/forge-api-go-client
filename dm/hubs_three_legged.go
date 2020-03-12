package dm

import (
	// "fmt"
	"time"
	"github.com/outer-labs/forge-api-go-client/oauth"
)

type HubAPI3L struct {
	Auth        oauth.ThreeLeggedAuth
	BearerToken *oauth.Bearer
	HubAPIPath  string
	TokenExpireTime time.Time
}

func NewHubAPI3LWithCredentials(
	auth oauth.ThreeLeggedAuth,
	bearer *oauth.Bearer,
) *HubAPI3L {
	return &HubAPI3L{
		Auth:        auth,
		BearerToken: bearer,
		HubAPIPath:  "/project/v1/hubs",
		TokenExpireTime: time.Now(),
	}
}

// Hub functions for use with 3legged authentication
func (a *HubAPI3L) GetHubsThreeLegged() (result ForgeResponseArray, err error) {
	if err = a.refreshTokenIfRequired(); err != nil {
		return
	}

	path := a.Auth.Host + a.HubAPIPath
	return getHubs(path, a.BearerToken.AccessToken)
}

func (a *HubAPI3L) GetHubDetailsThreeLegged(hubKey string) (result ForgeResponseObject, err error) {
	if err = a.refreshTokenIfRequired(); err != nil {
		return
	}

	path := a.Auth.Host + a.HubAPIPath
	return getHubDetails(path, hubKey, a.BearerToken.AccessToken)
}

func (a *HubAPI3L) ListProjectsThreeLegged(hubKey string) (result ForgeResponseArray, err error) {
	if err = a.refreshTokenIfRequired(); err != nil {
		return
	}

	path := a.Auth.Host + a.HubAPIPath
	return listProjects(path, hubKey, "", "", "", "", a.BearerToken.AccessToken)
}

func (a *HubAPI3L) GetProjectDetailsThreeLegged(hubKey, projectKey string) (result ForgeResponseObject, err error) {
	if err = a.refreshTokenIfRequired(); err != nil {
		return
	}

	path := a.Auth.Host + a.HubAPIPath
	return getProjectDetails(path, hubKey, projectKey, a.BearerToken.AccessToken)
}

func (a *HubAPI3L) GetTopFoldersThreeLegged(hubKey, projectKey string) (result ForgeResponseArray, err error) {
	if err = a.refreshTokenIfRequired(); err != nil {
		return
	}

	path := a.Auth.Host + a.HubAPIPath
	return getTopFolders(path, hubKey, projectKey, a.BearerToken.AccessToken)
}

func (a *HubAPI3L) refreshTokenIfRequired() error {
	
	// Check if token has expired
	now := time.Now()
	expiryTime := a.TokenExpireTime
	if now.Before(expiryTime){
		return nil
	}

	refreshedBearer, err := a.Auth.RefreshToken(a.BearerToken.RefreshToken, "data:read")
	if err != nil {
		return err
	}

	// Refresh "now" and add new token expiration time to API struct along with new credentials
	now = time.Now()
	newExpiryTime := now.Add(time.Second * time.Duration(refreshedBearer.ExpiresIn))
	a.TokenExpireTime = newExpiryTime

	a.BearerToken.AccessToken = refreshedBearer.AccessToken
	a.BearerToken.ExpiresIn = refreshedBearer.ExpiresIn
	a.BearerToken.RefreshToken = refreshedBearer.RefreshToken
	a.BearerToken.TokenType = refreshedBearer.TokenType

	return nil
}
