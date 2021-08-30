---
layout: page
title: "Follow the Blog"
nofollow: true
---

<script async type="module" src="/assets/api.js"></script>

Here's your options for receiving updates about new blog posts:

## Option 1: Email

Email is by far my preferred option for notifying followers of new posts.

The entire email list system for this blog, from storing subscriber email
addresses to the email server which sends the notifications out, has been
designed from scratch and is completely self-hosted in my living room.

I solemnly swear that:

* You will never receive an email from this blog except to notify of a new post.

* Your email will never be provided or sold to anyone else for any reason.

With all that said, if you'd like to receive an email everytime a new blog post
is published then input your email below and smash that subscribe button!

<style>

#emailStatus.success {
    color: green;
}

#emailStatus.fail {
    color: red;
}

</style>

<input type="email" placeholder="name@host.com" id="emailAddress" />
<input class="button-primary" type="submit" value="Subscribe" id="emailSubscribe" />
<span id="emailStatus"></span>

<script>

const emailAddress = document.getElementById("emailAddress");
const emailSubscribe = document.getElementById("emailSubscribe");
const emailSubscribeOrigValue = emailSubscribe.value;
const emailStatus = document.getElementById("emailStatus");

emailSubscribe.onclick = async () => {

    const api = await import("/assets/api.js");

    emailSubscribe.disabled = true;
    emailSubscribe.className = "";
    emailSubscribe.value = "Please hold...";
    emailStatus.innerHTML = '';

    try {

        if (!window.isSecureContext) {
            throw "The browser environment is not secure.";
        }

        await api.call('POST', '/api/mailinglist/subscribe', {
            body: { email: emailAddress.value },
            requiresPow: true,
        });

        emailStatus.className = "success";
        emailStatus.innerHTML = "Verification email sent (check your spam folder)";

    } catch (e) {
        emailStatus.className = "fail";
        emailStatus.innerHTML = e;

    } finally {
        emailSubscribe.disabled = false;
        emailSubscribe.className = "button-primary";
        emailSubscribe.value = emailSubscribeOrigValue;
    }

};

</script>

## Option 2: RSS

RSS is the classic way to follow any blog. It comes from a time before
aggregators like reddit and twitter stole the show, when people felt capable to
manage their own content feeds. We should use it again.

To follow over RSS give any RSS reader the following URL...

<a href="{{site.url}}/feed.xml">{{site.url}}/feed.xml</a>

...and posts from this blog will show up in your RSS feed as soon as they are
published. There are literally thousands of RSS readers out there. Here's some
recommendations:

* [Google Chrome Browser Extension](https://chrome.google.com/webstore/detail/rss-feed-reader/pnjaodmkngahhkoihejjehlcdlnohgmp)

* [spaRSS](https://f-droid.org/en/packages/net.etuldan.sparss.floss/) is my
  preferred android RSS reader, but you'll need to install
  [f-droid](https://f-droid.org/) on your device to use it (a good thing to do
  anyway, imo).

* [NetNewsWire](https://ranchero.com/netnewswire/) is a good reader for
  iPhone/iPad/Mac devices, so I'm told. Their homepage description makes a much
  better sales pitch for RSS than I ever could.

## Option 3: Twitter

New posts are automatically published to [my Twitter](https://twitter.com/{{
site.twitter_username }}). Simply follow me there and pray the algorithm smiles
upon my tweets enough to show them to you! :pray: :pray: :pray:

