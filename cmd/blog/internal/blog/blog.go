package blog

import (
	"net/http"

	"github.com/rs/zerolog"
	"go.seankhliao.com/mono/cmd/blog/internal"
	"go.seankhliao.com/mono/internal/envconf"
	"go.seankhliao.com/mono/internal/httpsvc"
)

var _ httpsvc.HTTPSvc = &Server{}

type Server struct {
	log zerolog.Logger

	host     string
	gtm      string
	notFound []byte

	mux *http.ServeMux
}

func (s Server) Desc() string {
	return `blog engine`
}

func (s Server) Help() string {
	return `
BLOG_GTM
        google tag manager id
BLOG_HOST
        hostname
`
}

func (s *Server) Init(log zerolog.Logger) error {
	s.log = log
	s.host = envconf.String("BLOG_HOST", "seankhliao.com")
	s.gtm = envconf.String("BLOG_GTM", "GTM-TLVN7D6")
	s.mux = http.NewServeMux()

	s.registerRedirects(s.mux)
	s.registerStatic(s.mux, internal.StaticFS)
	s.renderAndRegister(s.mux, internal.ContentFS)
	return nil
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	s.log.Debug().Str("user_agent", r.UserAgent()).Str("referrer", r.Referer()).Str("url", r.URL.String()).Msg("requested")
	s.mux.ServeHTTP(rw, r)
}
