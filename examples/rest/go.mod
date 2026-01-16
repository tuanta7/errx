module examples/rest

go 1.25.3

require (
	github.com/tuanta7/errx v0.0.0
	golang.org/x/text v0.33.0
)

replace github.com/tuanta7/errx => ../..

require (
	golang.org/x/sys v0.40.0 // indirect
	google.golang.org/grpc v1.78.0 // indirect
)
