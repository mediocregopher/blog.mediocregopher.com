---
title: >-
    Why Do We Have WiFi Passwords?
description: >-
    A possible UX improvement.
tags: tech
---

It's been longer than I'd like since the last post, and unfortunately I don't
have a ton that I can actually show for it. A lot of time has been spent on
cryptic-net, which is coming along great and even has a proper storage mechanism
now! But it also still has some data specific to our own network baked into the
code, so it can't be shown publicly yet.

-----

Since I don't have much I _can_ show, I thought I'd spend a post diving into a
thought I had the other day: **why do we have wifi passwords?**

The question is a bit facetious. Really what I want to ask is the adjacent
question: why do we use usernames _and_ passwords for wifi networks? The
question doesn't make much sense standing alone though, so it wouldn't do as a
title.

In any case, what I'm proposing is that the vast majority of people don't need a
username/password authentication mechanism to secure their wifi network in a
practical way. Rather, most people could get along just fine with a secret token
mechanism.

In the case of wifi networks, a secret token system might be better named a
secret _name_ mechanism. Using this mechanism a router would not broadcast its
own name to be discovered by the user's device, but rather the user inputs the
name into their device themselves. Existing hidden wifi networks work in this
way already, except they also require a password.

I'm not going to look at this from a technical or cryptographical perspective.
Hidden wifi networks work already, I assume that under the hood this wouldn't be
appreciably different. Instead I'd like to highlight how this change affects the
user experience of joining a wifi network.

The current experience is as follows:

* USER discovers the network name and password through external means.
* USER opens "add new wifi network" page on their device.
* USER finds network name in network list, possibly waiting or scrolling if
  there are many networks.
* USER selects the network name.
* USER inputs password into text box.
* USER is connected to the wifi.

What could this look like if the network name was secret and there was no
password? There'd be no network list, so the whole process is much slimmer:

* USER discovers the secret network name through external means.
* USER opens "add new wifi network" page on their device.
* USER inputs secret name into text box.
* USER is connected to the wifi.

The result is a 33% reduction in number of steps, and a 50% reduction in number
of things the user has to know. The experience is virtually the same across all
other axis.

So the upside of this proposal is clear, a far better UX, but what are the
downsides? Losing a fun avenue of self-expression in the form of wifi names is
probably the most compelling one I've thought of. There's also corporate
environments to consider (as one always must), where it's more practical to
remove users from the network in a targeted way, by revoking accounts, vs
changing the password for everyone anytime a user needs to be excluded.

Corporate offices can keep their usernames and passwords, I guess, and we
should come up with some other radio-based graffiti mechanism in any case. Let's
just get rid of these pointless extra steps!

-----

That's the post. Making this proposal into reality would require a movement far
larger than I care to organize, so we're just going to put this whole thing in
the "fun, pointless yak-shave" bucket and move along. If you happen to know the
architect of the next wifi protocol maybe slip this their way? Or just copy it
and take the credit yourself, that's fine by me.

What's coming next? I'm taking a break from cryptic to catch up on some house
keeping in the self-hosted arena. I've got a brand new password manager I'd like
to try, as well as some motivation to finish getting my own email server
properly set up (it can currently only send mail). At some point I'd like to get
this blog gemini-ified too. Plus there's some services running in their
vestigial docker containers on my server still, that needs to be remedied.

And somewhere in there I have to move too.
