package svc

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"runtime/debug"

	"github.com/nats-io/nats.go"
	natsproto "github.com/nats-io/nats.go/encoders/protobuf"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.seankhliao.com/mono/internal/flagwrap"
	"google.golang.org/grpc"
)

// S is the common interface that all services must implement
type S interface {
	// Register is called immediately, used to set up any flags to be read
	// and set up any default values from a zero value
	Register(Register) error
	// Init is called after flags have been parsed
	Init(Init) error
}

// SHttp is an optional interface for services to implement if they want to serve over http.
type SHTTP interface {
	S
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type SGRPC interface {
	S
	RegistergRPC(*grpc.Server)
}

type SBack interface {
	S
	Start() error
	Stop() error
}

type Register struct {
	Flags *flag.FlagSet
}

func newRegister() Register {
	return Register{
		Flags: flag.NewFlagSet("", flag.ContinueOnError),
	}
}

type Init struct {
	GRPCDialOpt   []grpc.DialOption
	HTTPTransport http.RoundTripper
	Logger        zerolog.Logger
	NATSConn      *nats.Conn
	NATSProtobuf  *nats.EncodedConn
}

func newInit(natsURL string, log zerolog.Logger) (Init, error) {
	var init Init

	init.GRPCDialOpt = []grpc.DialOption{
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
		grpc.WithUserAgent(userAgent()),
	}

	init.HTTPTransport = http.RoundTripper(http.DefaultTransport.(*http.Transport).Clone())
	init.HTTPTransport = newHttpTransport(init.HTTPTransport)
	init.HTTPTransport = otelhttp.NewTransport(init.HTTPTransport)

	init.Logger = log

	if natsURL != "" {
		var err error
		init.NATSConn, err = nats.Connect(natsURL)
		if err != nil {
			return init, fmt.Errorf("setup nats base conn: %w", err)
		} else {
			init.NATSProtobuf, err = nats.NewEncodedConn(init.NATSConn, natsproto.PROTOBUF_ENCODER)
			if err != nil {
				return init, fmt.Errorf("setup nats protobuf conn: %w", err)
			}
		}
	}

	return init, nil
}

type O struct {
	Args   []string
	Envs   []string
	Writer io.Writer
	Logger *zerolog.Logger
	Done   chan struct{}
}

func (o *O) init() {
	if o.Args == nil {
		o.Args = os.Args
	}
	if o.Envs == nil {
		o.Envs = os.Environ()
	}
	if o.Writer == nil {
		o.Writer = os.Stdout
	}
	if o.Logger == nil {
		log := zerolog.New(os.Stderr).With().Timestamp().Logger()
		o.Logger = &log
	}
	if o.Done == nil {
		o.Done = make(chan struct{})
	}
}

func Run(s S, o *O) int {
	if o == nil {
		o = &O{}
	}
	o.init()
	defer func() { close(o.Done) }()

	err := run(s, o)
	if err != nil {
		o.Logger.Error().Err(err).Msg("run svc")
		return 1
	}

	return 0
}

func run(s S, o *O) error {
	sHTTP, isHTTP := s.(SHTTP)
	sGRPC, isGRPC := s.(SGRPC)
	sBack, isBack := s.(SBack)

	// Register
	register := newRegister()
	err := s.Register(register)
	if err != nil {
		return fmt.Errorf("s.Register: %w", err)
	}

	var grpcSvr *grpcsvr
	if isGRPC {
		grpcSvr = newGrpcsvr(*o.Logger)
		grpcSvr.register(register.Flags, isHTTP)
	}
	var httpSvr *httpsvr
	if isHTTP {
		httpSvr = newHttpsvr(*o.Logger)
		httpSvr.register(register.Flags, isGRPC)
	}
	var natsURL, otlpURL string
	register.Flags.StringVar(&natsURL, "nats.url", "", "NATS address, ex nats://localhost:4222")
	register.Flags.StringVar(&otlpURL, "otlp.url", "", "OTLP address, ex grpc://localhost:4318")

	err = flagwrap.Parse(register.Flags, o.Args[1:])
	if errors.Is(err, flag.ErrHelp) {
		return nil
	} else if err != nil {
		return fmt.Errorf("parse flags: %w", err)
	}

	// Init
	init, err := newInit(natsURL, *o.Logger)
	if err != nil {
		return fmt.Errorf("setup init clients: %w", err)
	}

	var otelClient *otelclient
	if otlpURL != "" {
		otelClient, err = newOtelClient(otlpURL, init.GRPCDialOpt, *o.Logger)
		if err != nil {
			return fmt.Errorf("setup otel: %w", err)
		}
	}

	err = s.Init(init)
	if err != nil {
		return fmt.Errorf("s.Init: %w", err)
	}

	// Run
	rg := newRungroup(*o.Logger)

	if otelClient != nil {
		rg.add(otelClient.start, otelClient.stop)
	}
	if isBack {
		rg.add(sBack.Start, sBack.Stop)
	}
	if isGRPC {
		sGRPC.RegistergRPC(grpcSvr.svr)
		rg.add(grpcSvr.start, grpcSvr.stop)
	}
	if isHTTP {
		if init.NATSProtobuf != nil {
			httpSvr.pub = init.NATSProtobuf
		}
		httpSvr.sethandler(sHTTP)
		rg.add(httpSvr.start, httpSvr.stop)
	}

	return rg.wait()
}

func serviceName() (long, short, version string) {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		panic("no build info")
	}
	long = bi.Main.Path
	short = path.Base(long)
	version = bi.Main.Version
	return
}

func userAgent() string {
	long, short, version := serviceName()
	email := fmt.Sprintf("mono-internal-svc+%s@liao.dev", short)
	return fmt.Sprintf("%s/%s (%s, %s)", short, version, long, email)
}

type httpTransport struct {
	base      http.RoundTripper
	userAgent string
}

func newHttpTransport(base http.RoundTripper) http.RoundTripper {
	return &httpTransport{
		base:      base,
		userAgent: userAgent(),
	}
}

func (h *httpTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("user-agent", h.userAgent)
	return h.base.RoundTrip(r)
}
