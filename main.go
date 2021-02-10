package main

import (
	"gerrit.wikimedia.org/labs/tools/registry-admission-webhook/server"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

// Config is the general configuration of the webhook via env variables
type Config struct {
	ListenOn string   `default:"0.0.0.0:8080"`
	TLSCert  string   `default:"/etc/webhook/certs/cert.pem"`
	TLSKey   string   `default:"/etc/webhook/certs/key.pem"`
	Domains  []string `default:"toolforge.org,wmflabs.org,wmcloud.org,toolsbeta.wmflabs.org,toolsbeta.wmcloud.org"`
	Debug    bool     `default:"true"`
}

func main() {
	config := &Config{}
	envconfig.Process("", config)

	if config.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	logrus.Infoln(config)
	ingac := server.IngressAdmission{Domains: config.Domains}
	s := server.GetAdmissionValidationServer(&ingac, config.TLSCert, config.TLSKey, config.ListenOn)
	s.ListenAndServeTLS("", "")
}
