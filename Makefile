

help:
	@echo "GenFract -- fractal generator -- view at ip:4000"
	@echo ""
	@echo "build    -- build docker image"
	@echo "toreg    -- push image to docker hub"
	@echo "deploy   -- deploy as deployment and loadbalanced service"
	@echo "redo     -- k8s redeploy"
	@echo "delete   -- delete deployment & service"
	@echo "mkpod    -- make a k8s pod"
	@echo "delpod   -- delete k8s pod"
	@echo "drun     -- run local docker image"
	@echo "logs     -- show logs for deployment (each pod)"

# any local configuration
-include .loccfg.mk

REGISTRY ?= larryrau

build:
	go install genfract.go

image:
	GOOS=linux GOARCH=amd64 go build genfract.go
	docker build -t $(REGISTRY)genfract -f Dockerfile .

toreg: image
	docker push $(REGISTRY)genfract

deploy: xx
	kubectl apply -f xx
	kubectl apply -f k8s/service.yaml
	@rm xx

redo:
	kubectl rollout restart deployment/genfract

delete:
	kubectl delete service/genfract
	kubectl delete deployment/genfract

mkpod:
	kubectl run genfract --image=$(REGISTRY)genfract --port=4000
	kubectl get pods

delpod:
	kubectl delete pod/genfract

drun:
	docker run -it -p 4000:4000 $(REGISTRY)genfract

irun:
	docker run -it -p 4000:4000 -u root --entrypoint /bin/sh $(REGISTRY)genfract

logs:
	kubectl logs --tail=5 -l name=genfract

xx:
	-@rm xx
	cat k8s/deployment.yaml |(REGISTRY=$(REGISTRY)  envsubst) >> xx
