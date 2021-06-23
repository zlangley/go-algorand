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

	"H4sIAAAAAAAC/+x9a3PcOJLgX8HVboQfV6ySH90z1kXHntru7tG12+2wNHN3a+ncKDKrCiMS4BCgpGqd",
	"/vtFJgASJMGq0mO917H7yVYRj0QiX8hMJG4mqSpKJUEaPTm8mZS84gUYqOgvnqaqliYRGf6VgU4rURqh",
	"5OTQf2PaVEKuJtOJwF9LbtaT6UTyAto22H86qeAftaggmxyaqobpRKdrKDgObDYltm5Guk5WKnFDHNkh",
	"jt9Nbrd84FlWgdZDKH+V+YYJmeZ1BsxUXGqe4ifNroRZM7MWmrnOTEimJDC1ZGbdacyWAvJMz/wi/1FD",
	"tQlW6SYfX9JtC2JSqRyGcL5VxUJI8FBBA1SzIcwolsGSGq25YTgDwuobGsU08Cpds6WqdoBqgQjhBVkX",
	"k8PPEw0yg4p2KwVxSf9dVgC/Q2J4tQIzOZ/GFrc0UCVGFJGlHTvsV6Dr3GhGbWmNK3EJkmGvGful1oYt",
	"gHHJPv34lr169eoNLqTgxkDmiGx0Ve3s4Zps98nhJOMG/OchrfF8pSous6Rp/+nHtzT/iVvgvq241hBn",
	"liP8wo7fjS3Ad4yQkJAGVrQPHerHHhGmaH9ewFJVsOee2MaPuinh/P+uu5Jyk65LJaSJ7Aujr8x+jsqw",
	"oPs2GdYA0GlfIqYqHPTzQfLm/ObF9MXB7T99Pkr+1f35zavbPZf/thl3BwaiDdO6qkCmm2RVASduWXM5",
	"xMcnRw96reo8Y2t+SZvPCxL1ri/DvlZ0XvK8RjoRaaWO8pXSjDsyymDJ69wwPzGrZY5iCkdz1M6EZmWl",
	"LkUG2RSl79VapGuWcm2HoHbsSuQ50mCtIRujtfjqtjDTbYgShOte+KAF/f+LjHZdOzAB1yQNkjRXGhKj",
	"dqgnr3G4zFioUFpdpe+mrNjpGhhNjh+ssiXcSaTpPN8wQ/uaMa4ZZ141TZlYso2q2RVtTi4uqL9bDWKt",
	"YIg02pyOHkXmHUPfABkR5C2UyoFLQp7nuyHK5FKs6go0u1qDWTudV4EuldTA1OLvkBrc9v9x8usHpir2",
	"C2jNV/CRpxcMZKqy8T12k8Y0+N+1wg0v9Krk6UVcXeeiEBGQf+HXoqgLJutiARXul9cPRrEKTF3JMYDs",
	"iDvorODXw0lPq1qmtLnttB1DDUlJ6DLnmxk7XrKCX393MHXgaMbznJUgMyFXzFzLUSMN594NXlKpWmZ7",
	"2DAGNyzQmrqEVCwFZKwZZQskbppd8Ah5N3hayyoAxw8yCk4zyw5wJFxHaAZZF7+wkq8gIJkZ+6uTXPTV",
	"qAuQjYBjiw19Kiu4FKrWTacRGGnq7ea1VAaSsoKliNDYiUMHSg/bxonXwhk4qZKGCwkZSl4CWhmwkmgU",
	"pmDC7YeZoYpecA3fvh5T4O3XPXd/qfq7vnXH99ptapRYlozoRfzqGDZuNnX673H4C+fWYpXYnwcbKVan",
	"qEqWIic183fcP4+GWpMQ6CDCKx4tVpKbuoLDM/kc/2IJOzFcZrzK8JfC/vRLnRtxIlb4U25/eq9WIj0R",
	"qxFkNrBGT1PUrbD/4HhxcYzg1jm3a4wxWNPg0pIqclRjIhjFSqiQemjtBB9KA1VCRWOO0UE47XbmMtfR",
	"U817pS7qMsR42jk2Lzbs+N3Y7HbMu3LOUXPWDo89p9f+KHTXHua6obQRIEc3t+TY8AI2FSC0PF3SP9dL",
	"Ini+rH7Hf8oyj206cpizBMhr4bwZn9xv+BNttD204CgipZ2ak34/vAkA+ucKlpPDyT/NW1fO3H7Vczcu",
	"zng7nRy14zz+TG1Pu77eSav9zIS0u0NNp/bQ+vjw4KhRSMiS7sHwfa7Si3vBUFbIaEbYfVzgOENOoeHZ",
	"GngGFcu44bP21GcNwRF6p45/oX50jIMqooN/pf/wnOFn5EJuvH2JtrXQaGWqwBOWoUlqFZ2dCRuQqaxY",
	"Ya1QhtbjnaB8205uNUgj8j87tJz3R4vszg/W8GXUwy8Cl94ea48WqrofvfQIQbL2sM44jtqY57jy7s5S",
	"07pMHH4iBr9t0Buo9Y8OxWqIof7wMVx1sHBi+L8BFjSO+hhY6A702FhQRSlyeAR+XXO9Hi4CLbBXL9nJ",
	"X46+efHyy8tvvkUToqzUquIFW2wMaPbU6RWmzSaHZ8OVkYCvcxMf/dvX/ojXHXcnhgjgZux9OOoUUDJY",
	"jDHr0EDo3lWbqpaPgEKoKlVF7AYiHaNSlSeXUOmocfPRtWCuBcohezDo/W6hZVdcO7sGMlbLDKpZDPN4",
	"ECSVbqDQuxSFHfr0Wra4cQPyquKbwQ7Y9UZW5+bdZ0+6yPfHD42GXGKuJctgUa9CHcWWlSoYZxl1JIH4",
	"QWVwYrip9SNIgXawFhjciBAEvlC1YZxJlSFDY+O4fBhxtpKXh5xTJhQ5Zm31zwLQfE95vVobhmalim1t",
	"2zHhqd2UhHSFHjmbNk4F28pOZx15eQU827AFgGRq4Q6A7mhKi+TkNzI+JOSkUwtWc2jpwFVWKgWtIUtc",
	"/GsnaL6d3WWzBU8EOAHczMK0Ykte3RNYowzPdwBKbWLgNuaEOzUPod5v+m0b2J883EZe4SHYUgHaLsjd",
	"ORgYQ+GeOLmEik6P/6b75ye57/bV5Uhsx2ngU1Eg+zLJpdKQKpnp6GA51ybZxbbYqGMm4AoCTolxKg08",
	"4sF4z7WxPgQhMzIZrbiheagPTTEO8KhGwZH/5pXJcOwU5aTUtW40i67LUlUGstgaJFxvmesDXDdzqWUw",
	"dqO+jGK1hl0jj2EpGN8hy67EIogb7wDwTrbh4ihegHpgE0VlB4gWEdsAOfGtAuyG/u0RQPB80fQkwhG6",
	"RzmNU3060UaVJfKfSWrZ9BtD04ltfWT+2rYdEhc3rVzPFODsxsPkIL+ymLWRjTVH245GZgW/QN1Elpr1",
	"JQxhRmZMtJApJNsoH9nyBFuFLLCDSUeMZBc7DWbrMUePfqNEN0oEO3ZhbMEjFvtH66I/bb1Dj2C0vAPD",
	"Ra4bw6SJA7SzUMign86BVmQFKUiTb5BWl6IqbNSN1Jn2v1mzJ3Oz2PhSy34yYxVc8SrzLYanpWAxiZAZ",
	"XMelK+/4RjK4ZiIO9LKZWRiW+piYDAeYRRndRhnTXGkhV4kNX+5Sak3U8YlmtRROgV1B5eBaQuXUrvHh",
	"u8QoH+LbBsc2VDjnzH2QgF3j01rg7G7pWJSXPiAjFiKtFLfBW0Rqb4GsgoIjdBRGdGp/fM5tyH5rv/tY",
	"svfhh7QbH9fT66iEaUj0ak2bhaK2j8SQ6vFoCxrGFrLK1YLnCRr8kGSQm52uNzxIwDtqifpapcPuXZDP",
	"zj7n2dnZOXuPbelsAewCNnMKqbN0zeUK2jhHyC/21ADXkNahaumhca+DoPOVdqHvHgWnk1KpPGmOvP24",
	"zEDd9PF+IdILyBjKK2IxpwWfdHcIJ2FPkcR1E7m6Wm+8CVmWICF7NmPsSDIoSrNx/pWexdObXD4x2+a/",
	"plmzmoLoXDJa5OxMxl0bNgT/QJ7yw2znJJuT9sCp7CDbJzLXI/GXil9RBAmHi/LnVu/oCfUMVN9AowdE",
	"ZaHYx4fwEyVq8c4ui4yOI6120/WiEJStFTSbouT0AfThCV+YGWOnJDvwgKXhEiqeUyqK9o5joVkh8KCu",
	"6zQFyA7PZNKBJFWFm/hp+18rls7qg4NXwA6e9ftog+aqO0taHuj3/Y4dTO0nQhf7jp1NziaDkSoo1CVk",
	"9jwW0rXttXPY/9KMeyZ/HQhmVvCNPcl5XmS6Xi5FKizSc4VyfaV6VqdU9AUqBA9QzWomzJRUGWGUrHW7",
	"Ly0DTqLW02P4fCKjop2OqhSlnQ+bdmlHM7jmKa6Sk5DZWIugobOhEWRUmYQDRF3QW2Z0QQDdkeP35Luh",
	"PLcOiO3wnfZcEB10BOQ62227D5ARhWAf9j9ipcJdFy5BymfR5EKbAZDOHUERoIYgI0pnxv63qlnKiX/L",
	"2kBztlMVHZjoII0zkI71czpLrcUQ5FCA9RDRl+fP+wt//tztudBsCVc+qxAb9tHx/LllAqXNgzmgR5rX",
	"xxEDihzzqE0jmeBrrteznU56Gncv33ww9PE7PyExk9akYnDhlVLLR1ityK6jNgtcx1bqdo7cbU80K/lm",
	"1LwuEcBIOhlUFzn58tWyR5HMyb+1KHHINvVlY6CTNvt/nv7L4eej5F958vtB8ua/zs9vXt8+ez748eXt",
	"d9/93+5Pr26/e/Yv/xwzXrQRi3jc5y9crxFSJzmu5bG0kVu0PMlht3F+ALX82nD3SAw302M+WNI+RPcx",
	"tiECTQnabKK5kzY35BE0TZt/RSTFdXiItqlhiOEgIcVZ9M7K6sXYuYZtx57GMRaOt4CVkHEPaLqG9IL8",
	"mjsdvR2hWqKVRYICeLpm7TAxTRWamT0VZI2Rwl/BGdBrtvgS9/QibJDzUqPCdO7eQuS5cK4kRIA0LFMI",
	"47vv2ZWqLtAIWnOZ5c6XAf+oQcfd9BeoUh82tZCMRrnDrKiZtkxqwwFbp95/LiK9+DwlVFpoJHZHoItN",
	"4G4NKQu15QLnWUIFMo05kvvZEJx8wn5fO5gOEeAB3IepA4bt57ac1GWZbx6DjWkgVoHzF+hOqEPbr2oZ",
	"ZsA7LaI32kAxjBbarl9GPBmfvKdxsHFK5kJCUigJm+ilLyHhF/oY33Y0MUY6k7E31rfvie3A3wOrO89e",
	"e/hA/JLkDkyKj00+/iNsfn/cXqA4zP0nLwXkJeMszQWFwZTUpqpTcyY5Odp7x+geWfjwwXjo5a1vEo/1",
	"REIxbqgzyTXisHG/RxMIlhARPz8C+AiMrlcr0L1jNVsCnEnXSkhymtJc5JVI7IaVUFGmx8y2xJPkkucU",
	"KfodKsUWtelqGUpRtidjG7XGaZhankluWA5cG/aLkKfXNJz3kHmakWBI6HssjHj4QIIWOokbRT/Zr2Qb",
	"ueWvnZ1E98XsZ287fG1jzsMey091kB+/c8fa43d0dmnj1QPYv1oQsxAyiRIZ6p5CSLqH0aMt9hRPYJ6A",
	"nrWRb7frZ9JcSySkS56LjJv7kUNfxA140XJHj2o6G9GLSfm1nsfcZSuVlDy9oFyyyUqYdb2YpaqY++P8",
	"fKWao/0841AoSd+yOS/FHPXw/PLFjqPVA+QVi4ir2+nESR396FmrbuDYgvpzNtFg/7dR7MlPP5yyudsp",
	"/cRml9uhgyzjiAfGXebuOANx8fY2qL1OcCbP5DtYCinw++GZzLjh8wXXItXzWkP1Pc+5TGG2UuyQuSHf",
	"ccPJh9yLjI1d2Cb/voOmrBe5SNlFqIpb1hwLrJydfUYCOTs7H+SODBWnmyoerKIJkith1qo2iYsujvuh",
	"deeU4+I622adMje2pUgXvXTjjwTQylInQUQlvvyyzHH5ARlqRp0o95hpoyovBFEyOp847u8H5bJnKn7l",
	"r6jVGjT7reDlZyHNOUuc//aoLClcQ/GS35ysQZrclLB/zKUFsR0sdkiihVuDCq5NxZOSr0BHl2+Al7T7",
	"pKgL8ojnOaNundiRz7ykodoFbI0RBHDcOS+eFndie/lgaHwJ9Im2kNqgdGpjW/fdLxzqLypHIrv3dgVj",
	"RHepNusEeTu6KjwTZn5nmlukK5TJzZFdrCQygbtwuwB7rIaMAvkU65p2uvt0KafhvOgQ2t6RtenvdJHL",
	"H9TqMuPOBuBy07+wosEYf43oE1zA5lS198DuckMFD/Y2OJ0gzYwxKlFqoIyQWEO29QHu3ua7XAUKIJcl",
	"szFae7PAk8VhQxe+zzgjWw35CEwc9W94NGyh95JXEURY4h9BwT0WiuM9iPSjEWFeGZGKsrm/tUeM+WOn",
	"Dw6yS7lE1Yla9rXGQKhHhZhtnMQdaWdnnwG/4H4gD/UzE/1MNkJgk04Y1VlxhLvIIciO0I6zeUVGl1+2",
	"LRwxBlqcSqCSrVb3YHQxEpoPa5fmIy7b5B5y3+6jaHcmVyAVeYeQ6IZRBc6bwyUfjWiPXnA8DpLqgnvz",
	"zfVFL9j6zDBtrrLaEjb+mqO/2+gvNE6md7qcOJ24PO/YdihJVkYGOay4C+BSBrlPHrKgPdHBBiEcvy6X",
	"uZDAklh+HtdapcLm9LSy3M0BaIQ+Z8w6eNjeI8TIOACbIl80MPugQt6Uq7sAKUFQqIz7sSlmFvwNuyNH",
	"bS0hZ97uNEOHsqNloml719du49ALNZ1ERdLYCaHTitkmCxgcqWIkiqJp6JcZen805EDqOOlI1uQi5q1D",
	"qwKIDE98t+DYwJ6KJSr5Z0EAtIKV0AbaczNyq3cEfV3fxaUykCxFpU1CR/bo8rDRj5qMwR+xaVz8dFDF",
	"bDESkcWlD017AZskE3kd320378/vcNoPzflJ14sL2JCSoZDHgornoBbqTI9ttkxtc1S3Lvi9XfB7/mjr",
	"3Y+WsClOXCllenP8QaiqJ0+2MVOEAGPEMdy1UZRuES9BVt1QtgRnMpv7R3mCs21eg2GU7K6ZiaOS144U",
	"XUtg6G5dhU1gtTmqQe2Z4X2pER7gZSmy694Z3o46EoInA/4Ohrq1+CNh5Ukz2A4MBOf1WEp+Bd7nYLc0",
	"0Jm2itAgbXk3ZvrJ0oFACKcS2tfAiwQ1YZNQVukuXJ0Cz3+Gzd+wLS1ncjudPOzIH8O1G3EHrj822xvF",
	"M/my7RGw48G7I8p5WVbqkueJc4yMkWalLh1pUnPvR/nKoi5+/D794ej9Rwc+ZWEDr1zy8bZVUbvyD7Mq",
	"PBHHMpBPA88IWav+7GwNsWDzm7oAoTPFJ4x3bDmUYo64LHu1jrKAFZ1zZRkPqe10lTifnl3iFt8elI1r",
	"rz0RW89e15vHL7nI/VHUQ7s7wf1eUqGTIf9Qr2CYLv+o4mbA3XHuaKlrh0wK59qS8VPYomKaKdlPEkQT",
	"kk64RKoF3yAFWef0UDjJukiQ/RKdizTutpALjcQhrc8XGzNqPGKM4oi1GAkhyFoEY2EzvUe0rAdkMEcU",
	"meRS2oK7hXLVYGsp/lEDExlIg58qlzTcYVTkS38PZqhO43du3MDu2k0z/ENsDBxqzLogILYbGKGHOXLj",
	"yx84/UIb1zj+EDgG7xCoCmccqMQtQSZHH46abbR/3fUUh8Vbh/IPCcMW+tpdOda7LdYW0JE5opVgR7XF",
	"0bimoLtU++uIViUQuKEysPntPNcqMkwtr7i0hR2xn8Wh661dAhn2ulIVXUDWEI3SC50sK/U7xE+yS9yo",
	"SB6zQyWZi9R7FrnY2ReijVemLdnr8RvCMUraY5Zc8JF1A4kjHE5UHrjO6WKGd3BxacnaFqHshK/jzBGm",
	"nMzt+C1zOJgHaTo5v1rwWLkjNKgQpqM2SNNxxRnFfGe/C7q5j+RoL4j3NG2FvbVbQtVeNhhWiLincfTH",
	"IvkMUlHwPG4lZYT9bupqJlbCVvKsNQSlIt1AtgSypSJXbtOGwVrUHC/ZwTQoRut2IxOXQotFDtTixbTN",
	"8aVbo+FNUpcYZUCatabmL/dovq5lVkFm1toiVivWGLD2gqD3fS/AXAFIdkDtXrxhT8nrr8UlPEMsOltk",
	"cvjiDaWl2D8OYsrOlezdJlcyEiz/0wmWOB1T2MOOgUrKjTqL3iC3ddbHRdgWbrJd9+Elaumk3m5eKrjk",
	"K4hHc4sdMNm+tJvkNOzhhTK/WQbaVGrDhInPD4ajfBpJTUPxZ8FgbRI1FRdWBdJTWwfSTuqHsxWHXekz",
	"D5f/SCGW0t8b7B2Yv66D2Ory2KopEPaBF9BF65RxW2iBrj66Ah1OIM7YsS/XQrXgmhJwFjc4Fy6dTDrc",
	"Qip5JaShQ1RtlsmfWbrmFU9R/M3GwP2y+Pb1CMjfvkage9Wv5N3WsOfJMiihChqqyzgWqxEK9oaB68ue",
	"SiWTAoVD9qzN6gwYLFqDShmex/NTvHDupydtH3pfWxJHSUYpp+5QDg+E7oNoSG4Z8IFU1axnlLTq7aR1",
	"50XemcjqKr7TvEaw/vrpvdP9hapiJbVaJmzuEphKwCVlvcTxjWM+EK1VPs6rdRShD1nIHXHaz0JtTOzG",
	"7vEcFrO0v69Fnv2tzR3v3R+quEzX0SDCAjt+aaseN3BZ7orfHuJSQh4dziqlL155RdTr39W+8xRC7tm2",
	"f9vELre3uBbwLpgeKD8holeYHCcIsdpNpm2yr/KVyhjN01b+aalkeGGeajlKg3T6lueERJ7nvy4nh5+3",
	"exmaXqoouMwmt9P+LvNqNXKdi1er2t5ONYqVnF7+cdp+WUt3fZ3nefxmgGsRH7rpjycc3l5mbw2R1AG+",
	"RyoALiCYcEjo5yH6HCIQLIFgFUL6Y5BDzOaDeyCiwdm+Sbmn9koFTRS+4bAgdnc3tpW/eTG2vmkz88gE",
	"9PG/MZdG8QSPqE/wwP4EUfkkfvIfGcw7g1uXVrS7e+4oTibOGuhlGnYKFsBICcpBEr1dN8F7voUJjqV9",
	"l+GxmECrukpHbhnYbwxFe5NFuTdxuoHjNOmrW9qLd5F6HfTB5iKT/xfP0rayJQOZ2TuUzNa3QLA6FQro",
	"BCiKOre33SFbQeUCE3WZK55NGY5z+sPRe2Zn1a5WE9VVoMqaK1srpSOYeqwQVP67S/GYsdTl/cfZnkuJ",
	"q9aGSn9pw4sydisFW5z6BnT1JYyF0NEoxM6MvbOnUu152E7S1ghizXTOgiIxj/8xhqdrYv7O4Whci+1f",
	"EtYrGh283dBUmW+K99myN0b5qrC2KOyUKRQhV0Lb92fgEroXYZpbYY7w/cWY7vKqWkpLKVHhse3W4n3Q",
	"7oGzCS8+XBKFrIf4O54QLPPetULuiWX5WA2NfrndwaMN9jZ9U5PcvyuWcqmkSKmCRfDiTQOye8tmH8Nx",
	"j2Ifff3qWdxxaIS5okV+m5Q6h8XRsr9eEJ6MCOHwK26qpQ77Jz3xQE7KFRjtJBtkU1/I2fkYhdTgijHS",
	"s0aBnFRVJz5LEjIa8m/Lsd2RjCgtfuSk+SN+++BcCpTKeiEkHVMc2lzWrPUC0lMbBs82wrCVgsAiC9f0",
	"GfvMqCxDBtfnM/80B41hw5u4bBvLHw515CP7LpKObd9iW0ahzPbnTgq+nfSoLN2kUTOi2eFYKepRBEci",
	"tIkPkQXIbcYPR9tCbltTckifIqHBJQX0oSQ9PCCMkRJnP1zyvHYVGqhSkk2Fi16dFDICxnshoX04JqIg",
	"0qhKoI0hfh3pp9OKG3uq20umnQLPKYofE2jauLDGQ4fqbTChhNbo5xjfxrYg+YjgaBq0ZzEuN817NUjd",
	"gTHxlh7Kcogclhcnq8oZURklO/cKjscEBwpuX6q/qwCGbDC0iWx3tDfhrppo7JJYJjTXGopFHknvfNd8",
	"DIruUx75YkP/xsp2jK/AJX3cuyAidbyzfbm9OGGOe59osbrnrrT9H3FbejwQ7lGM+n9AsRLeqx3UCrOC",
	"p7n2Sqltyj+BQoeK5sJWl2ZJ0EX9MO1rFtuPPePvUkxJNI4kuH5qKzpwK31t3GoszTUdzcrmxl25MJxt",
	"qxJqH5OIjWBzZOwjFvbFyqh3cCwvxqbF4OdB7/3shoEVNuLcCBDqE66GAP3sszlZyYULyrYsMsSsy/se",
	"ZuLvkxHabnB/ES6bmgaJreSeyc978d4QSxHGDtPWdpDnRQel9pZkz5JUFTwyagMVekfUDhPy9l0erYMo",
	"ptYwXOfeG9DB7Qju90F8KxciJa9G2dks9mHn+GUz7E7yxCLEX4ccSpOvJg06b+C4eWO7/rcx74E9IY/4",
	"nns4rUWe7drcTiShLTdCvnIfPvl3KXjyxXo+h+zmaj/cRfH3N4EQE1lrZ/JgqiBGsEd4wHWLBAPIIZvW",
	"lTAbyiv1lqb4Er2v8xNI9xKQe1ityc5xySG2cJcLRa2a1u07kT8p+zRSgeYvmYKGimD+cM2LMgfHF989",
	"WfwJXv35dXbw6sWfFn8++OYghdffvDk44G9e8xdvXr2Al3/+5vUBvFh++2bxMnv5+uXi9cvX337zJn31",
	"+sXi9bdv/vTEv4FoAW3fF/xfVBUoOfp4nJxSqa92a0rxM2xsHRAkY19hhKfEiVBwkZMvn376757DZqkq",
	"gnfl3a8TF6KcrI0p9eF8fnV1NQu7zFdUlj0xqk7Xcz/PsObox+PGQWuT0WhHm5cmbaTRkcIRffv0w8kp",
	"O/p4PGsJZnI4OZgdzF5QIa8SJC/F5HDyin4i7lnTvs8dsU0Ob26nk/kaeG7W7o8CTCVS/0lf8dUKqpkr",
	"tYI/Xb6ce//O/MZ58G+3fetmwLkrlUGH4E7+/KZT1z8Lx6Ub6/Mbnx0YfLLv1sxvyH00+vu8+xZovE0H",
	"1BtzLbLbua/B6Xq4NyLmN+2jLbeWg3KIeQd8Mem2ORWJprfstP0VmcbnxQjdfeOnoYDjDHcee71tHrAJ",
	"rj8dfv4P+gL/ee+9z5cHB//BXi58fccVb7V5O2fESK2k73nGfPyJ5n7x9eY+lnRLEoUes0L9djr55muu",
	"/lgiyfOcUcsgm3G49X+VF1JdSd8SNXBdFLzaeDbWHaHgn6UiOc9XmorXV+ISD/7n9DpCLOA3Ilzoicg7",
	"Cxd69/I/hcvXEi5/jAdBX96Rwf/4K/5PcfpHE6cnVtztL069KeeSMvR84SIcgw83HUuz+/s85Xk+v/Fp",
	"REE7mzsxt3Vl2599KYPh/f6uKT0m7N05iz0lJ62Eq2cu/8IOG6kV0cS6VWYdMr7uYC8FZzZQBp/coJ2y",
	"JD/DRu/SDKdrYL+54ROR/Ub3IijyMWWqYr/xPA9+o/px/swwiyuStn7AuBYZcH40hwzA39KgdE/3tAZq",
	"yAvwlSYsDjrR0WFCQVuldgkw9g6/LeYZikZH2y8ODg5iyYV9mJ3zyEJMt2KuVJLDJeTDrR4DoldwYoCx",
	"LdOfdiuuhnVCwkN/hOrouZUFtKVDYpDRqN3iF3eB7p2STwy74sK98BUkj9n3Sgth2AKWil4yNXUlXeJ5",
	"o3xiQEmV4JAxWNqLaw+1Cv54T2XcbpGiel2bTF3JccFF12557u6t0E2SxtdhFPMDNJJqxvy7+/mGlZW6",
	"FBkwTrlSqjatMwo7+xpSvReBmiqHKyFpAuJymsVe0OLB7QBX3n4oBE8cZB/sc5o9uRejHwdjnO9jTP9Q",
	"WhpaMFv3qvWBzG+CP27nN83KAz3WlCjr/D1HDkGz2b4unBBCh+4XAzyfu6Sa3q829B382H0kKPLrvLki",
	"Hf3YdyrFvjp/jm/UenND7yhtbOMX/XyO+0NXVdyet86+w/mcws1rpc2cEkS7jsDw43mzJTeeUPzW3J7f",
	"/r8AAAD//3SoFlcqlgAA",
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
