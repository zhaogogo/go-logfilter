package topology

type Input interface {
	ReadEvent() chan map[string]interface{}
	Shutdown()
}

type Processer interface {
	Process(map[string]interface{}) map[string]interface{}
}

type Output interface {
	Emit(map[string]interface{})
}
