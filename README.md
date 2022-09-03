# Chaos Monkey

This is a simplistic implementation of https://github.com/Netflix/chaosmonkey.

This program runs inside a Kubernetes cluster, in a given namespace and deletes a random pod.

```
❯ ./pod-chaos-monkey --help
pod-chaos-monkey is a CLI tool used to kill a random pod in the given namespace. 
    
    The default namespace is "default", but it can be changed via CLI args.

Usage:
  pod-chaos-monkey [flags]

Flags:
      --label-selector string   Label selector to skip some pods from being deleted (default "metadata.name!=pod-chaos-monkey")
  -h, --help                    help for pod-chaos-monkey
  -n, --namespace string        Target namespace (default "default")
```

## How to deploy in a cluster

In the [deployment](deployment) folder you may find some Kubernetes resources needed to run this program in your cluster.

- [pod-chaos-monkey.yaml](deployment/pod-chaos-monkey.yaml) contains the cronjob definition. Feel free to change the `spec.schedule` according to your needs. By default it runs once per minute. It is also possible to configure the target namespace and label selector. Just update the values in the `args`.
- [service-account.yaml](deployment/service-account.yaml) contains the service account, role and role binding used for this cronjob. The service account has `list` and `destroy` pods privileges for all the namespaces in the cluster.
- [example.yaml](deployment/example.yaml) deploys 4 replicas of a Nginx container. Used to see the PodChaosMonkey in action. They are deployed in the "workloads" namespace.

Then you just have to connect to your cluster and run:

```
kubectl apply -f deployment
```

## Local development

Spin up a Kubernetes cluster using your favorite tool. For example, you can use Minikube:

```
minikube start
```

You can then deploy the changes quickly by running `./deploy.sh`. Make sure to uncomment `imagePullPolicy` in [pod-chaos-monkey.yaml](deployment/pod-chaos-monkey.yaml) in order to use the local version of the Docker image.


## TODO list

- Add support for running the program outside the cluster for easier debugging.
