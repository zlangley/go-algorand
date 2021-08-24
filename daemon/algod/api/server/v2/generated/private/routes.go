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
	"AP//BUThIXaaAAA=",
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
