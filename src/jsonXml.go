package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"io/ioutil"
)


type Data struct {
	XMLName    xml.Name      `xml:"data" json:"-"`
	ChartConfig []chartConfig `xml:"chartConfig" json:"chartConfig"`

}

type chartConfig struct {
	XMLName   xml.Name `xml:"chartConfig" json:"-"`
	Title string `xml:"title,attr" json:"title"`
	Type string `xml:"type,attr" json:"type"`
	Columns[] *chartColumn `xml:"Columns>element" json:"columns"`
	Rows[]  *chartRow `xml:"rows>element" json:"rows"`


}

type chartRow struct {
	RowName string `xml:"rowName" json:"rowName"`
	RowValue string `xml:"rowValue" json:"rowValue"`
}

type chartColumn struct {
	ColName string `xml:"colType" json:"colType"`
	ColValue string `xml:"colName" json:"colName"`
}


func loadConfg() {

	//rawXmlData := "<data> <chartConfig name='Personas' title='inserte algo interesante aqui ! :D' type='Pie'> <Columns> <element> <colName>variable1</colName> <colType>string</colType> </element> <element> <colName>value</colName> <colType>number</colType> </element> </Columns> <rows> <element> <rowName>test</rowName> <rowValue>500</rowValue> </element> <element> <rowName>testII</rowName> <rowValue>500</rowValue> </element> </rows> </chartConfig> <chartConfig name='otracosa' title='otra cosa que insertr ! :D' type='bar'> <Columns> <element> <colName>variable2</colName> <colType>string</colType> </element> <element> <colName>value2</colName> <colType>number</colType> </element> </Columns> <rows> <element> <rowName>test</rowName> <rowValue>500</rowValue> </element> <element> <rowName>testII</rowName> <rowValue>500</rowValue> </element> </rows> </chartConfig> </data>"

	xmlFile , _ := os.Open("chartConfig.xml")

	b, _ := ioutil.ReadAll(xmlFile)
	var data Data
	xml.Unmarshal(b, &data)



	jsonData, _ := json.Marshal(data)
	fmt.Println(string(jsonData))

}
