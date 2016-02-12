

$( document ).ready(function($) {
   var eventLogCounter = 0;

   function log(iMessage, iClass){      
      if (eventLogCounter % 10 === 0){
         $("#log").empty();
      }
      $("#log").append('<table><tr><td><div class="'+ iClass + '">&nbsp;&nbsp;&nbsp;</div></td><div><td> &nbsp;&nbsp;' + eventLogCounter + ' - ' + iMessage + '</td></tr></table></div>');
      eventLogCounter++ ;
   }

   function connectWebSocket(){
      if ("WebSocket" in window){
         //var ws = new WebSocket("ws://soundtouch.redsofa.ca/entry");
         var ws = new WebSocket("ws://localhost:8081/entry");
            
         ws.onopen = function(){
            log("Opened Socket Connection", "orange");
            addNode("SocketOpen");
         }  

         ws.onmessage = function (evt){ 
            var received_msg = evt.data;
            jsonDoc = JSON.parse(received_msg);
            parseXml(jsonDoc)
         };
      
         ws.onclose = function(){                   
            log("Closed Socket Connection");
            addNode("SocketClose");
         };
      }else{
         log("WebSocket NOT supported by your Browser!");
         alert("WebSocket NOT supported by your Browser!");
      }
   }
   
   //To add nodes and groups...
   var groupCount = 0;
   var nodeCount = 0;

   //Create instance of GraphStorage
   var storage = new GraphStore();

   //Create instance of Vis
   var myVis = new Vis("#visDiv");

   //Set graph storage
   myVis.setGraphStorage(storage);

   var addNode = function(iNodeType){
     var node1 = nodeCount++;
     storage.addNode(node1, iNodeType);
     myVis.update();
   };


   function parseXml(iXml){
      xmlDoc = $.parseXML( iXml.body );

      $xml = $( xmlDoc );

      $volumeUpdated = $xml.find("volumeUpdated" );
      $presetUpdated = $xml.find("preset");
      $recentsUpdated = $xml.find("recentsUpdated");
      $connectionStatusUpdated = $xml.find("connectionStateUpdated");

      if ($volumeUpdated.length ){
         $targetVolume = $volumeUpdated.find("targetvolume");
         log("Volume Change to : " + $targetVolume.text(), "lightblue");
         addNode("VolChange");         
      }else if ($presetUpdated.length) {
         $preset = $presetUpdated.find("itemName");          
         log("Preset Change : " + $preset.text(), "green");
         addNode("PresetChange");
      }else if ($recentsUpdated.length){       
         $recents = $recentsUpdated.find("recents");
         log("Recents List Update", "grey");
         addNode("RecentsUpdate");
      }else if ($connectionStatusUpdated.length) {
         log("Heartbeat from SoundTouch", "red");
         addNode("Heartbeat");
      }
      //console.log(xmlDoc);
   } 

   connectWebSocket();



});



