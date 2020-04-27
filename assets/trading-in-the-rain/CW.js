function CW(resource) {
  this.conn = new WebSocket('wss://stream.cryptowat.ch/connect?apikey=GPDLXH702E1NAD96OSBO');
  this.conn.binaryType = "arraybuffer";

  this.conn.onopen = () => {
    console.log("CW websocket connected");
    if (this.onconnect) this.onconnect();
  }

  let decoder = new TextDecoder();
  this.conn.onmessage = (msg) => {
    let d = JSON.parse(decoder.decode(msg.data));

    // The server will always send an AUTHENTICATED signal when you establish a valid connection
    // At this point you can subscribe to resources
    if (d.authenticationResult && d.authenticationResult.status === 'AUTHENTICATED') {
      if (this.onauth) this.onauth();
      this.conn.send(JSON.stringify({
        subscribe: {
          subscriptions: [
            {streamSubscription: {resource: resource}},
          ],
        }
      }));
      return;
    }

    // Market data comes in a marketUpdate
    // In this case, we're expecting trades so we look for marketUpdate.tradesUpdate
    if (!d.marketUpdate || !d.marketUpdate.tradesUpdate) {
      return;
    }

    let trades = d.marketUpdate.tradesUpdate.trades;
    for (let i in trades) {
      trades[i].price = parseFloat(trades[i].priceStr);
      trades[i].volume = parseFloat(trades[i].amountStr);
    }
    if (this.ontrades) this.ontrades(trades);
  }

  this.close = () => this.conn.close();
}
