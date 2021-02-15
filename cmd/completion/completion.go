/*
Package completion manages all completion related operations.

Copyright 2021 Michael Bungenstock

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License. You may obtain a copy of the
License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied. See the License for the
specific language governing permissions and limitations under the License.
*/
package completion

import (
	"os"

	"github.com/spf13/cobra"
)

// NewCommand creates an initialized command object
func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
		Long: `To load completions:
	
	Bash:
	
	$ source <(aepctl completion bash)
	
	# To load completions for each session, execute once:
	Linux:
	  $ aepctl completion bash > /etc/bash_completion.d/aepctl
	MacOS:
	  $ aepctl completion bash > /usr/local/etc/bash_completion.d/aepctl
	
	Zsh:
	
	# If shell completion is not already enabled in your environment you will need
	# to enable it.  You can execute the following once:
	
	$ echo "autoload -U compinit; compinit" >> ~/.zshrc
	
	# To load completions for each session, execute once:
	$ aepctl completion zsh > "${fpath[1]}/_aepctl"
	
	# You will need to start a new shell for this setup to take effect.
	
	Fish:
	
	$ aepctl completion fish | source
	
	# To load completions for each session, execute once:
	$ aepctl completion fish > ~/.config/fish/completions/aepctl.fish

	PowerShell:

	PS> aepctl completion powershell | Out-String | Invoke-Expression

	# To load completions for every new session, run:
	PS> aepctl completion powershell > aepctl.ps1
	# and source this file from your PowerShell profile.
	`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				_ = cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				_ = cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				_ = cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				_ = cmd.Root().GenPowerShellCompletion(os.Stdout)
			}
		},
	}
	return cmd
}
