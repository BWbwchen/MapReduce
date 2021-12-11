package worker

type MrContext struct {
	Chan chan KV
}

func newMrContext() MrContext {
	return MrContext{
		Chan: make(chan KV, 100),
	}
}

func (mc *MrContext) EmitIntermediate(key string, value string) {
	mc.Chan <- newKV(key, value)
}

func (mc *MrContext) Emit(key string, value string) {
	mc.Chan <- newKV(key, value)
}
