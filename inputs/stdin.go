package inputs

import (
	"bufio"
	"github.com/zhaogogo/go-logfilter/encoding"
	"os"
	"time"

	"k8s.io/klog/v2"
)

type StdinInput struct {
	config  map[string]interface{}
	decoder encoding.Codec

	scanner  *bufio.Scanner
	messages chan []byte

	stop bool
}

func init() {
	Register("stdin", newStdinInput)
}

func newStdinInput(config map[string]interface{}) Input {
	var codertype string = "json"
	if v, ok := config["codec"]; ok {
		codertype = v.(string)
	}
	decoder, ok := encoding.GetCodec(codertype)
	if !ok {
		klog.Fatalf("decoder类型不支持")
	}
	p := &StdinInput{

		config:   config,
		decoder:  decoder,
		scanner:  bufio.NewScanner(os.Stdin),
		messages: make(chan []byte, 10),
	}

	return p
}

func (p *StdinInput) ReadOneEvent() map[string]interface{} {
	if p.scanner.Scan() {
		t := p.scanner.Bytes()
		msg := make([]byte, len(t))
		copy(msg, t)
		event := make(map[string]interface{})
		if err := p.decoder.Unmarshal(msg, &event); err != nil {
			klog.V(1).Error(err, "event unmarshal错误")
		}
		return event
	}
	if err := p.scanner.Err(); err != nil {
		klog.Errorf("stdin scan error: %v", err)
	} else {
		// EOF here. when stdin is closed by C-D, cpu will raise up to 100% if not sleep
		time.Sleep(time.Millisecond * 1000)
	}
	return nil
}

func (p *StdinInput) Shutdown() {
	// what we need is to stop emit new event; close messages or not is not important
	p.stop = true
}
