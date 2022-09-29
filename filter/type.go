package filter

type Cluster struct {
	Name	string
	Labels	map[string]string
}

type Filter struct {
	Name		string
        Labels		[]*FilterLabels
        Op		string
}

type FilterLabels struct {
        Key		string
        Value		string
        Op		string
}

