package worker

type KV struct {
	Key   string
	Value string
}

func newKV(k string, v string) KV {
	return KV{
		Key:   k,
		Value: v,
	}
}

type byKey []KV

func (a byKey) Len() int           { return len(a) }
func (a byKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byKey) Less(i, j int) bool { return a[i].Key < a[j].Key }
