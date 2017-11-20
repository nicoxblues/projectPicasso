window.onload= function(){
    document.getElementById("subTitulo").innerHTML  += ", resolucion de pantalla local en  fullScreen = " +  window.screen.availWidth +  "X" +   window.screen.availHeight;

};

function sendData(){

    //  var chHandler =  new chartHandler();

    toggleFullScreen();

    var pixelRatio  = window.devicePixelRatio;
    var posX = document.querySelector('[name="xPos"]').value;
    var posY = document.querySelector('[name="yPos"]').value;


    var mainDiv = document.getElementById("mainPage");
    mainDiv.style.display = 'none';

    var imgElement  = document.getElementById("imageElemnet");
    var  ChartElement = document.getElementById("chart_frame");

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
        var firstLoad = true;
        var socket = new WebSocket(url);
        //   container = $("#container")
        socket.onopen = function() {
            console.log("open")
        };

        socket.onmessage = function (e) {

            console.log("mensaje");

            //toggleFullScreen();
            var data = e.data;

            if (firstLoad) {


                ChartElement.setAttribute("src", data) ;
                ChartElement.style.display = 'block';
                imgElement.style.display = 'none';
                firstLoad = false;
            }else if (data.indexOf("Charts") !== -1) {
                imgElement.setAttribute("src", "data:image/jpg;base64,");
                ChartElement.style.display = 'block';
                imgElement.style.display = 'none';

            }else {
                imgElement.setAttribute("src", "data:image/jpg;base64," + data);
                ChartElement.style.display = 'none';
                imgElement.style.display = 'block';

            }

        };



        socket.onclose = function () {
            console.log("cerrar")
        };
        return socket;
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