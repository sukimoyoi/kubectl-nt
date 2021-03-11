package resourcerelationer

import (
	"fmt"

	"github.com/sukimoyoi/kubectl-nt/pkg/client"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type StorageClass struct {
	*Resource
}

func NewStorageClass() *StorageClass {
	return &StorageClass{
		Resource: NewResource(),
	}
}

func (sc *StorageClass) GenerateParents(obj *unstructured.Unstructured) error {
	var (
		pvs          []ResourceRelationer
		pvsUntracked *[]unstructured.Unstructured
		err          error
	)
	sc.SetParamsWithObj(obj)

	// set parents
	if sc.parents == nil {
		rk := client.Client.RKMap[R_PV]

		if pvsUntracked, err = client.ListResources(&rk.GroupVersionResource, ""); err != nil {
			return fmt.Errorf("failed generate resource relation %w:", err)
		}
		for _, v := range *pvsUntracked {
			pv := NewPersistentVolume()
			pv.SetParamsWithObj(&v)
			pvs = append(pvs, pv)
		}
		sc.MakeRelationParents(pvs)
	}
	return nil
}
