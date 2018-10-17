// Package forge is an opinionated Autodesk Forge SDK for the Go programming language.
//
// The Forge SDK for Go provides APIs that developers can use to build Go applications
// that use Autodesk Forge Services such as
// Data Management, Model Derivative, Reality Capture and others.
//
// The SDK removes the complexity of coding directly against a web service
// interface and it hides a lot of the lower-level plumbing, such as authentication.
//
// Getting More Information
//
// Checkout the https://developer.autodesk.com/ portal for overviews, tutorials and
// detailed documentation for each Autodesk Forge Service.
//
// Checkout LearnForge http://learnforge.autodesk.io for a step-by-step tutorial on
// building a Forge powered web application in different language, including Go using this library.
//
// Overview of SDK's Packages
//
// The SDK is composed of several parts, corresponding to each Forge Service, but all of them
// are relying on OAuth service for 2-legged and 3-legged authentication necessary to access
// Forge Services.
//   * oauth - provides common shared types such as Config, Logger,
//     and utilities to make working with API parameters easier.
package forge

