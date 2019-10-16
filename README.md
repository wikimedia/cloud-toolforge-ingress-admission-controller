# Ingress Validation Webhook for Toolforge

This is a [Kubernetes Admission Validation Webhook](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#what-are-admission-webhooks) deployed to check that
users are not setting ingress values that could interfere with other users.  Since the webhook will compile with the versions it uses, it will work
as long as the admission v1beta1 API is valid. When that changes, this webhook should
begin to consume the new admission/v1 API whenever it happens.

## Use and development

This is pending adaption to Toolforge.  Currently it depends on local docker images and it
can be built and deployed on Kubernetes by insuring any node it is expected to run on
has access to the image it uses.  The image will need to be in a registry most likely when deployed.

It was developed using [Go Modules](https://github.com/golang/go/wiki/Modules), which will
validate the hash of every imported library during build.  At this time, it depends on
these external go libraries:

	* github.com/kelseyhightower/envconfig
	* github.com/sirupsen/logrus
	* k8s.io/api
	* k8s.io/apimachinery

To build on minikube and launch, follow these steps:

* `eval $(minikube docker-env)`
* `docker build -t ingress-admission:latest .`

That creates the image on minikube's docker daemon.  Then to launch the service:

* `./get-cert.sh`  <-- creates a CSR and a secret with the TLS cert and key
* `./ca-bundle.sh` <-- places the correct ca-bundle in the service.yaml file
* `kubectl create -f service.yaml`

As long as a suitable image can be placed where needed on toolforge, which can be done locally if
node affinity is used or some similar mechanism to prevent it being needed on every
spun-up node, the last three steps are likely all that is needed to bootstrap.

## Testing

At the top level, run `go test ./...` to capture all tests.  If you need to see output
or want to examine things more, use `go test -test.v ./...`

