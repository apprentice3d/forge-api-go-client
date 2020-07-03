package md

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/outer-labs/forge-api-go-client/oauth"
)

var (
	// TranslationSVFPreset specifies the minimum necessary for translating a generic (single file, uncompressed)
	// model into svf.
	TranslationSVFPreset = TranslationParams{
		Output: OutputSpec{
			Destination: DestSpec{"us"},
			Formats: []FormatSpec{
				FormatSpec{
					"svf",
					[]string{"2d", "3d"},
				},
			},
		},
	}
)

// API struct holds all paths necessary to access Model Derivative API
type ModelDerivativeAPI3L struct {
	Auth            	oauth.ThreeLeggedAuth
	Token          		TokenRefresher
	ModelDerivativePath string
}

// NewAPIWithCredentials returns a Model Derivative API client with default configurations
func NewAPI3LWithCredentials(auth oauth.ThreeLeggedAuth, token *oauth.RefreshableToken,) *ModelDerivativeAPI3L {
	return &ModelDerivativeAPI3L{
		Auth:            auth,
		Token:           token,
		ModelDerivativePath: "/modelderivative/v2/designdata",
	}
}

func (a ModelDerivativeAPI3L) GetManifest3L(urn string) (result ManifestResult, err error) {
	if err = a.Token.RefreshTokenIfRequired(a.Auth); err != nil {
		return
	}

	path := a.Auth.Host + a.ModelDerivativePath
	result, err = getManifest(path, urn, a.Token.Bearer().AccessToken)

	return
}

func (a ModelDerivativeAPI3L) GetMetadata3L(urn string) (result MetadataResult, err error) {
	if err = a.Token.RefreshTokenIfRequired(a.Auth); err != nil {
		return
	}

	path := a.Auth.Host + a.ModelDerivativePath
	result, err = getMetadata(path, urn, a.Token.Bearer().AccessToken)

	return
}

func (a ModelDerivativeAPI3L) GetObjectTree3L(urn string, viewId string) (status int, result TreeResult, err error) {
	if err = a.Token.RefreshTokenIfRequired(a.Auth); err != nil {
		return
	}

	path := a.Auth.Host + a.ModelDerivativePath
	status, result, err = getObjectTree(path, urn, viewId, a.Token.Bearer().AccessToken)

	return
}

func (a ModelDerivativeAPI3L) GetPropertiesStream3L(urn string, viewId string) (status int,
	result io.ReadCloser, err error) {
	if err = a.Token.RefreshTokenIfRequired(a.Auth); err != nil {
		return
	}

	path := a.Auth.Host + a.ModelDerivativePath
	status, result, err = getPropertiesStream(path, urn, viewId, a.Token.Bearer().AccessToken)
	return
}

func (a ModelDerivativeAPI3L) GetThumbnail3L(urn string) (reader io.ReadCloser, err error) {
	if err = a.Token.RefreshTokenIfRequired(a.Auth); err != nil {
		return
	}

	path := a.Auth.Host + a.ModelDerivativePath
	reader, err = getThumbnail(path, urn, a.Token.Bearer().AccessToken)

	return
}