// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package command

import "github.com/spf13/cobra"

// OptionsSet is an interface that defines the functions options set need
// to implement to make them reusable.
type OptionsSet interface {
	AddFlags(*cobra.Command)
	Validate() error
}
