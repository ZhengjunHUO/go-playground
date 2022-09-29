package filter

import (
	"testing"
	"reflect"
)

var c1= &Cluster{
	Name:	"huo_infra",
	Labels:	map[string]string{
		"client": "huo.io",
		"env": "infra",
	},
}

var c2 = &Cluster{
	Name:	"huo_dev",
	Labels:	map[string]string{
		"client": "huo.io",
		"env": "dev",
	},
}

var c3 = &Cluster{
	Name:	"huo_market",
	Labels:	map[string]string{
		"client": "huo.io",
		"env": "market",
	},
}

var c4 = &Cluster{
	Name:	"foo_business",
	Labels:	map[string]string{
		"client": "foo",
		"env": "business",
	},
}

var c5 = &Cluster{
	Name:	"bar_preprod",
	Labels:	map[string]string{
		"client": "bar",
		"env": "preprod",
	},
}

var c6 = &Cluster{
	Name:	"not_relevant",
	Labels:	map[string]string{
		"foo1": "bar1",
		"foo2": "bar2",
		"foo3": "bar3",
	},
}

var clusters1 = []*Cluster{
	c1,
	c2,
	c3,
	c4,
	c5,
	c6,
}

var f1 = &Filter{
	Name:	"Only huo.io filter",
        Op:	"OR",
        Labels:	[]*FilterLabels{
		&FilterLabels{
			Key: "client",
			Value: "huo.io",
			Op: "=",
		},
	},
}

var f2 = &Filter{
	Name:	"Not huo.io filter",
        Op:	"OR",
        Labels:	[]*FilterLabels{
		&FilterLabels{
			Key: "client",
			Value: "huo.io",
			Op: "!=",
		},
	},
}

var f3 = &Filter{
	Name:	"Choose huo.io's infra filter",
        Op:	"AND",
        Labels:	[]*FilterLabels{
		&FilterLabels{
			Key: "client",
			Value: "huo.io",
			Op: "=",
		},
		&FilterLabels{
			Key: "env",
			Value: "infra",
			Op: "=",
		},
	},
}

var f4 = &Filter{
	Name:	"Choose foo filter",
        Op:	"AND",
        Labels:	[]*FilterLabels{
		&FilterLabels{
			Key: "client",
			Value: "huo.io",
			Op: "!=",
		},
		&FilterLabels{
			Key: "env",
			Value: "business",
			Op: "=",
		},
	},
}

var f5 = &Filter{
	Name:	"Exclude huo.io dev filter",
        Op:	"OR",
        Labels:	[]*FilterLabels{
		&FilterLabels{
			Key: "client",
			Value: "huo.io",
			Op: "!=",
		},
		&FilterLabels{
			Key: "env",
			Value: "market",
			Op: "=",
		},
		&FilterLabels{
			Key: "env",
			Value: "infra",
			Op: "=",
		},
	},
}

func TestFilterClusters(t *testing.T) {
	fs := []*Filter{f1, f2, f3, f4, f5}
	wanted := [][]string{
		[]string{
			c1.Name,
			c2.Name,
			c3.Name,
		},
		[]string{
			c4.Name,
			c5.Name,
		},
		[]string{
			c1.Name,
		},
		[]string{
			c4.Name,
		},
		[]string{
			c1.Name,
			c3.Name,
			c4.Name,
			c5.Name,
		},
	}

	for i := range fs {
		rslt := FilterClusters(clusters1, fs[i])
		if !reflect.DeepEqual(wanted[i], rslt) {
			t.Errorf("Filtered result [%s]: %v, but wait for %v", fs[i].Name, rslt, wanted[i] )
		}
	}
}
