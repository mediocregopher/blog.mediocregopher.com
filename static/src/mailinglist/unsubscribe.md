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
    const unsubToken = urlParams.get('unsubToken');

    if (!unsubToken) {
        setErr("No unsubscribe token provided");
        return;
    }

    const unsubscribeForm = new FormData();
    unsubscribeForm.append('unsubToken', unsubToken);

    const unsubscribeReq = new Request('/api/mailinglist/unsubscribe', {
        method: 'POST',
        body: unsubscribeForm,
    });

    const res = await fetch(unsubscribeReq)
        .then(response => response.json());

    if (res.error) {
        setErr(res.error);
        return;
    }

    resultSpan.className = "success";
    resultSpan.innerHTML = "You have been unsubscribed! Please go on about your day.";

})();

</script>

