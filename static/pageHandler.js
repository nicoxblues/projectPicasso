window.onload= function(){
    document.getElementById("subTitulo").innerHTML  += ", resolucion de pantalla local en  fullScreen = " +  window.screen.availWidth +  "X" +   window.screen.availHeight;

};

function sendData(){

    var chHandler =  new chartHandler();

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
            try{
                chHandler.jsonData =  JSON.parse (data);
                if (chHandler._googleChart) {
                    chHandler.fillData();
                    chHandler.draw();
                }

            }catch(ex){
                if (data.indexOf("Charts") !== -1) {
                    imgElement.setAttribute("src", "data:image/jpg;base64,");
                    ChartElement.style.display = 'block';
                    imgElement.style.display = 'none';

                } else {
                    imgElement.setAttribute("src", "data:image/jpg;base64," + data);
                    ChartElement.style.display = 'none';
                    imgElement.style.display = 'block';
                }


            }

        };

        socket.onclose = function () {
            console.log("cerrar")
        };
        return socket;
    }




}

function chartHandler(){

    this.gooleChart = [];
    var self = this;
    this.drawChart = function(){
        self.gooleChart['pie'] = google.visualization.PieChart;
        self.gooleChart['bar'] = google.visualization.BarChart;
        self.gooleChart['Column'] = google.visualization.ColumnChart;


        // Create the data table.


        self.fillData();

        // Instantiate and draw our chart, passing in some options.
        var chartDiv = document.getElementById('chart_div');

        self._googleChart = new self.gooleChart[self.jsonData.type](chartDiv);


        self._googleChart.draw(self.data, self.option);

    };


    // Load the Visualization API and the corechart package.
    google.charts.load('current', {'packages':['corechart']});

    // Set a callback to run when the Google Visualization API is loaded.
    google.charts.setOnLoadCallback(this.drawChart);




};

chartHandler.prototype.draw = function () {
    if (this._googleChart)
        this._googleChart.draw(this.data, this.option);

};

chartHandler.prototype.fillData = function (jsonData) {


    this.data = new google.visualization.DataTable();
    var columns  = this.jsonData.columns;

    for ( var i = 0; i < columns.length; i++){
        this.data.addColumn(columns[i].colType, columns[i].colName);



    }

    this.data.addRow([columns[0].colName, parseInt(this.jsonData.rows[0].rowValue)]);
    this.data.addRow([columns[1].colName, parseInt(this.jsonData.rows[1].rowValue)]);




    // Set chart options
    var options = {'title':this.jsonData.title,
        'width': window.screen.availWidth,
        'height':window.screen.availHeight,


        animation:{
            duration: 1000,
            easing: 'out'},
        vAxis: {minValue:0, maxValue:700}
    };

    this.option = options

};



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