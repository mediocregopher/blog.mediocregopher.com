function SeriesComposer(resource, rainCanvas, color) {
  this.rainCanvas = rainCanvas;
  this.color = color;

  this.priceDist = new Distributor(200);
  this.volumeDist = new Distributor(200);
  this.musicBox = new MusicBox(this.priceDist, this.volumeDist);

  this.enabled = false;
  this.setEnabled = (enabled) => this.enabled = enabled;
  this.getEnabled = () => { return this.enabled; }

  this.totalTrades = 0;
  this.getTotalTrades = () => { return this.totalTrades; }

  this.cw = new CW(resource);
  this.cw.ontrades = (trades) => {
    if (this.totalTrades > 0 && this.enabled) {
      let priceVols = {}; // sum volumes by price, for deduplication
      for (let i in trades) {
        let price = trades[i].price, volume = trades[i].volume;
        if (!priceVols[price]) priceVols[price] = 0;
        priceVols[price] += volume;
      }

      trades = []; // overwrite trades with deduplicated ones.
      for (let price in priceVols) {
        let volume = priceVols[price];
        let intensity = this.volumeDist.distribute(volume, 0, 1);
        this.rainCanvas.newDrop({
          x: this.priceDist.distribute(price, 0, 1),
          y: Math.random(),
          intensity: intensity,
          color: this.color,
        });

        trades.push({price: price, volume: volume});
      }

      this.musicBox.playTrades(trades);
    }

    for (let i in trades) {
      this.priceDist.add(trades[i].price);
      this.volumeDist.add(trades[i].volume);
    }

    this.totalTrades += trades.length;
    if (this.ontrades) this.ontrades(trades);
  };
}
