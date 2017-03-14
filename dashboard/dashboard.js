"use strict";

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

function drawTodayChart() {
    var ctx        = document.getElementById("chart-today");
    var todayChart = new Chart(ctx, {
        type: 'doughnut',
        data: {
            labels: ["Up", "Down"],
            datasets: [{
                label: 'Uptime Today',
                data: [90, 10],
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

function drawLifetimeChart() {
    var ctx        = document.getElementById("chart-lifetime");
    var todayChart = new Chart(ctx, {
        type: 'doughnut',
        data: {
            labels: ["Up", "Down"],
            datasets: [{
                label: 'Lifetime Uptime',
                data: [95, 5],
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

function drawInstantChart() {
    var ctx          = document.getElementById("chart-instant");
    var instantChart = Chart.Line(ctx, {
        data: {
            labels:  ['2017-03-11 00:00', '2017-03-11 01:00', '2017-03-11 02:00', '2017-03-11 03:00',
                '2017-03-11 04:00', '2017-03-11 05:00', '2017-03-11 06:00', '2017-03-11 07:00',
                '2017-03-11 08:00', '2017-03-11 09:00', '2017-03-11 10:00', '2017-03-11 11:00',
                '2017-03-11 12:00', '2017-03-11 13:00', '2017-03-11 14:00', '2017-03-11 15:00',
                '2017-03-11 16:00', '2017-03-11 17:00', '2017-03-11 18:00', '2017-03-11 19:00',
                '2017-03-11 20:00', '2017-03-11 21:00', '2017-03-11 22:00', '2017-03-11 23:00',
                '2017-03-11 24:00' ],
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
                    data: [1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 0, 1],
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

function drawAverageChart() {
    var ctx          = document.getElementById("chart-30days");
    var instantChart = Chart.Line(ctx, {
        data: {
            labels: ['2017-03-01', '2017-03-02', '2017-03-03', '2017-03-04', '2017-03-05', '2017-03-06',
                '2017-03-07', '2017-03-08', '2017-03-09', '2017-03-10', '2017-03-11', '2017-03-12', '2017-03-13',
                '2017-03-14', '2017-03-15', '2017-03-16', '2017-03-17', '2017-03-18', '2017-03-19', '2017-03-20',
                '2017-03-21', '2017-03-22', '2017-03-23', '2017-03-24', '2017-03-25', '2017-03-26', '2017-03-27',
                '2017-03-28', '2017-03-29', '2017-03-30'],
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
                    data: [95, 97, 97, 96, 95, 97, 95, 97, 97, 96,
                        97, 96, 96, 95, 94, 95, 94, 94, 94, 96, 94, 97,
                        94, 94, 96, 97, 95, 94, 94, 97],
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

drawTodayChart();
drawLifetimeChart();
drawInstantChart();
drawAverageChart();
