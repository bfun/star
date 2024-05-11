package star

import (
	"encoding/xml"
	"path"
)

var esadmin ESAdmin

func init() {
	esadmin = ParseESAdminFile()
}

type ESAdmin struct {
	XMLName  xml.Name         `xml:"ESAdmin"`
	DtaParms []ESAdminDtaParm `xml:"DtaParmTab>DtaParm"`
}

type ESAdminDtaParm struct {
	DtaName string `xml:"DtaName,attr"`
	Ports   string `xml:"DtaMchTab>DtaMch>IPTab>Item>Port,attr"`
}

func ParseESAdminFile() ESAdmin {
	fullpath := path.Join(getRootDir(), "etc/ESAdmin.xml")
	decoder := getStarFileDecoder(fullpath)
	var v ESAdmin
	err := decoder.Decode(&v)
	if err != nil {
		panic(err)
	}
	return v
}
