package main

import (
	"encoding/xml"
	"os"
	"io/ioutil"
	"encoding/json"


)

type imageDefaultConfig struct {
	DefaultResolutionX int `xml:"defResolutionX"`
	DefaultResolutionY int `xml:"defResolutionY"`

	DefaultResizeX int `xml:"defResizeX"`
	DefaultResizeY int `xml:"defResizeY"`
}


type Data struct {
	XMLName     xml.Name      `xml:"data" json:"-"`
	ChartConfig []chartConfig `xml:"chartConfig" json:"chartConfig"`
	jsonByte    []byte
}

type chartConfig struct {
	XMLName xml.Name        `xml:"chartConfig" json:"-"`
	ChartID string          `xml:"id,attr" json:"id"`
	Title   string          `xml:"title,attr" json:"title"`
	Type    string          `xml:"type,attr" json:"type"`
	SvgRoot string          `xml:"svgRoot,attr" json:"svgRoot"`
	Columns [] *chartColumn `xml:"Columns>element" json:"columns"`
	Rows    []  *chartRow   `xml:"rows>element" json:"rows"`
}





type chartRow struct {
	RowName string `xml:"rowName" json:"rowName"`
	RowValue string `xml:"rowValue" json:"rowValue"`
}

type chartColumn struct {
	ColName string `xml:"colType" json:"colType"`
	ColValue string `xml:"colName" json:"colName"`
}


type serverConfiguration struct {
	ConfigImg          imageDefaultConfig
	ChartConfiguration Data
	currentChartNode   int
	currentChart       *chartConfig

}


func (servConf *serverConfiguration) nextChart() *chartConfig {

	if servConf.currentChartNode == len(servConf.ChartConfiguration.ChartConfig){
		servConf.currentChartNode = 0
	}

	servConf.currentChart = &servConf.ChartConfiguration.ChartConfig[servConf.currentChartNode]
	servConf.currentChartNode++





	return servConf.currentChart



}

func (servConf *serverConfiguration) parseConfigForDataAdmin(frm *[]formData) {

	formList := make([]formData, len(servConf.ChartConfiguration.ChartConfig))

	for i, chart := range servConf.ChartConfiguration.ChartConfig {
		formList[i].FormID = chart.ChartID
		formList[i].FormTitle = chart.Title
		formList[i].InputName = make([]string, len(chart.Columns))
		for j, col := range chart.Columns {
			formList[i].InputName[j] = col.ColValue
		}

		*frm = append(*frm, formList[i])

	}

}

func (servConf *serverConfiguration) getCurrentChart() *chartConfig {
	return servConf.currentChart

}

func (servConf *serverConfiguration) loadConfig() {
	xmlFile , _ := os.Open("chartConfig.xml")

	b, _ := ioutil.ReadAll(xmlFile)
	var data Data
	var chConfig Data
	xml.Unmarshal(b, &data)

	jsonServerConfig, _ := json.Marshal(data)
	json.Unmarshal(jsonServerConfig,&chConfig )

	servConf.ChartConfiguration = chConfig
	servConf.ChartConfiguration.jsonByte = jsonServerConfig


}




