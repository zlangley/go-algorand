// Package private provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/algorand/oapi-codegen DO NOT EDIT.
package private

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"github.com/algorand/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Aborts a catchpoint catchup.
	// (DELETE /v2/catchup/{catchpoint})
	AbortCatchup(ctx echo.Context, catchpoint string) error
	// Starts a catchpoint catchup.
	// (POST /v2/catchup/{catchpoint})
	StartCatchup(ctx echo.Context, catchpoint string) error

	// (POST /v2/register-participation-keys/{address})
	RegisterParticipationKeys(ctx echo.Context, address string, params RegisterParticipationKeysParams) error

	// (POST /v2/shutdown)
	ShutdownNode(ctx echo.Context, params ShutdownNodeParams) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// AbortCatchup converts echo context to params.
func (w *ServerInterfaceWrapper) AbortCatchup(ctx echo.Context) error {

	validQueryParams := map[string]bool{
		"pretty": true,
	}

	// Check for unknown query parameters.
	for name, _ := range ctx.QueryParams() {
		if _, ok := validQueryParams[name]; !ok {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unknown parameter detected: %s", name))
		}
	}

	var err error
	// ------------- Path parameter "catchpoint" -------------
	var catchpoint string

	err = runtime.BindStyledParameter("simple", false, "catchpoint", ctx.Param("catchpoint"), &catchpoint)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter catchpoint: %s", err))
	}

	ctx.Set("api_key.Scopes", []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.AbortCatchup(ctx, catchpoint)
	return err
}

// StartCatchup converts echo context to params.
func (w *ServerInterfaceWrapper) StartCatchup(ctx echo.Context) error {

	validQueryParams := map[string]bool{
		"pretty": true,
	}

	// Check for unknown query parameters.
	for name, _ := range ctx.QueryParams() {
		if _, ok := validQueryParams[name]; !ok {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unknown parameter detected: %s", name))
		}
	}

	var err error
	// ------------- Path parameter "catchpoint" -------------
	var catchpoint string

	err = runtime.BindStyledParameter("simple", false, "catchpoint", ctx.Param("catchpoint"), &catchpoint)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter catchpoint: %s", err))
	}

	ctx.Set("api_key.Scopes", []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.StartCatchup(ctx, catchpoint)
	return err
}

// RegisterParticipationKeys converts echo context to params.
func (w *ServerInterfaceWrapper) RegisterParticipationKeys(ctx echo.Context) error {

	validQueryParams := map[string]bool{
		"pretty":           true,
		"fee":              true,
		"key-dilution":     true,
		"round-last-valid": true,
		"no-wait":          true,
	}

	// Check for unknown query parameters.
	for name, _ := range ctx.QueryParams() {
		if _, ok := validQueryParams[name]; !ok {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unknown parameter detected: %s", name))
		}
	}

	var err error
	// ------------- Path parameter "address" -------------
	var address string

	err = runtime.BindStyledParameter("simple", false, "address", ctx.Param("address"), &address)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter address: %s", err))
	}

	ctx.Set("api_key.Scopes", []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params RegisterParticipationKeysParams
	// ------------- Optional query parameter "fee" -------------
	if paramValue := ctx.QueryParam("fee"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "fee", ctx.QueryParams(), &params.Fee)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter fee: %s", err))
	}

	// ------------- Optional query parameter "key-dilution" -------------
	if paramValue := ctx.QueryParam("key-dilution"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "key-dilution", ctx.QueryParams(), &params.KeyDilution)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter key-dilution: %s", err))
	}

	// ------------- Optional query parameter "round-last-valid" -------------
	if paramValue := ctx.QueryParam("round-last-valid"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "round-last-valid", ctx.QueryParams(), &params.RoundLastValid)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter round-last-valid: %s", err))
	}

	// ------------- Optional query parameter "no-wait" -------------
	if paramValue := ctx.QueryParam("no-wait"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "no-wait", ctx.QueryParams(), &params.NoWait)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter no-wait: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.RegisterParticipationKeys(ctx, address, params)
	return err
}

// ShutdownNode converts echo context to params.
func (w *ServerInterfaceWrapper) ShutdownNode(ctx echo.Context) error {

	validQueryParams := map[string]bool{
		"pretty":  true,
		"timeout": true,
	}

	// Check for unknown query parameters.
	for name, _ := range ctx.QueryParams() {
		if _, ok := validQueryParams[name]; !ok {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unknown parameter detected: %s", name))
		}
	}

	var err error

	ctx.Set("api_key.Scopes", []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params ShutdownNodeParams
	// ------------- Optional query parameter "timeout" -------------
	if paramValue := ctx.QueryParam("timeout"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "timeout", ctx.QueryParams(), &params.Timeout)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter timeout: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.ShutdownNode(ctx, params)
	return err
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}, si ServerInterface, m ...echo.MiddlewareFunc) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.DELETE("/v2/catchup/:catchpoint", wrapper.AbortCatchup, m...)
	router.POST("/v2/catchup/:catchpoint", wrapper.StartCatchup, m...)
	router.POST("/v2/register-participation-keys/:address", wrapper.RegisterParticipationKeys, m...)
	router.POST("/v2/shutdown", wrapper.ShutdownNode, m...)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"t59v/38AAAD//wcatH9akwAA",
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file.
func GetSwagger() (*openapi3.Swagger, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromData(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error loading Swagger: %s", err)
	}
	return swagger, nil
}
