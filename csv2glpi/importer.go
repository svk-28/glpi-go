package main

import (
	"encoding/csv"
	"strings"
	// "encoding/json"
	"flag"
	"fmt"
	"glpi"
	"os"
	"reflect"
	"strconv"
	"unicode/utf8"
)

var (
	glpiServerAPI   string
	glpiUserToken   string
	glpiAppToken    string
	glpiAssetsType  string
	glpiVersion     bool
	operation       string
	inFileName      string
	inFormat        string
	inCSVSeparator  string
	inCSVComment    string
	outFileName     string
	outFormat       string
	outCSVSeparator string
	outCSVComment   string
)

func CSVToGLPI(glpiAPI *glpi.Session, fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comment, _ = utf8.DecodeRuneInString(inCSVComment)
	reader.Comma, _ = utf8.DecodeRuneInString(inCSVSeparator)
	reader.TrimLeadingSpace = true

	headers, err := reader.Read()
	if err != nil {
		fmt.Println(err)
	}

	for {
		record, err := reader.Read()
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(record)
		var (
			recordMap glpi.GlpiPhone
		)

		// recordMap.Name = "test"
		// Read the headers and find field name in structure
		for fieldIndex, fieldName := range headers {
			s := reflect.ValueOf(&recordMap).Elem()
			typeOfT := s.Type()
			for i := 0; i < s.NumField(); i++ {
				// if the field is found in the structure, we assign it a value of the appropriate type
				if strings.EqualFold(strings.ToLower(fieldName), strings.ToLower(typeOfT.Field(i).Name)) {
					switch fmt.Sprint(typeOfT.Field(i).Type) {
					case "int":
						convValue, _ := strconv.Atoi(record[fieldIndex])
						s.Field(i).SetInt(int64(convValue))
					case "string":
						s.Field(i).SetString(record[fieldIndex])
					case "bool":
						convValue, _ := strconv.ParseBool(record[fieldIndex])
						s.Field(i).SetBool(convValue)
					}
				}
			}
		}
		// fmt.Println(recordMap)
		fmt.Println("Added", glpiAssetsType, "with ID:", glpiAPI.ItemOperation("add", glpiAssetsType, recordMap))
	}
}

func ExportFromGLPI(glpiAPI *glpi.Session, fileName string) {
	fmt.Println(fileName)
	// file, err := os.Open(fileName)
	// if err != nil {
	// 	panic(err)
	// }
	// defer file.Close()

	// reader := csv.NewReader(file)
	// reader.Comma = outCSVSeparator
}

func main() {

	flag.StringVar(&glpiServerAPI, "glpi-server-api", "", "Glpi instance API URL. GLPI_SERVER_URL")
	flag.StringVar(&glpiUserToken, "glpi-user-token", "", "Glpi user API token. GLPI_USER_TOKEN")
	flag.StringVar(&glpiAppToken, "glpi-app-token", "", "Glpi aplication API token. GLPI_APP_TOKEN")
	flag.StringVar(&glpiAssetsType, "glpi-assets-type", "", "Glpi assets type (aka Computer, Printer, Phone, e.t.c.)")
	flag.BoolVar(&glpiVersion, "glpi-version", false, "Get the Glpi version")

	flag.StringVar(&operation, "operation", "", "import - import data to GLPI from file\n\texport - export data from GLPI to file\n\t")

	flag.StringVar(&inFileName, "in-file", "", "the name of the file containing the imported data")
	flag.StringVar(&inFormat, "in-format", "csv", "input file data format")
	flag.StringVar(&inCSVSeparator, "in-csv-separator", ";", "input csv-file separator")
	flag.StringVar(&inCSVComment, "in-csv-comment", "#", "input csv-file comments symbol")

	flag.StringVar(&outFileName, "out-file", "", "the name of the file containing the exported data")
	flag.StringVar(&outFormat, "out-format", "json", "Output format: 'go' - go map, 'json' - json string, 'csv' - common separated value")
	flag.StringVar(&outCSVSeparator, "out-csv-separator", ";", "output csv-file separator")
	flag.StringVar(&outCSVComment, "out-csv-comment", "#", "output csv-file comments symbol")
	flag.Parse()

	//-----------------------------------------------------------
	if glpiServerAPI == "" && os.Getenv("GLPI_SERVER_API") == "" {
		fmt.Println("Send error: make sure environment variables `GLPI_SERVER_API`, or used with '-glpi-server-api' argument")
		os.Exit(1)
	} else if glpiServerAPI == "" && os.Getenv("GLPI_SERVER_API") != "" {
		glpiServerAPI = os.Getenv("GLPI_SERVER_API")
	}

	if glpiUserToken == "" && os.Getenv("GLPI_USER_TOKEN") == "" {
		fmt.Println("Send error: make sure environment variables `GLPI_USER_TOKEN`, or used with '-glpi-user-token' argument")
		os.Exit(1)
	} else if glpiUserToken == "" && os.Getenv("GLPI_USER_TOKEN") != "" {
		glpiUserToken = os.Getenv("GLPI_USER_TOKEN")
	}

	if glpiAppToken == "" && os.Getenv("GLPI_APP_TOKEN") == "" {
		fmt.Println("Send error: make sure environment variables `GLPI_APP_TOKEN`, or used with '-glpi-app-token' argument")
		os.Exit(1)
	} else if glpiAppToken == "" && os.Getenv("GLPI_APP_TOKEN") != "" {
		glpiAppToken = os.Getenv("GLPI_APP_TOKEN")
	}

	glpiSession, err := glpi.NewSession(glpiServerAPI, glpiUserToken, glpiAppToken)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = glpiSession.InitSession()
	if err != nil {
		fmt.Println(err)
		return
	}

	if glpiVersion {
		fmt.Println("GLPI version:", glpiSession.Version())
	}

	switch operation {
	case "import":
		if inFileName == "" || glpiAssetsType == "" {
			fmt.Println("Setting the input file name and glpi assets type, use the '-in-file' and '-glpi-assets-type' argument")
		} else {
			CSVToGLPI(glpiSession, inFileName)
		}
	case "export":
		if outFileName == "" {
			fmt.Println("Setting the output file name, use the '-out-file' argument")
		} else {
			ExportFromGLPI(glpiSession, outFileName)
		}
	}
	glpiSession.KillSession()
}
