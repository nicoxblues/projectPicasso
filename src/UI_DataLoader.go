package main

import (
	"html/template"
	"net/http"

)

type formData struct {

	FormTitle            string
	DeviceIdentification string
	InputName            []string
	FormID string


}

type pageData struct {
	NumClientConnected int
	Forms [] formData

}

func mainPageLoader(w http.ResponseWriter,hanlder  *ClientHandler )  {



	tmpl, _ := template.ParseFiles("static/template/manager.html")
	var data []formData
	hanlder.serverConf.parseConfigForDataAdmin(&data)
	pData := pageData{len(hanlder.clients),data}


	//data := pageData{Forms:}



	tmpl.Execute(w, pData)


}

func  updateChart(chartID string, connection *ClientHandler ){

	chart :=  connection.charConf[chartID]

	chart.Rows[0].RowValue = "50"
	chart.Rows[1].RowValue = "300"

	for cli := range connection.clients{
		if cli.clientID == chartID {
			cli.socket.WriteJSON(chart)
		}
	}





}
