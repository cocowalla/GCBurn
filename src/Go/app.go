package main

import (
	bt "./burn_test"
	. "./common"
	st "./speed_test"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"
)

func main() {
	ramSizeGb := GetRamSize()
	var durationSecFlag = flag.Int64("d", 10, "Test pass duration, seconds")
	var staticSetSizeGbFlag = flag.String("m", strconv.Itoa(ramSizeGb), "Static set size, GB")
	var threadCountFlag = flag.String("t", "", "Number of threads to use")
	var maxSizeFlag = flag.String("s", "", "Max. object size")
	var _ = flag.String("l", "", "Latency mode (ignored for Go)")
	flag.Parse()
	bt.DefaultDuration = time.Duration(*durationSecFlag * int64(time.Second))
	bt.DefaultMaxSize = ParseRelative(*maxSizeFlag, bt.DefaultMaxSize, true)
	staticSetSizeGb := int(ParseRelative(*staticSetSizeGbFlag, float64(ramSizeGb), true))
	ThreadCount = int(ParseRelative(*threadCountFlag, float64(ThreadCount), true))

	args := fmt.Sprintf("%+v", os.Args[1:])
	fmt.Printf("Launch parameters: %v\n", args[1:len(args)-1])
	fmt.Printf("Software:\n")
	fmt.Printf("  Runtime:         Go\n")
	fmt.Printf("    Version:       %v\n", runtime.Version())
	fmt.Printf("  OS:              %v %v\n", GetOSVersion(), runtime.GOARCH)
	fmt.Printf("Hardware:\n")
	fmt.Printf("  CPU:             %v\n", GetCpuModelName())
	coreCountAddon := ""
	if runtime.NumCPU() != ThreadCount {
		coreCountAddon = fmt.Sprintf(" (using %v during the test)", ThreadCount)
	}
	fmt.Printf("  CPU core count:  %v%v\n", runtime.NumCPU(), coreCountAddon)
	fmt.Printf("  RAM size:        %v GB\n", ramSizeGb)
	runtime.NumCPU()

	fmt.Printf("Warming up...\n")
	speedTester := st.NewWarmupSpeedTester()
	speedTester.Run()
	speedTester = nil
	burnTester := bt.NewWarmupBurnTester(1 * int64(GB))
	burnTester.Run()
	burnTester = nil
	fmt.Printf("  Done.\n")
	fmt.Println()

	if staticSetSizeGb != ramSizeGb {
		title := fmt.Sprintf("--- Static set = %v GB (%.2f %% RAM) ---",
			staticSetSizeGb, float64(staticSetSizeGb)*100/float64(ramSizeGb))
		runBurnTest(title, int64(staticSetSizeGb)*GB)
	} else {
		runSpeedTest("--- Raw allocation (w/o holding what's allocated) ---")
		runBurnTest("--- Stateless server (no static set) ---", 0)
		runBurnTest("--- Worker / typical server (static set = 20% RAM) ---", int64(ramSizeGb)*GB/5)
		runBurnTest("--- Caching / compute server (static set = 50% RAM) ---", int64(ramSizeGb)*GB/2)
	}
}

func runSpeedTest(title string) {
	fmt.Println(title)
	fmt.Println()
	speedTester := st.NewSpeedTester()
	speedTester.Run()
}

func runBurnTest(title string, ramSize int64) {
	fmt.Println(title)
	fmt.Println()
	burnTester := bt.NewBurnTester(ramSize)
	burnTester.Run()
}
