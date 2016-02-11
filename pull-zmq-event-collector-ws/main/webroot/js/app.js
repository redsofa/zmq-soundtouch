/**************************************************************** 
*  Very simple graph storage class for illustration purposes 
*/
function GraphStore(){
  this.nodes = Array();
}

GraphStore.prototype.addNode = function (id, iNodeType) {
  this.nodes.push({"id":id, "nodeType":iNodeType});
};


GraphStore.prototype.getNodes = function(){
  return this.nodes;
};

/****************************************************************/


/****************************************************************/
//Basic Force directed graph class for illustration purposes
function Vis(iSelector){
  this.w = 500;
  this.h = 500;
  this.selector = iSelector;
  this.graphStore = null;
  this.svg = null;
  this.force = null;
}

Vis.prototype.createSVG = function(){
  var retVal = d3.select(this.selector)
    .append("svg")
    .attr("width", this.w)
    .attr("height", this.h)
    .attr("border", 1)
    .attr("id", "svg");

   retVal.append("svg:g");

   var borderPath = retVal.append("rect")
               .attr("x", 0)
               .attr("y", 0)
               .attr("height", this.h)
               .attr("width", this.w)
               .style("stroke", "black")
               .style("fill", "none")
               .style("stroke-width", 1);

    return retVal;

};

Vis.prototype.createForceLayout = function(){
  return d3.layout.force()
    .size([this.w, this.h])
    .gravity(.05)
    .distance(50)
    .nodes(this.graphStore.getNodes());
};

Vis.prototype.setGraphStorage = function(iStorage){
  this.graphStore = iStorage;
};

Vis.prototype.recreateContainer = function(){
  d3.select('#svg').remove();
  this.svg = this.createSVG();
  this.force = this.createForceLayout();
}

Vis.prototype.update =function(){ 
  var sRef = this;

  this.recreateContainer();


  var forceDragBehavior = function(){
    var retVal = sRef.force.drag()
      .on("dragstart", function(d){
        
      })
      .on("drag", function(d){
           
      })
      .on("dragend", function(d){

      });
    return retVal;
  };


  var createNodes = function(){
    var node = sRef.svg.append("g")
      .selectAll("circle.node")
      .data(sRef.graphStore.getNodes(), function(d) { 
          return d.id;
      });

    var nodeEnter = node.enter()
      .append("svg:circle")
      .attr("class", "node")
      .attr('cursor', 'move')
      .attr("cx", function(d) { return d.x; })
      .attr("cy", function(d) { return d.y; })
      .attr("fill", function(d){        
        if (d.nodeType == "VolChange"){
          return "lightBlue";
        }else if(d.nodeType == "Heartbeat"){
          return "red" ;
        }else if(d.nodeType == "PresetChange"){
          return "green" ;
        }else if(d.nodeType == "SocketOpen"){
          return "orange" ;
        }else if(d.nodeType == "SocketClose"){
          return "blue" ;      
        }else if(d.nodeType == "RecentsUpdate"){
          return "grey" ;
        }else{
          return "black"
        }
      })
      .style("stroke", "black")      
      .style("stroke-width", 1)
      .style("opacity", 0.7)
      .attr("r", 10)
      .attr("id",function(d) { 
        return "Node;" + d.id;
      })
      .call(forceDragBehavior());

    node.exit().remove();

    return node;
  };

  var onForceNodeTick = function(e){
    //When you fire tick, the previous positions of nodes are used to determine 
    //the new positions, and are stored in the px and py attributes. 
    //If you change these attributes, running tick will update 
    //the x and y attributes to these values.    
    node.attr("cx", function(d) { return d.x; })
        .attr("cy", function(d) { return d.y; });
  };

  this.force.on("tick", function(e) {
      onForceNodeTick(e);
  });

  var node = createNodes();


  function setAllNodesFixed(){
    sRef.graphStore.getNodes().forEach(function (e){
        e.fixed = true;
    });
  };

  this.startForceLayout(500);
  setAllNodesFixed();
};

Vis.prototype.startForceLayout = function(iIterationCount){
    this.force.start();
    var n = iIterationCount;
    for (var i = n * n; i > 0; --i) {
      this.force.tick();
    }
    this.force.stop();
}



$( document ).ready(function($) {
   var eventLogCounter = 0;

   function log(iMessage){
      eventLogCounter++ ;
      if (eventLogCounter % 10 === 0){
         $("#log").empty();
      }
      $("#log").append("<p>" + eventLogCounter + " - " + iMessage + "</p>")
   }

   function connectWebSocket(){
      if ("WebSocket" in window){
         //var ws = new WebSocket("ws://soundtouch.redsofa.ca/entry");
         var ws = new WebSocket("ws://localhost:8081/entry");
            
         ws.onopen = function(){
            log("Opened Socket Connection");
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
         log("Volume Change to : " + $targetVolume.text());
         addNode("VolChange");         
      }else if ($presetUpdated.length) {
         $preset = $presetUpdated.find("itemName");          
         log("Preset Change : " + $preset.text());
         addNode("PresetChange");
      }else if ($recentsUpdated.length){       
         $recents = $recentsUpdated.find("recents");
         log("Recents List Update");
         addNode("RecentsUpdate");
      }else if ($connectionStatusUpdated.length) {
         log("Heartbeat from SoundTouch");
         addNode("Heartbeat");
      }
      //console.log(xmlDoc);
   } 



   connectWebSocket();



});



