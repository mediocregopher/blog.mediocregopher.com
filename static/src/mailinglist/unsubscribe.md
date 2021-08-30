---
layout: page
title: ""
nofollow: true
---

<script async type="module" src="/assets/api.js"></script>

<style>
#result.success { color: green; }
#result.fail { color: red; }
</style>

<span id="result"></span>

<script>

(async () => {

    const resultSpan = document.getElementById("result");
    
    try {
        const urlParams = new URLSearchParams(window.location.search);
        const unsubToken = urlParams.get('unsubToken');

        if (!unsubToken) throw "No unsubscribe token provided";

        const api = await import("/assets/api.js");

        await api.call('POST', '/api/mailinglist/unsubscribe', {
            body: { unsubToken },
        });

        resultSpan.className = "success";
        resultSpan.innerHTML = "You have been unsubscribed! Please go on about your day.";

    } catch (e) {
        resultSpan.className = "fail";
        resultSpan.innerHTML = e;
    }

})();

</script>

