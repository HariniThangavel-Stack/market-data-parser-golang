package model

type Options struct {
	OutputFormat string
	FilePath     string
}

type InputFileObj struct {
	Path string
	Name string
}

type Keyvalue map[string]interface{}
