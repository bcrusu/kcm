package util

import (
	"bytes"
	"text/template"
)

func GenerateTextTemplate(templateStr string, params interface{}) []byte {
	t := template.New("t")

	if _, err := t.Parse(templateStr); err != nil {
		panic(err)
	}

	buffer := &bytes.Buffer{}
	if err := t.ExecuteTemplate(buffer, "t", params); err != nil {
		panic(err)
	}

	return buffer.Bytes()
}
