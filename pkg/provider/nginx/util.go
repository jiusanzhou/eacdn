package nginx

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"go.zoe.im/x/sh"

	"go.zoe.im/eacdn/pkg/core"
)

const (
	siteConfTpl = `
upstream {{ .Name }}_upstreams {
	{{ range .Upstreams }}server {{ .Dial }} weight={{ .Weight }};
	{{ end }}
}

server {
	listen {{ .Port }}; # need to change?
	listen [::]:{{ .Port }};  # need to change?
	server_name {{ .Host }};
	{{ if .Root }}
	root        {{ .Root }};
	{{ end }}
	# TODO: load more cache rules

	location / {
		proxy_pass http://{{ .Name }}_upstreams; # schema need to be changed?
		{{ if .HealthCheck }}health_check;{{ end }}
	}
}

# TODO: tls, cache, waf
`
)

var (
	siteConfTplElem *template.Template
)

func init() {
	siteConfTplElem = template.Must(template.New("_nginx_tpl").Funcs(sprig.TxtFuncMap()).Parse(siteConfTpl))
}

func (s *nginx) genRunCmd(c string, args ...string) string {
	if s.Sudo {
		c = "sudo " + c
	}
	return strings.Join(append([]string{c}, args...), " ")
}

func (s *nginx) genSiteConFile(name string) string {
	// replace . to -
	return filepath.Join(s.nginxSitesConfDir, name+".conf")
}

func (s *nginx) genSiteConfContent(site *core.Site) ([]byte, error) {
	var buf bytes.Buffer
	err := siteConfTplElem.Execute(&buf, site)
	return buf.Bytes(), err
}

func (s *nginx) reload() error {
	// TODO: read error from stderr
	err := sh.Run(s.genRunCmd(s.nginxCmd, "-t"), sh.StdIO(nil, nil, nil))
	if err != nil {
		return err
	}

	return sh.Run(s.genRunCmd(s.nginxCmd, "-s reload"))
}

func parseVersion(s string) (string, error) {
	parts := strings.Split(strings.TrimSpace(s), "/")
	if len(parts) != 2 {
		return "", fmt.Errorf("nginx version layout not corrent: %s", s)
	}
	return parts[1], nil
}
