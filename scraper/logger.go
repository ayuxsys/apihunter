package scraper

import (
	"apihunter/must"
	"encoding/json"
	"io"
	"log"
	"os"
)

var (
	opLogger = log.New(os.Stdout, "", 0)
)

func (s *Scraper) logOutput(output Output) {
	switch s.Opts.OutputMode {
	case JsonLOp:
		res := must.Eval(json.Marshal(output))
		opLogger.Println(string(res))
	case MarkDownOp:
		opLogger.Printf("## %s\n", output.Query) // todo log organisation need to change core fn
		for _, res := range output.Results {
			opLogger.Printf("- [%s](%s)", res.Title, res.Href)
		}
		opLogger.Print("\n")
	}
}

// SetupLogTarget configures the logger to write logs to both the specified file
// and standard output. It returns a cleanup function that closes the log file.
func (s *Scraper) SetupLogTarget(filePath string) func() {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	multiWriter := io.MultiWriter(file, os.Stdout)
	opLogger.SetOutput(multiWriter)

	return func() { file.Close() }
}
