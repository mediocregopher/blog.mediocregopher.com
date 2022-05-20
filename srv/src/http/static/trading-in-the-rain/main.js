const tradesToWaitFor = 5;

function hexToRgb(hex) {
  var result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
  return [
    parseInt(result[1], 16),
    parseInt(result[2], 16),
    parseInt(result[3], 16),
  ];
}

const colorPalette = [
  "#28D2EE",
  "#ED778E",
  "#6557DC",
  "#EEE386",
  "#B55AA0",
];

// Use https://api.cryptowat.ch/markets/<exchange>
// or https://api.cryptowat.ch/pairs (for "all")
const markets = {
  "kraken":[
    {
      name: "BTCUSD",
      resource: "markets:87:trades",
    },
    {
      name: "BTCEUR",
      resource: "markets:86:trades",
    },
    {
      name: "BTCEUR",
      resource: "markets:96:trades",
    },
    {
      name: "ETHEUR",
      resource: "markets:97:trades",
    },
    {
      name: "BCHUSD",
      resource: "markets:146:trades",
    },
  ],
  "bitfinex":[
    {
      name: "BTCUSD",
      resource: "markets:1:trades",
    },
    {
      name: "ETHUSD",
      resource: "markets:4:trades",
    },
    {
      name: "BSVUSD",
      resource: "markets:5558:trades",
    },
    {
      name: "BTCEUR",
      resource: "markets:415:trades",
    },
    {
      name: "XRPUSD",
      resource: "markets:25:trades",
    },
  ],
  "all": [
    {
      name: "BTCUSD",
      resource: "instruments:9:trades",
    },
    {
      name: "ETHUSD",
      resource: "instruments:125:trades",
    },
    {
      name: "LTCUSD",
      resource: "instruments:138:trades",
    },
    {
      name: "EOSUSD",
      resource: "instruments:4:trades",
    },
    {
      name: "XRPUSD",
      resource: "instruments:160:trades",
    },
  ],
};

const exchange = "all";

function fillMarketP() {
  let marketsEl = document.getElementById("markets");
  for (let i in markets[exchange]) {
    let name = markets[exchange][i].name;
    let color = colorPalette[i];
    if (i > 0) marketsEl.innerHTML += "</br>";
    marketsEl.innerHTML += `<strong style="color: ${color}; font-size: 2rem;">${name}</strong>`;
  }
}

function run() {
  document.getElementById("button").style.display = "none";

  let progress = document.getElementById("progress");
  progress.innerHTML = "Connecting to Cryptowat.ch...";

  let canvas = document.getElementById("rainCanvas");
  let rainCanvas = new RainCanvas(canvas);

  let modalHidden = false;
  for (let i in markets[exchange]) {
    let seriesComposer = new SeriesComposer(
      markets[exchange][i].resource,
      rainCanvas, 
      hexToRgb(colorPalette[i]),
    );

    seriesComposer.cw.onconnect = () => {
      progress.innerHTML = "Preloading a few trades before continuing.";
    };

    // wait for each series to rech tradesToWaitFor before letting it begin.
    // Hide the modal when the first series is enabled.
    seriesComposer.ontrades = (trades) => {
      if (!modalHidden && seriesComposer.getTotalTrades() < tradesToWaitFor) {
        progress.innerHTML += "."; // indicate that _something_ is happening
        return;
      }
      if (!modalHidden) {
        let modal = document.getElementById("tradingInRainModal");
        modal.style.display = "none";
        modalHidden = true;
      }
      seriesComposer.setEnabled(true);
      seriesComposer.ontrades = undefined;
    };
  }
}

function autorun() {
  loadMIDI();
}
