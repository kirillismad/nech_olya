package main

import (
	"context"
	"notesv1/cmd"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	cobra.CheckErr(cmd.RootCmd.ExecuteContext(ctx))
}
