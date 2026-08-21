// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package command

import (
	"testing"
)

func TestLongFlag(t *testing.T) {
	t.Parallel()
	const flagName = "myflag"
	for _, tt := range []struct {
		name   string
		expect string
		conf   *OptionsSetConfig
	}{
		{"normal", flagName, &OptionsSetConfig{
			FlagPrefix: "",
			Flags: map[string]FlagConfig{
				flagName: {Long: flagName},
			},
		}},
		{"normal", "test-myflag", &OptionsSetConfig{
			FlagPrefix: "test",
			Flags: map[string]FlagConfig{
				flagName: {Long: flagName},
			},
		}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.expect != tt.conf.LongFlag(flagName) {
				t.Fail()
			}
		})
	}
}
