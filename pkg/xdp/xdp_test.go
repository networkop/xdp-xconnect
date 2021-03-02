package xdp

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

func TestDiff(t *testing.T) {
	tests := []struct {
		running   map[string]string
		candidate map[string]string
		new       []string
		changed   []string
		removed   []string
	}{
		{
			running: map[string]string{
				"intf-1": "intf-2",
				"intf-2": "intf-1",
			},
			candidate: map[string]string{
				"intf-1": "intf-3",
				"intf-3": "intf-1",
			},
			new:     []string{"intf-3"},
			changed: []string{"intf-1"},
			removed: []string{"intf-2"},
		},
		{
			running: map[string]string{
				"intf-1": "intf-2",
				"intf-2": "intf-1",
			},
			candidate: map[string]string{
				"intf-1": "intf-2",
				"intf-2": "intf-1",
			},
			new:     nil,
			changed: nil,
			removed: nil,
		},
		{
			running: map[string]string{
				"intf-1": "intf-2",
				"intf-2": "intf-1",
			},
			candidate: map[string]string{
				"intf-1": "intf-2",
				"intf-2": "intf-1",
				"intf-3": "intf-4",
				"intf-4": "intf-3",
			},
			new:     []string{"intf-3", "intf-4"},
			changed: nil,
			removed: nil,
		},
		{
			running: map[string]string{
				"intf-1": "intf-2",
				"intf-2": "intf-1",
				"intf-3": "intf-4",
				"intf-4": "intf-3",
			},
			candidate: map[string]string{
				"intf-1": "intf-4",
				"intf-4": "intf-1",
				"intf-3": "intf-2",
				"intf-2": "intf-3",
			},
			new:     nil,
			changed: []string{"intf-1", "intf-2", "intf-3", "intf-4"},
			removed: nil,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("TestDiff_%d", i), func(t *testing.T) {

			new, changed, removed := confDiff(tt.running, tt.candidate)

			sort.Strings(new)
			sort.Strings(changed)
			sort.Strings(removed)

			if !reflect.DeepEqual(new, tt.new) {
				t.Errorf("#%d NEW wanted %v, got: %v", i, tt.new, new)
			}
			if !reflect.DeepEqual(changed, tt.changed) {
				t.Errorf("#%d CHANGED wanted %v, got: %v", i, tt.changed, changed)
			}
			if !reflect.DeepEqual(removed, tt.removed) {
				t.Errorf("#%d REMOVED wanted %v, got: %v", i, tt.removed, removed)
			}
		})
	}
}

func TestSymm(t *testing.T) {
	tests := []struct {
		input  map[string]string
		output map[string]string
	}{
		{
			input: map[string]string{
				"intf-1": "intf-2",
			},
			output: map[string]string{
				"intf-1": "intf-2",
				"intf-2": "intf-1",
			},
		},
		{
			input: map[string]string{
				"intf-2": "intf-2",
			},
			output: map[string]string{
				"intf-2": "intf-2",
			},
		},
		{
			input: map[string]string{
				"intf-1": "intf-2",
				"intf-2": "intf-3",
			},
			output: map[string]string{
				"intf-1": "intf-2",
				"intf-2": "intf-1",
			},
		},
		{
			input: map[string]string{
				"intf-3": "intf-2",
				"intf-1": "intf-2",
			},
			output: map[string]string{
				"intf-1": "intf-2",
				"intf-2": "intf-1",
			},
		},
		{
			input:  map[string]string{},
			output: map[string]string{},
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("TestSymm_%d", i), func(t *testing.T) {

			output := makeSymm(tt.input)

			if !reflect.DeepEqual(tt.output, output) {
				t.Errorf("#%d wanted %v, got: %v", i, tt.output, output)
			}
		})
	}
}
