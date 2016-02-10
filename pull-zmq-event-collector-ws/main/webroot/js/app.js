$( document ).ready(function($) {

   function refreshDiv(iMessage){
      $("#msg").empty();
      $("#msg").text(iMessage)
   }

   function log(iMessage){
      $("#log").text(iMessage)
   }

   function connectWS(){
      if ("WebSocket" in window){
         //var ws = new WebSocket("ws://soundtouch.redsofa.ca/entry");

         var ws = new WebSocket("ws://localhost:8081/entry");
            
         ws.onopen = function(){
            log("connected")
         }  

         ws.onmessage = function (evt) 
         { 
            var received_msg = evt.data;
            jsonDoc = JSON.parse(received_msg);
            parseXml(jsonDoc)

         };
      
         ws.onclose = function()
         {                   
            alert("Connection is closed..."); 
         };
      }else{
         alert("WebSocket NOT supported by your Browser!");
      }
   }
   

   function parseXml(iXml){
      xmlDoc = $.parseXML( iXml.body );

      $xml = $( xmlDoc );

      $volumeUpdated = $xml.find("volumeUpdated" );
      $presetUpdated = $xml.find("preset");
      $recentsUpdated = $xml.find("recentsUpdated");

      if ( $volumeUpdated.length ){
         console.log("Volume Change");
         $targetVolume = $volumeUpdated.find("targetvolume");
         refreshDiv("Target Volume Change : " + $targetVolume.text()) 
         console.log(xmlDoc);
         return;
      }else if ($presetUpdated.length) {
         console.log("Preset Updated");
         $preset = $presetUpdated.find("itemName");
         refreshDiv("Preset Change : " + $preset.text()) 
         console.log(xmlDoc);
         return;
      }else if ($recentsUpdated.length){
         console.log("Recents Updated");
         $recents = $recentsUpdated.find("recents");
         refreshDiv("Recents List : " + $recents.text()) 
         console.log(xmlDoc);
         return;
      }

      console.log(xmlDoc)

   } 

   connectWS();

});



