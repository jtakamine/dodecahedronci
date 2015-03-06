package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"net/http"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"
	"time"
)

func req(url string, method string, response interface{}) (err error) {
	return reqBody(url, method, struct{}{}, response)
}

func reqBody(url string, method string, body interface{}, response interface{}) (err error) {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(response)
	if err != nil {
		return err
	}

	return nil
}

func newAction(inner func(string, *cli.Context)) (fn func(*cli.Context)) {
	return func(c *cli.Context) {
		endpt := strings.Trim(c.GlobalString("endpoint"), " \t\n")
		if endpt == "" {
			msg := "No endpoint specified, defaulting to http://localhost:8000. If this is not the correct (dodeccontrol) endpoint, please set the DODEC_ENDPOINT environment variable or use the global \"endpoint\" option: dodec-cli -endpoint <endpoint> [command]."
			fmt.Println(msg)
			endpt = "http://localhost:8000"
		}
		endpt = cleanEndpt(endpt)

		inner(endpt, c)
	}
}

func cleanEndpt(endpt string) (cleanEndpt string) {
	cleanEndpt = strings.TrimSuffix(endpt, "/")
	if !strings.HasPrefix(cleanEndpt, "http://") && !strings.HasPrefix(cleanEndpt, "https://") {
		cleanEndpt = "http://" + cleanEndpt
	}

	if !(strings.Count(cleanEndpt, ":") == 2) {
		cleanEndpt += ":8000"
	}
	cleanEndpt += "/"

	return cleanEndpt
}

func printObj(obj interface{}) {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, '\t', 0)

	for i := 0; i < v.NumField(); i++ {
		name := strings.ToUpper(t.Field(i).Name)
		val := fmt.Sprint(v.Field(i).Interface())
		val = serialize(val)
		fmt.Fprintln(w, name+":\t"+val)
	}

	err := w.Flush()
	if err != nil {
		panic(err)
	}
}

func printRows(rows interface{}, header bool) {
	rowsV := reflect.ValueOf(rows)
	if rowsV.Len() == 0 {
		if header {
			fmt.Println("No results!")
		}
		return
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, '\t', 0)

	if header {
		colNames := []string{}
		t := reflect.TypeOf(rowsV.Index(0).Interface())
		for i := 0; i < t.NumField(); i++ {
			n := strings.ToUpper(t.Field(i).Name)
			colNames = append(colNames, n)
		}
		fmt.Fprintln(w, strings.Join(colNames, "\t"))
	}

	for i := 0; i < rowsV.Len(); i++ {
		r := rowsV.Index(i).Interface()
		v := reflect.ValueOf(r)
		colVals := []string{}
		for i := 0; i < v.NumField(); i++ {
			val := fmt.Sprint(v.Field(i).Interface())
			val = serialize(val)
			colVals = append(colVals, val)
		}
		fmt.Fprintln(w, strings.Join(colVals, "\t"))
	}

	err := w.Flush()
	if err != nil {
		panic(err)
	}
}

func printLogs(logs []Log, header bool) {
	type row struct {
		Timestamp string
		Message   string
	}

	rows := []row{}
	for _, l := range logs {
		r := row{
			Timestamp: l.Created.Format(time.RFC3339),
			Message:   l.Message,
		}
		rows = append(rows, r)
	}

	printRows(rows, header)
}

func serialize(s string) (serializedS string) {
	serializedS = s
	serializedS = strings.Replace(serializedS, "\r", "\\r", -1)
	serializedS = strings.Replace(serializedS, "\n", "\\n", -1)

	return serializedS
}

func requireArg(args []string, cmdName string, argName string, argExampleVal string) (val string) {
	if len(args) == 0 {
		fmt.Printf("Please supply the %s argument. E.g.,\n\tdodec-cli [global options] %s \"%s\"\n", argName, cmdName, argExampleVal)
		os.Exit(1)
	}

	val = args[0]
	return val
}
