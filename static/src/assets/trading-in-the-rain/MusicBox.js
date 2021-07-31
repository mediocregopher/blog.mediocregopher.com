var midiLoaded = false;
function loadMIDI() {
  MIDI.loadPlugin({
    soundfontUrl: "/assets/trading-in-the-rain/soundfont/",
    instrument: "acoustic_grand_piano",
    onprogress: (state, progress) => {
      console.log("MIDI loading...", progress*100, "%");
    },
    onsuccess: () => {
      console.log("MIDI is ready to be used");
      MIDI.setVolume(0, 127);
      midiLoaded = true;
    },
  });
}

function MusicBox(priceDist, volumeDist) {
  this.priceDist = priceDist;
  this.volumeDist = volumeDist;

  // clamp the keyboard so we're not using the very low notes, they don't sound
  // good.
  const noteRange = {
    //low: 21,
    //low: 36, // C2
    low: 60, // C4, middle C
    high: 108
  };


  function makeScale(tpl) {
    tplObj = {};
    for (i in tpl) {
      tplObj[tpl[i]] = true;
    }
  
    let scale = [];
    for (let note=noteRange.low; note<=noteRange.high; note++) {
      let key = MIDI.noteToKey[note].replace(/\d+$/, "");
      if (tplObj[key]) {
        scale.push(note);
      }
    }
    return scale;
  }

  //this.scale = makeScale(["C", "D", "E", "F", "G", "A", "B"]); // cMajor
  //this.scale = makeScale(["D", "E", "Gb", "G", "A", "Db"]); //dMajor
  //this.scale = makeScale(["C", "D", "E", "G", "A"]); // cMajor pentatonic
  this.scale = makeScale(["F", "G", "A", "C", "D"]); // fMajor pentatonic

  this.playNote = (note, holdFor) => {
    if (!midiLoaded) return;
    let velocity = 127;
    MIDI.noteOn(0, note, velocity, 0);
    MIDI.noteOff(0, note, holdFor);
  };

  this.playTrades = (trades) => {
    if (this.priceDist.length == 0) return;
    for (let i in trades) {
      let noteIdx = this.priceDist.distribute(trades[i].price, 0, this.scale.length-1);
      noteIdx = Math.round(noteIdx);

      let holdFor = 0.25 + this.volumeDist.distribute(trades[i].volume, 0, 1.75);
      let note = this.scale[noteIdx];
      this.playNote(note, holdFor);
    }
  };
}
