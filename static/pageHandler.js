window.onload= function(){
    document.getElementById("subTitulo").innerHTML  += ", resolucion de pantalla local en  fullScreen = " +  window.screen.availWidth +  "X" +   window.screen.availHeight;

};

function sendData(){
    chartHandler();

    toggleFullScreen();

    var pixelRatio  = window.devicePixelRatio;
    var posX = document.querySelector('[name="xPos"]').value;
    var posY = document.querySelector('[name="yPos"]').value;


    var mainDiv = document.getElementById("mainPage");
    mainDiv.style.display = 'none';
    var imgElement  = document.getElementById("imageElemnet");
    var  ChartElement = document.getElementById("chart_div");

    var fullScreenHeight = window.screen.availHeight; // * pixelRatio;
    var fullScreenWidth = window.screen.availWidth; //* pixelRatio;

    var url = "ws://"  + document.location.host + "/ws?height=" + fullScreenHeight + "&width=" + fullScreenWidth + "&coordenadasX=" + posX  + "&coordenadasY=" +  posY ;

    var ws;
    if (window.WebSocket === undefined) {

        return null;
    } else {
        ws = initWS();
    }
    function initWS() {
        var socket = new WebSocket(url);
        //   container = $("#container")
        socket.onopen = function() {
            console.log("open")
        };

        socket.onmessage = function (e) {

            console.log("mensaje");
            //toggleFullScreen();
            var data = e.data;
            imgElement.setAttribute("src","data:image/jpg;base64," + data);
            ChartElement.style.display = 'none';
            imgElement.style.display = 'block';
        };

        socket.onclose = function () {
            console.log("cerrar")
        };
        return socket;
    }




}

function chartHandler(){

    // Load the Visualization API and the corechart package.
    google.charts.load('current', {'packages':['corechart']});

    // Set a callback to run when the Google Visualization API is loaded.
    google.charts.setOnLoadCallback(drawChart);

    var jsonChartData ={
        "chartConfig": {
            "chart":
                {
                    "-type": "pie",
                    "-name": "Personas",
                    "-title": "inserte algo interesante aqui ! :D",
                    "Columns":
                        [
                            {
                                "colType": "string",
                                "colName": "variable1"
                            },
                            {
                                "colType": "number",
                                "colName": "value"
                            }
                        ],
                    "rows":
                        [
                            {
                                "rowName": "test",
                                "rowValue": 500
                            },
                            {
                                "rowName": "testII",
                                "rowValue": 500
                            }
                        ]
                }
        }


    };


    // Callback that creates and populates a data table,
    // instantiates the pie chart, passes in the data and
    // draws it.
    function drawChart() {

        gooleChart = [];
        gooleChart['pie'] = google.visualization.PieChart;
        gooleChart['bar'] = google.visualization.BarChart;
        gooleChart['Column'] = google.visualization.ColumnChart;


        // Create the data table.
        var data = new google.visualization.DataTable();

        data.addColumn('string', 'N');
        data.addColumn('number', 'Value');



        data.addRow(['V', 300]);
        data.addRow(['C', 100]);

        // Set chart options
        var options = {'title':'titulo del grafico, no se , algo! ',
            'width': window.screen.availWidth,
            'height':window.screen.availHeight,


            	animation:{
                    duration: 1000,
                    easing: 'out',},
            vAxis: {minValue:0, maxValue:700}
        };



        /*var button = document.getElementById('b1');

                      // Disabling the button while the chart is drawing.
      button.disabled = true;
      google.visualization.events.addListener(chart, 'ready',
          function() {
            button.disabled = false;
          });*/





        // Instantiate and draw our chart, passing in some options.
        var chartDiv = document.getElementById('chart_div');

        var chart = new  gooleChart['Column'](chartDiv);
        chart.draw(data, options);


        chartDiv.onclick =  function() {
            data = new google.visualization.DataTable();
            var columns = jsonChartData.chartConfig.chart.Columns;
            for (i = 0; i < columns.length; i++){
                data.addColumn(columns[i].colType, columns[i].colName);
            }
            var rows = jsonChartData.chartConfig.chart.rows;

            for (i = 0; i < rows.length; i++){
                data.addRow([rows[i].rowName,rows[i].rowValue]) //data.addRow(['V', 500]);

            }
            //data.addColumn('number', 'Value');


            chart.draw(data, options);
        }



    }
}

function toggleFullScreen() {
    var doc = window.document;
    var docEl = doc.documentElement;

    var requestFullScreen = docEl.requestFullscreen || docEl.mozRequestFullScreen || docEl.webkitRequestFullScreen || docEl.msRequestFullscreen;
    var cancelFullScreen = doc.exitFullscreen || doc.mozCancelFullScreen || doc.webkitExitFullscreen || doc.msExitFullscreen;

    if(!doc.fullscreenElement && !doc.mozFullScreenElement && !doc.webkitFullscreenElement && !doc.msFullscreenElement) {
        requestFullScreen.call(docEl);
    }
    else {
        cancelFullScreen.call(doc);
    }
}