package main

import (
	"github.com/acidlemon/httpload-go/httpload"
	"flag"
	"fmt"
	"os"
	"io"
	"bufio"
)

func main() {
	var parallel *int = flag.Int("parallel", 10, "parallel")
	var seconds *int = flag.Int("seconds", 10, "seconds")
	var keepalive *bool = flag.Bool("keepalive", false, "keepalive")
	flag.Parse()

	urlfile := flag.Arg(0)

	if len(urlfile) == 0 {
		fmt.Println("give me url_file...")
		os.Exit(1)
	}

	fmt.Println("option: parallel=", *parallel, ", seconds=", *seconds)

	conf := new(httpload.Config)
	conf.Parallel = *parallel
	conf.Seconds = *seconds
	conf.Urls = openUrlFile(urlfile)
	conf.KeepAlive = *keepalive

	httpload.Start(*conf)


}


func openUrlFile(path string) []string {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("open file error('-';)")
		os.Exit(2)
	}

	urllist := make([]string, 0, 100)

	r := bufio.NewReaderSize(f, 4*1024)
	for {
		line, isPrefix, err := r.ReadLine()

		if err == io.EOF {
			break
		}

		if len(line) < 3 {
			continue
		}

		if len(line) > 0 && err != nil {
			fmt.Fprintf(os.Stderr, "cat: error %v\n", err)
			os.Exit(1)
		}

		if isPrefix {
			fmt.Fprintf(os.Stderr, "ReadLine returned prefix\n")
			os.Exit(1)
		}


		urllist = append(urllist, string(line))
	}

	return urllist
}

