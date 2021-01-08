var colors = [
    "rgb(255, 138, 128)",
    "rgb(213, 0, 0)",
    "rgb(118, 255, 3)",
    "rgb(130, 177, 255)",
    "rgb(41, 98, 255)"
]

var typeNames = [
    "Owner With Request", "Owner Terminal", "Idle", "Waiter With Request", "Waiter Terminal"
]


var labelsvg = d3.select("svg#legend")
var labelwidth = +labelsvg.attr("width")
var labelheight = +labelsvg.attr("height")

function addLabel() {

    //TODO:Refactor
    let lastPos = -20

    for (let index = 0; index < 5; index++) {
        lastPos = lastPos + 30
        addLabelNode(typeNames[index], colors[index], lastPos)
    }

    labelsvg.append("line").merge(link)
        .attr("class", "link")
        .attr('marker-end', 'url(#arrowhead)')
        .attr("x1", 10)
        .attr("y1", 160)
        .attr("x2", 100)
        .attr("y2", 160);
    labelsvg.append("text").attr("x", 110).attr("y", 160).text("Link").style("font-size", "15px").attr("alignment-baseline", "middle")

}

function addLabelNode(text, color, y) {
    labelsvg.append("circle").attr("cx", 10).attr("cy", y).attr("r", 8).style("fill", color)
    labelsvg.append("text").attr("x", 30).attr("y", y + 5).text(text).style("font-size", "15px").attr("alignment-baseline", "middle")
}

var svg = d3.select("svg#display"),
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
    links = [],
    queue_elements = [],
    current_owner = "";

var simulation = d3.forceSimulation(nodes)
    .force("charge", d3.forceManyBody().strength(-50))
    .force("link", d3.forceLink().id(function (d) {
        return d.MyAddress;
    }).distance(100))
    .force("x", d3.forceX())
    .force("y", d3.forceY())
    .alphaTarget(1)
    .on("tick", ticked);

/*

 */

var g = svg.append("g").attr("transform", "translate(" + width / 2 + "," + height / 2 + ")")
var link = g.append("g").attr("stroke", "#000").attr("stroke-width", 1.5).selectAll(".link")
var node = g.append("g").selectAll(".node")

restart();

function restart() {


    node = node.data(nodes)

    node.exit().remove();

    node = node.enter()
        .append("circle").merge(node).attr("fill", function (d) {
            return colors[d.Type]
        }).attr("r", 8)
        .on("mouseover", function (d) {

            //TODO:Refactor VVVVVVV

            d3.select(this).style("stroke", "black")
            let parent = d3.select(this.parentNode)
            parent.append("text")
                .attr("x", d.x + 4)
                .attr("y", d.y + 4)
                .attr("fill", "black")
                .text(`Address: ${d.MyAddress}`)
                .style("font-size", "15px")

            if (d.Link != "") {
                parent.append("text")
                    .attr("x", d.x + 4)
                    .attr("y", d.y + 4 + 15)
                    .attr("fill", "black")
                    .text(`Link -> ${d.Link}`)
                    .style("font-size", "15px")
            }

        }).on("mouseout", function (d) {
            d3.select(this).style("stroke", "none")
            d3.select(this.parentNode).selectAll("text").remove()
        }).on("dblclick", function (d) {
            remoteRequest(d)

        })

        .call(d3.drag()
            .on("start", dragstarted)
            .on("drag", dragged)
            .on("end", dragended));


    node.exit().remove();

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
    node.attr("cx", function (d) {
        return d.x;
    })
        .attr("cy", function (d) {
            return d.y;
        })


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
}

function mouseout(d) {
    d3.select("body").selectAll('div.tooltip').remove();
}

function remoteRequest(d) {
    fetch(`http://${d.MyAddress}/remoteRequest`)
        .then(response => response.json())
        .then(json => console.log(json))
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
    d3.json("/data", function (error, data) {
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

function updateQueue() {
    var queueTable = document.getElementById("queue");
    queueTable.innerHTML = ""
    for (var i = 0; i < queue_elements.length; i++) {
        var newRow = queueTable.insertRow(i)
        var newCell = newRow.insertCell(0)
        newCell.innerHTML = queue_elements[i]
        newCell.style.backgroundColor = stringToColour(queue_elements[i])
    }

    var ownerTable = document.getElementById("owner_table")
    ownerTable.innerText = ""
    if (current_owner != "") {
        var newRow = ownerTable.insertRow(0)
        var newCell = newRow.insertCell(0)
        newCell.innerHTML = current_owner
        newCell.style.backgroundColor = stringToColour(current_owner)

    }

}


function getQueue() {
    d3.json("/queue", function (error, data) {
        if (data.queue_nodes == null) {
            queue_elements = []
        } else {
            queue_elements = data.queue_nodes
        }
        if (data.owner == null) {
            current_owner = []
        } else {
            current_owner = data.owner
        }

    })

}

var stringToColour = function (str) {
    for (var i = 0, hash = 0; i < str.length; hash = str.charCodeAt(i++) + ((hash << 5) - hash)) ;
    color = Math.floor(Math.abs((Math.sin(hash) * 10000) % 1 * 16777216)).toString(16);
    return '#' + Array(6 - color.length + 1).join('0') + color;
}


addLabel()
getData()
d3.interval(function () {
    getData()
    getQueue()
    restart()
    updateQueue()
}, 200)
