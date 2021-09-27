package apollo

import (
	"log"
	"time"

	"github.com/apolloconfig/agollo/v4/storage"
	"github.com/micro/go-micro/v2/config/encoder"
	"github.com/micro/go-micro/v2/config/source"
)

type watcher struct {
	enc  encoder.Encoder
	name string
	ch   chan *source.ChangeSet
	exit chan bool
}

func (w *watcher) OnChange(event *storage.ChangeEvent) {
	kv := map[string]interface{}{}
	for k, v := range event.Changes {
		kv[k] = v.NewValue
	}

	kv = convert(kv)

	b, err := w.enc.Encode(kv)
	if err != nil {
		log.Println(err)
		return
	}

	scs := &source.ChangeSet{
		Data:      b,
		Format:    w.enc.String(),
		Source:    w.name,
		Timestamp: time.Now(),
	}
	scs.Checksum = scs.Sum()
	w.ch <- scs
}

func (w *watcher) OnNewestChange(event *storage.FullChangeEvent) {}

func (w *watcher) Next() (*source.ChangeSet, error) {
	select {
	case cs := <-w.ch:
		return cs, nil
	case <-w.exit:
		return nil, source.ErrWatcherStopped
	}
}

func (w *watcher) Stop() error {
	select {
	case <-w.exit:
		return nil
	default:
	}
	return nil
}

func NewWatcher(name string, e encoder.Encoder) *watcher {
	return &watcher{
		enc:  e,
		name: name,
		ch:   make(chan *source.ChangeSet),
		exit: make(chan bool),
	}
}
