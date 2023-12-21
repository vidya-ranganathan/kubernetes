package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Kluster struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Spec KlusterSpec
}

type KlusterSpec struct {
	Name    string
	Region  string
	Version string

	NodePools []NodePool
}

type NodePool struct {
	Size  string
	Name  string
	Count int
}
