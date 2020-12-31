package nginx

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"go.zoe.im/eacdn/pkg/core"
	"go.zoe.im/eacdn/pkg/provider"
	"go.zoe.im/x/sh"
)

const (
	nginxConfName = "nginx.conf"
)

type nginx struct {
	provider.Options

	// for different os???
	// or auto search
	Command     string `json:"command,omitempty" yaml:"command"`
	ConfDir     string `json:"conf_dir,omitempty" yaml:"conf_dir"`
	ConfDirName string `json:"conf_dir_name,omitempty" yaml:"conf_dir_name"`

	// load the real nginx command
	nginxCmd string
	// check the nginx version
	nginxVerion string
	// find the nginx servered config nginx file directorry
	nginxConfDir string

	// -t, -s reload
}

func (s *nginx) Init(opts ...provider.Option) error {
	// make sure nginx command exits
	if s.Command == "" {
		var buf bytes.Buffer
		// check if exits, search with which
		sh.Run("which nginx", sh.StdIO(nil, &buf, nil))
		s.nginxCmd = strings.TrimSpace(string(buf.Bytes()))
		if s.nginxCmd == "" {
			return errors.New("Can't auto load nginx command")
		}
	} else {
		// check file exits
		if _, err := os.Stat(s.Command); os.IsNotExist(err) {
			return err
		}
		s.nginxCmd = s.Command
	}

	// get the nginx version
	var buf bytes.Buffer
	err := sh.Run(s.nginxCmd+" -v", sh.StdIO(nil, nil, &buf))
	if err != nil {
		return fmt.Errorf("pre-check nginx error: %v", err)
	}

	s.nginxVerion, err = parseVersion(string(buf.Bytes()))
	if err != nil {
		return err
	}

	buf.Reset()

	// TODO: make sure nginx process alive

	// get the conf base direcrtory
	// use nginx -t to get conf file path

	if s.nginxConfDir == "" {
		err = sh.Run(s.nginxCmd+" -t", sh.StdIO(nil, nil, &buf))
		if err != nil {
			return fmt.Errorf("get conf with nginx -t error: %v", err)
		}

		// nginx: configuration file /etc/nginx/conf/nginx.conf test is successful
		// nginx: configuration file /etc/nginx/conf/nginx.conf test failed
		// last line

		hr := string(buf.Bytes()[bytes.LastIndex(bytes.TrimSpace(buf.Bytes()), []byte("\n")):][26:])
		hr = hr[:strings.Index(hr, " ")]
		fmt.Println("====>", hr)
	}

	log.Printf("I Init nginx success: %s %s %s\n", s.nginxCmd, s.nginxVerion, s.nginxConfDir)

	return nil
}

func (s nginx) Create(site *core.Site, opts ...provider.Option) error {
	log.Println("D Create site ====>", site)
	// modify nginx conf add or update

	// then reload
	return nil
}

func (s nginx) Update(site *core.Site, opts ...provider.Option) error {

	return nil
}

func init() {
	provider.Register(&nginx{}, "nginx")
}
