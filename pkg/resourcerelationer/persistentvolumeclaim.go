package resourcerelationer

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

type PersistentVolumeClaim struct {
	*Resource
}

func NewPersistentVolumeClaim() *PersistentVolumeClaim {
	return &PersistentVolumeClaim{
		Resource: NewResource(),
	}
}

func (pvc *PersistentVolumeClaim) GenerateChildren(obj *unstructured.Unstructured) error {
	// set children
	if pvc.children == nil {
		pvName := obj.Object["spec"].(map[string]interface{})["volumeName"].(string)
		pv := NewPersistentVolume()
		pv.SetParams(R_PV, pvName, "")
		pvc.MakeRelationChildren([]ResourceRelationer{pv})
	}
	return nil
}
