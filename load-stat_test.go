package main

import (
	"reflect"
	"testing"
	"time"
)

func TestStatsBookkeeper_AddDoneActions(t *testing.T) {
	s := statsBookkeeper{}
	now := time.Now()

	s.addDoneActions(23, now.Add(-time.Minute*16), now.Add(-time.Minute*14-time.Second*58))

	assertDeepEqual(t, s.history, []doneActions{
		{23, now.Add(-time.Minute * 16), now.Add(-time.Minute*14 - time.Second*58)},
	})

	time.Sleep(time.Second * 3)
	s.addDoneActions(42, now.Add(-time.Minute*4), now.Add(-time.Minute*3))

	assertDeepEqual(t, s.history, []doneActions{
		{42, now.Add(-time.Minute * 4), now.Add(-time.Minute * 3)},
	})

	time.Sleep(time.Second * 1)
	s.addDoneActions(777, now.Add(-time.Minute*2), now.Add(-time.Minute))

	assertDeepEqual(t, s.history, []doneActions{
		{42, now.Add(-time.Minute * 4), now.Add(-time.Minute * 3)},
		{777, now.Add(-time.Minute * 2), now.Add(-time.Minute)},
	})
}

func TestStatsBookkeeper_QueryLoad(t *testing.T) {
	now := time.Now()

	assertDeepEqual(t, (&statsBookkeeper{}).queryLoad(), [3]float64{0.0, 0.0, 0.0})

	assertDeepEqual(t, (&statsBookkeeper{history: []doneActions{
		{22, now.Add(-time.Minute * 17), now.Add(-time.Minute * 16)},
		{23, now.Add(-time.Minute * 14), now.Add(-time.Minute * 6)},
		{42, now.Add(-time.Minute * 4), now.Add(-time.Minute)},
		{777, now.Add(-time.Second * 50), now.Add(-time.Second * 10)},
		{778, now.Add(time.Minute), now.Add(time.Minute * 2)},
	}}).queryLoad(), [3]float64{777.0, 777.0 + 42.0, 777.0 + 42.0 + 23.0})

	{
		load15 := (&statsBookkeeper{history: []doneActions{
			{23, now.Add(-time.Minute * 24), now.Add(-time.Minute * 6)},
		}}).queryLoad()[2]
		if !(11 <= load15 && load15 <= 12) {
			t.Errorf("11 <= %#v && %#v <= 12", load15, load15)
		}
	}

	{
		load15 := (&statsBookkeeper{history: []doneActions{
			{23, now.Add(-time.Minute * 14), now.Add(time.Minute * 14)},
		}}).queryLoad()[2]
		if !(11 <= load15 && load15 <= 12) {
			t.Errorf("11 <= %#v && %#v <= 12", load15, load15)
		}
	}

	load15 := (&statsBookkeeper{history: []doneActions{
		{23, now.Add(-time.Minute * 45 / 2), now.Add(time.Minute * 15 / 2)},
	}}).queryLoad()[2]
	if !(11 <= load15 && load15 <= 12) {
		t.Errorf("11 <= %#v && %#v <= 12", load15, load15)
	}
}

func assertDeepEqual(t *testing.T, a, b interface{}) {
	t.Helper()

	if !reflect.DeepEqual(a, b) {
		t.Errorf("reflect.DeepEqual(%#v, %#v)", a, b)
	}
}
