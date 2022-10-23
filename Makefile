

help:
	@echo "GenFract -- fractal generator -- view at ip:4000"
	@echo ""
	@echo "build    -- build docker image"
	@echo "toreg    -- push image to docker hub"
	@echo "drun     -- run local docker image"
	@echo "** k8s helpers **"
	@echo "deploy   -- deploy as deployment and loadbalanced service"
	@echo "redo     -- k8s redeploy"
	@echo "delete   -- delete deployment & service"
	@echo "mkpod    -- make a k8s pod"
	@echo "delpod   -- delete k8s pod"

	@echo "logs     -- show logs for deployment (each pod)"

# any local configuration
-include .loccfg.mk

REGISTRY ?= larryrau/


build:
	go install genfract.go

image:
	GOOS=linux GOARCH=amd64 go build genfract.go
	docker build -t $(REGISTRY)genfract -f Dockerfile .

toreg: image
	docker push $(REGISTRY)genfract

drun:
	docker run -it -p 4000:4000 $(REGISTRY)genfract

irun:
	docker run -it -p 4000:4000 -u root --entrypoint /bin/sh $(REGISTRY)genfract


#
# k8s helpers
# note: kk -> use of envsubst to replace use of env in k8s files
#

kdeploy: kk
	kubectl apply -f kk
	kubectl apply -f k8s/service.yaml
	@rm kk

krestart:
	kubectl rollout restart deployment/genfract

kdelete:
	kubectl delete service/genfract
	kubectl delete deployment/genfract

kmkpod:
	kubectl run genfract --image=$(REGISTRY)genfract --port=4000
	kubectl get pods

kdelpod:
	kubectl delete pod/genfract

klogs:
	kubectl logs --tail=5 -l name=genfract

kk:
	-@rm kk
	cat k8s/deployment.yaml |(REGISTRY=$(REGISTRY)  envsubst) >> kk
