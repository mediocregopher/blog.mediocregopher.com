---
layout: page
title: ""
nofollow: true
---

<style>
#result.success { color: green; }
#result.fail { color: red; }
</style>

<span id="result"></span>

<script>

(async () => {
    const resultSpan = document.getElementById("result");
    
    function setErr(errStr) {
        resultSpan.className = "fail";
        resultSpan.innerHTML = errStr;
    }

    const urlParams = new URLSearchParams(window.location.search);
    const subToken = urlParams.get('subToken');

    if (!subToken) {
        setErr("No subscription token provided");
        return;
    }

    const finalizeForm = new FormData();
    finalizeForm.append('subToken', subToken);

    const finalizeReq = new Request('/api/mailinglist/finalize', {
        method: 'POST',
        body: finalizeForm,
    });

    const res = await fetch(finalizeReq)
        .then(response => response.json());

    if (res.error) {
        setErr(res.error);
        return;
    }

    resultSpan.className = "success";
    resultSpan.innerHTML = "Your email subscription has been finalized! Please go on about your day.";

})();

</script>
