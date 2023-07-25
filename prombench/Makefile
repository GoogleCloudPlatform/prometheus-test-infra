INFRA_CMD                   ?= ../infra/infra

PROVIDER                    ?= gke

PROMETHEUS_IMAGE_REPOSITORY ?= gke.gcr.io/prometheus-engine/prometheus

PROMETHEUS_IMAGE_VERSION    ?= v2.35.0-gmp.7-gke.0

NODE_RETENTION_TIME_PERIOD  ?= "3 days"

.PHONY: deploy clean
deploy: node_create resource_apply
# GCP sometimes takes longer than 30 tries when trying to delete nodes
# if k8s resources are not already cleared
clean: resource_delete node_delete

cluster_create:
	${INFRA_CMD} ${PROVIDER} cluster create \
		-v ZONE:${ZONE} -v GKE_PROJECT_ID:${GKE_PROJECT_ID} \
		-v EKS_WORKER_ROLE_ARN:${EKS_WORKER_ROLE_ARN} -v EKS_CLUSTER_ROLE_ARN:${EKS_CLUSTER_ROLE_ARN} \
		-v EKS_SUBNET_IDS:${EKS_SUBNET_IDS} -v SEPARATOR:${SEPARATOR} \
		-v CLUSTER_NAME:${CLUSTER_NAME} -v PR_NUMBER:${PR_NUMBER} \
		-f manifests/cluster_${PROVIDER}.yaml

cluster_resource_apply:
	${INFRA_CMD} ${PROVIDER} resource apply \
		-v ZONE:${ZONE} -v GKE_PROJECT_ID:${GKE_PROJECT_ID} \
		-v EKS_WORKER_ROLE_ARN:${EKS_WORKER_ROLE_ARN} -v EKS_CLUSTER_ROLE_ARN:${EKS_CLUSTER_ROLE_ARN} \
		-v EKS_SUBNET_IDS:${EKS_SUBNET_IDS} -v SEPARATOR:${SEPARATOR} \
		-v CLUSTER_NAME:${CLUSTER_NAME} -v PR_NUMBER:${PR_NUMBER} -v DOMAIN_NAME:${DOMAIN_NAME} -v RELEASE:${RELEASE} \
		-v GRAFANA_ADMIN_PASSWORD:${GRAFANA_ADMIN_PASSWORD} \
		-v SERVICEACCOUNT_CLIENT_EMAIL:${SERVICEACCOUNT_CLIENT_EMAIL} \
		-v OAUTH_TOKEN="$(printf ${OAUTH_TOKEN} | base64 -w 0)" \
		-v WH_SECRET="$(printf ${WH_SECRET} | base64 -w 0)" \
		-v GITHUB_ORG:${GITHUB_ORG} -v GITHUB_REPO:${GITHUB_REPO} \
		-v PROMETHEUS_IMAGE_REPOSITORY:${PROMETHEUS_IMAGE_REPOSITORY} \
		-v PROMETHEUS_IMAGE_VERSION:${PROMETHEUS_IMAGE_VERSION} \
		-f manifests/cluster-infra

cluster_delete:
	${INFRA_CMD} ${PROVIDER} cluster delete \
		-v ZONE:${ZONE} -v GKE_PROJECT_ID:${GKE_PROJECT_ID} \
		-v EKS_WORKER_ROLE_ARN:${EKS_WORKER_ROLE_ARN} -v EKS_CLUSTER_ROLE_ARN:${EKS_CLUSTER_ROLE_ARN} \
		-v EKS_SUBNET_IDS:${EKS_SUBNET_IDS} -v SEPARATOR:${SEPARATOR} \
		-v CLUSTER_NAME:${CLUSTER_NAME} -v PR_NUMBER:${PR_NUMBER} \
		-f manifests/cluster_${PROVIDER}.yaml

