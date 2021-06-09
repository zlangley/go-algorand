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
