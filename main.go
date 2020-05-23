package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/igorwwwwwwwwwwwwwwwwwwww/log2hist/hist"
)

// TODO: support non group-by mode
// TODO: scaling factor (e.g. bytes => GiB)

func main() {
	var err error

	m := make(map[string]*hist.Histogram)

	r := os.Stdin
	if len(os.Args) > 1 {
		file := os.Args[1]
		r, err = os.Open(file)
		if err != nil {
			log.Fatalf("could not open file %v: %v", file, err)
		}
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			continue
		}

		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			// TODO: support any whitespace here
			// TODO: custom separator to support keys containing spaces (or just
			//       use the last space of the line to split)
			log.Fatalf("warning: ignoring bad value, expected k <space> value, got %v", parts)
		}

		key := parts[0]
		rawval := parts[1]
		val, err := strconv.ParseUint(rawval, 10, 64)
		if err != nil {
			log.Fatalf("warning: ignoring bad value, expected int, got %v", rawval)
		}

		h, ok := m[key]
		if !ok {
			h = hist.New()
			m[key] = h
		}

		err = h.Record(val)
		if err != nil {
			log.Printf("warning: ignoring bad value %v, got error: %v", val, err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// TODO: top-k mode for high-cardinality keys

	// sort keys by total count
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		ki := keys[i]
		kj := keys[j]
		mi := m[ki]
		mj := m[kj]
		if mi.Count == mj.Count {
			return ki > kj
		}
		return mi.Count > mj.Count
	})

	for _, k := range keys {
		fmt.Println(k)
		fmt.Println(m[k])
	}
}
