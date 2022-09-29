package filter

func FilterClusters(clusters []*Cluster, filter *Filter) []string {
	rslt := []string{}
        if filter.Op == "OR" {
                for _, cluster := range clusters {
                        match := false
                        for _, label := range filter.Labels {
                                if label.Op == "=" {
                                        if v, exist := cluster.Labels[label.Key]; exist && v == label.Value {
                                                match = true
                                                break
                                        }
                                }else{
                                        if v, exist := cluster.Labels[label.Key]; exist && v != label.Value {
                                                match = true
                                                break
                                        }
                                }
                        }

                        if match {
                                rslt = append(rslt, cluster.Name)
                        }
                }

        }else{
                // case filter.Op == "AND"
                for _, cluster := range clusters {
                        match := true
                        for _, label := range filter.Labels {
                                if label.Op == "=" {
                                        if v, exist := cluster.Labels[label.Key]; !(exist && v == label.Value) {
                                                match = false
                                                break
                                        }
                                }else{
                                        if v, exist := cluster.Labels[label.Key]; !(exist && v != label.Value) {
                                                match = false
                                                break
                                        }
                                }
                        }

                        if match {
                                rslt = append(rslt, cluster.Name)
                        }
                }
        }

	return rslt
}

/*
func main() {
	fs := []*Filter{f1, f2, f3, f4, f5}

	for _, f := range fs {
		rslt := FilterClusters(clusters1, f)
		fmt.Printf("%s: %+v\n", f.Name, rslt)
	}
}
*/
