package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PullMergeConflictThenRebase = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pull with a merge strategy, get a conflict, then choose to rebase instead",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file", "content1")
		shell.Commit("one")
		shell.UpdateFileAndAdd("file", "content2")
		shell.Commit("two")
		shell.EmptyCommit("three")

		shell.CloneIntoRemote("origin")

		shell.SetBranchUpstream("master", "origin/master")

		shell.HardReset("HEAD^^")
		shell.UpdateFileAndAdd("file", "content4")
		shell.Commit("four")

		shell.SetConfig("pull.rebase", "false")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.Pull)

		// The merge pull conflicts. Choose "rebase instead" to abort the merge
		// and retry the pull with a rebase.
		t.ExpectPopup().Menu().
			Title(Equals("Conflicts!")).
			Select(Contains("Rebase instead")).
			Confirm()

		// The rebase hits the same conflict. Resolve it and verify the history
		// ends up linear rather than with a merge commit.
		t.Common().AcknowledgeConflicts()

		t.Views().Files().
			IsFocused().
			Lines(
				Contains("UU").Contains("file"),
			).
			PressEnter()

		t.Views().MergeConflicts().
			IsFocused().
			TopLines(
				Contains("<<<<<<< HEAD"),
				Contains("content2"),
				Contains("======="),
				Contains("content4"),
				Contains(">>>>>>>"),
			).
			SelectNextItem().
			PressPrimaryAction()

		t.Common().ContinueOnConflictsResolved("rebase")

		t.Views().Status().Content(Equals("↑1 repo → master"))

		t.Views().Commits().
			Focus().
			Lines(
				Contains("four").IsSelected(),
				Contains("three"),
				Contains("two"),
				Contains("one"),
			)
	},
})