node_create:
	${INFRA_CMD} ${PROVIDER} nodes create \
		-v ZONE:${ZONE} -v GKE_PROJECT_ID:${GKE_PROJECT_ID} \
		-v EKS_WORKER_ROLE_ARN:${EKS_WORKER_ROLE_ARN} -v EKS_CLUSTER_ROLE_ARN:${EKS_CLUSTER_ROLE_ARN} \
		-v EKS_SUBNET_IDS:${EKS_SUBNET_IDS} \
		-v CLUSTER_NAME:${CLUSTER_NAME} -v PR_NUMBER:${PR_NUMBER} \
		-f manifests/prombench/nodes_${PROVIDER}.yaml

resource_apply:
	$(INFRA_CMD) ${PROVIDER} resource apply \
		-v ZONE:${ZONE} -v GKE_PROJECT_ID:${GKE_PROJECT_ID} \
		-v CLUSTER_NAME:${CLUSTER_NAME} \
		-v PR_NUMBER:${PR_NUMBER} -v RELEASE:${RELEASE} -v DOMAIN_NAME:${DOMAIN_NAME} \
		-v GITHUB_ORG:${GITHUB_ORG} -v GITHUB_REPO:${GITHUB_REPO} \
		-v PROMETHEUS_IMAGE_REPOSITORY:${PROMETHEUS_IMAGE_REPOSITORY} \
		-v PROMETHEUS_IMAGE_VERSION:${PROMETHEUS_IMAGE_VERSION} \
		-v NODE_RETENTION_TIME_PERIOD:${NODE_RETENTION_TIME_PERIOD} \
		-f manifests/prombench/benchmark

# Required because namespace and cluster-role are not part of the created nodes
resource_delete:
	$(INFRA_CMD) ${PROVIDER} resource delete \
		-v ZONE:${ZONE} -v GKE_PROJECT_ID:${GKE_PROJECT_ID} \
		-v CLUSTER_NAME:${CLUSTER_NAME} -v PR_NUMBER:${PR_NUMBER} \
		-f manifests/prombench/benchmark/1c_cluster-role-binding.yaml \
		-f manifests/prombench/benchmark/1a_namespace.yaml \
		-f manifests/prombench/benchmark/7_delete_old_nodes.yaml

node_delete:
	$(INFRA_CMD) ${PROVIDER} nodes delete \
		-v ZONE:${ZONE} -v GKE_PROJECT_ID:${GKE_PROJECT_ID} \
		-v EKS_WORKER_ROLE_ARN:${EKS_WORKER_ROLE_ARN} -v EKS_CLUSTER_ROLE_ARN:${EKS_CLUSTER_ROLE_ARN} \
		-v EKS_SUBNET_IDS:${EKS_SUBNET_IDS} \
		-v CLUSTER_NAME:${CLUSTER_NAME} -v PR_NUMBER:${PR_NUMBER} \
		-f manifests/prombench/nodes_${PROVIDER}.yaml

all_nodes_running:
	$(INFRA_CMD) ${PROVIDER} nodes check-running \
		-v ZONE:${ZONE} -v GKE_PROJECT_ID:${GKE_PROJECT_ID} \
		-v EKS_WORKER_ROLE_ARN:${EKS_WORKER_ROLE_ARN} -v EKS_CLUSTER_ROLE_ARN:${EKS_CLUSTER_ROLE_ARN} \
		-v EKS_SUBNET_IDS:${EKS_SUBNET_IDS} -v SEPARATOR:${SEPARATOR} \
		-v CLUSTER_NAME:${CLUSTER_NAME} -v PR_NUMBER:${PR_NUMBER} \
		-f manifests/prombench/nodes_${PROVIDER}.yaml	

all_nodes_deleted:
	$(INFRA_CMD) ${PROVIDER} nodes check-deleted \
		-v ZONE:${ZONE} -v GKE_PROJECT_ID:${GKE_PROJECT_ID} \
		-v EKS_WORKER_ROLE_ARN:${EKS_WORKER_ROLE_ARN} -v EKS_CLUSTER_ROLE_ARN:${EKS_CLUSTER_ROLE_ARN} \
		-v EKS_SUBNET_IDS:${EKS_SUBNET_IDS} -v SEPARATOR:${SEPARATOR} \
		-v CLUSTER_NAME:${CLUSTER_NAME} -v PR_NUMBER:${PR_NUMBER} \
		-f manifests/prombench/nodes_${PROVIDER}.yaml
