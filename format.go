package star

import (
	"encoding/xml"
	"fmt"
	"log"
	"path"
	"regexp"
	"strings"
	"sync"
)

type FormatTab struct {
	XMLName xml.Name `xml:"FormatTab"`
	Formats []Format `xml:"Format"`
}
type Format struct {
	FmtName string       `xml:"FmtName,attr"`
	FmtType string       `xml:"FmtType,attr"`
	Items   []FormatItem `xml:"ItemTab>Item"`
	SubFmts []string
}
type FormatItem struct {
	ItemType  string `xml:"ItemType,attr"`
	ItemIgnr  string `xml:"ItemIgnr,attr"`
	ElemName  string `xml:"ElemName,attr"`
	XmlType   string `xml:"XmlType,attr"`
	XmlName   string `xml:"XmlName,attr"`
	SubName   string `xml:"SubName,attr"`
	SubExpr   string `xml:"SubExpr"`
	ConstData string `xml:"ConstData"`
}

func trimFormatCDATA(formats map[string]Format) {
	re2 := regexp.MustCompile(`.*?\?(.*?):(.*)`)
	re3 := regexp.MustCompile(`.*\?\(.*\?(.*?):(.*)\):(.*)`)
	for kf, vf := range formats {
		for ki, vi := range vf.Items {
			if vi.ItemType == "fmt" {
				if vi.SubName == "" {
					panic(kf + "/fmt/SubName empty")
				}
				vf.SubFmts = append(vf.SubFmts, vi.SubName)
				continue
			}
			if vi.ItemType == "expr" {
				s := strings.TrimSpace(vi.SubExpr)
				if s != "" {
					s = strings.Replace(s, " ", "", -1)
					s = strings.Replace(s, `"`, "", -1)
					if strings.Contains(s, "?") && strings.Contains(s, ":") {
						fs := re3.FindStringSubmatch(s)
						if len(fs) > 0 {
							if len(fs) != 4 {
								panic(kf + s)
							}
						} else {
							fs = re2.FindStringSubmatch(s)
							if len(fs) != 3 {
								panic(kf + s)
							}
						}
						for _, f := range fs[1:] {
							vf.SubFmts = append(vf.SubFmts, f)
						}
					} else {
						vf.SubFmts = append(vf.SubFmts, s)
					}
				}
				vf.Items[ki].SubExpr = s
			} else {
				vf.Items[ki].SubExpr = ""
			}
			vf.Items[ki].ConstData = strings.TrimSpace(vi.ConstData)
		}
		formats[kf] = vf
	}
}

func formatArrayToMap(formats []Format, m map[string]Format) {
	for _, v := range formats {
		m[v.FmtName] = v
	}
}
func parseOneFormatXml(fileName string, ch chan Format, wg *sync.WaitGroup) {
	defer wg.Done()
	fullPath := path.Join(getRootDir(), fileName)
	decoder := getStarFileDecoder(fullPath)
	var v FormatTab
	err := decoder.Decode(&v)
	if err != nil {
		panic(err)
	}
	for _, f := range v.Formats {
		ch <- f
	}
}

func ParseAllFormatXml(wg *sync.WaitGroup) {
	defer wg.Done()
	m := make(map[string]Format)
	files := getFormatFiles()
	ch := make(chan Format, 1024*1024)
	wg2 := new(sync.WaitGroup)
	for _, f := range files {
		wg2.Add(1)
		go parseOneFormatXml(f, ch, wg2)
	}
	wg2.Wait()
	close(ch)
	for f := range ch {
		m[f.FmtName] = f
	}
	log.Print("fmt count:", len(m))
	trimFormatCDATA(m)
	FMTMAP = m
	log.Print("Format.xml parse success")
}

func getVarFormatName(dta, svc, format string) string {
	if !strings.Contains(format, "+") {
		return format
	}
	RIG := "RIG($stdmsgtype+$stdprocode,10)"
	if strings.Contains(format, RIG) {
		format = strings.Replace(format, RIG, svc, 1)
	}
	CBS := "$CBS_FORMAT"
	if strings.Contains(format, CBS) {
		format = strings.Replace(format, CBS, svc, 1)
	}
	SVC := "$__SVCNAME"
	if strings.Contains(format, SVC) {
		format = strings.Replace(format, SVC, svc, 1)
	}
	SDTA := "$NESB_SDTA_NAME"
	if strings.Contains(format, SDTA) {
		to := dta
		d, ok := DTAMAP[dta]
		if !ok {
			panic(dta + svc + format)
		}
		if d.NESB_SDTA_NAME != "" {
			to = d.NESB_SDTA_NAME
		}
		s, ok := d.Services[svc]
		if !ok {
			panic(dta + svc + format)
		}
		if s.NESB_SDTA_NAME != "" {
			to = s.NESB_SDTA_NAME
		}
		format = strings.Replace(format, SDTA, to, 1)
	}
	DDTA := "$NESB_DDTA_NAME"
	if strings.Contains(format, DDTA) {
		to := dta
		d, ok := DTAMAP[dta]
		if !ok {
			panic(dta + svc + format)
		}
		if d.NESB_DDTA_NAME != "" {
			to = d.NESB_DDTA_NAME
		}
		format = strings.Replace(format, DDTA, to, 1)
	}
	return strings.ReplaceAll(format, "+", "")
}
func findElemsInFormat(dta, svc, format string) map[string]map[string]string {
	fm := make(map[string]map[string]string)
	em := make(map[string]string)
	f, ok := FMTMAP[format]
	if !ok {
		fmt.Println(dta, svc, format, "format not found")
		return nil
	}
	for _, v := range f.Items {
		if v.ItemType != "item" || v.ItemIgnr == "yes" {
			continue
		}
		em[v.XmlName] = v.ElemName
	}
	fm[format] = em
	for _, sub := range f.SubFmts {
		sub2 := getVarFormatName(dta, svc, sub)
		if sub2 != sub {
			fmt.Printf("getVarFormatName %v.%v %v -> %v\n", dta, svc, sub, sub2)
		}
		fm2 := findElemsInFormat(dta, svc, sub2)
		for k, v := range fm2 {
			fm[k] = v
		}
	}
	return fm
}

func findElemsInFormat2(dta, svc, format string) map[string]string {
	m := make(map[string]string)
	f, ok := FMTMAP[format]
	if !ok {
		fmt.Println(dta, svc, format, "format not found")
		return nil
	}
	constId := 1
	for _, v := range f.Items {
		if v.ItemType != "item" || v.ItemIgnr == "yes" {
			continue
		}
		if len(v.ConstData) > 0 {
			k := fmt.Sprintf("常量 %d: [%s]", constId, v.ConstData)
			m[k] = v.XmlName
			constId++
		}
		if v.XmlType != "tag" {
			continue
		}
		m[v.ElemName] = v.XmlName
	}
	for _, sub := range f.SubFmts {
		sub2 := getVarFormatName(dta, svc, sub)
		if sub2 != sub {
			fmt.Printf("getVarFormatName %v.%v %v -> %v\n", dta, svc, sub, sub2)
		}
		m2 := findElemsInFormat2(dta, svc, sub2)
		for k, v := range m2 {
			m[k] = v
		}
	}
	return m
}
