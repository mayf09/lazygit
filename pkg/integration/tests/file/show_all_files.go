package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ShowAllFiles = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Toggling between showing only changed files and showing all files",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file-a", "a")
		shell.CreateFileAndAdd("file-b", "b")
		shell.CreateFileAndAdd("dir/file-c", "c")
		shell.Commit("initial commit")

		shell.UpdateFile("file-a", "a changed")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Equals(" M file-a"),
			)

		t.Views().Files().
			Press(keys.Files.ToggleShowAllFiles).
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  ▶ dir"),
				Equals("   M file-a"),
				Equals("     file-b"),
			)

		t.Views().Files().
			Press(keys.Files.ToggleShowAllFiles).
			Lines(
				Equals(" M file-a"),
			)
	},
})
