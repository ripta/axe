package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/klog"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	"github.com/ripta/axe/pkg/kubelogs"
	"github.com/ripta/axe/pkg/ui"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	if err := start(context.Background()); err != nil {
		fmt.Printf("Error: %+v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func start(ctx context.Context) error {
	logger := log.New(os.Stderr, "", log.Lmicroseconds)

	root := &cobra.Command{
		Use:           "axe",
		Short:         "Split and display logs in more manageable chunks",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	kcf := genericclioptions.NewConfigFlags(true)
	kcf.AddFlags(root.PersistentFlags())

	{
		logflags := flag.NewFlagSet("dummy", flag.ExitOnError)
		klog.InitFlags(logflags)
		root.PersistentFlags().AddGoFlagSet(logflags)
	}

	f := cmdutil.NewFactory(kcf)
	root.RunE = run(logger, f)

	return root.ExecuteContext(ctx)
}

func run(logger *log.Logger, f cmdutil.Factory) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		cs, err := f.KubernetesClientSet()
		if err != nil {
			return err
		}

		m := kubelogs.NewManager(logger, cs, 1*time.Second, 3*time.Minute)
		u := ui.New(logger, m)

		nss, _, err := f.ToRawKubeConfigLoader().Namespace()
		if err != nil {
			return err
		}

		for _, ns := range strings.Split(nss, ",") {
			m.Watch(strings.TrimSpace(ns))
		}

		return u.Run(ctx)
	}
}
