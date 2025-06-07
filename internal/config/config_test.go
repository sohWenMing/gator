package config

import (
	"reflect"
	"testing"
)

func TestReadConfig(t *testing.T) {
	type test struct {
		name          string
		path          string
		expected      Config
		isErrExpected bool
	}

	tests := []test{
		{
			"basic test should pass",
			"./testconfig.json",
			Config{
				"postgres://example",
				"",
			},
			false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Read(test.path)

			switch test.isErrExpected {
			case true:
				if err == nil {
					t.Errorf("expected error, didn't get one")
				}
			default:
				if err != nil {
					t.Errorf("didn't expect error, got %v", err)
				}
			}

			if !reflect.DeepEqual(*got, test.expected) {
				t.Errorf("\ngot: %s\nwant: %s", got.String(),
					test.expected.String())
			}
		})
	}

}
