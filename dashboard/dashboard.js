"use strict";

// Dashboard settings
const COLOR_BACKGROUND_RED    = 'rgba(255, 99, 132, 0.2)';
const COLOR_BACKGROUND_ORANGE = 'rgba(252, 138, 7, 0.2)';
const COLOR_BACKGROUND_BLUE   = 'rgba(54, 162, 235, 0.2)';
const COLOR_BACKGROUND_GREEN  = 'rgba(18, 190, 53, 0.2)';
const COLOR_BORDER_RED        = 'rgba(255,99,132,1)';
const COLOR_BORDER_ORANGE     = 'rgba(252, 138, 7, 1)';
const COLOR_BORDER_BLUE       = 'rgba(54, 162, 235, 1)';
const COLOR_BORDER_GREEN      = 'rgba(18, 190, 53, 1)';

const COLOR_BACKGROUND_UP   = COLOR_BACKGROUND_GREEN;
const COLOR_BORDER_UP       = COLOR_BORDER_GREEN;
const COLOR_BACKGROUND_DOWN = COLOR_BACKGROUND_ORANGE;
const COLOR_BORDER_DOWN     = COLOR_BORDER_ORANGE;

// API settings
const API_URL = "http://your.hostname.here/api/"

////////////////////
// Dashboard drawing
////////////////////

function makePercentages(up, down) {
    var pctUp   = parseFloat(up)/(up+down)*100;
    var pctDown = parseFloat(down)/(up+down)*100;

    return { up: pctUp.toFixed(2), down: pctDown.toFixed(2) };
}

function drawTodayChart(up, down) {
    var data = makePercentages(up, down);

    var ctx        = document.getElementById("chart-today");
    var todayChart = new Chart(ctx, {
        type: 'doughnut',
        data: {
            labels: ["Up", "Down"],
            datasets: [{
                label: 'Uptime Today',
                data: [data.up, data.down],
                backgroundColor: [
                    COLOR_BACKGROUND_UP,
                    COLOR_BACKGROUND_DOWN
                ],
                borderColor: [
                    COLOR_BORDER_UP,
                    COLOR_BORDER_DOWN
                ],
                borderWidth: 1
            }]
        },
        options: {
            title: {
                display: true,
                text: 'Uptime Today'
            },
            maintainAspectRatio: false,
            cutoutPercentage: 33
        }
    });
}

function drawLifetimeChart(up, down) {
    var data = makePercentages(up, down)

    var ctx        = document.getElementById("chart-lifetime");
    var todayChart = new Chart(ctx, {
        type: 'doughnut',
        data: {
            labels: ["Up", "Down"],
            datasets: [{
                label: 'Lifetime Uptime',
                data: [data.up, data.down],
                backgroundColor: [
                    COLOR_BACKGROUND_UP,
                    COLOR_BACKGROUND_DOWN
                ],
                borderColor: [
                    COLOR_BORDER_UP,
                    COLOR_BORDER_DOWN
                ],
                borderWidth: 1
            }]
        },
        options: {
            title: {
                display: true,
                text: 'Lifetime Uptime'
            },
            maintainAspectRatio: false,
            cutoutPercentage: 33
        }
    });
}

