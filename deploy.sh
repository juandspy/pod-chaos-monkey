# Needed to pull the local image from Minikube
eval $(minikube docker-env)
# Build the image using the local version
docker build -t juandspy/pod-chaos-monkey:latest .
# Update all the resources
kubectl apply -f deployment