package main

import (
	"bufio"
	"os"
	"path"
	"regexp"
	"strings"
)

type TagType struct {
	Name   string
	GetReq map[string]string
	GetRes map[string]string
	SetReq map[string]string
	SetRes map[string]string
}

var RE_TAGTYPE_GETREQ *regexp.Regexp = regexp.MustCompile("get_req:{(.*?)}")
var RE_TAGTYPE_GETRES *regexp.Regexp = regexp.MustCompile("get_res:{(.*?)}")
var RE_TAGTYPE_SETREQ *regexp.Regexp = regexp.MustCompile("set_req:{(.*?)}")
var RE_TAGTYPE_SETRES *regexp.Regexp = regexp.MustCompile("set_res:{(.*?)}")

func tag_type() (tags map[string]TagType) {
	tags = map[string]TagType{}
	fileName := path.Join(getRootDir(), "etc/enum/tag_type.txt")
	file, err := os.Open(fileName)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		line = strings.ReplaceAll(line, " ", "")
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "enumData1") {
			continue
		}
		kv := strings.Split(line, "\t")
		if len(kv) != 2 {
			continue
		}
		var item TagType
		item.Name = kv[0]
		item.getReq(kv[1])
		item.getRes(kv[1])
		item.setReq(kv[1])
		item.setRes(kv[1])
		tags[item.Name] = item
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return
}

func (t *TagType) getReq(line string) {
	t.GetReq = map[string]string{}
	vs := RE_TAGTYPE_GETREQ.FindStringSubmatch(line)
	if vs == nil {
		return
	}
	if len(vs) != 2 {
		panic(line)
	}
	couples := strings.Split(vs[1], ",")
	for _, c := range couples {
		te := strings.Split(c, ":")
		if len(te) != 2 {
			panic(line)
		}
		t.GetReq[te[1]] = te[0]
	}
}
func (t *TagType) getRes(line string) {
	t.GetRes = map[string]string{}
	vs := RE_TAGTYPE_GETRES.FindStringSubmatch(line)
	if len(vs) != 2 {
		panic(line)
	}
	couples := strings.Split(vs[1], ",")
	for _, c := range couples {
		te := strings.Split(c, ":")
		if len(te) != 2 {
			panic(line)
		}
		t.GetRes[te[1]] = te[0]
	}
}
func (t *TagType) setReq(line string) {
	t.SetReq = map[string]string{}
	vs := RE_TAGTYPE_SETREQ.FindStringSubmatch(line)
	if len(vs) != 2 {
		panic(line)
	}
	couples := strings.Split(vs[1], ",")
	for _, c := range couples {
		te := strings.Split(c, ":")
		if len(te) != 2 {
			panic(line)
		}
		t.SetReq[te[1]] = te[0]
	}
}
func (t *TagType) setRes(line string) {
	t.SetRes = map[string]string{}
	vs := RE_TAGTYPE_SETRES.FindStringSubmatch(line)
	if len(vs) != 2 {
		panic(line)
	}
	couples := strings.Split(vs[1], ",")
	for _, c := range couples {
		te := strings.Split(c, ":")
		if len(te) != 2 {
			panic(line)
		}
		t.SetRes[te[1]] = te[0]
	}
}
