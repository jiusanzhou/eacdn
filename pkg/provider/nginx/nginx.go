package nginx

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"go.zoe.im/eacdn/pkg/core"
	"go.zoe.im/eacdn/pkg/provider"
	"go.zoe.im/x/sh"
)

const (
	nginxConfName = "nginx.conf"

	nginxDynamicSitesTplStart = `    # Load dynamic sites by EaCDN. DO NOT EDIT ANY THING!`
)

var (
	nginxDynamicSitesTpl = nginxDynamicSitesTplStart + "\n    include %s/*.conf;"
)

type nginx struct {
	provider.Options

	// for different os or auto search
	Command string `json:"command,omitempty" yaml:"command"`
	// nginx default conf located dir
	ConfDir string `json:"conf_dir,omitempty" yaml:"conf_dir"`
	// custom sites conf dir
	SitesConfDir string `json:"sites_conf_dir,omitempty" yaml:"sites_conf_dir"`
	// if we need to use sudo to run nginx command
	Sudo bool `json:"sudo,omitempty" yaml:"sudo"`

	// load the real nginx command
	nginxCmd string
	// check the nginx version
	nginxVerion string
	// find the nginx servered config nginx file directorry
	nginxConfDir  string
	nginxConfFile string
	// custom sites conf dir
	nginxSitesConfDir string

	// -t, -s reload

	_confMode os.FileMode
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
	err := sh.Run(s.genRunCmd(s.nginxCmd, "-v"), sh.StdIO(nil, nil, &buf))
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
		err = sh.Run(s.genRunCmd(s.nginxCmd, "-t"), sh.StdIO(nil, nil, &buf))
		if err != nil {
			return fmt.Errorf("get conf with nginx -t error: %v, %s", err, string(buf.Bytes()))
		}

		// nginx: configuration file /etc/nginx/conf/nginx.conf test is successful
		// nginx: configuration file /etc/nginx/conf/nginx.conf test failed
		// last line

		hr := string(buf.Bytes()[bytes.LastIndex(bytes.TrimSpace(buf.Bytes()), []byte("\n")):][27:])
		s.nginxConfFile = hr[:strings.Index(hr, " ")]

		s.nginxConfDir = filepath.Dir(s.nginxConfFile)
	} else {
		// conf file is dir  "nginx.conf"
		s.nginxConfFile = filepath.Join(s.nginxConfDir, "nginx.conf")
	}

	// finally check version and conf
	if s.nginxVerion == "" || s.nginxConfDir == "" {
		return fmt.Errorf("Init nginx failed, version: %s conf: %s", s.nginxVerion, s.nginxConfFile)
	}

	if s.SitesConfDir == "" {
		// set with default dir
		s.nginxSitesConfDir = filepath.Join(s.nginxConfDir, "eacdn-sites")
	}

	if _, err := os.Stat(s.SitesConfDir); os.IsNotExist(err) {
		// create directorry
		// why not in current processor to create dir?
		log.Println("I Create sites conf dir", s.nginxSitesConfDir)
		sh.Run(s.genRunCmd("mkdir", "-p", s.nginxSitesConfDir))
	}

	// TODO: make sure the nginx process is running

	// check if we need to add loader

	confdata, err := ioutil.ReadFile(s.nginxConfFile)
	if err != nil {
		log.Println("E Read nginx conf error:", err)
		return err
	}
	confstr := string(confdata)

	// TODO: if site dir not same remove the old one and replace it
	if strings.Index(confstr, nginxDynamicSitesTplStart) < 0 {
		// not exits, add dynamic load sites
		conflines := fmt.Sprintf(nginxDynamicSitesTpl, s.nginxSitesConfDir)
		// sed -i -z '0,/\n\n/s/\n\n/\n%s\n\n/' /etc/nginx/nginx.conf

		// replace \n\n and write to a new file
		tmpconf := filepath.Join(os.TempDir(), "nginx.conf")

		confstat, _ := os.Stat(s.nginxConfFile)
		s._confMode = confstat.Mode()
		err = ioutil.WriteFile(tmpconf, []byte(strings.Replace(confstr, "http {", "http {\n\n"+conflines+"\n", 1)), s._confMode)
		if err != nil {
			log.Println("E Create temp nginx conf error:", err)
			return err
		}

		// copy to
		if err := sh.Run(s.genRunCmd("cp", tmpconf, s.nginxConfFile)); err != nil {
			log.Println("E Copy nginx conf error:", err)
			return err
		}

		log.Println("I Modidy conf add EaCDN sites loader success")
	} else {
		log.Println("I No need to add EaCDN sites loader")
	}

	log.Printf("I Init nginx success, %s version: %s conf: %s\n", s.nginxCmd, s.nginxVerion, s.nginxConfFile)

	return s.reload()
}

// Create or update the site
func (s nginx) Create(site *core.Site, opts ...provider.Option) error {

	// TODO: lint the site
	if err := site.Lint(); err != nil {
		log.Println("E Site", site.Name, " lint error", err)
		return err
	}

	conf := s.genSiteConFile(site.Name)

	// if disabled we need to remove the site conf
	if _, err := os.Stat(conf); site.Disabled && err != nil && !os.IsNotExist(err) {
		// remove this file
		os.Remove(conf)
		return s.reload()
	}

	confcontent, err := s.genSiteConfContent(site)
	if err != nil {
		log.Println("E Generate site", site.Name, "conf content error", err)
		return err
	}

	tmpconf := filepath.Join(os.TempDir(), site.Host)
	// create the site conf
	err = ioutil.WriteFile(tmpconf, confcontent, 0644)
	if err != nil {
		log.Println("E Write conf", conf, "error", err)
		return err
	}

	// use copy
	err = sh.Run(s.genRunCmd("cp", "-rf", tmpconf, conf))
	if err != nil {
		fmt.Println("E Copy site", site.Name, "conf error", err)
		return err
	}

	// TODO: check the site available
	log.Println("I Enable site", site.Name, conf, "success")

	return s.reload()
}

func (s nginx) Update(site *core.Site, opts ...provider.Option) error {

	return nil
}

func init() {
	provider.Register(&nginx{}, "nginx")
}
