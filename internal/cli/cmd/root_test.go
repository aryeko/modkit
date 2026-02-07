package cmd

import "testing"

func TestExecuteHelp(t *testing.T) {
	origArgs := rootCmd.Args
	origOut := rootCmd.OutOrStdout()
	origErr := rootCmd.ErrOrStderr()
	t.Cleanup(func() {
		rootCmd.SetArgs(nil)
		rootCmd.SetOut(origOut)
		rootCmd.SetErr(origErr)
		rootCmd.Args = origArgs
	})

	rootCmd.SetArgs([]string{"--help"})
	if err := Execute(); err != nil {
		t.Fatalf("Execute() failed: %v", err)
	}
}
