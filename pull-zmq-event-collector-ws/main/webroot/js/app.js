/**************************************************************** 
*  Very simple graph storage class for illustration purposes 
*/
function GraphStore(){
  this.nodes = Array();
  this.links = Array();
}

GraphStore.prototype.addNode = function (id, nodeGroup) {
  this.nodes.push({"id":id, "nodeGroup":nodeGroup});
};

GraphStore.prototype.addLink = function (source, target, value) {
  this.links.push({"source":this.findNode(source),"target":this.findNode(target),"value":value});
};

GraphStore.prototype.findNode = function(id) {
  for (var i in this.nodes) {
      if (this.nodes[i]["id"] === id){ 
        return this.nodes[i];
      }
  };
};

GraphStore.prototype.getNodes = function(){
  return this.nodes;
};

GraphStore.prototype.getLinks = function(){
  return this.links;
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
    .linkDistance(50)
    .nodes(this.graphStore.getNodes())
    .links(this.graphStore.getLinks());
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

  var getHullGroups = function(){
    //groupBy nodeGroup
    var retVal = d3.nest().key(function(d){
      return d.nodeGroup;
    }).entries(sRef.graphStore.getNodes());
    return retVal;
  };

  var groups = getHullGroups();

  var hullPath = function (d) {
      if (d.values.length > 2) {
          return "M" + d3.geom.hull(d.values.map(function (i) {
                  return [i.x, i.y];
              }))
              .join("L")
              + "Z";
      } else {
          var x = d.values[0].x;
          var y = d.values[0].y;
          var x2 = parseFloat(x) + 1;
          var y2 = parseFloat(y) + 1;
          var x3 = parseFloat(x) - 1;
          var y3 = parseFloat(y) - 1;
          var path = "M" + x + "," + y + "L" + x2 + "," + y2 + "L" + x3 + "," + y3;

          for (var i = 1; i < d.values.length; i++) {
              var a = d.values[i].x;
              var b = d.values[i].y;
              path += "L" + a + "," + b;
          }
          path += "Z";
          return path;
      }
  };

  //Credit : http://jsfiddle.net/FEM3e/7/
  var updateHulls = function(){
    sRef.svg.selectAll("path")
      .data(groups)
      .attr("d", hullPath)
      .enter()
      .insert("path", "g")
      .attr('class','hull')
      .attr("id", function (d,i) { 
        return "path_" + i; 
      })
      .attr('cursor', 'crosshair')
      .style("fill", "lightblue")
      .style("stroke", "lightblue")
      .style("stroke-width", 57)
      .style("stroke-linejoin", "round")
      .style("opacity", .2)
      .call(hullDragBehavior())
      .on('dblclick', function(d){
          if (d3.event.defaultPrevented){
              return;
          } 
          console.log("Double Click");
          console.log(d);
      });
  };

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

  var createLinks = function(){
    var link = sRef.svg.selectAll("link")
      .data(sRef.graphStore.getLinks(), function(d) { 
          return d.source.id + "-" + d.target.id; 
      });

    var linkEnter = link.enter()
      .append("line")
      .attr("class", "link")
      .style('stroke', 'red')
      .style("stroke-width", 2);

    link.exit().remove();
    return link;
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

  var hullDragBehavior  = function(e) {
    var retVal = d3.behavior.drag()
      .on('dragstart', function (d) {
          d3.event.sourceEvent.stopPropagation(); // silence other listeners
          d3.select(this).style("stroke", "red");
      })
      .on('drag', function (d, i) {
          d3.selectAll('.d3-tip').remove();
          var nodeGroup = parseInt(d.key);
          var dx = d3.event.dx;
          var dy = d3.event.dy;
          sRef.moveNodeGroup(nodeGroup, dx, dy);
      })
      .on('dragend',function(d){
          d3.event.sourceEvent.stopPropagation(); // silence other listeners
          d3.select(this).style("stroke", "lightblue");          
      });

      return retVal;
    };

  this.force.on("tick", function(e) {
      updateHulls(e);
      onForceNodeTick(e);
      onForceLinkTick();
  });

  var link = createLinks();
  var node = createNodes();

  var onForceLinkTick = function(){
    link.attr("x1", function (d) {
          return d.source.x;
       })
       .attr("y1", function (d) {
          return d.source.y;
       })
       .attr("x2", function (d) {
          return d.target.x;
       })
       .attr("y2", function (d) {
          return d.target.y;
       });
  };

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

Vis.prototype.moveNodeGroup = function(iGroup, dx, dy){
  
  var filteredNodes = d3.selectAll(".node").filter(function(n){
      return n.nodeGroup === iGroup;
  });

  filteredNodes.attr("cx", function (n){
      n.px = n.px + dx;
      n.x = n.x + dx;
      return n.x;
    }).attr("cy", function(n){
      n.py = n.py + dy;
      n.y = n.y + dy;
      return n.y;
    });

  this.startForceLayout(1);
};



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
            log("Opened Socket Connection");
            addData();
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
   
   //To add nodes and groups...
   var groupCount = 0;
   var nodeCount = 0;

   //Create instance of GraphStorage
   var storage = new GraphStore();

   //Create instance of Vis
   var myVis = new Vis("#visDiv");

   //Set graph storage
   myVis.setGraphStorage(storage);

   var addData = function(){
     var group = groupCount++;
     var node1 = nodeCount++;
     var node2 = nodeCount++
     storage.addNode(node1, group);
     storage.addNode(node2, group);
     storage.addLink(node1, node2, {"value": 12, "nodeGroup": group});
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
         addData();         
      }else if ($presetUpdated.length) {
         $preset = $presetUpdated.find("itemName");          
         log("Preset Change : " + $preset.text());
         addData();
      }else if ($recentsUpdated.length){       
         $recents = $recentsUpdated.find("recents");
         log("Recents List Update");
         addData();
      }else if ($connectionStatusUpdated.length) {
         log("Heartbeat from SoundTouch");
         addData();
      }
      //console.log(xmlDoc);
   } 

   connectWS();



});



