package inputs

import (
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	"github.com/zhaogogo/go-logfilter/core"
	"github.com/zhaogogo/go-logfilter/field"
	"github.com/zhaogogo/go-logfilter/metrics"
)

func NewInputCell(inputType string, input Input, cellConfig map[string]interface{}, process core.Processer) (*InputCell, error) {
	i := &InputCell{
		name:   inputType,
		input:  input,
		config: cellConfig,
		//stop:         false,
		shutdownChan: make(chan struct{}, 1),
	}

	if process == nil {
		return nil, errors.New("process Node is nil")
	} else {
		i.process = process
	}
	p, err := metrics.NewPrometheusCounter(cellConfig)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	i.prometheusCounter = p
	if add_fields, ok := cellConfig["add_fields"]; ok && add_fields != nil {
		i.addFields = make(map[field.FieldSetter]field.ValueRender)
		for k, v := range add_fields.(map[string]interface{}) {
			fieldSetter := field.NewFieldSetter(k)
			if fieldSetter == nil {
				log.Fatal().Msgf("input fieldSetter构建失败", k)
			}
			i.addFields[fieldSetter] = field.GetValueRender(v)
		}
	} else {
		i.addFields = nil
	}

	return i, nil
}

type InputCell struct {
	name   string
	input  Input
	config map[string]interface{}
	//stop              bool
	shutdownChan      chan struct{}
	shutdownWhenNil   bool
	prometheusCounter prometheus.Counter
	process           core.Processer
	//exit              func()

	addFields map[field.FieldSetter]field.ValueRender
}

func (i *InputCell) SetShutdownWhenNil(shutdownWhenNil bool) {
	i.shutdownWhenNil = shutdownWhenNil
}

func (i *InputCell) Start(gid int) {
	//wg := sync.WaitGroup{}
	//wg.Add(worker)
	//for j := 0; j < worker; j++ {
	//	go func(gid int) {
	//		defer wg.Done()
	i.start(gid)
	//	}(j)
	//}
	//wg.Wait()
}

func (i *InputCell) Shutdown() {
	i.input.Shutdown()
}

func (i *InputCell) start(goid int) {
	eventCh := i.input.ReadEvent()
	log.Debug().Msgf("[%v]start inputCell event chan: %T %p\n", goid, eventCh, eventCh)
	for event := range eventCh {
		if i.prometheusCounter != nil {
			i.prometheusCounter.Inc()
		}
		if event == nil {
			log.Info().Msgf("received nil message.")
			//if i.stop {
			//	break
			//}
			if i.shutdownWhenNil {
				log.Info().Msgf("received nil message. shutdown...")
				//i.exit()
				break
			} else {
				continue
			}
		}
		for fs, v := range i.addFields {
			event = fs.SetField(event, v.Render(event), "", false)
		}
		log.Debug().Any("event", event).Msgf("[%v]ReadEvent成功", goid)
		// v, _ := json.Marshal(event)
		// fmt.Printf("res: [%v] %v\n", goid, string(v))
		i.process.Process(event)
	}
	log.Debug().Msgf("[%v]input cell %v read event stop, len: %v", goid, i.name, len(eventCh))
}
