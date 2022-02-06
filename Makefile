build:
	# build images
	echo "Building docker images..."
	# set minikube env to reference local images
	eval $(minikube docker-env); docker build -t api-service .
	eval $(minikube docker-env); docker build -t coindesk -f microservice.docker . --build-arg MAIN_ARG="-coindesk"
	eval $(minikube docker-env); docker build -t twitterscraper -f microservice.docker . --build-arg MAIN_ARG="-twitterscraper"
	echo "Completed building docker images..."
	
deploy:
	# deploy kubernetes
	echo "Deploying kubernetes microservices..."
	eval $(minikube docker-env); minikube kubectl -- apply -f coindesk-deployment.yaml
	eval $(minikube docker-env); minikube kubectl -- apply -f twitterscraper-deployment.yaml
	eval $(minikube docker-env); minikube kubectl -- apply -f api-service-deployment.yaml
	eval $(minikube docker-env); minikube kubectl -- get deployments
	eval $(minikube docker-env); minikube dashboard 
	echo "Done deploying kubernetes microservices."

clean:
	# clean go cache
	go clean -modcache
	# clean docker images	
	docker system prune -f
	echo "cleaning kubernetes default namespace"
	eval $(minikube docker-env); kubectl delete -f coindesk-deployment.yaml
	eval $(minikube docker-env); kubectl delete -f twitterscraper-deployment.yaml
	eval $(minikube docker-env); kubectl delete -f api-service-deployment.yaml
	eval $(minikube docker-env); kubectl -n default delete pod,svc --all