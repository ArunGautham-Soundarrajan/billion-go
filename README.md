# billion-go

My Go take on the [One Billion Row Challenge](https://github.com/gunnarmorling/1brc): compute min / mean / max per station over 1B temperature readings.

## Result

Full 1B rows, run natively on an M2 Max (12 cores): **99.2s → 4.6s**, about **21x faster**.

## How I got here

I started with a deliberately naive version (`bufio.Scanner`, `strings.Split`, `ParseFloat`, one map). Each step after that came from a profile. Timings below are full 1B rows **inside Docker (Linux)**, from `pprof`, so they run slower than native. The same v6 code takes 8.5s here and 4.6s natively on the Mac.

| Step | Change | Time |
|---|---|---|
| v2b | `Scanner.Bytes()` + `bytes.IndexByte`, no per-line allocation | 70.8s |
| v3 | Hand-written integer parser (tenths) | 61.6s |
| v4 | Own chunked read loop, cut at the last newline | 46.2s |
| v5 | Reader goroutine feeding one worker over a channel | 37.9s |
| v6 | Pool of 12 workers, one map each, merged at the end | 8.5s |

On the Mac, 12 workers (4.6s) beat 8 (5.8s) and 16 (5.0s).

## What I learned

- The macOS profiler is unreliable, so I profile in Linux and trust `time` for wall-clock.
- Removing per-line allocation was the first big win.
- A custom integer parser beat `ParseFloat`. Avoiding the string copy did not.
- Run-to-run noise is real (the same code gave 61.6s and 54.4s).

## Caveats

Mac numbers are wall-clock, Docker numbers are `pprof` durations, so compare within one setup only. Mac v3–v5 weren't timed. My generated data is uniform random, so it doesn't match the official expected output.

## Next

Parse temperatures without a loop, a custom hash table instead of the built-in map, per-worker `ReadAt`, and reusing chunk buffers.
