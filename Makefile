prep-k8s-test:
	kubectl create ns testing
	kubectl create service clusterip test-service --tcp=80:80 -n testing
clean-k8s-test:
	kubectl delete service test-service -n testing
	kubectl delete ns testing
