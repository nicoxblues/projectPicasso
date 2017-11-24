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
    var chartElement = document.getElementById("chart_frame");
    var resetElement = document.getElementById("reset_frame");

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

                setTimeout(function(){

                    chartElement.setAttribute("src", data) ;
                    chartElement.style.display = 'block';

                    imgElement.style.display = 'none';
                    firstLoad = false;}
                    ,300)


            }else if (data.indexOf("Charts") !== -1) {
                imgElement.setAttribute("src", "data:image/jpg;base64,");
                imgElement.style.display = 'none';
                resetElement.style.display = 'none';

                chartElement.style.display = 'block';


            }else if (data.indexOf("resetFinneg") !== -1) {

                chartElement.style.display = 'none';
                imgElement.style.display = 'none';


                resetElement.style.display = 'block';
                resetElement.contentWindow.location.reload();


            }else {
                imgElement.setAttribute("src", "data:image/jpg;base64," + data);
                chartElement.style.display = 'none';
                resetElement.style.display = 'none';
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