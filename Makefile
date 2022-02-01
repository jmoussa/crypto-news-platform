build:
	# build images
	eval $(minikube docker-env)
	docker build -t api-service .
	docker build -t coindesk-server -f coindesk-server .
	
deploy:
	# deploy kubernetes
	eval $(minikube docker-env)
	minikube kubectl -- apply -f server-deployment.yaml
	minikube kubectl -- apply -f api-service-deployment.yaml
	minikube kubectl -- get deployments

clean:
	# clean kubernetes
	eval $(minikube docker-env)
	kubectl -n default delete pod,svc --all
