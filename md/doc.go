/*
Package md provides wrappers for the Model Derivative V2 REST API.
- https://aps.autodesk.com/model-derivative-api-2d-3d-conversions
- https://aps.autodesk.com/en/docs/model-derivative/v2/developers_guide/

The API offers the following features:
- Translate CAD file into viewables for rendering in the Viewer.
- Extract design metadata and integrate it into your app.
- Extract selected parts of a design and export the set of geometries in OBJ format.
- Translate designs into different formats, such as STL and OBJ.

To-Do:
- Update `GetDerivative` to use the new download endpoint (using signedcookies)
- Implement all endpoints
*/
package md
