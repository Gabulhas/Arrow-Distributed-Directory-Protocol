var colors = [
    "rgb(255, 138, 128)",
    "rgb(213, 0, 0)",
    "rgb(118, 255, 3)",
    "rgb(41, 98, 255)",
    "rgb(130, 177, 255)"
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
    .force("charge", d3.forceManyBody().strength(-100))
    .force("link", d3.forceLink().id(function (d) {
        return d.id;
    }).distance(100))
    .force("x", d3.forceX())
    .force("y", d3.forceY())
    .alphaTarget(1)
    .on("tick", ticked);

var g = svg.append("g").attr("transform", "translate(" + width / 2 + "," + height / 2 + ")"),
    link = g.append("g").attr("stroke", "#000").attr("stroke-width", 1.5).selectAll(".link"),
    node = g.append("g").selectAll(".node")

restart
();


function restart() {

    node = node.data(nodes, function (d) {
        return d;
    });

    node.exit().remove();

    node = node.enter()
        .append("circle").attr("fill", function (d) {
            return colors[d.type]
        }).attr("r", 8).merge(node)
        .on("mouseover", function (d) {
            mouseover(d)

        }).on("mouseout", function (d) {
            mouseout(d)
        })

    ;

    node.append("title")
        .text(function (d) {
            return d.name;
        });


    link = link.data(links, function (d) {
        return d.source + "-" + d.target;
    });
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

function mouseover(d) {
    var div = d3.select("body").append("div")
        .attr("class", "tooltip")
        .style("opacity", 1)
        .style("left", (d.y + 120) + "px")
        .style("top", (d.x - 20) + "px")
        .html(
            "<table style='font-size: 10px; font-family: sans-serif;' >" +
            "<tr><td>Name: </td><td>" + d.name + "</td></tr>" +
            "<tr><td>Value: </td><td>" + d.id + "</td></tr>" +
            "</table>"
        );
    if (d.parent) mouseover(d.parent);
}

function mouseout(d) {
    d3.select("body").selectAll('div.tooltip').remove();
}

d3.json("/data", function (error, data) {
    if (error) throw error;


    //aqui ele deve só fazer alteraões, não fazer reset
    nodes = data.nodes
    links = data.links

    restart()
})

