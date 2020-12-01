package structstream

import (
	"sync"
	"time"

	"github.com/dgraph-io/ristretto"
)

type Buffer struct {
	cap    int
	mu     sync.RWMutex
	parser Transformer

	metas  []string
	raw    []string
	parsed *ristretto.Cache
}

type Structline struct {
	Complete bool
	Type     string

	Message   string
	Meta      string
	Priority  string
	Timestamp time.Time

	KV map[string]interface{}
}

func New(cap int, tr Transformer) (*Buffer, error) {
	cfg := &ristretto.Config{
		NumCounters: 10 * int64(cap),
		MaxCost:     200 * int64(cap),
		BufferItems: 64,
	}
	cache, err := ristretto.NewCache(cfg)
	if err != nil {
		return nil, err
	}

	return &Buffer{
		cap:    cap,
		mu:     sync.RWMutex{},
		parser: tr,

		metas:  make([]string, 0, cap),
		raw:    make([]string, 0, cap),
		parsed: cache,
	}, nil
}

func (b *Buffer) Append(meta, line string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.raw = append(b.raw, line)
	b.metas = append(b.metas, meta)
}

func (b *Buffer) Clear() {
	b.mu.Lock()
	b.mu.Unlock()
	b.raw = b.raw[:0]
	b.metas = b.metas[:0]
	b.parsed.Clear()
}

func (b *Buffer) GetAt(loc int) Structline {
	ss := b.GetRange(loc, loc)
	if ss == nil || len(ss) != 1 {
		return Structline{}
	}
	return ss[0]
}

func (b *Buffer) GetRange(st, fi int) []Structline {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if st < 0 {
		st = 0
	}
	if fi >= len(b.raw) {
		fi = len(b.raw) - 1
	}

	ss := make([]Structline, 0)
	for i := st; i <= fi; i++ {
		cs, ok := b.parsed.Get(i)
		if ok {
			ss = append(ss, cs.(Structline))
			continue
		}

		s, ok := b.parser(b.metas[i], b.raw[i])
		if !ok {
			continue
		}

		_ = b.parsed.Set(i, s, int64(len(b.raw[i])))
		ss = append(ss, s)
	}

	return ss
}

func (b *Buffer) Len() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.raw)
}
