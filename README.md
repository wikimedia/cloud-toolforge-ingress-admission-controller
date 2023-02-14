# Ingress Validation Webhook for Toolforge

This is a [Kubernetes Admission Validation Webhook](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#what-are-admission-webhooks) deployed to check that
users are not setting ingress values that could interfere with other users.  Since the webhook will compile with the versions it uses, it will work
as long as the admission v1 API is valid.

## Deploying to Toolforge
This project uses the [standard workflow](https://wikitech.wikimedia.org/wiki/Wikimedia_Cloud_Services_team/EnhancementProposals/Toolforge_Kubernetes_component_workflow_improvements):
1. Build the container image using the
    `wmcs.toolforge.k8s.component.build` cookbook.
2. Update the file for the project you're updating in `deployment/values`.
   Commit those changes to the repository and get it merged in Gerrit.
3. Use the `wmcs.toolforge.k8s.component.deploy` cookbook to deploy the updated
   image to the cluster.

## Local development
1. Start a local Toolforge cluster using [lima-kilo](https://gitlab.wikimedia.org/repos/cloud/toolforge/lima-kilo/).
2. Build the Docker image locally and load it to kind:
```shell-session
$ docker build -f Dockerfile -t ingress-admission:test . && kind load docker-image ingress-admission:test -n toolforge
```
3. Run the deploy script to start the service
```shell-session
$ ./deploy.sh local
```
4. After you've made changes, update the docker image and restart the running container:
```shell-session
$ docker build -f Dockerfile -t ingress-admission:test . && kind load docker-image ingress-admission:test -n toolforge && kubectl rollout restart -n ingress-admission deployment ingress-admission
```
