// This package provides functionality to create microservices.
// `gokit` provides following functionalities out of the box.

/*
	- Accepts any mux/router which is compaitle with `http.Handler`
	- Provides some default middlewares
	- Allows you to pass custom middlewares
	- Provides graceful shutdown out of the box
	- Initializes default middlewares automatically for convinence
	- Pass data between middlewares through `context.Context`
*/
package kit
