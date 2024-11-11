package inputs

import (
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	"github.com/zhaogogo/go-logfilter/core/topology"
	"github.com/zhaogogo/go-logfilter/pkg/field"
	"github.com/zhaogogo/go-logfilter/pkg/metrics"
)

func NewInputer(name string, input topology.Input, cellConfig map[string]interface{}, process topology.Process) (*Inputer, error) {
	var overwrite bool = false
	overwriteAny, ok := cellConfig["overwrite"]
	if ok {
		if overwriteBool, ok := overwriteAny.(bool); ok {
			overwrite = overwriteBool
		}
	}
	i := &Inputer{
		name:         name,
		input:        input,
		config:       cellConfig,
		shutdownChan: make(chan struct{}, 1),
		addFields:    field.NewAddFields(name, cellConfig),
		overwrite:    overwrite,
	}

	if process == nil {
		return nil, errors.New("process Node is nil")
	} else {
		i.process = process
	}
	p, err := metrics.NewPrometheusCounter(cellConfig)
	if err != nil {
		log.Fatal().Err(err)
	}
	i.prometheusCounter = p
	return i, nil
}

type Inputer struct {
	input  topology.Input
	name   string
	config map[string]interface{}
	//stop              bool
	shutdownChan      chan struct{}
	shutdownWhenNil   bool
	prometheusCounter prometheus.Counter
	process           topology.Process
	//exit              func()

	addFields    []map[field.FieldSetter]field.ValueRender
	overwrite    bool
	addFailedTag bool
}

func (i *Inputer) SetShutdownWhenNil(shutdownWhenNil bool) {
	i.shutdownWhenNil = shutdownWhenNil
}

func (i *Inputer) Start(gid int) {
	i.start(gid)
}

func (i *Inputer) Shutdown() {
	i.input.Shutdown()
}

func (i *Inputer) start(goid int) {
	eventCh := i.input.ReadEvent()
	log.Debug().Msgf("[%v]start inputCell event chan: %T %p\n", goid, eventCh, eventCh)
	for event := range eventCh {
		if i.prometheusCounter != nil {
			i.prometheusCounter.Inc()
		}
		if event == nil {
			log.Info().Msgf("received nil message.")
			if i.shutdownWhenNil {
				log.Info().Msgf("received nil message. shutdown...")
				//i.exit()
				break
			} else {
				continue
			}
		}
		for _, fs := range i.addFields {
			for fieldsetter, valuerender := range fs {
				event = fieldsetter.SetField(event, valuerender.Render(event), i.overwrite)
			}
		}
		log.Debug().Any("event", event).Msgf("[%v]ReadEvent成功", goid)
		// v, _ := json.Marshal(event)
		// fmt.Printf("res: [%v] %v\n", goid, string(v))
		i.process.Process(event)
	}
	log.Info().Msgf("[%v]input cell %v read event stop, len: %v", goid, i.name, len(eventCh))
}
