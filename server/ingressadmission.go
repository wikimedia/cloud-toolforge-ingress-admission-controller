package server

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	"k8s.io/api/admission/v1beta1"
	netv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IngressAdmission type is where the project is stored and the handler method is linked
type IngressAdmission struct {
	Domains []string
}

// HandleAdmission is the logic of the whole webhook, really.  This is where
// the decision to allow a Kubernetes ingress update or create or not takes place.
func (ing *IngressAdmission) HandleAdmission(review *v1beta1.AdmissionReview) error {
	// logrus.Debugln(review.Request)
	req := review.Request
	var ingress netv1beta1.Ingress
	if err := json.Unmarshal(req.Object.Raw, &ingress); err != nil {
		logrus.Errorf("Could not unmarshal raw object: %v", err)
		review.Response = &v1beta1.AdmissionResponse{
			Result: &v1.Status{
				Message: err.Error(),
			},
		}
		return nil
	}
	logrus.Debugf("AdmissionReview for Kind=%v, Namespace=%v Name=%v (%v) UID=%v patchOperation=%v UserInfo=%v",
		req.Kind, req.Namespace, req.Name, ingress.Name, req.UID, req.Operation, req.UserInfo)

	// Whitelist kube-system
	if req.Namespace == "kube-system" {
		review.Response = &v1beta1.AdmissionResponse{
			Allowed: true,
			Result: &v1.Status{
				Message: "Welcome, admin!",
			},
		}
		return nil
	}
	domstr := strings.Join(ing.Domains, "|")

	for _, rule := range ingress.Spec.Rules {
		subdomRe := regexp.MustCompile(fmt.Sprintf("^%s\\.(%s)", req.Namespace[5:], domstr))
		pathRe := regexp.MustCompile(fmt.Sprintf("^tools\\.(%s)", domstr))
		toolPathRe := regexp.MustCompile(fmt.Sprintf("^/%s\\b", req.Namespace[5:]))
		logrus.Debugf("Found ingress host: %v", rule.Host)
		if rule.Host == "" || pathRe.MatchString(rule.Host) {
			for _, ingressPath := range rule.HTTP.Paths {
				logrus.Debugf("Found ingress path: %v", ingressPath.Path)
				if !toolPathRe.MatchString(ingressPath.Path) {
					review.Response = &v1beta1.AdmissionResponse{
						Allowed: false,
						Result: &v1.Status{
							Message: "Ingress path incorrect",
						},
					}
					return nil
				}
			}
		} else if !subdomRe.MatchString(rule.Host) {
			review.Response = &v1beta1.AdmissionResponse{
				Allowed: false,
				Result: &v1.Status{
					Message: "Ingress host must be <toolname>.toolforge.org or tools.wmflabs.org/<toolname>",
				},
			}
			return nil
		}
	}

	review.Response = &v1beta1.AdmissionResponse{
		Allowed: true,
		Result: &v1.Status{
			Message: "Welcome to the Toolforge!",
		},
	}
	return nil
}
