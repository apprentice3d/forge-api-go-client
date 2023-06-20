// Package forge is an opinionated APS (formerly `Forge`) API client.
//
// The forge package provides APIs that developers can use to build Go applications
// that use Autodesk Platform Services such as Data Management and Model Derivative.
//
// The API removes the complexity of coding directly against a web service
// interface, and it hides a lot of the lower-level plumbing, such as authentication.
//
// # Getting More Information
//
// Checkout the https://aps.autodesk.com/ portal for overviews, tutorials and
// detailed documentation for each Autodesk Platform Service.
//
// Checkout https://tutorials.autodesk.io/ for step-by-step tutorials on
// building APS powered web application in different language.
//
// # Overview of the forge APIs Packages
//
// The API is composed of several parts, corresponding to each Autodesk Platform Service,
// but all of them are relying on OAuth service for 2-legged and 3-legged authentication
// necessary to access Autodesk Platform Services.
//   - oauth - provides 2-legged and 3-legged authentication
//   - dm - provides access to Data Management service
//   - md - provides access to Model Derivative service
package forge
