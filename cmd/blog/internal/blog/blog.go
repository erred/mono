package blog

import (
	"net/http"

	"github.com/rs/zerolog"
	"go.seankhliao.com/mono/cmd/blog/internal"
	"go.seankhliao.com/mono/internal/httpsvc"
	"go.seankhliao.com/mono/internal/webstatic"
)

var _ httpsvc.HTTPSvc = &Server{}

type Server struct {
	log zerolog.Logger

	host     string
	gtm      string
	notFound []byte

	mux *http.ServeMux
}

func (s *Server) Init(init *httpsvc.Init) error {
	s.log = init.Log
	init.Flags.StringVar(&s.host, "blog.host", "seankhliao.com", "canonical host for blog")
	init.Flags.StringVar(&s.gtm, "blog.gtm", "GTM-TLVN7D6", "GTM id for analytics")
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
