package rules

import (
	"bytes"
	"embed"
	"io/ioutil"
	"text/template"

	"github.com/ViaQ/logerr/v2/kverrors"
)

const (
	// LokiRulesConfigName is the name of the rules config file in the configmap
	LokiRulesConfigName = "rules.yaml"
)

var (
	//go:embed loki-rules-config.yaml
	lokiRulesConfigYAMLTmplFile embed.FS

	lokiRulesConfigYAMLTmpl = template.Must(template.ParseFS(lokiRulesConfigYAMLTmplFile, "loki-rules-config.yaml"))
)

// Build builds a loki stack configuration files
func Build(opts Options) (string, error) {
	// Build loki config yaml
	w := bytes.NewBuffer(nil)
	err := lokiRulesConfigYAMLTmpl.Execute(w, opts)
	if err != nil {
		return "", kverrors.Wrap(err, "failed to create loki rules configuration")
	}
	cfg, err := ioutil.ReadAll(w)
	if err != nil {
		return "", kverrors.Wrap(err, "failed to read configuration from buffer")
	}

	return string(cfg), nil
}
