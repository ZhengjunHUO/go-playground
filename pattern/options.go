package pattern

const (
	defaultName = "huo's k8s"
	defaultCNI  = "cilium"
	defaultSize = 6
	defaultMgmt = false
	defaultBare = true
	defaultOvly = false
)

type K8sCluster struct {
	Name		string
	CNI		string
	Size		int
	IsManaged	bool
	IsBaremetal	bool
	IsOverlay	bool
}

type K8sClusterOption func(*K8sCluster) 

func K8sName(name string) K8sClusterOption {
	return func(k *K8sCluster) {
		k.Name = name
	}
}

func K8sCNI(cni string) K8sClusterOption {
	return func(k *K8sCluster) {
		k.CNI = cni
	}
}

func K8sSize(size int) K8sClusterOption {
	return func(k *K8sCluster) {
		k.Size = size
	}
}

func K8sManaged(isManaged bool) K8sClusterOption {
	return func(k *K8sCluster) {
		k.IsManaged = isManaged 
	}
}

func K8sBaremetal(isBaremetal bool) K8sClusterOption {
	return func(k *K8sCluster) {
		k.IsBaremetal = isBaremetal
	}
}

func K8sOverlay(isOverlay bool) K8sClusterOption {
	return func(k *K8sCluster) {
		k.IsOverlay = isOverlay
	}
}

func NewK8sCluster(kops ...K8sClusterOption) *K8sCluster {
	k8s := &K8sCluster {
		Name: defaultName,
		CNI:  defaultCNI,
		Size: defaultSize,
		IsManaged:   defaultMgmt,
		IsBaremetal: defaultBare,
		IsOverlay:   defaultOvly,
	}

	for _, v := range kops {
		v(k8s)
	}

	return k8s
}
