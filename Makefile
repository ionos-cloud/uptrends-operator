VERSION ?= 0.0.1

# ENVTEST_K8S_VERSION refers to the version of kubebuilder assets to be downloaded by envtest binary.
ENVTEST_K8S_VERSION = 1.25.0

# kustomize for deploy
KUSTOMIZE = go run sigs.k8s.io/kustomize/kustomize/v3

IMAGE_TAG_BASE ?= ghcr.io/ionos-cloud/uptrends-operator/operator

IMG ?= $(IMAGE_TAG_BASE):v$(VERSION)

# GO related variables.
GO ?= go
GO_RUN_TOOLS ?= $(GO) run -modfile ./tools/go.mod
GO_TEST = $(GO_RUN_TOOLS) gotest.tools/gotestsum --format pkgname
GO_RELEASER ?= $(GO_RUN_TOOLS) github.com/goreleaser/goreleaser

##@ Development

.PHONY: generate
generate:
	$(GO) generate ./...
	$(GO) run cmd/manifest/manifest.go --file manifests/crd/bases/operators.ionos-cloud.github.io_uptrends.yaml \
		--file manifests/install/service_account.yaml \
		--file manifests/install/cluster_role.yaml \
		--file manifests/install/cluster_role_binding.yaml \
		--file manifests/install/statefulset.yaml \
		--output manifests/install.yaml

##@ Build

.PHONY: fmt
fmt: ## Run go fmt against code.
	$(GO_RUN_TOOLS) mvdan.cc/gofumpt -w .

.PHONY: vet
vet: ## Run go vet against code.
	$(GO) vet ./...

.PHONY: lint
lint: ## Run lint.
	$(GO_RUN_TOOLS) github.com/golangci/golangci-lint/cmd/golangci-lint run --timeout 5m -c .golangci.yml

.PHONY: build
build: ## Build manager binary.
	$(GO_RELEASER) build --rm-dist --snapshot

##@ Deployment

deploy: generate ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	cd manifests/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build manifests/default | kubectl apply -f -

remove: generate ## Remove controller to the K8s cluster specified in ~/.kube/config.
	cd manifests/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build manifests/default | kubectl delete -f -

##@ Test

LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

ENVTEST ?= $(LOCALBIN)/setup-envtest

.PHONY: envtest
envtest: $(ENVTEST) ## Download envtest-setup locally if necessary.
$(ENVTEST): $(LOCALBIN)
	test -s $(LOCALBIN)/setup-envtest || GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest

.PHONY: test
test: envtest ## Run tests.
	mkdir -p .test/reports
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) --bin-dir $(LOCALBIN) -p path)" $(GO_TEST) --junitfile .test/reports/unit-test.xml -- -race ./... -count=1 -short -cover -coverprofile .test/reports/unit-test-coverage.out
