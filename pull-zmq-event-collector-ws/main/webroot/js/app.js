$( document ).ready(function($) {
   var eventLogCounter = 0;

   function log(iMessage){
      eventLogCounter++ ;
      if (eventLogCounter % 10 === 0){
         $("#log").empty();
      }
      $("#log").append("<p>" + eventLogCounter + " - " + iMessage + "</p>")
   }

   function connectWS(){
      if ("WebSocket" in window){
         //var ws = new WebSocket("ws://soundtouch.redsofa.ca/entry");
         var ws = new WebSocket("ws://localhost:8081/entry");
            
         ws.onopen = function(){
            log("Opened Socket Connection")
         }  

         ws.onmessage = function (evt){ 
            var received_msg = evt.data;
            jsonDoc = JSON.parse(received_msg);
            parseXml(jsonDoc)
         };
      
         ws.onclose = function(){                   
            log("Closed Socket Connection");
         };
      }else{
         log("WebSocket NOT supported by your Browser!");
         alert("WebSocket NOT supported by your Browser!");
      }
   }
   
   function parseXml(iXml){
      xmlDoc = $.parseXML( iXml.body );

      $xml = $( xmlDoc );

      $volumeUpdated = $xml.find("volumeUpdated" );
      $presetUpdated = $xml.find("preset");
      $recentsUpdated = $xml.find("recentsUpdated");
      $connectionStatusUpdated = $xml.find("connectionStateUpdated");

      if ($volumeUpdated.length ){
         $targetVolume = $volumeUpdated.find("targetvolume");
         log("Volume Change to : " + $targetVolume.text());         
      }else if ($presetUpdated.length) {
         $preset = $presetUpdated.find("itemName");          
         log("Preset Change : " + $preset.text());         
      }else if ($recentsUpdated.length){       
         $recents = $recentsUpdated.find("recents");
         log("Recents List Update");          
      }else if ($connectionStatusUpdated.length) {
         log("Heartbeat from SoundTouch");
      }
      //console.log(xmlDoc);
   } 

   connectWS();

});



