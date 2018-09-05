console.log("main.js started");

// For asynchronously loading the qr code library, loadQRLib returns a promise
// which will be resolved when the library is loaded.
var qrLibProm;
var loadQRLib = () => {
    if (qrLibProm) { return qrLibProm; }
    qrLibProm = new Promise((resolve, reject) => {
        console.log("loading qrcode.min.js");
        var script = document.createElement('script');
        script.src = "/assets/qrcode.min.js";
        script.async = true;
        script.onload = () => { resolve(); };
        document.querySelectorAll('head')[0].appendChild(script);
    });
    return qrLibProm;
};

document.addEventListener("DOMContentLoaded", () => {
    console.log("DOM loaded");

    var cryptoDisplay = document.querySelector('#crypto-display');
    var clearCryptoDisplay = () => {
        while(cryptoDisplay.firstChild) {
            cryptoDisplay.removeChild(cryptoDisplay.firstChild);
        }
    };

    console.log("setting up crypto buttons");
    document.querySelectorAll('.crypto').forEach((el) => {
        var href = el.href;
        el.href="#";
        el.onclick = () => {
            var parts = href.split(":");
            if (parts.length != 2) {
                console.error(el, "href not properly formatted");
                return;
            }
            var currency = parts[0];
            var address = parts[1];

            clearCryptoDisplay();

            var cryptoDisplayQR = document.createElement('div');
            cryptoDisplayQR.id = "crypto-display-qr";

            var cryptoDisplayAddr = document.createElement('div');
            cryptoDisplayAddr.id = "crypto-display-addr";
            cryptoDisplayAddr.innerHTML = '<span>'+currency + " address: " + address + '</span>';

            var cryptoDisplayX = document.createElement('div');
            cryptoDisplayX.id = "crypto-display-x";
            cryptoDisplayX.onclick = clearCryptoDisplay;
            cryptoDisplayX.innerHTML = '<span>X</span>';

            cryptoDisplay.appendChild(cryptoDisplayQR);
            cryptoDisplay.appendChild(cryptoDisplayAddr);
            cryptoDisplay.appendChild(cryptoDisplayX);

            loadQRLib().then(() => {
                new QRCode(cryptoDisplayQR, {
                    text: currency,
                    width: 512,
                    height: 512,
                });
            });
        };
    });
})
