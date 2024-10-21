package inputs

import (
	"bufio"
	"os"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/zhaogogo/go-logfilter/textcodec"
)

type StdinInput struct {
	one     sync.Once
	config  map[string]interface{}
	decoder textcodec.Decoder

	scanner *bufio.Scanner
	fifo    chan map[string]interface{}
	stop    bool
}

func init() {
	Register("stdin", newStdinInput)
}

func newStdinInput(config map[string]interface{}) Input {
	var (
		codertype    string = "plain"
		timestampkey        = "timestamp"
	)
	if v, ok := config["codec"]; ok {
		switch v.(type) {
		case string:
			codertype = v.(string)
		}
	}
	if v, ok := config["timestamp"]; ok {
		switch v.(type) {
		case string:
			timestampkey = v.(string)
		}
	}

	p := &StdinInput{
		config:  config,
		decoder: textcodec.NewDecoder(codertype, timestampkey),
		scanner: bufio.NewScanner(os.Stdin),
		fifo:    make(chan map[string]interface{}, 3),
		stop:    false,
	}

	return p
}

func (p *StdinInput) ReadEvent() chan map[string]interface{} {

	go p.read()

	return p.fifo
}

func (p *StdinInput) read() {
	for !p.stop {
		if p.scanner.Scan() {
			t := p.scanner.Bytes()
			msg := make([]byte, len(t))
			copy(msg, t)
			event := p.decoder.Decode(msg)
			p.fifo <- event
		}
		if err := p.scanner.Err(); err != nil {
			log.Error().Msgf("stdin scan error: %v", err)
		}

	}
	close(p.fifo)
	log.Info().Msg("stdin input plugin close channel")
	log.Info().Msg("stdin input plugin shutdown success")
}

func (p *StdinInput) Shutdown() {
	// what we need is to stop emit new event; close messages or not is not important
	log.Info().Msg("stdin plugin Shutdown...")
	p.stop = true
}
