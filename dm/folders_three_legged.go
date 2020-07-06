package dm

import (
	"github.com/outer-labs/forge-api-go-client/oauth"
)

type FolderAPI3L struct {
	Auth          oauth.ThreeLeggedAuth
	Token         TokenRefresher
	FolderAPIPath string
}

func NewFolderAPI3LWithCredentials(
	auth oauth.ThreeLeggedAuth,
	token TokenRefresher,
) *FolderAPI3L {
	return &FolderAPI3L{
		Auth:          auth,
		Token:         token,
		FolderAPIPath: "/data/v1/projects",
	}
}

// Three legged Folder api calls
func (a FolderAPI3L) GetFolderDetailsThreeLegged(projectKey, folderKey string) (result ForgeResponseObject, err error) {
	if err = a.Token.RefreshTokenIfRequired(a.Auth); err != nil {
		return
	}

	path := a.Auth.Host + a.FolderAPIPath
	return getFolderDetails(path, projectKey, folderKey, a.Token.Bearer().AccessToken)
}

func (a FolderAPI3L) GetFolderContentsThreeLegged(projectKey, folderKey string) (result ForgeResponseArray, err error) {
	if err = a.Token.RefreshTokenIfRequired(a.Auth); err != nil {
		return
	}

	path := a.Auth.Host + a.FolderAPIPath
	return getFolderContents(path, projectKey, folderKey, a.Token.Bearer().AccessToken)
}

func (a FolderAPI3L) GetItemDetailsThreeLegged(projectKey, itemKey string) (result ForgeResponseObject, err error) {
	if err = a.Token.RefreshTokenIfRequired(a.Auth); err != nil {
		return
	}

	path := a.Auth.Host + a.FolderAPIPath
	return getItemDetails(path, projectKey, itemKey, a.Token.Bearer().AccessToken)
}
