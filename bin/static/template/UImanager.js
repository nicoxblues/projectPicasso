const url = document.location.host;

function show(){

    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
            //      document.getElementById("demo").innerHTML = this.responseText;
        }
    };
    xhttp.open("POST", "/loadImage2", true);
    xhttp.send();


}

function send(obj){

    var nodes = document.querySelectorAll( "form#" + obj.form.id +" input[type=number]");

    for(var i=0;i<nodes.length;i++) {
        nodes[i].value



    }


}

function sendAll(){
    var formsCollection = document.getElementsByTagName("form");
    for(var i=0;i<formsCollection.length;i++)
    {


    }


}



function resetFinneg(){
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
            //    document.getElementById("demo").innerHTML = this.responseText;
        }
    };
    xhttp.open("GET", "/resetFinneg", true);
    xhttp.send();
}

function showCharts(){

    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
            //    document.getElementById("demo").innerHTML = this.responseText;
        }
    };
    xhttp.open("POST", "/showChart", true);
    xhttp.send();


}