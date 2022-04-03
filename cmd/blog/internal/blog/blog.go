package blog

import (
	"net/http"

	"github.com/rs/zerolog"
	"go.seankhliao.com/mono/cmd/blog/internal"
	"go.seankhliao.com/mono/internal/svc"
	"go.seankhliao.com/mono/internal/webstatic"
)

var _ svc.SHTTP = &Server{}

type Server struct {
	log zerolog.Logger

	host     string
	gtm      string
	notFound []byte

	mux *http.ServeMux
}

func (s *Server) Register(r svc.Register) error {
	r.Flags.StringVar(&s.host, "blog.host", "seankhliao.com", "canonical host for blog")
	r.Flags.StringVar(&s.gtm, "blog.gtm", "GTM-TLVN7D6", "GTM id for analytics")
	return nil
}

func (s *Server) Init(init svc.Init) error {
	s.log = init.Logger
	s.mux = http.NewServeMux()

	s.registerRedirects(s.mux)
	s.registerStatic(s.mux, internal.StaticFS)
	s.renderAndRegister(s.mux, internal.ContentFS)
	webstatic.Register(s.mux)
	return nil
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(rw, r)
}
