MINISHIFT_IP?=$(shell minishift ip)
MINISHIFT_HOSTNAME=minishift.local
MINISHIFT_HOSTNAME_REGEX='minishift\.local'
ETC_HOSTS=/etc/hosts

APP_NAMESPACE ?= $(NAMESPACE)
NAMESPACE ?= "devcluster"

.PHONY: login-as-admin
## Log in as system:admin
login-as-admin:
	$(Q)-echo "Logging using system:admin..."
	$(Q)-oc login -u system:admin

.PHONY: create-namespace
## Create the test namespace
create-namespace:
	$(Q)-echo "Creating Namespace"
	$(Q)-oc new-project $(NAMESPACE)
	$(Q)-echo "Switching to the namespace $(NAMESPACE)"
	$(Q)-oc project $(NAMESPACE)

.PHONY: use-namespace
## Log in as system:admin and enter the test namespace
use-namespace: login-as-admin
	$(Q)-echo "Using to the namespace $(NAMESPACE)"
	$(Q)-oc project $(NAMESPACE)

.PHONY: clean-namespace
## Delete the test namespace
clean-namespace:
	$(Q)-echo "Deleting Namespace"
	$(Q)-oc delete project $(NAMESPACE)

.PHONY: reset-namespace
## Delete an create the test namespace and deploy rbac there
reset-namespace: login-as-admin clean-namespace create-namespace

.PHONY: deploy-on-minishift
## Deploy DevCluster service on minishift
deploy-on-minishift: login-as-admin create-namespace build docker-image-dev
	$(Q)oc process -f ./deploy/devcluster.yaml \
        -p IMAGE=${IMAGE_DEV} \
        -p ENVIRONMENT=dev \
        -p NAMESPACE=${NAMESPACE} \
        | oc apply -f -

.PHONY: deploy-dev
## Deploy DevCluster on dev environment
deploy-dev: create-namespace build docker-image-dev
	$(Q)oc process -f ./deploy/devcluster.yaml \
        -p IMAGE=${IMAGE_DEV} \
        -p ENVIRONMENT=dev \
        -p NAMESPACE=${NAMESPACE} \
        | oc apply -f -
