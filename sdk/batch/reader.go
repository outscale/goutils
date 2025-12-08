package batch

import (
	"context"
	"maps"
	"slices"
	"time"

	"github.com/outscale/goutils/sdk/log"
)

type result[T any] struct {
	ok  *T
	err error
}

func resultOk[T any](t *T) result[T] {
	return result[T]{ok: t}
}

func resultError[T any](err error) result[T] {
	return result[T]{err: err}
}

type watcher[T any] struct {
	ctx   context.Context // the context of the waiter
	id    string
	until func(r *T) (ok bool, err error)
	resp  chan result[T]
}

type Batcher[T any] struct {
	interval time.Duration
	refresh  func(ctx context.Context, ids []string) (map[string]*T, error)

	in       chan watcher[T]
	watchers map[string]watcher[T]
}

func NewBatcher[T any](interval time.Duration,
	refresh func(ctx context.Context, ids []string) (map[string]*T, error)) *Batcher[T] {
	return &Batcher[T]{
		interval: interval,
		refresh:  refresh,

		in:       make(chan watcher[T]),
		watchers: make(map[string]watcher[T]),
	}
}

func (sw *Batcher[T]) Run(ctx context.Context) {
	t := time.NewTicker(sw.interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case in := <-sw.in:
			sw.watchers[in.id] = in
		case <-t.C:
			if len(sw.watchers) == 0 {
				continue
			}
			ids := slices.Collect(maps.Keys(sw.watchers))
			log.Default.Info(ctx, "Watching resources", "count", len(ids))
			rsrcs, err := sw.refresh(ctx, ids)
			if err != nil {
				log.Default.Error(ctx, err, "unable to check statuses")
				continue
			}
			for id, rsrc := range rsrcs {
				w, found := sw.watchers[id]
				if !found { // should not occur
					continue
				}
				ok, err := w.until(rsrc)
				switch {
				case ok:
					log.Default.Info(ctx, "Resource is ok", "id", id)
					sw.response(ctx, w, resultOk(rsrc))
					delete(sw.watchers, id)
					close(w.resp)
				case err != nil:
					log.Default.Info(ctx, "Resource is in error", "id", id)
					sw.response(ctx, w, resultError[T](err))
					delete(sw.watchers, id)
					close(w.resp)
				default:
					log.Default.Info(ctx, "Resource is not ready", "id", id)
				}
			}
		}
	}
}

func (sw *Batcher[T]) response(ctx context.Context, w watcher[T], res result[T]) {
	select {
	case <-ctx.Done():
	case <-w.ctx.Done(): // we do not want to block if the waiter context has been cancelled.
	case w.resp <- res: // send the response, would block without the previous test...
	}
}

func (sw *Batcher[T]) WaitUntil(ctx context.Context, id string, until func(r *T) (ok bool, err error)) (r *T, err error) {
	start := time.Now()
	defer func() {
		log.Default.Info(ctx, "End of wait", "success", err == nil, "duration", time.Since(start))
	}()
	resp := make(chan result[T], 1)
	w := watcher[T]{ctx: ctx, id: id, until: until, resp: resp}
	// send (unless context has been cancelled)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case sw.in <- w:
	}
	// receive (unless context has been cancelled)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-w.resp:
		return res.ok, res.err
	}
}

func (sw *Batcher[T]) Read(ctx context.Context, id string) (r *T, err error) {
	return sw.WaitUntil(ctx, id, func(r *T) (ok bool, err error) { return true, nil })
}
