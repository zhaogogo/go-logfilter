package topology

type Input interface {
	ReadEvent() chan map[string]interface{}
	Shutdown()
}

type Filter interface {
	Filter(map[string]interface{}) (map[string]interface{}, error)
}

type Process interface {
	Process(map[string]interface{}) map[string]interface{}
}

type Output interface {
	Emit(map[string]interface{})
}
