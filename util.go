package star

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"os"
	"path"
	"sort"
	"strings"
)

func getRootDir() string {
	return os.Getenv("FAPWORKDIR")
}

func fileScanner(filename string) (file *os.File, scanner *bufio.Scanner) {
	var err error
	file, err = os.Open(path.Join(getRootDir(), filename))
	if err != nil {
		panic(err)
	}
	scanner = bufio.NewScanner(file)
	return
}
func getGbFileDecoder(fileName string) *xml.Decoder {
	buf, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	decoder := xml.NewDecoder(transform.NewReader(bytes.NewReader(buf), simplifiedchinese.GB18030.NewDecoder()))
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		return transform.NewReader(input, simplifiedchinese.GB18030.NewDecoder()), nil
	}
	return decoder
}
func getStarFileDecoder(fileName string) *xml.Decoder {
	buf, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	olds := `<?xml version="1.0" encoding="ISO-8859-1" ?>`
	news := `<?xml version="1.0" encoding="UTF-8" ?>`
	buf = bytes.Replace(buf, []byte(olds), []byte(news), 1)
	r := transform.NewReader(bytes.NewReader(buf), simplifiedchinese.GB18030.NewDecoder())
	return xml.NewDecoder(r)
}

func isSVR(dta string) bool {
	return strings.Contains(dta, "_SVR") || strings.Contains(dta, "_SGET") || strings.Contains(dta, "_SPUT") || strings.Contains(dta, "_PAY") || strings.Contains(dta, "_CLS")
}

func isCLT(dta string) bool {
	return strings.Contains(dta, "_CLT") || strings.Contains(dta, "_CPUT") || strings.Contains(dta, "_CGET")
}

func getKeysInMap[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func getSomeKeysInMap[T any](m map[string]T, sub string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		if strings.Contains(k, sub) {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	return keys
}
