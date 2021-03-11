package resourcerelationer

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type PersistentVolume struct {
	*Resource
}

func NewPersistentVolume() *PersistentVolume {
	return &PersistentVolume{
		Resource: NewResource(),
	}
}

func (pv *PersistentVolume) GenerateParents(obj *unstructured.Unstructured) error {
	// set parents
	if pv.parents == nil {
		claimRef := obj.Object["spec"].(map[string]interface{})["claimRef"].(map[string]interface{})
		pvc := NewPersistentVolumeClaim()
		pvc.SetParams(R_PVC, claimRef["name"].(string), claimRef["namespace"].(string))
		pv.MakeRelationParents([]ResourceRelationer{pvc})
	}
	return nil
}

func (pv *PersistentVolume) GenerateChildren(obj *unstructured.Unstructured) error {
	// set children
	if pv.children == nil {
		storageClassName := obj.Object["spec"].(map[string]interface{})["storageClassName"].(string)
		sc := NewStorageClass()
		sc.SetParams(R_SC, storageClassName, "")
		pv.MakeRelationChildren([]ResourceRelationer{sc})
	}
	return nil
}
