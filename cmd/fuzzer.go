package cmd

import (
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/netrixframework/netrix/config"
	"github.com/netrixframework/netrix/strategies"
	"github.com/netrixframework/netrix/strategies/fuzzing"
	"github.com/netrixframework/redisraft-testing/util"
	"github.com/spf13/cobra"
)

func FuzzerCommand() *cobra.Command {
	var iterations int
	var netrixAddr string
	var tlcAddr string

	cmd := &cobra.Command{
		Use: "tlc-fuzzer",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunFuzzer(iterations, netrixAddr, tlcAddr)
		},
	}
	cmd.PersistentFlags().IntVarP(&iterations, "iterations", "i", 10, "Number of fuzzer iteration loops to run")
	cmd.PersistentFlags().StringVarP(&netrixAddr, "netrix-addr", "n", "localhost:7074", "Address to run the netrix server")
	cmd.PersistentFlags().StringVarP(&tlcAddr, "tlc-addr", "t", "localhost:2023", "Address to TLC Server to measure coverage")

	return cmd
}

func RunFuzzer(iterations int, netrixServerAddr, tlcAddr string) error {
	termCh := make(chan os.Signal, 1)
	signal.Notify(termCh, os.Interrupt, syscall.SIGTERM)

	fuzzingStrategy := fuzzing.NewFuzzStrategy(&fuzzing.FuzzStrategyConfig{
		Mutator:           fuzzing.NewDefaultMutator(),
		Guider:            fuzzing.NewTLCGuider(tlcAddr, util.TLCEventMapper()),
		Seed:              rand.NewSource(time.Now().UnixNano()),
		InitialPopulation: 4,
		MutationsPerTrace: 4,
	})

	driver := strategies.NewStrategyDriver(
		&config.Config{
			NumReplicas:   3,
			APIServerAddr: netrixServerAddr,
			LogConfig: config.LogConfig{
				Format: "json",
				Path:   "/local/snagendra/data/testing/raft/t/checker.log",
				Level:  "info",
			},
		},
		&util.RedisRaftMessageParser{},
		fuzzingStrategy,
		&strategies.StrategyConfig{
			Iterations:       iterations,
			IterationTimeout: 10 * time.Second,
		},
	)

	go func() {
		<-termCh
		driver.Stop()
	}()

	err := driver.Start()
	if err != strategies.ErrDriverQuit {
		return err
	}
	return nil
}
