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
