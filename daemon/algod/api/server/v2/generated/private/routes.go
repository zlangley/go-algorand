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
