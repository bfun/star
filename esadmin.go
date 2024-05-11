package star

import (
	"encoding/xml"
	"path"
)

type ESAdmin struct {
	XMLName  xml.Name         `xml:"ESAdmin"`
	DtaParms []ESAdminDtaParm `xml:"DtaParmTab>DtaParm"`
}

type ESAdminDtaParm struct {
	DtaName    string             `xml:"DtaName,attr"`
	DtaDesc    string             `xml:"DtaDesc,attr"`
	IPTabItems []ESAdminIPTabItem `xml:"DtaMchTab>DtaMch>IPTab>Item"`
	Nodes      []ESAdminDtaNode   `xml:"DtaNodeTab>DtaNode"`
}

type ESAdminIPTabItem struct {
	Port string `xml:"Port,attr"`
}

type ESAdminDtaNode struct {
	NodeName string             `xml:"NodeName,attr"`
	Port     ESAdminDtaNodePort `xml:"PORTDefs>port"`
}

type ESAdminDtaNodePort struct {
	NodeIP   string `xml:"NodeIP,attr"`
	NodePort string `xml:"NodePort,attr"`
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
