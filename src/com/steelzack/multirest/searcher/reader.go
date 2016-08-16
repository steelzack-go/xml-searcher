package searcher

import (
	"os"
	"bufio"
	"log"
	"fmt"
	"path/filepath"
	"strings"
	"bytes"
	"github.com/kokardy/saxlike"
)

type Config struct {
	CASSANDRA     struct {
					  PORT int
					  HOST string
				  }
	NETWORKFOLDER struct {
					  FOLDER string
				  }
}

type ReaderEA struct {
	keystorage KeyStorage
}
func ReadFileLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func (readerea ReaderEA)  ReadLines(path string, config Config) {
	lines, err := ReadFileLines(path)
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}
	var stringBuffer bytes.Buffer
	for _, line := range lines {
		stringBuffer.WriteString(line)
	}

	source := stringBuffer.String()
	r := bytes.NewReader([]byte(source))
	handler := new(PartialHandler)
	parser := saxlike.NewParser(r, handler)
	parser.SetHTMLMode()
	parser.Parse()
	readerea.keystorage.InsertKeys(handler.value1, handler.value2)
}

func (readerea ReaderEA) WalkThrough(config Config) (err error) {
	readerea.keystorage.OpenDatabase(config.CASSANDRA.HOST, config.CASSANDRA.PORT)
	readerea.keystorage.Init()
	fileList := []string{}
	err = filepath.Walk(string(filepath.Separator) + config.NETWORKFOLDER.FOLDER + string(filepath.Separator), func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() && strings.HasSuffix(path, readdocument) {
			readerea.ReadLines(path, config)
			fileList = append(fileList, path)
		}
		return nil
	})
	for _, file := range fileList {
		fmt.Println(file)
	}
	readerea.keystorage.session.Close()
	return
}
