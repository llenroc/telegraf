package models

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/influxdata/telegraf"

	"github.com/stretchr/testify/assert"
)

func TestMakeMetricNoFields(t *testing.T) {
	now := time.Now()
	ri := RunningInput{
		Config: &InputConfig{
			Name: "TestRunningInput",
		},
	}

	m := ri.MakeMetric(
		"RITest",
		map[string]interface{}{},
		map[string]string{},
		telegraf.Untyped,
		now,
	)
	assert.Nil(t, m)
}

// nil fields should get dropped
func TestMakeMetricNilFields(t *testing.T) {
	now := time.Now()
	ri := RunningInput{
		Config: &InputConfig{
			Name: "TestRunningInput",
		},
	}

	m := ri.MakeMetric(
		"RITest",
		map[string]interface{}{
			"value": int(101),
			"nil":   nil,
		},
		map[string]string{},
		telegraf.Untyped,
		now,
	)
	assert.Equal(
		t,
		fmt.Sprintf("RITest value=101i %d\n", now.UnixNano()),
		m.String(),
	)
}

// make an untyped, counter, & gauge metric
func TestMakeMetric(t *testing.T) {
	now := time.Now()
	ri := RunningInput{
		Config: &InputConfig{
			Name: "TestRunningInput",
		},
	}
	ri.SetDebug(true)
	assert.Equal(t, true, ri.Debug())
	ri.SetTrace(true)
	assert.Equal(t, true, ri.Trace())
	assert.Equal(t, "inputs.TestRunningInput", ri.Name())

	m := ri.MakeMetric(
		"RITest",
		map[string]interface{}{"value": int(101)},
		map[string]string{},
		telegraf.Untyped,
		now,
	)
	assert.Equal(
		t,
		fmt.Sprintf("RITest value=101i %d\n", now.UnixNano()),
		m.String(),
	)
	assert.Equal(
		t,
		m.Type(),
		telegraf.Untyped,
	)

	m = ri.MakeMetric(
		"RITest",
		map[string]interface{}{"value": int(101)},
		map[string]string{},
		telegraf.Counter,
		now,
	)
	assert.Equal(
		t,
		fmt.Sprintf("RITest value=101i %d\n", now.UnixNano()),
		m.String(),
	)
	assert.Equal(
		t,
		m.Type(),
		telegraf.Counter,
	)

	m = ri.MakeMetric(
		"RITest",
		map[string]interface{}{"value": int(101)},
		map[string]string{},
		telegraf.Gauge,
		now,
	)
	assert.Equal(
		t,
		fmt.Sprintf("RITest value=101i %d\n", now.UnixNano()),
		m.String(),
	)
	assert.Equal(
		t,
		m.Type(),
		telegraf.Gauge,
	)
}

func TestMakeMetricWithPluginTags(t *testing.T) {
	now := time.Now()
	ri := RunningInput{
		Config: &InputConfig{
			Name: "TestRunningInput",
			Tags: map[string]string{
				"foo": "bar",
			},
		},
	}
	ri.SetDebug(true)
	assert.Equal(t, true, ri.Debug())
	ri.SetTrace(true)
	assert.Equal(t, true, ri.Trace())

	m := ri.MakeMetric(
		"RITest",
		map[string]interface{}{"value": int(101)},
		nil,
		telegraf.Untyped,
		now,
	)
	assert.Equal(
		t,
		fmt.Sprintf("RITest,foo=bar value=101i %d\n", now.UnixNano()),
		m.String(),
	)
}

func TestMakeMetricFilteredOut(t *testing.T) {
	now := time.Now()
	ri := RunningInput{
		Config: &InputConfig{
			Name: "TestRunningInput",
			Tags: map[string]string{
				"foo": "bar",
			},
			Filter: Filter{NamePass: []string{"foobar"}},
		},
	}
	ri.SetDebug(true)
	assert.Equal(t, true, ri.Debug())
	ri.SetTrace(true)
	assert.Equal(t, true, ri.Trace())
	assert.NoError(t, ri.Config.Filter.Compile())

	m := ri.MakeMetric(
		"RITest",
		map[string]interface{}{"value": int(101)},
		nil,
		telegraf.Untyped,
		now,
	)
	assert.Nil(t, m)
}

