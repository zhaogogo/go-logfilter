package file

import (
	"context"
	"github.com/fsnotify/fsnotify"
	"github.com/zhaogogo/go-logfilter/internal/config"
	"k8s.io/klog/v2"
)

var _ config.Watcher = (*watcher)(nil)

type watcher struct {
	f  *file
	fw *fsnotify.Watcher

	ctx    context.Context
	cancel context.CancelFunc
}

func newWatcher(f *file) (config.Watcher, error) {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	if err := fw.Add(f.path); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &watcher{
		f:      f,
		fw:     fw,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

func (w *watcher) Next() ([]*config.KeyValue, error) {
	select {
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	case event := <-w.fw.Events:
		klog.V(1).Infof("fsnotify get event %v", event.String())
		//if event.Op == fsnotify.Rename {
		//	if _, err := os.Stat(event.Name); err == nil || os.IsExist(err) {
		//		if err := w.fw.Add(event.Name); err != nil {
		//			return nil, err
		//		}
		//	}
		//}
		if event.Has(fsnotify.Write) {
			return nil, nil
		}
		return nil, nil
	case err := <-w.fw.Errors:
		return nil, err
	}
}

func (w *watcher) Stop() error {
	w.cancel()
	return w.fw.Close()
}
