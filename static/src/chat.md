---
layout: page
---

<script async type="module" src="/assets/api.js"></script>

<style>
    #messages {
        max-height: 65vh;
        overflow: auto;
        padding-right: 2rem;
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

let messagesScrolledToBottom = true;
messagesEl.onscroll = () => {
    const el = messagesEl;
    messagesScrolledToBottom = el.scrollHeight == el.scrollTop + el.clientHeight;
};

function renderMessages(msgs) {

    messagesEl.innerHTML = '';

    msgs.forEach((msg) => {
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

    try {

        const api = await import("/assets/api.js");

        const history = await api.call("/api/chat/global/history");
        const msgs = history.messages;

        // history returns msgs in time descending, but we display them in time
        // ascending.
        msgs.reverse()

        const sinceID = (msgs.length > 0) ?  msgs[msgs.length-1].id : "";

        const ws = await api.ws("/api/chat/global/listen", {
            params: { sinceID },
        });

        while (true) {
            renderMessages(msgs);

            // If the user was previously scrolled to the bottom then keep them
            // there.
            if (messagesScrolledToBottom) {
                messagesEl.scrollTop = messagesEl.scrollHeight;
            }

            const msg = await ws.next();
            msgs.push(msg.message);
            renderMessages(msgs);
        }


    } catch (e) {
        e = `Failed to fetch message history: ${e}`
        setErr(e);
        console.error(e);
        return;
    }

})()

</script>

<style>
#append {
    border: 1px dashed #AAA;
    border-radius: 10px;
    padding: 2rem;
}

#append #appendBody {
    font-family: monospace;
}

#append #appendStatus {
    color: red;
}

</style>

<form id="append">
    <h5>New Message</h5>
    <div class="row">
        <div class="columns four">
            <input class="u-full-width" placeholder="Name" id="appendName" type="text" />
            <input class="u-full-width" placeholder="Secret" id="appendSecret" type="password" />
        </div>
        <div class="columns eight">
            <p>
                Your name is displayed alongside your message.

                Your name+secret is used to generate your userID, which is also
                displayed alongside your message.

                Other users can validate two messages are from the same person
                by comparing the messages' userID.
            </p>
        </div>
    </div>
    <div class="row">
        <div class="columns twelve">
            <textarea
                style="font-family: monospace"
                id="appendBody"
                class="u-full-width"
                placeholder="Well thought out statement goes here..."
                ></textarea>
        </div>
    </div>
    <div class="row">
        <div class="columns four">
            <input class="u-full-width button-primary" id="appendSubmit" type="button" value="Submit" />
        </div>
    </div>
    <span id="appendStatus"></span>
</form>

<script>

const append = document.getElementById("append");
const appendName = document.getElementById("appendName");
const appendSecret = document.getElementById("appendSecret");
const appendBody = document.getElementById("appendBody");
const appendSubmit = document.getElementById("appendSubmit");
const appendStatus = document.getElementById("appendStatus");

appendSubmit.onclick = async () => {

    const appendSubmitOrigValue = appendSubmit.value;

    appendSubmit.disabled = true;
    appendSubmit.className = "";
    appendSubmit.value = "Please hold...";

    appendStatus.innerHTML = '';

    try {

        const api = await import("/assets/api.js");

        await api.call('/api/chat/global/append', {
            body: {
                name: appendName.value,
                password: appendSecret.value,
                body: appendBody.value,
            },
            requiresPow: true,
        });

        appendBody.value = '';

    } catch (e) {

        appendStatus.innerHTML = e;

    } finally {
        appendSubmit.disabled = false;
        appendSubmit.className = "button-primary";
        appendSubmit.value = appendSubmitOrigValue;
    }
};

</script>