func TestMakeMetricWithDaemonTags(t *testing.T) {
	now := time.Now()
	ri := RunningInput{
		Config: &InputConfig{
			Name: "TestRunningInput",
		},
	}
	ri.SetDefaultTags(map[string]string{
		"foo": "bar",
	})
	ri.SetDebug(true)
	assert.Equal(t, true, ri.Debug())
	ri.SetTrace(true)
	assert.Equal(t, true, ri.Trace())

	m := ri.MakeMetric(
		"RITest",
		map[string]interface{}{"value": int(101)},
		map[string]string{},
		telegraf.Untyped,
		now,
	)
	assert.Equal(
		t,
		fmt.Sprintf("RITest,foo=bar value=101i %d\n", now.UnixNano()),
		m.String(),
	)
}

// make an untyped, counter, & gauge metric
func TestMakeMetricInfFields(t *testing.T) {
	inf := math.Inf(1)
	ninf := math.Inf(-1)
	now := time.Now()
	ri := RunningInput{
		Config: &InputConfig{
			Name: "TestRunningInput",
		},
	}
	ri.SetDebug(true)
	assert.Equal(t, true, ri.Debug())
	ri.SetTrace(true)
	assert.Equal(t, true, ri.Trace())

	m := ri.MakeMetric(
		"RITest",
		map[string]interface{}{
			"value": int(101),
			"inf":   inf,
			"ninf":  ninf,
		},
		map[string]string{},
		telegraf.Untyped,
		now,
	)
	assert.Equal(
		t,
		fmt.Sprintf("RITest value=101i %d\n", now.UnixNano()),
		m.String(),
	)
}

func TestMakeMetricAllFieldTypes(t *testing.T) {
	now := time.Now()
	ri := RunningInput{
		Config: &InputConfig{
			Name: "TestRunningInput",
		},
	}
	ri.SetDebug(true)
	assert.Equal(t, true, ri.Debug())
	ri.SetTrace(true)
	assert.Equal(t, true, ri.Trace())

	m := ri.MakeMetric(
		"RITest",
		map[string]interface{}{
			"a": int(10),
			"b": int8(10),
			"c": int16(10),
			"d": int32(10),
			"e": uint(10),
			"f": uint8(10),
			"g": uint16(10),
			"h": uint32(10),
			"i": uint64(10),
			"j": float32(10),
			"k": uint64(9223372036854775810),
			"l": "foobar",
			"m": true,
		},
		map[string]string{},
		telegraf.Untyped,
		now,
	)
	assert.Contains(t, m.String(), "a=10i")
	assert.Contains(t, m.String(), "b=10i")
	assert.Contains(t, m.String(), "c=10i")
	assert.Contains(t, m.String(), "d=10i")
	assert.Contains(t, m.String(), "e=10i")
	assert.Contains(t, m.String(), "f=10i")
	assert.Contains(t, m.String(), "g=10i")
	assert.Contains(t, m.String(), "h=10i")
	assert.Contains(t, m.String(), "i=10i")
	assert.Contains(t, m.String(), "j=10")
	assert.NotContains(t, m.String(), "j=10i")
	assert.Contains(t, m.String(), "k=9223372036854775807i")
	assert.Contains(t, m.String(), "l=\"foobar\"")
	assert.Contains(t, m.String(), "m=true")
}

func TestMakeMetricNameOverride(t *testing.T) {
	now := time.Now()
	ri := RunningInput{
		Config: &InputConfig{
			Name:         "TestRunningInput",
			NameOverride: "foobar",
		},
	}

	m := ri.MakeMetric(
		"RITest",
		map[string]interface{}{"value": int(101)},
		map[string]string{},
		telegraf.Untyped,
		now,
	)
	assert.Equal(
		t,
		fmt.Sprintf("foobar value=101i %d\n", now.UnixNano()),
		m.String(),
	)
}

func TestMakeMetricNamePrefix(t *testing.T) {
	now := time.Now()
	ri := RunningInput{
		Config: &InputConfig{
			Name:              "TestRunningInput",
			MeasurementPrefix: "foobar_",
		},
	}

	m := ri.MakeMetric(
		"RITest",
		map[string]interface{}{"value": int(101)},
		map[string]string{},
		telegraf.Untyped,
		now,
	)
	assert.Equal(
		t,
		fmt.Sprintf("foobar_RITest value=101i %d\n", now.UnixNano()),
		m.String(),
	)
}

func TestMakeMetricNameSuffix(t *testing.T) {
	now := time.Now()
	ri := RunningInput{
		Config: &InputConfig{
			Name:              "TestRunningInput",
			MeasurementSuffix: "_foobar",
		},
	}

	m := ri.MakeMetric(
		"RITest",
		map[string]interface{}{"value": int(101)},
		map[string]string{},
		telegraf.Untyped,
		now,
	)
	assert.Equal(
		t,
		fmt.Sprintf("RITest_foobar value=101i %d\n", now.UnixNano()),
		m.String(),
	)
}
