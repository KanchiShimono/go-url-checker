package main

import (
	"flag"
	"fmt"
	"github.com/KanchiShimono/go-healthcheck/checkers"
	"github.com/KanchiShimono/go-url-checker/infrastructure/datastore"
	"github.com/KanchiShimono/go-url-checker/repository"
	"os"
	"time"
)

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		panic("Too many or few args. This programe accept only a input file name.")
	}

	iFileName := flag.Arg(0)
	iFile, err := os.Open(iFileName)
	if err != nil {
		fmt.Println(err.Error())
		panic("Can not open file")
	}
	defer iFile.Close()

	var fileRepo repository.FileRepositoryReader = datastore.NewCSVRepositoryReader(iFile, '\t')
	conditions, err := fileRepo.ReadAll()
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed parse file")
	}

	var results []repository.ResultHTTPCheck
	for _, con := range conditions {
		checker := &checkers.HTTPChecker{
			URL:        con.URL,
			StatusCode: con.StatusCode,
			Timeout:    con.Timeout,
		}

		err = checker.Check()

		result := repository.ResultHTTPCheck{
			TimeStamp:   time.Now().Format("2006-01-02 15:04:05"),
			URL:         con.URL,
			Result:      err,
			Description: con.Description,
		}

		results = append(results, result)

		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed: ", result.URL.String())
		}

	}

	// Time stamp format for output results file is yyyymmdd_hhmmss
	ts := time.Now().Format("20060102_150405")
	oFileName := "health_check_" + ts + ".tsv"
	oFile, err := os.OpenFile(oFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println(err.Error())
		panic("Can not open file")
	}
	defer oFile.Close()

	var writeRepo repository.FileRepositoryWriter = datastore.NewCSVRepositoryWriter(oFile, '\t')
	err = writeRepo.WriteAll(results)
	if err != nil {
		fmt.Println(err.Error())
	}

	var writeStdout repository.FileRepositoryWriter = datastore.NewCSVRepositoryWriter(os.Stdout, ':')
	err = writeStdout.WriteAll(results)
	if err != nil {
		fmt.Println(err.Error())
	}

}
