package da

import "github.com/apprentice3d/forge-api-go-client/oauth"

// API struct holds all paths necessary to access Design Automation API
type API struct {
	Authenticator oauth.ForgeAuthenticator
	DesignAutomationPath string
	UploadAppURL string
}

// NewAPI returns a DesignAutomation API client with default configurations
func NewAPI(authenticator oauth.ForgeAuthenticator) API {
	return API{
		authenticator,
		"/da/us-east/v3",
		"https://dasprod-store.s3.amazonaws.com",
	}
}


// UserId gives you the id used to identify the user
func (api API) UserId() (nickname string, err error) {
	bearer, err := api.Authenticator.GetToken("code:all")
	if err != nil {
		return
	}
	path := api.Authenticator.GetHostPath() + api.DesignAutomationPath
	nickname, err = getUserID(path, bearer.AccessToken)

	return
}



// EngineList lists all available Engines.
func (api API) EngineList() (list EngineList, err error) {

	bearer, err := api.Authenticator.GetToken("code:all")
	if err != nil {
		return
	}
	path := api.Authenticator.GetHostPath() + api.DesignAutomationPath
	list, err = listEngines(path, bearer.AccessToken)

	return
}

// EngineDetails gives details on an engine providing it's id.
func (api API) EngineDetails(id string) (list EngineDetails, err error) {

	bearer, err := api.Authenticator.GetToken("code:all")
	if err != nil {
		return
	}
	path := api.Authenticator.GetHostPath() + api.DesignAutomationPath
	list, err = getEngineDetails(path, id, bearer.AccessToken)

	return
}


// CreateApp creates an app with given name and using specified engine
// 	name - should be unique and will be the appID
// 	engine - engineId to be used by this app (check EngineList)
func (api API) CreateApp(name, engine string) (app AppBundle, err error) {

	bearer, err := api.Authenticator.GetToken("code:all")
	if err != nil {
		return
	}
	path := api.Authenticator.GetHostPath() + api.DesignAutomationPath
	app, err = createApp(path, name, engine, bearer.AccessToken)

	app.authenticator = api.Authenticator
	app.path = path
	app.name = name
	app.uploadURL = api.UploadAppURL

	//WARNING: when an AppBundle is created, it is assigned an '$LATEST' alias
	// but this alias is not usable and if no other alias is created for this
	// appBundle, then the alias listing will fail.
	// Thus I decided to autoasign a "default" alias upon app creation
	go app.CreateAlias("default", 1)

	return
}




// AppList lists all available appbundles.
func (api API) AppList() (list AppList, err error) {

	bearer, err := api.Authenticator.GetToken("code:all")
	if err != nil {
		return
	}
	path := api.Authenticator.GetHostPath() + api.DesignAutomationPath
	list, err = listApps(path, bearer.AccessToken)

	return
}

// CreateActivity creates an activity given an app
// 	name - should be unique and will be the appID
// 	engine - engineId to be used by this app (check EngineList)
func (api API) CreateActivity(config ActivityConfig) (activity Activity, err error) {

	bearer, err := api.Authenticator.GetToken("code:all")
	if err != nil {
		return
	}
	path := api.Authenticator.GetHostPath() + api.DesignAutomationPath
	activity, err = createActivity(path, config, bearer.AccessToken)

	activity.authenticator = api.Authenticator
	activity.path = path
	activity.name = config.ID

	//WARNING: when an Activity is created, it is assigned an '$LATEST' alias
	// but this alias is not usable and if no other alias is created for this
	// appBundle, then the alias listing will fail.
	// Thus I decided to autoasign a "default" alias upon app creation
	go activity.CreateAlias("default", 1)

	return
}




















// AppDelete will delete the app with specified id
//func (api API) AppDelete(id string) (err error) {
//
//	bearer, err := api.GetToken("code:all")
//	if err != nil {
//		return
//	}
//	path := api.Host + api.DesignAutomationPath
//	err = deleteApp(path, id, bearer.AccessToken)
//
//	return
//}




