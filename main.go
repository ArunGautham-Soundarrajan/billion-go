package main

import (
	"bytes"
	"fmt"
	"os"
	"runtime/pprof"
	"sort"
)

type weatherData struct {
	min   int
	max   int
	sum   int
	count int
}

func newWeatherData(temperature int) *weatherData {
	return &weatherData{
		min:   temperature,
		max:   temperature,
		sum:   temperature,
		count: 1,
	}
}

func main() {
	f, _ := os.Create("cpu.prof")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	// Actual implementation start
	file, err := os.Open("measurements.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	output := make(map[string]*weatherData)

	// chunk of constant size to read from the file
	buf := make([]byte, 16<<20)
	leftover := 0

	for {
		// file.Read() reads the file and filles the buf we passed and it
		// keeps track where it last read. We pass in leftover in case our buffer
		// has partial lines which we want to keep for the next read at the top and
		// only fill the rest and not rewrite the entire buf and erase the leftover
		n, err := file.Read(buf[leftover:])
		total := leftover + n
		if total == 0 {
			break
		}

		// this is where we find the last full line and we copy the rest to the buf
		// which copy the trailing partial line to the front of the buf so the next
		// read continues it
		last := bytes.LastIndexByte(buf[:total], '\n')
		processChunk(buf[:last+1], output)

		leftover = copy(buf, buf[last+1:total])

		if err != nil {
			break
		}
	}

	sortMap(output)
}

func processChunk(chunk []byte, m map[string]*weatherData) {
	for len(chunk) > 0 {
		// We have a chunk of bytes, we need to process line by line
		// IndexByte, returns the index of the first occurance of the
		// given byte, in our case '\n', so we just take the slice before
		// that index, which gives us a single line and we do this until
		// the chunk becomes zero
		newLine := bytes.IndexByte(chunk, '\n')
		line := chunk[:newLine]
		chunk = chunk[newLine+1:]

		// We use the same index byte to split the byte into station and temperature
		i := bytes.IndexByte(line, ';')
		stationName := line[:i]
		temperature := line[i+1:]
		temp := customTemperatureParser(temperature)
		insertToMap(stationName, temp, m)
	}
}

func insertToMap(station []byte, temperature int, m map[string]*weatherData) {
	val, ok := m[string(station)]
	if ok {
		val.min = min(val.min, temperature)
		val.max = max(val.max, temperature)
		val.sum += temperature
		val.count++
	} else {
		newVal := newWeatherData(temperature)
		m[string(station)] = newVal
	}
}

func sortMap(m map[string]*weatherData) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		v := m[k]

		mean := float64(v.sum) / float64(v.count) / 10
		fmt.Printf("%s=%.1f/%.1f/%.1f\n", k, float64(v.min)/10, mean, float64(v.max)/10)
	}
}

func customTemperatureParser(b []byte) int {
	negative := false
	if b[0] == '-' {
		negative = true
		b = b[1:]
	}
	n := 0
	for _, c := range b {
		if c == '.' {
			continue
		}

		n = n*10 + int(c-'0')
	}

	if negative {
		return -n
	}
	return n
}
