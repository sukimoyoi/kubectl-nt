package cmd

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/disiqueira/gotree"
	"github.com/sukimoyoi/kubectl-nt/pkg/client"
	rr "github.com/sukimoyoi/kubectl-nt/pkg/resourcerelationer"
	"k8s.io/client-go/dynamic"
)

func NeoTree(args []string) error {
	var namespace string

	if len(args) == 2 {
		namespace = ""
	}
	if err := GetTree(args[0], args[1], namespace); err != nil {
		return fmt.Errorf("Error: %w", err)
	}
	return nil
}

func GetTree(kind, name, namespace string) error {
	if err := client.InitKubeClient(); err != nil {
		return fmt.Errorf("failed to init NeoTreeClient: %w", err)
	}

	resourceKind := client.Client.RKMap[kind]
	if resourceKind == nil {
		return fmt.Errorf("failed to get basic resource information. '%s' may not exist in the target cluster", kind)
	}

	r, err := rr.NewResourceRelationerWithBasicInfo(resourceKind.GroupVersionResource.Resource, name, namespace)
	if err != nil {
		return err
	}
	// fmt.Println("ParsedResourcekind: ", r.GetKind(), "ResourceName: ", r.GetName(), "ResourceNamespace: ", r.GetNamespace())

	var RQ rr.ResourceQueue
	if err := GetRelation(r); err != nil {
		return fmt.Errorf("failed to get resource relation: %w", err)
	}

	// Get parents relation tree
	parents := r.GetParents()
	if parents != nil {

		RQ.PushSlice(parents)
		for RQ.Len() > 0 {
			var elem rr.ResourceRelationer
			elem = RQ.Pop()
			// fmt.Println("ParentsLen:", RQ.Len(), "Elem:", elem.GetName(), elem.GetKind(), elem.GetNamespace())
			if err := GetRelation(elem); err != nil {
				return fmt.Errorf("failed to get parent resources: %w", err)
			}
			RQ.PushSlice(elem.GetParents())
		}

		// Print parents
		parentsRoot := gotree.New(r.GetFormalName() + " (Parents)")
		AddParentsBranchRecursively(parentsRoot, r)
		fmt.Println(parentsRoot.Print())
	}

	// Get children relation tree
	children := r.GetChildren()
	if children != nil {
		RQ.PushSlice(children)
		for RQ.Len() > 0 {
			var elem rr.ResourceRelationer
			elem = RQ.Pop()
			// fmt.Println("ChildrenLen:", RQ.Len(), "Elem:", elem.GetName(), elem.GetKind(), elem.GetNamespace())
			if err := GetRelation(elem); err != nil {
				return fmt.Errorf("failed to get chlid resources: %w", err)
			}
			RQ.PushSlice(elem.GetChildren())
		}

		// Print children
		childrenRoot := gotree.New(r.GetFormalName() + " (Children)")
		AddChildrenBranchRecursively(childrenRoot, r)
		fmt.Println(childrenRoot.Print())
	}
	return nil
}

func AddParentsBranchRecursively(root gotree.Tree, r rr.ResourceRelationer) {
	elements := r.GetParents()
	if elements == nil {
		return
	}
	for _, elem := range elements {
		branch := root.Add(elem.GetFormalName())
		AddParentsBranchRecursively(branch, elem)
	}
}

func AddChildrenBranchRecursively(root gotree.Tree, r rr.ResourceRelationer) {
	elements := r.GetChildren()
	if elements == nil {
		return
	}
	for _, elem := range elements {
		branch := root.Add(elem.GetFormalName())
		AddChildrenBranchRecursively(branch, elem)
	}
}

func GetRelation(r rr.ResourceRelationer) error {
	var err error

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)

	go func() {
		time.Sleep(time.Second * 10)
		println("Timeout")
		cancel()
	}()

	var ri dynamic.ResourceInterface

	rkm := client.Client.RKMap[r.GetKind()]

	if rkm.Namespaced {
		ns := r.GetNamespace()

		if ns == "" {
			ns, _, err = client.Client.ConfigFlags.ToRawKubeConfigLoader().Namespace()
			if err != nil {
				return fmt.Errorf("namespaced resource but 'namespace' was not given and can't get namespace set in current context: %w", err)
			}
		}

		ri = client.Client.DYN.Resource(rkm.GroupVersionResource).Namespace(ns)
	} else {
		ri = client.Client.DYN.Resource(rkm.GroupVersionResource)
	}

	obj, err := ri.Get(ctx, r.GetName(), metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get resources: %w", err)
	}

	r.SetParamsWithObj(obj)
	if err := r.GenerateParents(obj); err != nil {
		return err
	}
	if err := r.GenerateChildren(obj); err != nil {
		return err
	}

	return nil
}
