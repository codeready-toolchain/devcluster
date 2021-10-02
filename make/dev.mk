ETC_HOSTS=/etc/hosts

APP_NAMESPACE ?= $(NAMESPACE)
NAMESPACE ?= "devcluster"

.PHONY: create-namespace
## Create the test namespace
create-namespace:
	$(Q)-echo "Creating Namespace"
	$(Q)-oc new-project $(NAMESPACE)
	$(Q)-echo "Switching to the namespace $(NAMESPACE)"
	$(Q)-oc project $(NAMESPACE)

.PHONY: use-namespace
## Log in as system:admin and enter the test namespace
use-namespace:
	$(Q)-echo "Using to the namespace $(NAMESPACE)"
	$(Q)-oc project $(NAMESPACE)

.PHONY: clean-namespace
## Delete the test namespace
clean-namespace:
	$(Q)-echo "Deleting Namespace"
	$(Q)-oc delete project $(NAMESPACE)

.PHONY: reset-namespace
## Delete an create the test namespace and deploy rbac there
reset-namespace: clean-namespace create-namespace

.PHONY: deploy
## Deploy DevCluster
deploy: create-namespace build docker-image-dev docker-push-dev apply-resources print-route

.PHONY: apply-resources
## Apply DevCluster resources
apply-resources:
    # Try to create a secret. Will fail if it already exists. So we don't override the existing one.
	$(Q)-oc process -f ./deploy/secret.yaml \
        -p NAMESPACE=${NAMESPACE} \
        | oc create -f -
    # Try to create a config map. Will fail if it already exists. So we don't override the existing one.
	$(Q)-oc process -f ./deploy/configmap.yaml \
        -p NAMESPACE=${NAMESPACE} \
        | oc create -f -
	$(Q)oc process -f ./deploy/devcluster.yaml \
        -p IMAGE=${IMAGE_DEV} \
        -p NAMESPACE=${NAMESPACE} \
        | oc apply -f -

.PHONY: print-route
print-route:
	@echo "------------------------------------------------------------------"
	@echo "Deployment complete! Waiting for the devcluster service route."
	@echo -n "."
	@while [[ -z `oc get routes devcluster -n ${NAMESPACE} 2>/dev/null` ]]; do \
		if [[ $${NEXT_WAIT_TIME} -eq 100 ]]; then \
            echo ""; \
            echo "The timeout of waiting for the service route has been reached. Try to run 'make print-route' later or check the deployment logs"; \
            exit 1; \
		fi; \
		echo -n "."; \
		sleep 1; \
	done
	@echo ""
	$(eval ROUTE = $(shell oc get routes devcluster -n ${NAMESPACE} -o=jsonpath='{.spec.host}'))
	@echo Access the Landing Page here: https://${ROUTE}
	@echo "------------------------------------------------------------------"