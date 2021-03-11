package client

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"k8s.io/client-go/dynamic"
)

type ResourceKind struct {
	GroupVersionResource schema.GroupVersionResource
	Namespaced           bool
}

type KubeClient struct {
	ConfigFlags *genericclioptions.ConfigFlags
	DYN         dynamic.Interface
	RKMap       map[string]*ResourceKind
}

var Client *KubeClient

func InitKubeClient() error {
	Client = &KubeClient{}

	kubeConfig, err := newRestconfig()
	if err != nil {
		return fmt.Errorf("failed to get restconfig: %w", err)
	}

	cf := genericclioptions.NewConfigFlags(true)
	Client.ConfigFlags = cf

	dc, err := cf.ToDiscoveryClient()
	if err != nil {
		return fmt.Errorf("failed to create discoveryclient: %w", err)
	}

	apiResourceList, err := dc.ServerPreferredResources()
	if err != nil {
		return fmt.Errorf("failed to get apiresourcelist: %w", err)
	}
	dyn, err := dynamic.NewForConfig(kubeConfig)
	if err != nil {
		return fmt.Errorf("failed to construct dynamic client: %w", err)
	}

	Client.DYN = dyn
	Client.RKMap = map[string]*ResourceKind{}

	for _, apiRL := range apiResourceList {
		gv, err := schema.ParseGroupVersion(apiRL.GroupVersion)
		if err != nil {
			return fmt.Errorf("failed to parse api resource group & version: %w", err)
		}

		for _, apiR := range apiRL.APIResources {
			if !contains(apiR.Verbs, "list") {
				continue
			}
			merge(Client.RKMap, getResourceKindMap(apiR, gv))
		}
	}

	// for k, v := range Client.RKMap {
	// 	fmt.Println("RsourceKind:", k, "Namespaced:", v.Namespaced, "APIGroup&APIVersion:", v.GroupVersionResource)
	// 	// fmt.Printf("address of slice %p add of Arr %p \n", &k, &v.GroupVersionResource)
	// }

	return nil
}

func newRestconfig() (*rest.Config, error) {
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		return nil, err
	}
	return kubeConfig, nil
}

func getResourceKindMap(apiR metav1.APIResource, gv schema.GroupVersion) map[string]*ResourceKind {
	aNs := apiR.ShortNames
	aNs = append(aNs, apiR.Name, apiR.Name[:len(apiR.Name)-1], apiR.Kind)

	rkm := map[string]*ResourceKind{}
	rk := ResourceKind{
		GroupVersionResource: schema.GroupVersionResource{
			Version:  gv.Version,
			Group:    gv.Group,
			Resource: apiR.Name,
		},
		Namespaced: apiR.Namespaced,
	}
	for _, v := range aNs {
		rkm[v] = &rk
	}
	return rkm
}

func ListResources(gvr *schema.GroupVersionResource, namespace string) (*[]unstructured.Unstructured, error) {
	var out []unstructured.Unstructured
	var next string
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		time.Sleep(time.Second * 10)
		println("Timeout")
		cancel()
	}()

	var intf dynamic.ResourceInterface
	if namespace == "" {
		intf = Client.DYN.Resource(*gvr)
	} else {
		intf = Client.DYN.Resource(*gvr).Namespace(namespace)
	}
	resp, err := intf.List(ctx, metav1.ListOptions{
		Limit:    250,
		Continue: next,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get resources: %w", err)
	}

	for {
		out = append(out, resp.Items...)
		next = resp.GetContinue()
		if next == "" {
			break
		}
	}

	// for i, v := range out {
	// 	fmt.Println("Number: ", i)
	// 	fmt.Println(v)
	// }

	return &out, nil
}

func contains(s []string, e string) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}

func merge(rKM1, rKM2 map[string]*ResourceKind) {
	for k, v := range rKM2 {
		rKM1[k] = v
	}
}
