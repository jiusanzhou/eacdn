package core

import (
	"errors"
	"net/url"
	"strings"

	"go.zoe.im/eacdn/pkg/utils"
)

// Site describes an HTTP server.
// Site define a instance of site need to be added in cdn
// [scheme:][//[userinfo@]host][/]path[?query][#fragment]
type Site struct {
	Provider string `json:"provider,omitempty" yaml:"provider"`

	Name string `json:"name,omitempty" yaml:"name"`

	// If true, this site disabled.
	Disabled bool `json:"disabled,omitempty" yaml:"disabled"`

	// Scheme string
	// User       *url.Userinfo // username and password information
	Host string `json:"host,omitempty" yaml:"host"` // host or host:port
	// Path       string    // path (relative paths may omit leading slash)
	// RawPath    string    // encoded path hint (see EscapedPath method)
	Port int `json:"port,omitempty" yaml:"port"`

	// TODO: do we need to move those timeouts to upstream?

	// How long to allow a read from a client's upload. Setting this
	// to a short, non-zero value can mitigate slowloris attacks, but
	// may also affect legitimately slow clients.
	ReadTimeout utils.Duration `json:"read_timeout,omitempty" yaml:"read_timeout"`

	// ReadHeaderTimeout is like ReadTimeout but for request headers.
	ReadHeaderTimeout utils.Duration `json:"read_header_timeout,omitempty" yaml:"read_header_timeout"`

	// WriteTimeout is how long to allow a write to a client. Note
	// that setting this to a small value when serving large files
	// may negatively affect legitimately slow clients.
	WriteTimeout utils.Duration `json:"write_timeout,omitempty" yaml:"write_timeout"`

	// Upstream host or url, only schema and host are usable
	Upstreams UpstreamPool `json:"upstreams,omitempty" yaml:"upstreams"`

	// TODO: root directory
	Root string `json:"root,omitempty" yaml:"root"`

	// AutoHTTPS configures or disables automatic HTTPS within this server.
	// HTTPS is enabled automatically and by default when qualifying names
	// are present in a Host matcher and/or when the server is listening
	// only on the HTTPS port.
	AutoHTTPS *AutoHTTPSConfig `json:"auto_https,omitempty" yaml:"auto_https"`

	// // How to handle TLS connections. At least one policy is
	// // required to enable HTTPS on this server if automatic
	// // HTTPS is disabled or does not apply.
	// TLSConnPolicies caddytls.ConnectionPolicies `json:"tls_connection_policies,omitempty"`

	// TODO: cache, waf

	HealthCheck bool `json:"health_check,omitempty" yaml:"health_check"`
}

// AutoHTTPSConfig is used to disable automatic HTTPS
// or certain aspects of it for a specific server.
// HTTPS is enabled automatically and by default when
// qualifying hostnames are available from the config.
type AutoHTTPSConfig struct {
	// If true, automatic HTTPS will be entirely disabled.
	Disabled bool `json:"disabled,omitempty" yaml:"disabled"`

	// If true, only automatic HTTP->HTTPS redirects will
	// be disabled.
	DisableRedir bool `json:"disable_redirect,omitempty" yaml:"disable_redirect"`

	// TODO: use different https provider
}

// UpstreamPool is a collection of upstreams.
type UpstreamPool []*Upstream

// Upstream bridges this proxy's configuration to the
// state of the backend host it is correlated with.
type Upstream struct {
	// The [network address](/docs/conventions#network-addresses)
	// to dial to connect to the upstream. Must represent precisely
	// one socket (i.e. no port ranges). A valid network address
	// either has a host and port or is a unix socket address.
	//
	// Placeholders may be used to make the upstream dynamic, but be
	// aware of the health check implications of this: a single
	// upstream that represents numerous (perhaps arbitrary) backends
	// can be considered down if one or enough of the arbitrary
	// backends is down. Also be aware of open proxy vulnerabilities.
	Dial string `json:"dial,omitempty" yaml:"dial"`

	// TODO: add more field for trip
	// Weight ...
	Weight int `json:"weight,omitempty" yaml:"weight"`

	// sets the number of unsuccessful attempts to communicate with
	// the server that should happen in the duration set by the
	// fail_timeout parameter to consider the server unavailable
	// for a duration also set by the fail_timeout parameter.
	MaxFails int `json:"max_fails,omitempty" yaml:"max_fails"`

	// the time during which the specified number of unsuccessful
	// attempts to communicate with the server should happen to
	// consider the server unavailable;
	// and the period of time the server will be considered unavailable.
	// By default, the parameter is set to 10 seconds.
	FailTimeout string `json:"fail_timeout,omitempty" yaml:"fail_timeout"`
}

// Lint check the upstream
func (u *Upstream) Lint() error {

	uri, err := url.Parse(u.Dial)
	if err != nil {
		return err
	}

	if uri.Scheme == "" {
		uri.Scheme = "http"
	}

	return nil
}

// Lint check the site
func (s *Site) Lint() error {

	if s.Root == "" && len(s.Upstreams) == 0 {
		return errors.New("must offer a root or upstream")
	}

	if s.Port == 0 {
		s.Port = 80
	}

	if s.Name == "" {
		s.Name = strings.ReplaceAll(s.Host, ".", "-")
	}

	for _, u := range s.Upstreams {
		if err := u.Lint(); err != nil {
			return err
		}
	}

	return nil
}
