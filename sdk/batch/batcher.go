/*
SPDX-FileCopyrightText: 2025 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package batch

import (
	"context"
	"errors"
	"reflect"
	"slices"
	"time"

	"github.com/outscale/goutils/sdk/log"
)

var ErrNotFound = errors.New("not found")

type result[R any] struct {
	result *R
	err    error
}

func resultOk[R any](t *R) result[R] {
	return result[R]{result: t}
}

func resultError[R any](err error) result[R] {
	return result[R]{err: err}
}

type watcher[Q, R any] struct {
	ctx   context.Context // the context of the waiter
	query Q
	until func(r *R) (ok bool, err error)
	resp  chan result[R]
}

type resultFn[Q, R any] func(query Q) (*R, bool)
type batcher[Q, R any] struct {
	interval time.Duration
	refresh  func(ctx context.Context, queries []Q) (resultFn[Q, R], error)
	merge    func(query Q, queries []Q) ([]Q, bool)
	in       chan watcher[Q, R]
	batches  []batch[Q, R]
}

type batch[Q, R any] struct {
	query    []Q
	watchers []watcher[Q, R]
}

func newBatcher[Q, R any](interval time.Duration,
	refresh func(ctx context.Context, queries []Q) (resultFn[Q, R], error),
	merge func(query Q, queries []Q) ([]Q, bool),
) *batcher[Q, R] {
	return &batcher[Q, R]{
		interval: interval,
		refresh:  refresh,
		merge:    merge,
		in:       make(chan watcher[Q, R]),
	}
}

func (b *batcher[Q, R]) Run(ctx context.Context) {
	t := time.NewTicker(b.interval)
	defer t.Stop()
LOOPBATCHER:
	for {
		select {
		case <-ctx.Done():
			return
		case in := <-b.in:
			for i := range b.batches {
				if newQuery, ok := b.merge(in.query, b.batches[i].query); ok {
					b.batches[i].query = newQuery
					b.batches[i].watchers = append(b.batches[i].watchers, in)
					continue LOOPBATCHER
				}
			}
			b.batches = append(b.batches, batch[Q, R]{
				query:    []Q{in.query},
				watchers: []watcher[Q, R]{in},
			})
		case <-t.C:
			needClean := false
			for i := range b.batches {
				batch := &b.batches[i]
				log.Default.Info(ctx, "Watching resources", "count", len(batch.query))
				result, err := b.refresh(ctx, batch.query)
				if err != nil {
					log.Default.Error(ctx, err, "unable to check statuses")
					continue
				}
				var left []watcher[Q, R]
				for _, w := range batch.watchers {
					res, found := result(w.query)
					if !found {
						log.Default.Info(ctx, "Resource is not found", "id", w.query)
						b.response(ctx, w, resultError[R](ErrNotFound))
						close(w.resp)
						continue
					}
					ok, err := w.until(res)
					switch {
					case ok:
						log.Default.Info(ctx, "Resource is ok", "id", w.query)
						b.response(ctx, w, resultOk(res))
						close(w.resp)
					case err != nil:
						log.Default.Info(ctx, "Resource is in error", "id", w.query)
						b.response(ctx, w, resultError[R](err))
						close(w.resp)
					default:
						log.Default.Info(ctx, "Resource is not ready", "id", w.query)
						left = append(left, w)
					}
				}
				batch.watchers = left
				if len(left) == 0 {
					needClean = true
				}
			}
			if needClean {
				var newbatches []batch[Q, R]
				for _, batch := range b.batches {
					if len(batch.watchers) > 0 {
						newbatches = append(newbatches, batch)
					}
				}
				b.batches = newbatches
			}
		}
	}
}

func (b *batcher[Q, R]) response(ctx context.Context, w watcher[Q, R], res result[R]) {
	select {
	case <-ctx.Done():
	case <-w.ctx.Done(): // we do not want to block if the waiter context has been cancelled.
	case w.resp <- res: // send the response, would block without the previous test...
	}
}

// BatcherByID batches all reads by ID in a single Read call.
type BatcherByID[R any] struct {
	*batcher[string, R]
}

// WaitUntil repeatedly reads the resource until the until func returns either true or an error.
func (b *BatcherByID[R]) WaitUntil(ctx context.Context, id string, until func(r *R) (ok bool, err error)) (r *R, err error) {
	start := time.Now()
	defer func() {
		log.Default.Info(ctx, "End of wait", "success", err == nil, "duration", time.Since(start))
	}()
	resp := make(chan result[R], 1)
	w := watcher[string, R]{ctx: ctx, query: id, until: until, resp: resp}
	// send (unless context has been cancelled)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case b.in <- w:
	}
	// receive (unless context has been cancelled)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-w.resp:
		return res.result, res.err
	}
}

// Read reads a resource.
func (b *BatcherByID[R]) Read(ctx context.Context, id string) (r *R, err error) {
	return b.WaitUntil(ctx, id, func(r *R) (ok bool, err error) { return true, nil })
}

func NewBatcherByID[R any](interval time.Duration,
	refresh func(ctx context.Context, ids []string) (resultFn[string, R], error),
) *BatcherByID[R] {
	return &BatcherByID[R]{
		batcher: newBatcher(interval, refresh,
			func(query string, queries []string) ([]string, bool) { // merge
				if slices.Contains(queries, query) {
					return queries, true
				}
				return append(queries, query), true
			},
		),
	}
}

// BatcherSameQuery batches all queries of type Q that are equal in a single call.
type BatcherSameQuery[Q, R any] struct {
	*batcher[Q, R]
}

// Read returns all resources matching the query Q, without pagination.
func (b *BatcherSameQuery[Q, R]) Read(ctx context.Context, query Q) (r *R, err error) {
	start := time.Now()
	defer func() {
		log.Default.Info(ctx, "End of wait", "success", err == nil, "duration", time.Since(start))
	}()
	resp := make(chan result[R], 1)
	w := watcher[Q, R]{ctx: ctx, query: query, until: func(_ *R) (ok bool, err error) { return true, nil }, resp: resp}
	// send (unless context has been cancelled)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case b.in <- w:
	}
	// receive (unless context has been cancelled)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-w.resp:
		return res.result, res.err
	}
}

func NewBatcherSameQuery[Q, R any](interval time.Duration,
	refresh func(ctx context.Context, queries []Q) (resultFn[Q, R], error),
) *BatcherSameQuery[Q, R] {
	return &BatcherSameQuery[Q, R]{
		batcher: newBatcher(interval, refresh, func(query Q, queries []Q) ([]Q, bool) { // merge
			if len(queries) == 0 {
				return nil, false
			}
			if reflect.DeepEqual(query, queries[0]) {
				return append(queries, query), true
			}
			return nil, false
		}),
	}
}
