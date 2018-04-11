package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/dlavalli/golangchest/logstash/smshubbeat/beater"
)

func main() {
	err := beat.Run("smshubbeat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
