package main

import (
	"gerrit.wikimedia.org/cloud/toolforge/ingress-admission-controller/server"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

// Config is the general configuration of the webhook via env variables
type Config struct {
	ListenOn string   `default:"0.0.0.0:8080"`
	TLSCert  string   `default:"/etc/webhook/certs/tls.crt"`
	TLSKey   string   `default:"/etc/webhook/certs/tls.key"`
	Domains  []string `default:"toolforge.org,toolsbeta.wmflabs.org,toolsbeta.wmcloud.org"`
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
