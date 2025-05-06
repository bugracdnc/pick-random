package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"path/filepath"
	"slices"

	"github.com/spf13/cobra"
)

func resolvePath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}
	return abs
}

var (
	count         int
	unique        bool
	verbose       bool
	humanReadable bool
)

var rootCmd = &cobra.Command{
	Use:   "pick-random [path]",
	Short: "Pick one or more random files from a directory",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := resolvePath(args[0])
		entries := getEntries(path)

		if verbose {
			fmt.Printf("Picking from %s/", resolvePath(path))
			if unique {
				fmt.Println("Picking unique items")
			}
			if count != 1 {
				fmt.Printf("Picking %d items\n", count)
			}
			if humanReadable {
				fmt.Println("Printing human-readable results")
			}
			fmt.Println()
		}

		if count >= len(entries) && unique {
			if verbose {
				fmt.Printf("Count (%d) is equal to or greater than available items (%d)\nDefaulting to available items (%d)\n\n", count, len(entries), len(entries))
			}
			for i, e := range entries {
				if humanReadable {
					fmt.Printf("Picked \"%s\" from %d options\r\n", e, len(entries)-i)
				} else {
					fmt.Println(e)
				}
			}
			return
		}

		for range count {
			nth := rand.IntN(len(entries))
			picked := entries[nth]
			if humanReadable {
				fmt.Printf("Picked \"%s\" from %d options\r\n", picked, len(entries))
			} else {
				fmt.Println(picked)
			}
			if unique {
				entries = slices.Delete(entries, nth, nth+1)
			}
		}
	},
}

func init() {
	rootCmd.Flags().IntVarP(&count, "count", "c", 1, "Number of items to pick")
	rootCmd.Flags().BoolVarP(&unique, "unique", "u", false, "Ensure selections are unique")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.Flags().BoolVarP(&humanReadable, "human-readable", "r", false, "Make output human-readable")
}

func getEntries(path string) []string {
	stat, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}
	if !stat.IsDir() {
		log.Fatalf("Path is not a directory: %s", path)
	}

	dirEntries, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	var entries []string
	for _, e := range dirEntries {
		name := e.Name()

		if name[0] == '.' {
			continue
		}

		if e.IsDir() {
			name += "/"
		}
		entries = append(entries, name)
	}

	return entries
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
