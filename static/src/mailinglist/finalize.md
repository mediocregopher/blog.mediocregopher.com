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
        const subToken = urlParams.get('subToken');

        if (!subToken) throw "No subscription token provided";

        const api = await import("/assets/api.js");

        await api.call('/api/mailinglist/finalize', {
            body: { subToken },
        });

        resultSpan.className = "success";
        resultSpan.innerHTML = "Your email subscription has been finalized! Please go on about your day.";

    } catch (e) {
        resultSpan.className = "fail";
        resultSpan.innerHTML = e;
    }

})();

</script>
