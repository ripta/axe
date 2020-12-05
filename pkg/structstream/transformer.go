package structstream

import (
	"encoding/json"
	"strings"
	"time"
)

type Transformer func(string, string) (Structline, bool)

func CombineTransformers(strict bool, ts ...Transformer) Transformer {
	return func(meta, in string) (Structline, bool) {
		for _, t := range ts {
			s, ok := t(meta, in)
			if ok {
				if !strict || s.Complete {
					return s, true
				}
			}
		}
		return Structline{}, false
	}
}

func GlogTransformer(meta, in string) (Structline, bool) {
	s := Structline{
		Type:      "glog",
		Meta:      meta,
		Timestamp: time.Now(),
		KV:        make(map[string]interface{}),
	}
	if len(in) < 22 {
		return s, false
	}

	v, ok := map[byte]string{
		'I': "INFO",
		'W': "WARNING",
		'E': "ERROR",
		'F': "FATAL",
	}[in[0]]
	if !ok {
		return s, false
	}

	s.Priority = v

	segs := strings.SplitN(in, "] ", 2)
	if len(segs) != 2 {
		return s, false
	}

	s.Message = segs[1]

	fields := strings.Fields(segs[0])
	if len(fields) != 4 {
		return s, false
	}

	s.KV["pid"] = fields[2]
	s.KV["fileline"] = fields[3]

	cyear := time.Now().Format("2006")
	cdate := fields[0][1:]
	ctime := fields[1]
	if len(cdate) == 4 {
		t, err := time.Parse("20060102 15:04:05.999999999", cyear+cdate+" "+ctime)
		if err != nil {
			return s, true
		}
		s.Timestamp = t
	} else if len(cdate) == 8 {
		t, err := time.Parse("20060102 15:04:05.999999999", cdate+" "+ctime)
		if err != nil {
			return s, true
		}
		s.Timestamp = t
	}

	s.Complete = true
	return s, true
}

func JSONTransformer(meta, in string) (Structline, bool) {
	s := Structline{
		Type:      "json",
		Meta:      meta,
		Timestamp: time.Now(),
		KV:        make(map[string]interface{}),
	}
	if err := json.Unmarshal([]byte(in), &s.KV); err != nil {
		return s, false
	}

	s.Complete = true

	tryFields(s.KV, []string{"ts", "timestamp", "@ts", "@timestamp"}, func(vs string) bool {
		for _, format := range []string{time.RFC822Z, time.RFC1123Z, time.RFC3339, time.RFC3339Nano} {
			t, err := time.Parse(format, vs)
			if err != nil {
				continue
			}

			s.Timestamp = t
			return true
		}
		return false
	})

	tryFields(s.KV, []string{"priority", "prio", "@priority", "@prio"}, func(vs string) bool {
		s.Priority = vs
		return true
	})

	tryFields(s.KV, []string{"msg", "message", "mesg"}, func(vs string) bool {
		s.Message = vs
		return true
	})

	return s, true
}

func PassthruTransformer(meta, in string) (Structline, bool) {
	s := Structline{
		Type:     "passthru",
		Complete: true,

		Message:   in,
		Meta:      meta,
		Timestamp: time.Now(),
	}
	return s, true
}

func tryFields(kv map[string]interface{}, fields []string, fn func(string) bool) {
	if kv == nil {
		return
	}

	for _, field := range fields {
		vi, ok := kv[field]
		if !ok {
			continue
		}
		vs, ok := vi.(string)
		if !ok {
			continue
		}
		if fn(vs) {
			delete(kv, field)
			return
		}
	}
}
