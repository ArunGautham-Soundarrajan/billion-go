build:
  go build -o billion-go .

run: build
  time ./billion-go > out.txt 

profile:
    docker run --rm -v "{{justfile_directory()}}":/app -w /app golang:latest \
      sh -c 'go build -o billion-linux . && ./billion-linux > /dev/null'
    go tool pprof -top cpu.prof
