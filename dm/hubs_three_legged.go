package dm

import "github.com/outer-labs/forge-api-go-client/oauth"

type HubAPI3L struct {
	Auth       oauth.ThreeLeggedAuth
	Token      TokenRefresher
	HubAPIPath string
}

func NewHubAPI3LWithCredentials(
	auth oauth.ThreeLeggedAuth,
	token TokenRefresher,
) *HubAPI3L {
	return &HubAPI3L{
		Auth:       auth,
		Token:      token,
		HubAPIPath: "/project/v1/hubs",
	}
}

// Hub functions for use with 3legged authentication
func (a *HubAPI3L) GetHubsThreeLegged() (result ForgeResponseArray, err error) {
	if err = a.Token.RefreshTokenIfRequired(a.Auth); err != nil {
		return
	}

	path := a.Auth.Host + a.HubAPIPath
	return getHubs(path, a.Token.Bearer().AccessToken)
}

func (a *HubAPI3L) GetHubDetailsThreeLegged(hubKey string) (result ForgeResponseObject, err error) {
	if err = a.Token.RefreshTokenIfRequired(a.Auth); err != nil {
		return
	}

	path := a.Auth.Host + a.HubAPIPath
	return getHubDetails(path, hubKey, a.Token.Bearer().AccessToken)
}

func (a *HubAPI3L) ListProjectsThreeLegged(hubKey string) (result ForgeResponseArray, err error) {
	if err = a.Token.RefreshTokenIfRequired(a.Auth); err != nil {
		return
	}

	path := a.Auth.Host + a.HubAPIPath
	return listProjects(path, hubKey, "", "", "", "", a.Token.Bearer().AccessToken)
}

func (a *HubAPI3L) GetProjectDetailsThreeLegged(hubKey, projectKey string) (result ForgeResponseObject, err error) {
	if err = a.Token.RefreshTokenIfRequired(a.Auth); err != nil {
		return
	}

	path := a.Auth.Host + a.HubAPIPath
	return getProjectDetails(path, hubKey, projectKey, a.Token.Bearer().AccessToken)
}

func (a *HubAPI3L) GetTopFoldersThreeLegged(hubKey, projectKey string) (result ForgeResponseArray, err error) {
	if err = a.Token.RefreshTokenIfRequired(a.Auth); err != nil {
		return
	}

	path := a.Auth.Host + a.HubAPIPath
	return getTopFolders(path, hubKey, projectKey, a.Token.Bearer().AccessToken)
}
