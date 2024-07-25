package inputs

import (
	"bufio"
	"github.com/zhaogogo/go-logfilter/textcodec"
	"os"
	"sync"
	"time"

	"k8s.io/klog/v2"
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
	var codertype string = "plain"
	if v, ok := config["codec"]; ok {
		codertype = v.(string)
	}
	p := &StdinInput{

		config:  config,
		decoder: textcodec.NewDecoder(codertype),
		scanner: bufio.NewScanner(os.Stdin),
		fifo:    make(chan map[string]interface{}, 1),
		stop:    false,
	}

	return p
}

func (p *StdinInput) ReadEvent() chan map[string]interface{} {
	p.one.Do(func() {
		go p.read()
	})
	return p.fifo
}

func (p *StdinInput) read() {
	for {
		if p.stop {
			close(p.fifo)
			break
		}
		if p.scanner.Scan() {
			t := p.scanner.Bytes()
			msg := make([]byte, len(t))
			copy(msg, t)
			event := p.decoder.Decode(msg)
			p.fifo <- event
		}
		if err := p.scanner.Err(); err != nil {
			klog.Errorf("stdin scan error: %v", err)
		} else {
			// EOF here. when stdin is closed by C-D, cpu will raise up to 100% if not sleep
			time.Sleep(time.Millisecond * 1000)
		}

	}
	klog.Infof("stdin input plugin shutdown success")
}

func (p *StdinInput) Shutdown() {
	// what we need is to stop emit new event; close messages or not is not important
	klog.Infof("stdin plugin Shutdown")
	p.stop = true
}
