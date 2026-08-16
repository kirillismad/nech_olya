package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"notesv1/cmd"

	"github.com/spf13/cobra"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	cobra.CheckErr(cmd.RootCmd.ExecuteContext(ctx))
}
