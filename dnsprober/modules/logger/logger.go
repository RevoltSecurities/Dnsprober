package logger

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
)

const (
	tcl = "\r\x1b[2K"
)

type Logger struct {
	NoColor bool
	aurora  aurora.Aurora
}

func New(noColor bool) *Logger {
	return &Logger{
		NoColor: noColor,
		aurora:  aurora.NewAurora(noColor),
	}
}

func (l *Logger) Logger(message, level string) {
	var leveler aurora.Value
	switch level {
	case "info":
		leveler = l.aurora.Bold(l.aurora.Blue("INFO"))
	case "warn":
		leveler = l.aurora.Bold(l.aurora.Yellow("WRN"))
	case "verbose":
		leveler = l.aurora.Bold(l.aurora.Green("VERBOSE"))
	case "error":
		leveler = l.aurora.Bold(l.aurora.Red("ERROR"))
	default:
		leveler = l.aurora.Bold(l.aurora.Blue(level))
	}
	formatted := fmt.Sprintf("%s%s%s %s\n",
		l.aurora.Bold(l.aurora.White("[")),
		leveler,
		l.aurora.Bold(l.aurora.White("]: ")),
		l.aurora.Bold(l.aurora.White(message)),
	)
	fmt.Fprint(os.Stderr, formatted)
}

func (l *Logger) StdinLogger(message string) {
	fmt.Println(fmt.Sprintf("%s%s", tcl, message))
}

func (l *Logger) Bolder(message string) string {
	return fmt.Sprintf("%s", l.aurora.Bold(l.aurora.White(message)))
}

func (l *Logger) Loader(message, level string) string {
	var leveler aurora.Value
	switch level {
	case "info":
		leveler = l.aurora.Bold(l.aurora.Blue("INFO"))
	case "warn":
		leveler = l.aurora.Bold(l.aurora.Red("WRN"))
	case "verbose":
		leveler = l.aurora.Bold(l.aurora.Green("VERBOSE"))
	default:
		leveler = l.aurora.Bold(l.aurora.Blue(level))
	}
	formatted := fmt.Sprintf("%s%s%s %s\n",
		l.aurora.Bold(l.aurora.White("[")),
		leveler,
		l.aurora.Bold(l.aurora.White("]: ")),
		l.aurora.Bold(l.aurora.White(message)),
	)
	return formatted
}

func (l *Logger) StdLogger(message, level string) {
	fmt.Fprint(os.Stderr, tcl)
	switch level {
	case "info":
		fmt.Fprintf(os.Stderr, "%s%s%s%s%s\n",
			tcl,
			l.aurora.Bold(l.aurora.White("[")),
			l.aurora.Bold(l.aurora.Blue("INFO")),
			l.aurora.Bold(l.aurora.White("]: ")),
			l.aurora.Bold(l.aurora.White(message)),
		)
	case "warn":
		fmt.Fprintf(os.Stderr, "%s%s%s%s%s\n",
			tcl,
			l.aurora.Bold(l.aurora.White("[")),
			l.aurora.Bold(l.aurora.Red("WRN")),
			l.aurora.Bold(l.aurora.White("]: ")),
			l.aurora.Bold(l.aurora.White(message)),
		)
	case "error":
		fmt.Fprintf(os.Stderr, "%s%s%s%s%s\n",
			tcl,
			l.aurora.Bold(l.aurora.White("[")),
			l.aurora.Bold(l.aurora.Red("ERR")),
			l.aurora.Bold(l.aurora.White("]: ")),
			l.aurora.Bold(l.aurora.White(message)),
		)
	case "vuln":
		fmt.Fprintf(os.Stderr, "%s%s%s%s%s\n",
			tcl,
			l.aurora.Bold(l.aurora.White("[")),
			l.aurora.Bold(l.aurora.Green("VULN")),
			l.aurora.Bold(l.aurora.White("]: ")),
			l.aurora.Bold(l.aurora.White(message)),
		)
	default:
		fmt.Fprintf(os.Stderr, "%s%s%s%s%s\n",
			tcl,
			l.aurora.Bold(l.aurora.White("[")),
			l.aurora.Bold(l.aurora.Blue(level)),
			l.aurora.Bold(l.aurora.White("]: ")),
			l.aurora.Bold(l.aurora.White(message)),
		)
	}
}

func (l *Logger) VLogger(level, toolname, version string) {
	if level == "latest" {
		fmt.Fprintf(os.Stderr, "%s%s%s%s%s%s\n",
			l.aurora.Bold(l.aurora.White("[")),
			l.aurora.Bold(l.aurora.Blue("Version")),
			l.aurora.Bold(l.aurora.White("]: ")),
			l.aurora.Bold(l.aurora.White(fmt.Sprintf("%s current version %s (", toolname, version))),
			l.aurora.Bold(l.aurora.Green("latest")),
			l.aurora.Bold(l.aurora.White(")")),
		)
	} else {
		fmt.Fprintf(os.Stderr, "%s%s%s%s%s%s\n",
			l.aurora.Bold(l.aurora.White("[")),
			l.aurora.Bold(l.aurora.Blue("Version")),
			l.aurora.Bold(l.aurora.White("]: ")),
			l.aurora.Bold(l.aurora.White(fmt.Sprintf("%s current version %s (", toolname, version))),
			l.aurora.Bold(l.aurora.Red("outdated")),
			l.aurora.Bold(l.aurora.White(")")),
		)
	}
}

func (l *Logger) Colorizer(banner, color string) aurora.Value {
	switch color {
	case "blue":
		return l.aurora.Bold(l.aurora.Blue(banner))
	case "white":
		return l.aurora.Bold(l.aurora.White(banner))
	case "magenta":
		return l.aurora.Bold(l.aurora.Magenta(banner))
	case "green":
		return l.aurora.Bold(l.aurora.Green(banner))
	case "cyan":
		return l.aurora.Bold(l.aurora.Cyan(banner))
	case "yellow":
		return l.aurora.Bold(l.aurora.Yellow(banner))
	default:
		return l.aurora.Bold(l.aurora.Blue(banner))
	}
}
