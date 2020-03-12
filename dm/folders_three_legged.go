package dm

import (
	"time"
	"github.com/outer-labs/forge-api-go-client/oauth"
)

type FolderAPI3L struct {
	Auth        oauth.ThreeLeggedAuth
	BearerToken *oauth.Bearer
	HubAPIPath  string
}

func NewFolderAPI3LWithCredentials(
	auth oauth.ThreeLeggedAuth,
	bearer *oauth.Bearer,
) *FolderAPI3L {
	return &FolderAPI3L{
		Auth:        auth,
		BearerToken: bearer,
		HubAPIPath:  "/project/v1/hubs",
	}
}

// Three legged Folder api calls
func (a FolderAPI3L) GetFolderDetailsThreeLegged(bearer oauth.Bearer, projectKey, folderKey string) (result ForgeResponseObject, err error) {
	if err = a.refreshTokenIfRequired(); err != nil {
		return
	}

	path := a.Auth.Host + a.HubAPIPath
	return getFolderDetails(path, projectKey, folderKey, a.BearerToken.AccessToken)
}

func (a FolderAPI3L) GetFolderContentsThreeLegged(bearer oauth.Bearer, projectKey, folderKey string) (result ForgeResponseArray, err error) {
	if err = a.refreshTokenIfRequired(); err != nil {
		return
	}
	
	path := a.Auth.Host + a.HubAPIPath
	return getFolderContents(path, projectKey, folderKey, a.BearerToken.AccessToken)
}

func (a FolderAPI3L) GetItemDetailsThreeLegged(bearer oauth.Bearer, projectKey, itemKey string) (result ForgeResponseObject, err error) {
	if err = a.refreshTokenIfRequired(); err != nil {
		return
	}

	path := a.Auth.Host + a.HubAPIPath
	return getItemDetails(path, projectKey, itemKey, a.BearerToken.AccessToken)
}

func (a *FolderAPI3L) refreshTokenIfRequired() error {
	
	now := time.Now()
	expiryTime := a.BearerToken.ExpireTime
	
	if now.Before(expiryTime){
		return nil
	}
	
	refreshedBearer, err := a.Auth.RefreshToken(a.BearerToken.RefreshToken, "data:read")
	if err != nil {
		return err
	}

	newExpiryTime := now.Add(time.Second * time.Duration(refreshedBearer.ExpiresIn))
	expireTimeFormatted := newExpiryTime.Format(time.RFC3339)

	a.BearerToken.AccessToken = refreshedBearer.AccessToken
	a.BearerToken.ExpiresIn = refreshedBearer.ExpiresIn
	a.BearerToken.RefreshToken = refreshedBearer.RefreshToken
	a.BearerToken.TokenType = refreshedBearer.TokenType
	a.BearerToken.ExpireTime = expireTimeFormatted

	return nil	
}