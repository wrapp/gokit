// trace package provides a tracing http client. A trace client adds a unique request-id for outgoing
// http request. It requires a func as a parameter which generates the request-id. A service would want to
// generate a unique id for outgoing requests so that they can be traced in a large distributed systems
// where different services are communicating. This package uses requestidmw middleware to add the request-id.
// In addition to request-id this client also sets the User-Agent header for each outgoing request. This
// User-Agent is set to the value of SERVICE_NAME environment variable by default but it can be set to anything
// by calling SetUserAgent method of TraceClient.
package trace
