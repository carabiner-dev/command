// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package command

import (
	"testing"
)

func TestLongFlag(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name   string
		expect string
		conf   *OptionsSetConfig
	}{
		{"normal", "myflag", &OptionsSetConfig{
			FlagPrefix: "",
			Flags: map[string]FlagConfig{
				"myflag": {Long: "myflag"},
			},
		}},
		{"normal", "test-myflag", &OptionsSetConfig{
			FlagPrefix: "test",
			Flags: map[string]FlagConfig{
				"myflag": {Long: "myflag"},
			},
		}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.expect != tt.conf.LongFlag("myflag") {
				t.Fail()
			}
		})
	}
}
