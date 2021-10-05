package server

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	admissionv1 "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	AdmissionRequestFail = admissionv1.AdmissionReview{
		TypeMeta: v1.TypeMeta{
			Kind: "AdmissionReview",
		},
		Request: &admissionv1.AdmissionRequest{
			UID: "e911777a-c418-11e8-bbad-025000000001",
			Kind: v1.GroupVersionKind{
				Group: "networking.k8s.io", Version: "v1", Kind: "Ingress",
			},
			Namespace: "tool-tool2",
			Operation: "CREATE",
			Object: runtime.RawExtension{
				Raw: []byte(`{
					"apiVersion": "networking.k8s.io/v1",
					"kind": "Ingress",
					"metadata": {
					    "name": "tool2-ingress",
					    "namespace": "tool-tool2"
					},
					"spec": {
					    "rules": [
						{
						    "host": "tool1.wmflabs.org",
						    "http": {
							"paths": [
							    {
								"backend": {
								    "serviceName": "tool2-svc",
								    "servicePort": "8081"
								}
							    }
							]
						    }
						}
					    ]
					}
				    }`),
			},
		},
	}
	AdmissionRequestPass = admissionv1.AdmissionReview{
		TypeMeta: v1.TypeMeta{
			Kind: "AdmissionReview",
		},
		Request: &admissionv1.AdmissionRequest{
			UID: "e911857d-c318-11e8-bbad-025000000001",
			Kind: v1.GroupVersionKind{
				Group: "networking.k8s.io", Version: "v1", Kind: "Ingress",
			},
			Operation: "CREATE",
			Namespace: "tool-tool2",
			Object: runtime.RawExtension{
				Raw: []byte(`{
					"kind": "Ingress",
					"apiVersion": "networking.k8s.io/v1",
					"metadata": {
					    "name": "tool2-ingress",
					    "namespace": "tool-tool2",
					    "uid": "4b54be10-8d3c-11e9-8b7a-080027f5f85c",
					    "creationTimestamp": "2019-06-12T18:02:51Z"
					},
					"spec": {
					    "rules": [
						{
						    "host": "tool2.wmflabs.org",
						    "http": {
							"paths": [
							    {
								"backend": {
								    "service": {
										"name": "tool2-svc",
										"port": {
											"number": 8081
										}
									}
								}
							    }
							]
						    }
						}
					    ]
					}
				    }`),
			},
		},
	}
)

func decodeResponse(body io.ReadCloser) *admissionv1.AdmissionReview {
	response, _ := ioutil.ReadAll(body)
	review := &admissionv1.AdmissionReview{}
	codecs.UniversalDeserializer().Decode(response, nil, review)
	return review
}

func encodeRequest(review *admissionv1.AdmissionReview) []byte {
	ret, err := json.Marshal(review)
	if err != nil {
		logrus.Errorln(err)
	}
	return ret
}

func TestServeReturnsCorrectJson(t *testing.T) {
	inc := &IngressAdmission{}
	server := httptest.NewServer(GetAdmissionServerNoSSL(inc, ":8080").Handler)
	requestString := string(encodeRequest(&AdmissionRequestPass))
	myr := strings.NewReader(requestString)
	r, _ := http.Post(server.URL, "application/json", myr)
	review := decodeResponse(r.Body)
	t.Log(review.Response)
	if review.Request.UID != AdmissionRequestPass.Request.UID {
		t.Error("Request and response UID don't match")
	}
}
func TestHookFailsOnBadIngress(t *testing.T) {
	nsc := &IngressAdmission{}
	server := httptest.NewServer(GetAdmissionServerNoSSL(nsc, ":8080").Handler)
	requestString := string(encodeRequest(&AdmissionRequestFail))
	myr := strings.NewReader(requestString)
	r, _ := http.Post(server.URL, "application/json", myr)
	review := decodeResponse(r.Body)
	t.Log(review.Response)
	if review.Response.Allowed {
		t.Error("Allowed ingress that should not have been allowed!")
	}
}
func TestHookPassesOnRightIngress(t *testing.T) {
	nsc := &IngressAdmission{}
	server := httptest.NewServer(GetAdmissionServerNoSSL(nsc, ":8080").Handler)
	requestString := string(encodeRequest(&AdmissionRequestPass))
	myr := strings.NewReader(requestString)
	r, _ := http.Post(server.URL, "application/json", myr)
	review := decodeResponse(r.Body)
	t.Log(review.Response)
	if !review.Response.Allowed {
		t.Error("Failed to allow ingress should have been allowed!")
	}
}
