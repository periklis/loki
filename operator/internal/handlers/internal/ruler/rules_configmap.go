package ruler

import (
	"context"
	"crypto/sha1"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/grafana/loki/operator/internal/external/k8s"
	"github.com/grafana/loki/operator/internal/manifests"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GetRulesConfigMapSHA1 returns the SHA1 of all binary data of the rules configmap if any available or
// an empty string.
func GetRulesConfigMapSHA1(ctx context.Context, log logr.Logger, k k8s.Client, req ctrl.Request) string {
	var rulesCM corev1.ConfigMap
	key := client.ObjectKey{Name: manifests.RulesConfigMapName(req.Name), Namespace: req.Namespace}
	if err := k.Get(ctx, key, &rulesCM); err != nil {
		log.Error(err, "couldn't find")
		return ""
	}

	var c []byte
	for _, value := range rulesCM.Data {
		c = append(c, value...)
	}

	s := sha1.New()
	_, err := s.Write(c)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%x", s.Sum(nil))
}
