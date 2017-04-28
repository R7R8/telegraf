package sum

import (
	"testing"
	"time"

	"github.com/influxdata/telegraf/metric"
	"github.com/influxdata/telegraf/testutil"
)

var m1, _ = metric.New("m1",
	map[string]string{"foo": "bar"},
	map[string]interface{}{
		"a": int64(1),
		"b": int64(1),
		"c": int64(1),
		"d": int64(1),
		"e": int64(1),
		"f": float64(2),
		"g": float64(2),
		"h": float64(2),
		"i": float64(2),
		"j": float64(3),
	},
	time.Now(),
)
var m2, _ = metric.New("m1",
	map[string]string{"foo": "bar"},
	map[string]interface{}{
		"a":        int64(1),
		"b":        int64(3),
		"c":        int64(3),
		"d":        int64(3),
		"e":        int64(3),
		"f":        float64(1),
		"g":        float64(1),
		"h":        float64(1),
		"i":        float64(1),
		"j":        float64(1),
		"k":        float64(200),
		"ignoreme": "string",
		"andme":    true,
	},
	time.Now(),
)

func BenchmarkApply(b *testing.B) {
	minmax := NewSum()

	for n := 0; n < b.N; n++ {
		minmax.Add(m1)
		minmax.Add(m2)
	}
}

// Test two metrics getting added.
func TestSumWithPeriod(t *testing.T) {
	acc := testutil.Accumulator{}
	minmax := NewSum()

	minmax.Add(m1)
	minmax.Add(m2)
	minmax.Push(&acc)

	expectedFields := map[string]interface{}{
		"a": float64(2),
		"b": float64(4),
		"c": float64(4),
		"d": float64(4),
		"e": float64(4),
		"f": float64(3),
		"g": float64(3),
		"h": float64(3),
		"i": float64(3),
		"j": float64(4),
		"k": float64(200),
	}
	expectedTags := map[string]string{
		"foo": "bar",
	}
	acc.AssertContainsTaggedFields(t, "m1", expectedFields, expectedTags)
}

// Test two metrics getting added with a push/reset in between (simulates
// getting added in different periods.)
func TestSumDifferentPeriods(t *testing.T) {
	acc := testutil.Accumulator{}
	minmax := NewSum()

	minmax.Add(m1)
	minmax.Push(&acc)
	expectedFields := map[string]interface{}{
		"a": float64(1),
		"b": float64(1),
		"c": float64(1),
		"d": float64(1),
		"e": float64(1),
		"f": float64(2),
		"g": float64(2),
		"h": float64(2),
		"i": float64(2),
		"j": float64(3),
	}
	expectedTags := map[string]string{
		"foo": "bar",
	}
	acc.AssertContainsTaggedFields(t, "m1", expectedFields, expectedTags)

	acc.ClearMetrics()
	minmax.Reset()
	minmax.Add(m2)
	minmax.Push(&acc)
	expectedFields = map[string]interface{}{
		"a": float64(1),
		"b": float64(3),
		"c": float64(3),
		"d": float64(3),
		"e": float64(3),
		"f": float64(1),
		"g": float64(1),
		"h": float64(1),
		"i": float64(1),
		"j": float64(1),
		"k": float64(200),
	}
	expectedTags = map[string]string{
		"foo": "bar",
	}
	acc.AssertContainsTaggedFields(t, "m1", expectedFields, expectedTags)
}