function drawInstantChart(data) {
    var labels = data.map(d => { return d.timestamp });
    var states = data.map(d => { return d.up ? 1 : 0 });

    var ctx          = document.getElementById("chart-instant");
    var instantChart = Chart.Line(ctx, {
        data: {
            labels:  labels,
            datasets: [
                {
                    label: "Instant Uptime",
                    fill: true,
                    lineTension: 0.1, // ignored
                    steppedLine: true,
                    backgroundColor: COLOR_BACKGROUND_UP,
                    borderColor: COLOR_BORDER_UP,
                    borderCapStyle: 'butt',
                    borderDash: [],
                    borderDashOffset: 0.0,
                    borderJoinStyle: 'miter',
                    pointBorderColor: COLOR_BORDER_UP,
                    pointBackgroundColor: "#fff",
                    pointBorderWidth: 1,
                    pointHoverRadius: 5,
                    pointHoverBackgroundColor: COLOR_BACKGROUND_UP,
                    pointHoverBorderColor: COLOR_BORDER_UP,
                    pointHoverBorderWidth: 2,
                    pointRadius: 1,
                    pointHitRadius: 10,
                    data: states,
                    spanGaps: false,
                }
            ]
        },
        options: {
            title: {
                display: true,
                text: "Instant feed"
            },
            legend: {
                display: false
            },
            maintainAspectRatio: false,
            scales: {
                xAxes: [{
                    type: 'time',
                    time: {
                        unit: 'hour',
                        unitStepSize: 2,
                        displayFormats: {
                            hour: 'HH:mm'
                        }
                    },
                }],
                yAxes: [{
                    ticks: {
                        min: 0,
                        max: 1,
                        stepSize: 1,
                        callback: function(label, index, labels) {
                            return label == '1' ? 'Up' : 'Down';
                        }
                    }
                }]
            }
        }
    });
}

function drawAverageChart(data) {
    var labels  = data.map(d => { return d.day });
    var uptimes = data.map(d => { return makePercentages(d.up, d.down).up });

    var ctx          = document.getElementById("chart-30days");
    var instantChart = Chart.Line(ctx, {
        data: {
            labels: labels,
            datasets: [
                {
                    label: "Daily Uptime",
                    fill: true,
                    lineTension: 0.1,
                    backgroundColor: COLOR_BACKGROUND_UP,
                    borderColor: COLOR_BORDER_UP,
                    borderCapStyle: 'butt',
                    borderDash: [],
                    borderDashOffset: 0.0,
                    borderJoinStyle: 'miter',
                    pointBorderColor: COLOR_BORDER_UP,
                    pointBackgroundColor: "#fff",
                    pointBorderWidth: 1,
                    pointHoverRadius: 5,
                    pointHoverBackgroundColor: COLOR_BACKGROUND_UP,
                    pointHoverBorderColor: COLOR_BORDER_UP,
                    pointHoverBorderWidth: 2,
                    pointRadius: 1,
                    pointHitRadius: 10,
                    data: uptimes,
                    spanGaps: false,
                }
            ]
        },
        options: {
            title: {
                display: true,
                text: "30 Days Uptime"
            },
            maintainAspectRatio: false,
            legend: {
                display: false
            },
            scales: {
                xAxes: [{
                    type: 'time',
                    time: {
                        unit: 'day',
                        unitStepSize: 5,
                        displayFormats: {
                            day: 'MMM DD'
                        }
                    }
                }],
                yAxes: [{
                    display: true,
                    ticks: {
                        beginAtZero: true
                    },
                }]
            }
        }
    });
}

////////////
// API Calls
////////////

function getApiInstant() {
    var req = new XMLHttpRequest();
    req.open("GET", API_URL + "instant", true)
    req.setRequestHeader("Accept", "application/json");
    req.send();
    req.addEventListener("readystatechange", function() {
        if(req.readyState === XMLHttpRequest.DONE && req.status === 200) {
            var instantData = JSON.parse(req.responseText)['data'];
            drawInstantChart(instantData);
        }
    });
}

function getApiDaily() {
    var req = new XMLHttpRequest();
    req.open("GET", API_URL + "daily", true)
    req.setRequestHeader("Accept", "application/json");
    req.send();

    req.addEventListener("readystatechange", function() {
        if(req.readyState === XMLHttpRequest.DONE && req.status === 200) {
            var dailyData = JSON.parse(req.responseText)['data'];
            var now = moment();

            var sumUp = 0, sumDown = 0;
            for (var i in dailyData)
            {
                sumUp   += dailyData[i].up;
                sumDown += dailyData[i].down;

                if (moment(dailyData[i].day).isSame(now, 'd'))
                    drawTodayChart(dailyData[i].up, dailyData[i].down);
            }
            drawLifetimeChart(sumUp, sumDown);

            var last30DaysData = dailyData.filter(d => { return now.diff(d.day, 'days') < 30 });
            drawAverageChart(last30DaysData);
        }
    });
}

/////////
// Action
/////////

getApiInstant();
getApiDaily();
