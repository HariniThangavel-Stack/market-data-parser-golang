package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"market-data-parser-golang/model"
	"market-data-parser-golang/utility"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func MarketDataParser(options model.Options) {
	var inputFileObj model.InputFileObj
	var requiredFormat string
	if (model.Options{}) != options && options.FilePath != "" && options.OutputFormat != "" {
		inputFileObj = utility.GetFilePathAndNameFromOptions(options)
		requiredFormat = options.OutputFormat
	} else {
		inputFileObj = utility.GetFilePathAndNameFromCli(os.Args[1:])
		requiredFormat = os.Getenv("OUTPUT_FORMAT")
	}

	//Channels
	fileContent := make(chan string, 50)
	csvToJSONContent := make(chan []model.Keyvalue, 50)
	tableHeaders := make(chan []string, 50)
	done := make(chan struct{})

	//Goroutines
	if utility.IsInputFilePathExists(inputFileObj) {

		go readFile(inputFileObj, fileContent, done)
		go transformCsvToJSON(fileContent, csvToJSONContent, tableHeaders, done)
		csvToJSON := <-csvToJSONContent
		tHeaders := <-tableHeaders
		<-done
		<-done

		if requiredFormat == utility.HTML_FORMAT {
			utility.CreateOutputFolder()
			go fileWriter(inputFileObj, csvToJSON, tHeaders, done)
			<-done
		} else {
			for _, mapData := range csvToJSON {
				x := fmt.Sprint(mapData)
				bytes, _ := json.Marshal(x)
				temp := json.RawMessage(bytes)
				fmt.Printf("JSON returned successfully %#v \n", string(temp))
			}
		}
	} else {
		fmt.Println("File path and output format is required")
	}
}

func readFile(inputFileObj model.InputFileObj, fileContent chan string, done chan struct{}) {
	r, err := ioutil.ReadFile(inputFileObj.Path)
	if err != nil {
		log.Fatalf("Failed to read file: %s", err)
	}
	fmt.Println("File Read successfullly")
	fileContent <- string(r)
	done <- struct{}{}
	close(fileContent)
}

func transformCsvToJSON(fileContent chan string, csvToJSONContent chan []model.Keyvalue, tableHeaders chan []string, done chan struct{}) {
	headers := make([]string, 0)
	headerEndIdex := 0
	content := <-fileContent
	lines := strings.Split(content, "\n")
	for n, item := range lines {
		if utility.ExceuteRegexPattern(item, utility.HEADER_END_REGEX) && headerEndIdex == 0 {
			headers = append(headers, item)
			headerEndIdex = n
			break
		}
		headers = append(headers, item)
	}
	formatedHeaders := utility.GetFormatedHeader(headers)
	formatedBody := utility.GetFormatedBody(lines, headerEndIdex)
	formatedJson := utility.ConstructJson(formatedHeaders, formatedBody)
	tableHeaders <- formatedHeaders
	csvToJSONContent <- formatedJson
	done <- struct{}{}
	close(csvToJSONContent)
}

func fileWriter(inputFileObj model.InputFileObj, csvToJSON []model.Keyvalue, tHeaders []string, done chan struct{}) {
	w, err := os.Create(utility.OUTPUT_DIR_NAME + "/" + inputFileObj.Name + ".html")
	if err != nil {
		log.Fatalf("Failed creating file: %s", err)
	}
	defer w.Close()
	_, err = w.WriteString(utility.ConstructHtml(csvToJSON, tHeaders))
	if err != nil {
		log.Fatalf("Failed writing to file: %s", err)
	}
	fmt.Println("File Written successfullly")
	done <- struct{}{}
}
