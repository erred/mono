package blog

import (
	"fmt"
	"net/http"
	"net/url"
)

type redirectLink struct {
	code     int
	path     string
	redirect string
}

func (s *Server) registerRedirects(mux *http.ServeMux) {
	redirects := []redirectLink{
		skhlRedirect("angelco", "profile", "seankhliao"),
		skhlRedirect("github", "profile", "erred"),
		skhlRedirect("github", "profile", "sean-dbk"),
		skhlRedirect("github", "profile", "seankhliao"),
		skhlRedirect("github", "readme", "erred"),
		skhlRedirect("github", "readme", "sean-dbk"),
		skhlRedirect("github", "readme", "seankhliao"),
		skhlRedirect("github", "site", "erred"),
		skhlRedirect("github", "site", "sean-dbk"),
		skhlRedirect("github", "site", "seankhliao"),
		skhlRedirect("google", "profile", "seankhliao"),
		skhlRedirect("instagram", "profile", "seankhliao"),
		skhlRedirect("linkedin", "about", "seankhliao"),
		skhlRedirect("linkedin", "profile", "seankhliao"),
		skhlRedirect("twitter", "profile", "seankhliao"),
	}

	for _, r := range redirects {
		mux.Handle(r.path, http.RedirectHandler(r.redirect, r.code))
	}
}

func skhlRedirect(source, medium, campaign string) redirectLink {
	p := fmt.Sprintf("/%s-%s-%s", shorten(source), shorten(medium), shorten(campaign))
	u := utm("https://seankhliao.com/", source, medium, campaign)
	return redirectLink{http.StatusTemporaryRedirect, p, u}
}

func shorten(source string) string {
	m := map[string]string{
		// source
		"angelco":   "ac",
		"github":    "gh",
		"google":    "g",
		"instagram": "ig",
		"linkedin":  "li",
		"twitter":   "tw",
		// medium
		"profile": "p",
		"readme":  "r",
		"site":    "s",
		// campaign
		"seankhliao": "s",
		"sean-dbk":   "sd",
		"erred":      "er",
	}

	s, ok := m[source]
	if ok {
		return s
	}
	return source
}

func utm(base, source, medium, campaign string) string {
	v := url.Values{
		"utm_source":   []string{source},
		"utm_medium":   []string{medium},
		"utm_campaign": []string{campaign},
	}
	return base + "?" + v.Encode()
}
