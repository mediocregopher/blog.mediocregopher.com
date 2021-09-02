---
layout: page
---

<script async type="module" src="/assets/api.js"></script>

<style>
    #messages {
        max-height: 65vh;
        overflow: auto;
    }

    #messages .message {
        border: 1px solid #AAA;
        border-radius: 10px;
        margin-bottom: 1rem;
        padding: 2rem;
        overflow: auto;
    }

    #messages .message .title {
        font-weight: bold;
        font-size: 120%;
    }

    #messages .message .secondaryTitle {
        font-family: monospace;
        color: #CCC;
    }

    #messages .message p {
        font-family: monospace;
        margin: 1rem 0 0 0;
    }

</style>

<div id="messages"></div>

<span id="fail" style="color: red;"></span>

<script>

const messagesEl = document.getElementById("messages");

function renderMessages(msgs) {

    msgs = [...msgs].reverse();

    messagesEl.innerHTML = '';

    msgs.forEach((msg) => {
        console.log(msg);
        const el = document.createElement("div");
        el.className = "row message"

        const elWithTextContents = (tag, body) => {
            const el = document.createElement(tag);
            el.appendChild(document.createTextNode(body));
            return el;
        };

        const titleEl = document.createElement("div");
        titleEl.className = "title";
        el.appendChild(titleEl);

        const userNameEl = elWithTextContents("span", msg.userID.name);
        titleEl.appendChild(userNameEl);

        const secondaryTitleEl = document.createElement("div");
        secondaryTitleEl.className = "secondaryTitle";
        el.appendChild(secondaryTitleEl);

        const dt = new Date(msg.createdAt*1000);
        const dtStr
            = `${dt.getFullYear()}-${dt.getMonth()+1}-${dt.getDate()}`
            + ` ${dt.getHours()}:${dt.getMinutes()}:${dt.getSeconds()}`;

        const userIDEl = elWithTextContents("span", `userID:${msg.userID.id} @ ${dtStr}`);
        secondaryTitleEl.appendChild(userIDEl);

        const bodyEl = document.createElement("p");

        const bodyParts = msg.body.split("\n");
        for (const i in bodyParts) {
            if (i > 0) bodyEl.appendChild(document.createElement("br"));
            bodyEl.appendChild(document.createTextNode(bodyParts[i]));
        }

        el.appendChild(bodyEl);

        messagesEl.appendChild(el);
    });
}


(async () => {

    const failEl = document.getElementById("fail");

    setErr = (msg) => failEl.innerHTML = `${msg} (please refresh the page to retry)`;

    const api = await import("/assets/api.js");

    try {

        const history = await api.call("/api/chat/global/history");
        renderMessages(history.messages);

    } catch (e) {
        e = `Failed to fetch message history: ${e}`
        setErr(e);
        console.error(e);
        return;
    }

    //const ws = await api.ws("/api/chat/global/listen");

    //while (true) {
    //    const msg = await ws.next();
    //    console.log("got msg", msg);
    //}

})()

</script>
