package cbtests

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/dahernan/goHystrix"
)

func testDependency() error {
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
	result := rand.Intn(2)
	var err error
	if result == 1 {
		err = fmt.Errorf("Ack! Got an error!")
	}
	return err
}

func printLongStats(msg string, c *goHystrix.Command) {
	fmt.Printf("\n%s\n\n", msg)
	fmt.Println("\tRSY: Count Buckets:\t", c.HealthCounts().HealthCountsBucket)
	fmt.Println("\tRSY: Total calls:\t\t", c.HealthCounts().Total)
	fmt.Println("\tRSY: Successful calls:\t\t", c.HealthCounts().Success)
	fmt.Println("\tRSY: Failed calls:\t\t", c.HealthCounts().Failures)
	fmt.Println("\tRSY: Error Percentage:\t\t", c.HealthCounts().ErrorPercentage)
	fmt.Println("")
	fmt.Println("\tRSY: Stats Count:\t\t", c.Metric().Stats().Count(), "(only successes are counted)")
	fmt.Println("\tRSY: Max Exec Time:\t\t", c.Metric().Stats().Max()/1000000, "ms")
	fmt.Println("\tRSY: Min Exec Time:\t\t", c.Metric().Stats().Min()/1000000, "ms")
	fmt.Println("\tRSY: Mean Exec Time:\t\t", c.Metric().Stats().Mean()/1000000, "ms")
	fmt.Println("\tRSY: Stats 50 Percentile:\t", c.Metric().Stats().Percentile(0.50)/1000000, "ms")
	fmt.Println("\tRSY: Stats 90 Percentile:\t", c.Metric().Stats().Percentile(0.90)/1000000, "ms")

	fmt.Printf("\n\n\n")

}

func printStats(msg string, total, success, failures, consecFails int, errPercent float64) {
	fmt.Printf("\n%s\n\n", msg)
	fmt.Println("\tTotal calls:\t\t", total)
	fmt.Println("\tSuccessful calls:\t", success)
	fmt.Println("\tFailed calls:\t\t", failures)
	fmt.Println("\tConsec failed calls:\t", consecFails)
	fmt.Printf("\tError Percentage:\t %f\n", errPercent*100)

	fmt.Printf("\n\n\n")

}
