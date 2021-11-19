package worker

type MrContext struct {
	MapChan    chan KV
	ReduceChan chan KV
}

func newMrContext() MrContext {
	return MrContext{
		MapChan:    make(chan KV, 100),
		ReduceChan: make(chan KV, 100),
	}
}

func (mc *MrContext) EmitIntermediate(key string, value string) {
	mc.MapChan <- newKV(key, value)
}

func (mc *MrContext) Emit(key string, value string) {
	mc.ReduceChan <- newKV(key, value)
}
