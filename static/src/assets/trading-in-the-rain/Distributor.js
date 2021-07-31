function distribute(val, minOld, maxOld, minNew, maxNew) {
  let scalar = (val - minOld) / (maxOld - minOld);
  return minNew + ((maxNew - minNew) * scalar);
}

function Distributor(capacity) {
  this.cap = capacity;

  this.reset = () => {
    this.arr = [];
    this.arrSorted = [];
    this.length = 0;
  };
  this.reset();

  // add adds the given value into the series, shifting off the oldest value if
  // the series is at capacity.
  this.add = (val) => {
    this.arr.push(val);
    if (this.arr.length >= this.cap) this.arr.shift();
    this.arrSorted = this.arr.slice(); // copy array
    this.arrSorted.sort();
    this.length = this.arr.length;
  };

  // distribute finds where the given value falls within the series, and then
  // scales that into the given range (inclusive).
  this.distribute = (val, min, max) => {
    if (this.length == 0) throw "cannot locate within empty Distributor";

    let idx = this.length;
    for (i in this.arrSorted) {
      if (val < this.arrSorted[i]) {
        idx = i;
        break;
      }
    }

    return distribute(idx, 0, this.length, min, max);
  };
}

