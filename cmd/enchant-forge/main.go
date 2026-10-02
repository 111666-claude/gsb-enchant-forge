// Command enchant-forge 跑强化样例。
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"example.com/forge"
)

func newForge() *forge.Enhancer {
	return forge.New([]float64{1.0, 0.95, 0.4, 0.6, 0.1, 0.1, 0.1}, 8, 3, 3, 100)
}

// Run 执行一次命令行调用，返回退出码。
func Run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("enchant-forge", flag.ContinueOnError)
	flags.SetOutput(stderr)
	sample := flags.String("sample", "grade", "grade / degrade / protect / pity / dup / scan")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	switch *sample {
	case "grade":
		enhancer := newForge()
		enhancer.Register("sword", 1, 1000)
		level, success, _ := enhancer.Enhance("r1", "sword", 0.5, false)
		fmt.Fprintf(stdout, "level=%d success=%v\n", level, success)
	case "degrade":
		enhancer := newForge()
		enhancer.Register("sword", 5, 1000)
		level, _, _ := enhancer.Enhance("r1", "sword", 0.99, false)
		fmt.Fprintf(stdout, "level=%d\n", level)
	case "protect":
		enhancer := newForge()
		enhancer.Register("sword", 3, 1000)
		level, _, _ := enhancer.Enhance("r1", "sword", 0.99, true)
		fmt.Fprintf(stdout, "level=%d\n", level)
	case "pity":
		enhancer := forge.New([]float64{0, 0, 0, 0}, 8, 1, 3, 100)
		enhancer.Register("sword", 1, 1000)
		level := 0
		for index := 0; index < 4; index++ {
			level, _, _ = enhancer.Enhance(fmt.Sprintf("r-%d", index), "sword", 0.5, false)
		}
		fmt.Fprintf(stdout, "level=%d\n", level)
	case "dup":
		enhancer := newForge()
		enhancer.Register("sword", 0, 1000)
		enhancer.Enhance("same", "sword", 0.0, false)
		level, _, _ := enhancer.Enhance("same", "sword", 0.0, false)
		fmt.Fprintf(stdout, "level=%d gold=%d\n", level, enhancer.Gold("sword"))
	case "scan":
		enhancer := newForge()
		for index := 0; index < 3000; index++ {
			enhancer.Register(fmt.Sprintf("i-%d", index), 1, 1000)
		}
		for index := 0; index < 3000; index++ {
			enhancer.Enhance(fmt.Sprintf("r-%d", index), "i-2999", 0.5, false)
		}
		fmt.Fprintf(stdout, "scanned=%d\n", enhancer.Scanned())
	default:
		fmt.Fprintln(stderr, "需要 --sample grade|degrade|protect|pity|dup|scan")
		return 2
	}
	return 0
}

func main() {
	os.Exit(Run(os.Args[1:], os.Stdout, os.Stderr))
}
