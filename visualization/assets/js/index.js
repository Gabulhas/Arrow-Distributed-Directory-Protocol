var colors = [
    "rgb(255, 138, 128)",
    "rgb(213, 0, 0)",
    "rgb(118, 255, 3)",
    "rgb(130, 177, 255)",
    "rgb(41, 98, 255)"
]



var svg = d3.select("svg"),
    width = +svg.attr("width"),
    height = +svg.attr("height"),
    color = d3.scaleOrdinal(d3.schemeCategory10);

svg
    .append('defs')
    .append('marker')
    .attrs({
        'id': 'arrowhead',
        'viewBox': '-0 -5 10 10',
        'refX': 13,
        'refY': 0,
        'orient': 'auto',
        'markerWidth': 13,
        'markerHeight': 13,
        'xoverflow': 'visible'
    })
    .append('svg:path')
    .attr('d', 'M 0,-5 L 10 ,0 L 0,5')
    .attr('fill', '#999')
    .style('stroke', 'none');

var
    nodes = [],
    links = [];

var simulation = d3.forceSimulation(nodes)
    .force("charge", d3.forceManyBody().strength(-50))
    .force("link", d3.forceLink().id(function(d) {
        return d.MyAddress;
    }).distance(100))
    .force("x", d3.forceX())
    .force("y", d3.forceY())
    .alphaTarget(1)
    .on("tick", ticked);

/*

 */

var g = svg.append("g").attr("transform", "translate(" + width / 2 + "," + height / 2 + ")"),
    link = g.append("g").attr("stroke", "#000").attr("stroke-width", 1.5).selectAll(".link"),
    node = g.append("g").selectAll(".node")

restart();


function restart() {


    node = node.data(nodes)

    node.exit().remove();

    node = node.enter()
        .append("circle").merge(node).attr("fill", function(d) {
            return colors[d.Type]
        }).attr("r", 8)
        .on("mouseover", function(d) {
            d3.select(this).style("stroke", "black") // set the line colour
            mouseover(d)
        }).on("mouseout", function(d) {
            d3.select(this).style("stroke", "none") // set the line colour
            mouseout(d)
        })
        .call(d3.drag()
            .on("start", dragstarted)
            .on("drag", dragged)
            .on("end", dragended));


    link = link.data(links)
    link.exit().remove();
    link = link.enter().append("line").merge(link)
        .attr("class", "link")
        .attr('marker-end', 'url(#arrowhead)');


    simulation.nodes(nodes);
    simulation.force("link").links(links);
    simulation.alpha(1).restart();
}


function ticked() {
    node.attr("cx", function(d) {
            return d.x;
        })
        .attr("cy", function(d) {
            return d.y;
        })


    link.attr("x1", function(d) {
            return d.source.x;
        })
        .attr("y1", function(d) {
            return d.source.y;
        })
        .attr("x2", function(d) {
            return d.target.x;
        })
        .attr("y2", function(d) {
            return d.target.y;
        });
}

function mouseover(d) {
    var div = d3.select("body").append("div")
        .attr("class", "tooltip")
        .style("opacity", 1)
        .style("left", (d.y + 120) + "px")
        .style("top", (d.x - 20) + "px")
        .html(
            "<table class=\"nodeDetails\" style='font-size: 10px; font-family: sans-serif;' >" +
            "<tr><td>Address: </td><td>" + d.MyAddress + "</td></tr>" +
            "<tr><td>Link: </td><td>" + d.Link + "</td></tr>" +
            "<tr><td>Type: </td><td>" + d.Type + "</td></tr>" +
            "</table>"
        );
    if (d.parent) mouseover(d.parent);
}

function dragstarted(d) {
    if (!d3.event.active) simulation.alphaTarget(0.3).restart();
    d.fx = d.x;
    d.fy = d.y;
}

function dragged(d) {
    d.fx = d3.event.x;
    d.fy = d3.event.y;
}

function dragended(d) {
    if (!d3.event.active) simulation.alphaTarget(0);
    //d.fx = null;
    //d.fy = null;
}

function mouseout(d) {
    d3.select("body").selectAll('div.tooltip').remove();
}

function recolor() {


}

//Complexidade muito alta, mudar para set/map

function updateNodes(newNodes) {

    for (let i = nodes.length - 1; i >= 0; i--) {

        var found = newNodes.findIndex((findNode) => findNode.MyAddress == nodes[i].MyAddress)
        if (found != -1) {
            var foundNode = newNodes[found]
            nodes[i].Type = foundNode.Type
            nodes[i].Link = foundNode.Link

            newNodes.splice(found, 1)

        } else {
            nodes.splice(i, 1)
        }

    }
    newNodes.forEach((newNode) => nodes.push(newNode))

}

function getData() {
    d3.json("/data", function(error, data) {
        if (error) throw error;

        //nodes = data.nodes


        if (data.nodes == null) {
            nodes = []
        } else {
            //aqui ele deve só fazer alteraões, não fazer reset
            updateNodes(data.nodes)
        }
        if (data.links == null) {
            links = []
        } else {

            links = data.links
        }

    })
}


getData()
d3.interval(function() {
    getData()
    restart()
    d3.select("body").selectAll('div.tooltip').remove();


}, 300)