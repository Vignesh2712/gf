package zookeeper

import (
	"context"
	"errors"
	"github.com/go-zookeeper/zk"
	"github.com/gogf/gf/v2/net/gsvc"
	"golang.org/x/sync/singleflight"
	"path"
	"strings"
)

var _ gsvc.Watcher = (*watcher)(nil)

var ErrWatcherStopped = errors.New("watcher stopped")

type watcher struct {
	ctx       context.Context
	event     chan zk.Event
	conn      *zk.Conn
	cancel    context.CancelFunc
	prefix    string
	nameSpace string
	group     singleflight.Group
}

func newWatcher(ctx context.Context, nameSpace, prefix string, conn *zk.Conn) (*watcher, error) {
	w := &watcher{
		conn:      conn,
		event:     make(chan zk.Event, 1),
		nameSpace: nameSpace,
		prefix:    prefix,
	}
	w.ctx, w.cancel = context.WithCancel(ctx)
	go w.watch(w.ctx)
	return w, nil
}

func (w *watcher) Proceed() ([]gsvc.Service, error) {
	select {
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	case e := <-w.event:
		if e.State == zk.StateDisconnected {
			return nil, ErrWatcherStopped
		}
		if e.Err != nil {
			return nil, e.Err
		}
		return w.getServicesByPrefix()
	}
}

func (w *watcher) getServicesByPrefix() ([]gsvc.Service, error) {
	prefix := strings.TrimPrefix(strings.ReplaceAll(w.prefix, "/", "-"), "-")
	serviceNamePath := path.Join(w.nameSpace, prefix)
	instances, err, _ := w.group.Do(serviceNamePath, func() (interface{}, error) {
		servicesID, _, err := w.conn.Children(serviceNamePath)
		if err != nil {
			return nil, err
		}
		items := make([]gsvc.Service, 0, len(servicesID))
		for _, service := range servicesID {
			servicePath := path.Join(serviceNamePath, service)
			byteData, _, err := w.conn.Get(servicePath)
			if err != nil {
				return nil, err
			}
			item, err := unmarshal(byteData)
			if err != nil {
				return nil, err
			}
			svc, err := gsvc.NewServiceWithKV(item.Key, item.Value)
			if err != nil {
				return nil, err
			}
			items = append(items, svc)
		}
		return items, nil
	})
	if err != nil {
		return nil, err
	}
	return instances.([]gsvc.Service), nil
}

func (w *watcher) Close() error {
	w.cancel()
	return nil
}

func (w *watcher) watch(ctx context.Context) {
	prefix := strings.TrimPrefix(strings.ReplaceAll(w.prefix, "/", "-"), "-")
	serviceNamePath := path.Join(w.nameSpace, prefix)
	for {

		if w.conn.State() == zk.StateConnected || w.conn.State() == zk.StateHasSession {
			// 每次 watch 只有一次有效期 所以循环 watch
			_, _, ch, err := w.conn.ChildrenW(serviceNamePath)
			if err != nil {
				w.event <- zk.Event{Err: err}
			}
			select {
			case <-ctx.Done():
				return
			default:
				w.event <- <-ch
			}
		}
	}
}
