package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
)

const (
	markName  = "GOLANG_CLI_REMINDER"
	markValue = "1"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s <HH:MM> <text message>\n", os.Args[0])
		os.Exit(1)
	}

	fmt.Println(os.Args)

	now := time.Now()

	w := when.New(nil)
	w.Add(en.All...)
	w.Add(common.All...)

	t, err := w.Parse(os.Args[1], now)

	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	if t == nil {
		fmt.Println("Unable to parse time")
		os.Exit(2)
	}

	if now.After(t.Time) {
		fmt.Println("set a future time")
		os.Exit(3)
	}

	diff := t.Time.Sub(now)

	if os.Getenv(markName) == markValue {
		time.Sleep(diff)
		fmt.Println("we are now printing...")
		err = beeep.Alert("Reminder", strings.Join(os.Args[2:], " "), "")
		if err != nil {
			fmt.Println(err)
			os.Exit(4)
		}
	} else {
		fmt.Println("we are in else")
		cmd := exec.Command(os.Args[0], os.Args[1:]...)
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", markName, markValue))
		if err = cmd.Start(); err != nil {
			fmt.Println(err)
			os.Exit(5)
		}
		fmt.Println("Reminder will be displayed after ", diff.Round(time.Second))
		os.Exit(0)
	}
}