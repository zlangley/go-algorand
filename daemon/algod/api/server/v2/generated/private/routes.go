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
	"akjy2TOKJ83Np5v/FwAA//9nUHA0/JcAAA==",
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
