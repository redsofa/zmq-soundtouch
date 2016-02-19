/*
Copyright 2016 Rene Richard

This file is part of zmq-soundtouch.

zmq-soundtouch is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

zmq-soundtouch is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with zmq-soundtouch.  If not, see <http://www.gnu.org/licenses/>.
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
  d3.selectAll(".tt").remove();

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


    var tooltip = d3.select("body")
      .append("div")
      .attr("class", "tt")
      .style("position", "absolute")
      .style("z-index", "10")
      .style("visibility", "hidden");


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
      .on("mouseover", function(){return tooltip.style("visibility", "visible");})
      .on("mousemove", function(e){
        tooltip.text(e.id + " - " + e.nodeType);
        return tooltip.style("top", (event.pageY-10)+"px").style("left",(event.pageX+10)+"px");
      })
      .on("mouseout", function(){return tooltip.style("visibility", "hidden");})
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