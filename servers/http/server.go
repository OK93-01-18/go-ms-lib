package http

import (
	"net"
)

// Server - http server
type Server interface {
	Router

	// Listen serves HTTP requests from the given addr.
	Listen(string) error

	// Listener can be used to pass a custom listener.
	Listener(net.Listener) error

	// Shutdown gracefully shuts down the server without interrupting any active connections.
	// Shutdown works by first closing all open listeners and then waiting indefinitely for all connections to return to idle and then shut down.
	//
	// Make sure the program doesn't exit and waits instead for Shutdown to return.
	//
	// Shutdown does not close keepalive connections so its recommended to set ReadTimeout to something else than 0.
	Shutdown() error
}

type Router interface {

	// Get registers a route for GET methods that requests a representation
	// of the specified resource. Requests using GET should only retrieve data.
	Get(string, ...Handler) Router

	// Head registers a route for HEAD methods that asks for a response identical
	// to that of a GET request, but without the response body.
	Head(string, ...Handler) Router

	// Post registers a route for POST methods that is used to submit an entity to the
	// specified resource, often causing a change in state or side effects on the server.
	Post(string, ...Handler) Router

	// Options registers a route for OPTIONS methods that is used to describe the
	// communication options for the target resource.
	Options(string, ...Handler) Router

	// Delete registers a route for DELETE methods that deletes the specified resource.
	Delete(string, ...Handler) Router

	// Use registers a middleware route that will match requests
	// with the provided prefix (which is optional and defaults to "/").
	// This method will match all HTTP verbs: GET, POST, PUT, HEAD etc...
	Use(...interface{}) Router

	// Group is used for Routes with common prefix to define a new sub-router with optional middleware.
	//  api := app.Group("/api")
	//  api.Get("/users", handler)
	Group(string, ...Handler) Router
}

// Context represents the Context which hold the HTTP request and response.
type Context interface {
	// IP returns the remote IP address of the request.
	// Please use Config.EnableTrustedProxyCheck to prevent header spoofing, in case when your app is behind the proxy.
	IP() string

	// Hostname contains the hostname derived from the X-Forwarded-Host or Host HTTP header.
	// Returned value is only valid within the handler. Do not store any references.
	// Please use Config.EnableTrustedProxyCheck to prevent header spoofing, in case when your app is behind the proxy.
	Hostname() string

	// Query returns the query string parameter in the url.
	// Defaults to empty string "" if the query doesn't exist.
	// If a default value is given, it will return that value if the query doesn't exist.
	// Returned value is only valid within the handler. Do not store any references.
	Query(string, ...string) string

	// Set sets the response's HTTP header field to the specified key, value.
	Set(string, string)

	// Append the specified value to the HTTP response header field.
	// If the header is not already set, it creates the header with the specified value.
	Append(string, ...string)

	// Write appends p into response body.
	Write([]byte) (int, error)

	// Status sets the HTTP status for the response.
	// This method is chainable.
	Status(int) Context

	// GetReqHeaders returns the HTTP request headers.
	// Returned value is only valid within the handler. Do not store any references.
	GetReqHeaders() map[string]string

	// Request return Request interface struct of HTTP request
	Request() Request

	// Writef appends f & a into response body writer.
	Writef(string, ...interface{}) (int, error)

	// WriteString appends s to response body.
	WriteString(string) (int, error)

	// BodyParser binds the request body to a struct.
	// It supports decoding the following content types based on the Content-Type header:
	// application/json, application/xml, application/x-www-form-urlencoded, multipart/form-data
	// If none of the content types above are matched, it will return a ErrUnprocessableEntity error
	BodyParser(interface{}) error

	// Next executes the next method in the stack that matches the current route.
	Next() error

	// Redirect to the URL derived from the specified path, with specified status.
	// If status is not specified, status defaults to 302 Found.
	Redirect(string, int) error

	// Locals makes it possible to pass interface{} values under string keys scoped to the request
	// and therefore available to all following routes that match the request.
	Locals(string, ...interface{}) interface{}

	// Get returns the HTTP request header specified by field.
	// Field names are case-insensitive
	// Returned value is only valid within the handler. Do not store any references.
	// Make copies or use the Immutable setting instead.
	Get(string, ...string) string

	// Method contains a string corresponding to the HTTP method of the request: GET, POST, PUT and so on.
	Method(...string) string

	// Params is used to get the route parameters.
	// Defaults to empty string "" if the param doesn't exist.
	// If a default value is given, it will return that value if the param doesn't exist.
	// Returned value is only valid within the handler. Do not store any references.
	// Make copies or use the Immutable setting to use the value outside the Handler.
	Params(key string, defaultValue ...string) string
}

// Request - HTTP request
type Request interface {
	// GetContentLength returns content length
	GetContentLength() int

	// Body returns request body.
	Body() []byte

	// RequestURI returns request's URI.
	RequestURI() string
}

// Handler - handler of http request
type Handler func(context Context) error
