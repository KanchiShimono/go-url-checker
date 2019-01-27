package datastore

import (
	"encoding/csv"
	"errors"
	"github.com/KanchiShimono/go-url-checker/repository"
	"io"
	"net/url"
	"reflect"
	"strconv"
	"time"
)

// CSVRepositoryReader is implementation of interface FileRepositoryReader
type CSVRepositoryReader struct {
	Reader    io.Reader
	separator rune
}

// NewCSVRepositoryReader work as constructor of CSVRepositoryReader which read from csv file
func NewCSVRepositoryReader(r io.Reader, s rune) *CSVRepositoryReader {
	return &CSVRepositoryReader{
		Reader:    r,
		separator: s,
	}
}

// ReadAll function reads conditon of http health check
func (cr *CSVRepositoryReader) ReadAll() (hccons []repository.ConditionHTTPCheck, err error) {
	r := csv.NewReader(cr.Reader)
	r.Comma = cr.separator
	cons, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	// Skip header
	for _, con := range cons[1:] {
		c, err := parseCondition(con)
		if err != nil {
			return hccons, err
		}

		hccons = append(hccons, c)
	}

	return hccons, nil
}

// parseCondition parse strings from CSV to HTTP Check condition
func parseCondition(constrs []string) (con repository.ConditionHTTPCheck, err error) {
	if len(constrs) != 4 {
		return con, errors.New("Illegal file format. Input file must be consisted by URL, Expected Code, Timeout and Description")
	}

	u, err := url.Parse(constrs[0])
	if err != nil {
		return con, err
	}

	code, err := strconv.Atoi(constrs[1])
	if err != nil {
		return con, err
	}

	timeout, err := strconv.Atoi(constrs[2])
	if err != nil {
		return con, err
	}

	desc := constrs[3]

	con = repository.ConditionHTTPCheck{
		URL:         u,
		StatusCode:  code,
		Timeout:     time.Duration(timeout) * time.Second,
		Description: desc,
	}

	return con, nil
}

// CSVRepositoryWriter is implementation of interface FileRepositoryWriter
type CSVRepositoryWriter struct {
	Writer    io.Writer
	separator rune
}

// NewCSVRepositoryWriter work as constructor of CSVRepositoryWriter which write to csv file
func NewCSVRepositoryWriter(w io.Writer, s rune) *CSVRepositoryWriter {
	return &CSVRepositoryWriter{
		Writer:    w,
		separator: s,
	}
}

// WriteAll function writes results of http health check
func (cw *CSVRepositoryWriter) WriteAll(hcrlts []repository.ResultHTTPCheck) error {
	w := csv.NewWriter(cw.Writer)
	w.Comma = cw.separator

	// Get header of output results file from the struct that define health check result
	var header []string
	// We must pass interfaces or pointer to reflect.ValueOf()
	// If passed struct, ValueOf() would throw panic
	v := reflect.ValueOf(&repository.ResultHTTPCheck{}).Elem()
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		header = append(header, t.Field(i).Name)
	}

	err := w.Write(header)
	if err != nil {
		return err
	}

	for _, r := range hcrlts {
		rec := parseResult(r)
		err := w.Write(rec)
		if err != nil {
			return nil
		}
	}
	// File is wrote only after writer flush
	w.Flush()

	return nil
}

func parseResult(r repository.ResultHTTPCheck) []string {
	rlt := "OK"
	if r.Result != nil {
		rlt = r.Result.Error()
	}

	ts := r.TimeStamp
	u := r.URL.String()
	desc := r.Description
	return []string{ts, u, rlt, desc}
}
