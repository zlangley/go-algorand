// Copyright (C) 2019-2021 Algorand, Inc.
// This file is part of go-algorand
//
// go-algorand is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// go-algorand is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with go-algorand.  If not, see <https://www.gnu.org/licenses/>.


package main

import "github.com/spf13/cobra"

func init() {
	speculationCmd.AddCommand(newSpeculationCmd)
}

var speculationCmd = &cobra.Command{
	Use: "speculation",
	Short:  "Provides tools to control speculative execution",
	Long:  "Collection of commands to support the management of Layer 2 speculative execution",
	Args: validateNoPosArgsFn,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.HelpFunc()(cmd, args)
	},
}

var newSpeculationCmd = &cobra.Command{
	Use: "new",
	Args: validateNoPosArgsFn,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.HelpFunc()(cmd, args)
	},
}