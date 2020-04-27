function RainCanvas(canvasDOM) {
  this.canvas = canvasDOM;
  this.ctx = this.canvas.getContext("2d");

  this.drops = [];
  this.tick = 0;

  // drop: {x, y, intensity, color} (all in range [0, 1], except color which is
  // an array [r,g,b])
  this.newDrop = (newDrop) => {
    if (!document.hasFocus()) return;

    // scale intensity up a bit right off the bat. If the intensity was near 0
    // the drop wouldn't actually show up at all.
    newDrop.intensity = distribute(newDrop.intensity, 0, 1, 0.1, 1);
    newDrop.tick = this.tick;
    this.drops.push(newDrop);
  };

  // alpha isn't really alpha, it's used to determine line width, but it plays
  // the same role.
  this.drawDrop = (drop, alpha) => {
    let cW = this.canvas.width, cH = this.canvas.height;
    let minDim = Math.min(cW, cH);

    let tickDiff = this.tick - drop.tick;
    let radius = tickDiff * (minDim / 250);
    let x = distribute(drop.x, 0, 1, cW*0.1, cW*0.9);
    let y = distribute(drop.y, 0, 1, cH*0.1, cH*0.9);

    this.ctx.beginPath();
    this.ctx.arc(x, y, radius, 0, Math.PI * 2, false);
    this.ctx.closePath();

    // multiple lineWidth by alpha so that the line width drops over time in
    // correspondence with the opacity.
    this.ctx.lineWidth = distribute(drop.intensity, 0, 1, 2, 9) * alpha;

    let r = drop.color[0], g = drop.color[1], b = drop.color[2];
    this.ctx.strokeStyle = `rgba(${r}, ${g}, ${b}, 1)`;
    this.ctx.stroke();
  };

  let requestAnimationFrame = 
    window.requestAnimationFrame || 
    window.mozRequestAnimationFrame || 
    window.webkitRequestAnimationFrame || 
    window.msRequestAnimationFrame;

  this.doTick = () => {
    this.canvas.width = window.innerWidth;
    this.canvas.height = window.innerHeight;
    this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);

    let newDrops = [];
    for (let i in this.drops) {
      let drop = this.drops[i];
      let alpha = distribute(
        this.tick - drop.tick,
        0, 200 * drop.intensity,
        1, 0,
      );
      if (alpha <= 0) continue;

      this.drawDrop(drop, alpha);
      newDrops.push(drop);
    }
    this.drops = newDrops;

    this.tick++;
    requestAnimationFrame(this.doTick);
  };
  requestAnimationFrame(this.doTick);
}
