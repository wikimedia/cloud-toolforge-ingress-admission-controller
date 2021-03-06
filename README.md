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

NOTE: Before you run ca-bundle.sh on MacOS, read the comments in that file and adjust accordingly

* `./get-cert.sh`  <-- creates a CSR and a secret with the TLS cert and key
* `./ca-bundle.sh` <-- places the correct ca-bundle in the service.yaml file
* `kubectl apply -f service.yaml`

As long as a suitable image can be placed where needed on toolforge, which can be done locally if
node affinity is used or some similar mechanism to prevent it being needed on every
spun-up node, the last three steps are likely all that is needed to bootstrap.

## Testing

At the top level, run `go test ./...` to capture all tests.  If you need to see output
or want to examine things more, use `go test -test.v ./...`

## Deploying

Since this was designed for use in [Toolforge](https://wikitech.wikimedia.org/wiki/Portal:Toolforge "Toolforge Portal"), so the instructions here focus on that.

The version of docker on the builder host is very old, so the builder/scratch pattern in
the Dockerfile won't work.

* Build the container image locally and copy it to the docker-builder host (currently tools-docker-builder-06.tools.eqiad.wmflabs). `$ docker build . -t docker-registry.tools.wmflabs.org/ingress-admission:latest`
* Then copy it over by saving it and using scp to get it on the docker-builder host `$ docker save -o saved_image.tar docker-registry.tools.wmflabs.org/ingress-admission:latest`
* Use scp or similar to transfer saved_image.tar from your local host to the docker builder.
* Load it into docker after copying the tar file to the builder host: `root@tools-docker-builder-06:~# docker load -i /home/bstorm/saved_image.tar`
* Push the image to the internal repo: `root@tools-docker-builder-06:~# docker push docker-registry.tools.wmflabs.org/ingress-admission:latest`
* On a control plane node as root (or as a cluster-admin user), with a checkout of the repo there somewhere (in a home directory is probably great), as root or admin user on Kubernetes, run `root@tools-k8s-control-1:# ./get-cert.sh`
* Then run `root@tools-k8s-control-1:# ./ca-bundle.sh`, which will insert the right ca-bundle in the service.yaml manifest.
* Now run `root@tools-k8s-control-1:# kubectl apply -f service.yaml` to launch it in the cluster.

## Updating the certs

Certificates created with the Kubernetes API are valid for one year. When upgrading Kubernetes (or whenever necessary)
it is wise to rotate the certs for this service. To do so simply run (as cluster admin or root@control host) `root@tools-k8s-control-1:# ./get-cert.sh`. That will recreate the cert secret. Then delete the existing pods to ensure
that the golang web services are serving the new cert.
