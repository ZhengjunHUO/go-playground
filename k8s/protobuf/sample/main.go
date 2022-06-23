package main

import (
	"fmt"

	pf "github.com/ZhengjunHUO/playground/k8s/protobuf/protob"
	"github.com/golang/protobuf/proto"
)

func main() {
	k8sCluster, k8sStruct := &pf.K8SInfo{
		Name: "huo",
		Size: 6,
		Ismanaged: true,
		Cni: &pf.Cni{
			Name: "Cilium",
			IsOverlayed: false,
			IsDirectRouting: true,
		},
	}, &pf.K8SInfo{}

	// Encoding
	d, err := proto.Marshal(k8sCluster)
	if err != nil {
		fmt.Println("During marshaling: ", err)
		return
	}

	fmt.Println("Serialized data: ", d)

	// Decoding
	if err = proto.Unmarshal(d, k8sStruct); err != nil {
		fmt.Println("During marshaling: ", err)
		return
	}

	fmt.Printf("Deserialized structure: %v\n", k8sStruct)
}
