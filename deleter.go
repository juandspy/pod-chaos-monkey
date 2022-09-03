package main

import (
	"context"
	"errors"
	"math/rand"

	"github.com/rs/zerolog/log"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type PodDeleter struct {
	clientset kubernetes.Interface
	ctx       context.Context
}

func NewPodDeleter() PodDeleter {
	log.Debug().
		Msg("Initializing client")
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Error reading in cluster configuration")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Error creating the client")
	}

	return PodDeleter{
		clientset: clientset,
		ctx:       context.TODO(),
	}
}

func (p PodDeleter) DeleteRandomPod(namespace, labelSelector string) {
	log.Info().Str("namespace", namespace).Msg("Deleting a random pod")
	pods, err := p.listPods(namespace, labelSelector)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Error listing pods")
	}
	pod, err := pickRandomPod(pods)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Error picking pod")
	}

	err = p.deletePod(pod.Name, namespace)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Error deleting pod")
	}
	log.Info().Str("pod", pod.Name).Msg("Pod deleted")
}

func (p PodDeleter) listPods(namespace, labelSelector string) ([]v1.Pod, error) {
	log.Debug().
		Str("labelSelector", labelSelector).
		Msg("Listing pods")
	podList, err := p.clientset.CoreV1().
		Pods(namespace).
		List(p.ctx, metav1.ListOptions{
			LabelSelector: labelSelector,
		})
	log.Debug().
		Strs("pods", prettySlicePods(podList.Items)).
		Msgf("There are %d pods\n", len(podList.Items))
	if podList.Items == nil {
		return []v1.Pod{}, nil
	}
	return podList.Items, err
}

func (p PodDeleter) deletePod(pod, namespace string) error {
	log.Debug().
		Str("pod", pod).
		Msg("Deleting pod")
	return p.clientset.CoreV1().
		Pods(namespace).
		Delete(p.ctx, pod, metav1.DeleteOptions{})
}

func pickRandomPod(in []v1.Pod) (v1.Pod, error) {
	if len(in) == 0 {
		return v1.Pod{}, errors.New("no pods found")
	}
	randomIndex := rand.Intn(len(in))
	return in[randomIndex], nil
}

// prettySlicePods converts []v1.Pod into []string using each Pod.Name. This is useful for logging and debugging
func prettySlicePods(pods []v1.Pod) []string {
	out := []string{}
	for _, pod := range pods {
		out = append(out, pod.Name)
	}
	return out
}
