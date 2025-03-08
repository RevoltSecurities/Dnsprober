package handler

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/RevoltSecurities/Dnsprober/dnsprober/modules/banner"
	"github.com/RevoltSecurities/Dnsprober/dnsprober/modules/cli"
	"github.com/RevoltSecurities/Dnsprober/dnsprober/modules/gorunner"
	"github.com/RevoltSecurities/Dnsprober/dnsprober/modules/logger"
	"github.com/RevoltSecurities/Dnsprober/dnsprober/modules/utils"
	"github.com/RevoltSecurities/Dnsprober/dnsprober/modules/version"
)

var (
	Loggers = logger.New(true)
	V       = "v1.0.0"
)

func Version() {
	git, err := version.GitVersion()
	if err != nil {
		Loggers.StdLogger("Unable to get the latest version of the dnsprober", "warn")
	} else {
		if git == V {
			Loggers.VLogger("latest", "dnsprober", V)
		} else {
			Loggers.VLogger("outdated", "dnsprober", V)
		}
	}
}

func Init() {
	args := cli.CLI()
	if args.Exceptions != nil {
		Loggers.Logger("exception occured initalizing CLI arguments due to: %s", args.Exceptions.Error())
		os.Exit(1)
	}

	if !args.Silent {
		fmt.Fprintf(os.Stderr, "%s", banner.BannerGenerator("dnsprober"))
		fmt.Fprintf(os.Stderr, "\n              - %s", Loggers.Bolder("Revoltsecurities"))
		fmt.Println("\n")
		Version()
	}

	if args.Output != "" {
		if _, err := utils.CanWrite(args.Output); err != nil {
			Loggers.Logger(fmt.Sprintf("unable to create the output %s file due to: %s", args.Output, err.Error()), "warn")
			os.Exit(1)
		}
	}
	gorunners, err := gorunner.New(&args)
	if err != nil {
		Loggers.Logger(fmt.Sprintf("unable to start gorunner due to: %s", err.Error()), "warn")
		os.Exit(1)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for range c {
			Loggers.StdLogger("CTRL+C Pressed!", "warn")
			os.Exit(1)
		}
	}()
	if err := gorunners.Sprint(); err != nil {
		Loggers.Logger(fmt.Sprintf("unable to start sprinter module due to: %s", err.Error()), "warn")
		os.Exit(1)
	}
	gorunners.Down()
}
