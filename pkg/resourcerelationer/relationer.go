package resourcerelationer

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// support resources
const (
	R_PV  = "persistentvolumes"
	R_PVC = "persistentvolumeclaims"
	R_SC  = "storageclasses"
)

// ResourceRelationer ...
type ResourceRelationer interface {
	GetKind() string
	GetName() string
	GetNamespace() string
	GetFormalName() string
	GetParents() []ResourceRelationer
	GetChildren() []ResourceRelationer
	SetParent(ResourceRelationer)
	SetChild(ResourceRelationer)
	SetParents([]ResourceRelationer)
	SetChildren([]ResourceRelationer)
	SetParams(string, string, string)
	SetParamsWithObj(obj *unstructured.Unstructured)
	GenerateParents(*unstructured.Unstructured) error
	GenerateChildren(*unstructured.Unstructured) error
	// GenerateResourceRelation(*unstructured.Unstructured) error
}

// Resource ...
type Resource struct {
	ResourceRelationer
	kind      string
	name      string
	namespace string
	// parents controll or refer this resource (e.g. Deploy is RS's parent)
	parents []ResourceRelationer
	// children are controlled or referd by this resource (e.g. Pod is RS's child)
	children []ResourceRelationer
}

func (r *Resource) GetKind() string {
	return r.kind
}

func (r *Resource) GetName() string {
	return r.name
}

func (r *Resource) GetNamespace() string {
	return r.namespace
}

func (r *Resource) GetFormalName() string {
	return r.kind + "/" + r.name
}

func (r *Resource) GetParents() []ResourceRelationer {
	return r.parents
}

func (r *Resource) GetChildren() []ResourceRelationer {
	return r.children
}

func (r *Resource) SetParents(parents []ResourceRelationer) {
	r.parents = append(r.parents, parents...)
}

func (r *Resource) SetChildren(children []ResourceRelationer) {
	r.children = append(r.children, children...)
}

func (r *Resource) SetParent(parent ResourceRelationer) {
	r.parents = append(r.parents, parent)
}

func (r *Resource) SetChild(child ResourceRelationer) {
	r.children = append(r.children, child)
}

func (r *Resource) GenerateParents(*unstructured.Unstructured) error {
	return nil
}

func (r *Resource) GenerateChildren(*unstructured.Unstructured) error {
	return nil
}

// func (r *Resource) GenerateResourceRelation(*unstructured.Unstructured) error {
// 	panic("This is a fake function. Please override.")
// }

// func NewResourceWithoutRelation(kind, name, namespace string) *Resource {
// 	return &Resource{
// 		Kind: kind, Name: name, Namespace: namespace,
// 	}
// }

// func NewResourceWithParents(kind, name, namespace string, parents *[]*Resource) *ResourceRelationer {
// 	r := NewResourceWithoutRelation(kind, name, namespace)
// 	r.parents = *parents
// 	return r
// }

// func NewResourceWithChildren(kind, name, namespace string, children *[]*Resource) *ResourceRelationer {
// 	r := NewResourceWithoutRelation(kind, name, namespace)
// 	r.children = *children
// 	return r
// }

func (r *Resource) SetParams(kind, name, namespace string) {
	r.kind = strings.ToLower(kind)
	r.name = name
	r.namespace = namespace
}

func (r *Resource) SetParamsWithObj(obj *unstructured.Unstructured) {
	r.kind = strings.ToLower(obj.GetKind())
	r.namespace = obj.GetNamespace()
	r.name = obj.GetName()
}

func (r *Resource) MakeRelationParents(parents []ResourceRelationer) {
	r.SetParents(parents)
	for _, p := range parents {
		p.SetChild(r)
	}
}

func (r *Resource) MakeRelationChildren(children []ResourceRelationer) {
	r.SetChildren(children)
	for _, c := range children {
		c.SetParent(r)
	}
}

func NewResource() *Resource {
	return &Resource{}
}

func NewResourceRelationerWithBasicInfo(kind, name, namespace string) (ResourceRelationer, error) {
	var rc ResourceRelationer
	switch kind {
	case R_PV:
		rc = NewPersistentVolume()
	case R_PVC:
		rc = NewPersistentVolumeClaim()
	case R_SC:
		rc = NewStorageClass()
	default:
		return nil, fmt.Errorf("kubectl-nt doesn't supoort the resource '%s'", kind)
	}
	rc.SetParams(kind, name, namespace)
	return rc, nil
}
