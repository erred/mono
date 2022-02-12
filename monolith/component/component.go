package component

import (
	"context"
	"flag"
	"net/http"

	"go.seankhliao.com/mono/monolith/o11y"
	"go.seankhliao.com/mono/monolith/run"
	"google.golang.org/grpc"
)

type Component interface {
	Enabled() bool
	Register(*flag.FlagSet)
}

type HTTP interface {
	Component
	HTTP(context.Context, o11y.ToolProvider, *http.ServeMux)
}

type GRPC interface {
	Component
	GRPC(context.Context, o11y.ToolProvider, *grpc.Server)
}

type Background interface {
	Component
	Background(context.Context, o11y.ToolProvider)
	run.Runner
}
