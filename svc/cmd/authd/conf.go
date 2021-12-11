package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"go.seankhliao.com/mono/proto/authdpb"
	"google.golang.org/protobuf/encoding/prototext"
)

// fromConfig initializes the authorization strategies from a config file
func (s *Server) fromConfig() error {
	b, err := os.ReadFile(s.configFile)
	if err != nil {
		return fmt.Errorf("read config file=%s: %w", s.configFile, err)
	}
	var config authdpb.Config
	err = prototext.Unmarshal(b, &config)
	if err != nil {
		return fmt.Errorf("unmarshal config: %w", err)
	}

	s.allow = make(map[string][]*regexp.Regexp)
	for host, allowed := range config.Allowlist {
		res := make([]*regexp.Regexp, 0, len(allowed.PathRe))
		for i, rre := range allowed.PathRe {
			re, err := regexp.Compile(rre)
			if err != nil {
				return fmt.Errorf("allowlist host=%s re=%d: %w", host, i, err)
			}
			res = append(res, re)
		}
		s.allow[host] = res
	}

	s.tokens = make(map[string]map[string]string)
	for host, tokens := range config.Tokens {
		tokmap := make(map[string]string, len(tokens.Tokens))
		for _, token := range tokens.Tokens {
			tokmap[token.Token] = token.Id
		}
		s.tokens[host] = tokmap
	}

	s.passwds = make(map[string][]byte)
	sc := bufio.NewScanner(strings.NewReader(config.HtpasswdFile))
	for idx := 0; sc.Scan(); idx++ {
		ss := sc.Text()
		if ss == "" {
			break
		}
		i := strings.Index(ss, ":")
		if i < 0 {
			return fmt.Errorf("htpasswd_file idx=%d no user:pass")
		}
		s.passwds[ss[:i]] = []byte(ss[i+1:])
	}
	for user, pass := range config.Htpasswd {
		s.passwds[user] = []byte(pass)
	}

	return nil
}
