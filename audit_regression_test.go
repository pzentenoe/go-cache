package cache

import (
	"bytes"
	"encoding/gob"
	"errors"
	"math"
	"testing"
	"time"
)

type panicEncoder struct{}

func (panicEncoder) GobEncode() ([]byte, error) { panic("encode") }

func TestAuditSavePanicUnlocks(t *testing.T) {
	sc := newShardedCache(1, NoExpiration)
	sc.Set("x", panicEncoder{}, NoExpiration)
	if err := sc.Save(&bytes.Buffer{}); err == nil || err.Error() != errGobRegistration {
		t.Fatalf("unexpected error: %v", err)
	}
	if !sc.cs[0].mu.TryLock() {
		t.Fatal("shard remains locked after panic")
	}
	sc.cs[0].mu.Unlock()
	sc.Set("x", 1, NoExpiration)
}

func TestAuditGenericIntegerDeltas(t *testing.T) {
	cases := []struct {
		name    string
		initial any
		delta   int64
		dec     bool
		want    any
		err     error
	}{
		{"signed-wide-valid", int8(-128), 255, false, int8(127), nil},
		{"signed-wide-invalid", int8(0), 256, false, int8(0), ErrOverflow},
		{"signed-decrement-wide", int8(127), 255, true, int8(-128), nil},
		{"unsigned-negative", uint8(5), -3, false, uint8(2), nil},
		{"unsigned-negative-decrement", uint8(5), -3, true, uint8(8), nil},
		{"unsigned-wide-invalid", uint8(0), 256, false, uint8(0), ErrOverflow},
		{"min-delta", uint64(1) << 63, math.MinInt64, false, uint64(0), nil},
		{"min-decrement", uint64(0), math.MinInt64, true, uint64(1) << 63, nil},
		{"high-overflow", uint64(math.MaxUint64), -1, true, uint64(math.MaxUint64), ErrUnderflow},
		{"max-delta", int8(0), math.MaxInt64, false, int8(0), ErrOverflow},
		{"signed-min-valid", int64(0), math.MinInt64, false, int64(math.MinInt64), nil},
		{"signed-min-subtract", int64(-1), math.MinInt64, true, int64(math.MaxInt64), nil},
		{"signed-min-overflow", int64(0), math.MinInt64, true, int64(0), ErrUnderflow},
		{"unsigned-high-minus", uint64(math.MaxUint64), math.MaxInt64, true, uint64(1) << 63, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, c := range []ShardedCache{unexportedNewSharded(NoExpiration, 0, 2)} {
				c.Set("k", tc.initial, NoExpiration)
				var err error
				if tc.dec {
					err = c.Decrement("k", tc.delta)
				} else {
					err = c.Increment("k", tc.delta)
				}
				if !errors.Is(err, tc.err) {
					t.Errorf("error %v, want %v", err, tc.err)
				}
				got, _ := c.Get("k")
				if got != tc.want {
					t.Errorf("got %v (%T), want %v (%T)", got, got, tc.want, tc.want)
				}
			}
		})
	}
}

func TestAuditLoadDecodeBeforeApply(t *testing.T) {
	sc := newShardedCache(2, NoExpiration)
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(2); err != nil {
		t.Fatal(err)
	}
	if err := enc.Encode(map[string]Item{"new": {Object: "value"}}); err != nil {
		t.Fatal(err)
	}
	if err := sc.Load(&b); err == nil {
		t.Fatal("expected decode failure")
	}
	if _, ok := sc.Get("new"); ok {
		t.Fatal("partial load on decode failure")
	}
}

func TestAuditIntervalNonblocking(t *testing.T) {
	j := newJanitor(time.Second, func() {})
	j.setInterval(2 * time.Second)
	done := make(chan struct{})
	go func() { j.setInterval(3 * time.Second); close(done) }()
	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		<-j.updateInterval
		<-done
		t.Fatal("second update blocked")
	}
}

func TestAuditJanitorCallbackUpdates(t *testing.T) {
	for _, sharded := range []bool{false, true} {
		t.Run(map[bool]string{false: "cache", true: "sharded"}[sharded], func(t *testing.T) {
			var c interface {
				Set(string, any, time.Duration)
				OnEvicted(func(string, any))
				SetJanitorInterval(time.Duration)
				Close()
			}
			var j *janitor
			if sharded {
				sc := unexportedNewSharded(NoExpiration, time.Millisecond, 2)
				c = sc
				j = sc.janitor
			} else {
				cache := New(NoExpiration, time.Millisecond)
				c = cache
				j = cache.janitor
			}
			defer c.Close()
			done := make(chan struct{})
			c.OnEvicted(func(string, any) { c.SetJanitorInterval(time.Hour); c.SetJanitorInterval(2 * time.Hour); close(done) })
			c.Set("expired", 1, time.Nanosecond)
			select {
			case <-done:
			case <-time.After(2 * time.Second):
				t.Fatal("callback did not finish")
			}
			deadline := time.After(2 * time.Second)
			for {
				j.mu.Lock()
				interval := j.interval
				j.mu.Unlock()
				if interval == 2*time.Hour {
					break
				}
				select {
				case <-deadline:
					t.Fatalf("latest interval not applied: %v", interval)
				case <-time.After(time.Millisecond):
				}
			}
		})
	}
}

func TestAuditJanitorUpdateCloseRace(t *testing.T) {
	for _, sharded := range []bool{false, true} {
		for i := 0; i < 20; i++ {
			var c interface {
				SetJanitorInterval(time.Duration)
				Close()
			}
			if sharded {
				c = unexportedNewSharded(NoExpiration, time.Hour, 2)
			} else {
				c = New(NoExpiration, time.Hour)
			}
			done := make(chan struct{})
			go func() {
				defer close(done)
				for n := 0; n < 100; n++ {
					c.SetJanitorInterval(time.Duration(n))
					c.SetJanitorInterval(-1)
				}
			}()
			c.Close()
			<-done
		}
	}
}

func TestAuditGenericDeltaPreservesExpiryAndErrorText(t *testing.T) {
	c := New(NoExpiration, 0)
	c.Set("key", int8(-128), time.Hour)
	_, before, _ := c.GetWithExpiration("key")
	if err := c.Increment("key", 255); err != nil {
		t.Fatal(err)
	}
	value, after, ok := c.GetWithExpiration("key")
	if !ok || value != int8(127) || !before.Equal(after) {
		t.Fatalf("value=%v expiry=%v want %v", value, after, before)
	}
	err := c.Increment("key", 256)
	if !errors.Is(err, ErrOverflow) || err.Error() != "overflow would occur for key" {
		t.Fatalf("unexpected error %v", err)
	}
	value, after, ok = c.GetWithExpiration("key")
	if !ok || value != int8(127) || !before.Equal(after) {
		t.Fatalf("value=%v expiry=%v", value, after)
	}
	err = c.Decrement("key", -256)
	if !errors.Is(err, ErrUnderflow) || err.Error() != "underflow would occur for key" {
		t.Fatalf("unexpected error %v", err)
	}
}
