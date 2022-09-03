package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	testclient "k8s.io/client-go/kubernetes/fake"
)

const testNamespace = "test"

var testPod = v1.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "test",
		Namespace: testNamespace,
	},
}

type tests struct {
	name        string
	objects     []runtime.Object
	want        []v1.Pod
	expectError bool
}

func TestListPods(t *testing.T) {
	for _, test := range []tests{
		{
			name:    "if the namespace has no pods, return empty slice",
			objects: nil,
			want:    []v1.Pod{},
		},
		{
			name: "if the namespace has one pod, return a slice with the pod",
			objects: []runtime.Object{
				&testPod,
			},
			want: []v1.Pod{testPod},
		},
	} {
		t.Run(test.name, func(t *testing.T) {

			sut := PodDeleter{
				clientset: testclient.NewSimpleClientset(test.objects...),
			}

			got, err := sut.listPods(testNamespace, "")
			assert.NoError(t, err, "error listing pods")

			assert.Equal(t, test.want, got)
		})
	}
}

func TestDeletePod(t *testing.T) {
	for _, test := range []tests{
		{
			name:        "if the pod doesn't exist, return an error",
			objects:     nil,
			expectError: true,
		},
		{
			name: "if the pod exists, don't return an error",
			objects: []runtime.Object{
				&testPod,
			},
			expectError: false,
		},
	} {
		t.Run(test.name, func(t *testing.T) {

			sut := PodDeleter{
				clientset: testclient.NewSimpleClientset(test.objects...),
			}

			err := sut.deletePod(testPod.Name, testNamespace)
			if test.expectError {
				assert.Error(t, err, "expecting an error deleting pod")
			} else {
				assert.NoError(t, err, "got an error deleting pod")
			}

			pods, err := sut.listPods(testNamespace, "")
			assert.NoError(t, err, "error listing pods")
			assert.Zero(t, len(pods), "there is one pod alive and there shouldn't")
		})
	}
}

func TestPickRandomPod(t *testing.T) {
	t.Run("if in slice is empty, return an error", func(t *testing.T) {
		in := []v1.Pod{}
		got, err := pickRandomPod(in)
		assert.Equal(t, got, v1.Pod{}, "expected an empty pod")
		assert.Error(t, err, "expected an error if in slice is empty")
	})

	t.Run("if in slice is not empty, return any of the pods", func(t *testing.T) {
		in := []v1.Pod{testPod}
		got, err := pickRandomPod(in)
		assert.Equal(t, got, testPod, "expected pod with name %q", testPod.Name)
		assert.NoError(t, err, "unexpected error when picking a random pod")
	})
}
