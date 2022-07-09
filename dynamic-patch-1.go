package main

import (
	"context"
	"fmt"
	"os"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	sigyaml "sigs.k8s.io/yaml"
)

func main() {
	s, _ := os.ReadFile("D:/code/deploy.yaml")
	kubeconfig := "D:/code/config"
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	gr, _ := restmapper.GetAPIGroupResources(clientset.Discovery())
	mapper := restmapper.NewDiscoveryRESTMapper(gr)

	var resourceREST dynamic.ResourceInterface
	dyn, err := dynamic.NewForConfig(config)
	unstructuredObj := &unstructured.Unstructured{}

	_, groupVersionAndKind, err :=
		yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(s, nil, unstructuredObj)
	mapping, _ := mapper.RESTMapping(groupVersionAndKind.GroupKind(), groupVersionAndKind.Version)

	fmt.Printf("%T %v\n", mapping, mapping)
	fmt.Printf("%T %v\n", unstructuredObj, unstructuredObj)

	// 需要为 namespace 范围内的资源提供不同的接口
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		if unstructuredObj.GetNamespace() == "" {
			unstructuredObj.SetNamespace("default")
		}
		resourceREST = dyn.Resource(mapping.Resource).Namespace(unstructuredObj.GetNamespace())
	} else {
		resourceREST = dyn.Resource(mapping.Resource)
	}
	unstructuredYaml, _ := sigyaml.Marshal(unstructuredObj)

	_, getErr := resourceREST.Get(context.Background(), unstructuredObj.GetName(), metav1.GetOptions{})
	if getErr != nil {
		_, createErr := resourceREST.Create(context.Background(), unstructuredObj, metav1.CreateOptions{})
		if createErr != nil {
			panic(createErr)
		}
	}
	force := true
	d, err := resourceREST.Patch(context.Background(), unstructuredObj.GetName(), types.ApplyPatchType,
		unstructuredYaml,
		metav1.PatchOptions{
			FieldManager: unstructuredObj.GetName(),
			Force:        &force,
		})
	if err != nil {
		panic(err)
		return
	}
	fmt.Printf("%T %v\n", d, d)
}
