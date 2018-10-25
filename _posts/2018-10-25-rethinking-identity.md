---
title: Rethinking Identity
description: >-
    A more useful way of thinking about identity on the internet, and using that
    to build a service which makes our online life better.
---

In my view, the major social media platforms (Facebook, Twitter, Instagram,
etc...) are broken. They worked well at small scales, but billions of people are
now exposed to them, and [Murphy's Law][murphy] has come into effect. The weak
points in the platforms have been found and exploited, to the point where
they're barely usable for interacting with anyone you don't already know in
person.

[murphy]: https://en.wikipedia.org/wiki/Murphy%27s_law

On the other hand, social media, at its core, is a powerful tool that humans
have developed, and it's not one to be thrown away lightly (if it can be thrown
away at all). It's worthwhile to try and fix it. So that's what this post is
about.

A lot of moaning and groaning has already been done on how social media is toxic
for the average person. But the average person isn't doing anything more than
receiving and reacting to their environment. If that environment is toxic, the
person in it becomes so as well. It's certainly possible to filter the toxicity
out, and use a platform to your own benefit, but that takes work on the user's
part. It would be nice to think that people will do more than follow the path of
least resistance, but at scale that's simply not how reality is, and people
shouldn't be expected to do that work.

To identify what has become toxic about the platforms, first we need to identify
what a non-toxic platform would look like.

The ideal definition for social media is to give people a place to socialize
with friends, family, and the rest of the world. Defining "socialize" is tricky,
and probably an exercise only a socially awkward person who doesn't do enough
socializing would undertake. "Expressing one's feelings, knowledge, and
experiences to other people, and receiving theirs in turn" feels like a good
approximation. A platform where true socializing was the only activity would be
ideal.

Here are some trends on our social media which have nothing to do with
socializing: artificially boosted follower numbers on Instagram to obtain
product sponsors, shills in Reddit comments boosting a product or company,
russian trolls on Twitter spreading propaganda, trolls everywhere being dicks
and switching IPs when they get banned, and [that basketball president whose
wife used burner Twitter accounts to trash talk players][president].

[president]: https://www.nytimes.com/2018/06/07/sports/bryan-colangelo-sixers-wife.html

These are all examples of how anonymity can be abused on social media. I want
to say up front that I'm _not_ against anonymity on the internet, and that I
think we can have our cake and eat it too. But we _should_ acknowledge the
direct and indirect problems anonymity causes. We can't trust that anyone on
social media is being honest about who they are and what their motivation is.
This problem extends outside of social media too, to Amazon product reviews (and
basically any other review system), online polls and raffles, multiplayer games,
and surely many other other cases.

## Identity

To fix social media, and other large swaths of the internet, we need to rethink
identity. This process started for me a long time ago, when I watched [this TED
talk][identity], which discusses ways in which we misunderstand identity.
Crucially, David Birch points out that identity is not a name, it's more
fundamental than that.

[identity]: https://www.ted.com/talks/david_birch_identity_without_a_name

In the context of online platforms, where a user creates an account which
identifies them in some way, identity breaks down into 3 distinct problems
which are often conflated:

* Authentication: Is this identity owned by this person?
* Differentiation: Is this identity unique to this person?
* Authorization: Is this identity allowed to do X?

For internet platform developers, authentication has been given the full focus.
Blog posts, articles, guides, and services abound which deal with properly
hashing and checking passwords, two factor authentication, proper account
recovery procedure, etc... While authentication is not a 100% solved problem,
it's had the most work done on it, and the problems which this post deals with
are not affected by it.

The problem which should instead be focused on is differentiation.

## Differentiation

I want to make very clear, once more, that I am _not_ in favor of de-anonymizing
the web, and doing so is not what I'm proposing.

Differentiation is without a doubt the most difficult identity problem to solve.
It's not even clear that it's solvable offline. Take this situation: you are in
a room, and you are told that one person is going to walk in, then leave, then
another person will do the same. These two persons may or may not be the same
person. You're allowed to do anything you like to each person (with their
consent) in order to determine if they are the same person or not.

For the vast, vast majority of cases you can simply look with your eyeballs and
see if they are different people. But this will not work 100% of the time.
Identical twins are an obvious example of two persons looking like one, but a
malicious actor with a disguise might be one person posing as two. Biometrics
like fingerprints, iris scanning, and DNA testing fail for many reasons (the
identical twin case being one). You could attempt to give the first a unique
marking on their skin, but who's to say they don't have a solvent, which can
clean that marking off, waiting right outside the door?

The solutions and refutations can continue on pedantically for some time, but
the point is that there is likely not a 100% solution, and even the 90%
solutions require significant investment. Differentiation is a hard problem,
which most developers don't want to solve. Most are fine with surrogates like
checking that an email or phone number is unique to the platform, but these
aren't enough to stop a dedicated individual or organization.

### Roll Your Own Differentiation

If a platform wants to roll their own solution to the differentiation problem, a
proper solution, it might look something like this:

* Submit an image of your passport, or other government issued ID. This would
  have to be checked against the appropriate government agency to ensure the
  ID is legitimate.

* Submit an image of your face, alongside a written note containing a code given
  by the platform. Software to detect manipulated images would need to be
  employed, as well as reverse image searching to ensure the image isn't being
  reused.

* Once completed, all data needs to be hashed/fingerprinted and then destroyed,
  so sensitive data isn't sitting around on servers, but can still be checked
  against future users signing up for the platform.

* A dedicated support team would be needed to handle edge-cases and mistakes.

None of these is trivial, nor would I trust an up-and-coming platform which is
being bootstrapped out of a basement to implement any of them correctly.
Additionally, going through with this process would be a _giant_ point of
friction for a user creating a new account; they likely would go use a different
platform instead, which didn't have all this nonsense required.

### Differentiation as a Service

This is the crux of this post.

Instead of each platform rolling their own differentiation, what if there was a
service for it. Users would still have to go through the hassle described above,
but only once forever, and on a more trustable site. Then platforms, no matter
what stage of development they're at, could use that service to ensure that
their community of users is free from the problems of fake accounts and trolls.

This is what the service would look like:

* A user would have to, at some point, have gone through the steps above to
  create an account on the differentiation-as-a-service (DaaS) platform. This
  account would have the normal authentication mechanisms that most platforms
  do (password, two-factor, etc...).

* When creating an account on a new platform, the user would login to their DaaS
  account (similar to the common "login with Google/Facebook/Twitter" buttons).

* The DaaS then returns an opaque token, an effectively random string which
  uniquely identifies that user, to the platform. The platform can then check in
  its own user database for any other users using that token, and know if the
  user already has an account. All of this happens without any identifying
  information being passed to the platform.

Similar to how many sites outsource to Cloudflare to handle DDoS protection,
which is better handled en masse by people familiar with the problem, the DaaS
allows for outsourcing the problem of differentiation. Users are more likely to
trust an established DaaS service than a random website they're signing up for.
And signing up for a DaaS is a one-time event, so if enough platforms are using
the DaaS it could become worthwhile for them to do so.

Finally, since the DaaS also handles authentication, a platform could outsource
that aspect of identity management to it as well. This is optional for the
platform, but for smaller platforms which are just starting up it might be
worthwhile to save that development time.

### Traits of a Successful DaaS

It's possible for me to imagine a world where use of DaaS' is common, but
bridging the gap between that world and this one is not as obvious. Still, I
think it's necessary if the internet is to ever evolve passed being, primarily,
a home for trolls. There are a number of traits of an up-and-coming DaaS which
would aid it in being accepted by the internet:

* **Patience**: there is a critical mass of users and platforms using DaaS'
  where it becomes more advantageous for platforms to use the DaaS than not.
  Until then, the DaaS and platforms using it need to take deliberate but small
  steps. For example: making DaaS usage optional for platform users, and giving
  their accounts special marks to indicate they're "authentic" (like Twitter's
  blue checkmark); giving those users' activity higher weight in algorithms;
  allowing others to filter out activity of non-"authentic" users; etc... These
  are all preliminary steps which can be taken which encourage but don't require
  platform users to use a DaaS.

* **User-friendly**: most likely the platforms using a DaaS are what are going
  to be paying the bills. A successful DaaS will need to remember that, no
  matter where the money comes from, if the users aren't happy they'll stop
  using the DaaS, and platforms will be forced to switch to a different one or
  stop using them altogether. User-friendliness means more than a nice
  interface; it means actually caring for the users' interests, taking their
  privacy and security seriously, and in all other aspects being on their side.
  In that same vein, competition is important, and so...

* **No country/government affiliation**: If the DaaS was to be run by a
  government agency it would have no incentive to provide a good user
  experience, since the users aren't paying the bills (they might not even be in
  that country). A DaaS shouldn't be exclusive to any one government or country
  anyway. Perhaps it starts out that way, to get off the ground, but ultimately
  the internet is a global institution, and is healthiest when it's connecting
  individuals _around the world_. A successful DaaS will reach beyond borders
  and try to connect everyone.

Obviously actually starting a DaaS would be a huge undertaking, and would
require proper management and good developers and all that, but such things
apply to most services.

## Authorization

The final aspect of identity management, which I haven't talked about yet, is
authorization. This aspect deals with what a particular identity is allowed to
do. For example, is an identity allowed to claim they have a particular name, or
are from a particular place, or are of a particular age? Other things like
administration and moderation privileges also fall under authorization, but they
are generally defined and managed within a platform.

A DaaS has the potential to help with authorization as well, though with a giant
caveat. If a DaaS were to not fingerprint and destroy the user's data, like
their name and birthday and whatnot, but instead store them, then the following
use-case could also be implemented:

* A platform wants to know if a user is above a certain age, let's say. It asks
  the DaaS for that information.

* The DaaS asks the user, OAuth style, whether the user is ok with giving the
  platform that information.

* If so, the platform is given that information.

This is a tricky situation. It adds a lot of liablity for the user, since their
raw data will be stored with the DaaS, ripe for hacking. It also places a lot of
trust with the DaaS to be responsible with users' data and not go giving it out
willy-nilly to others, and instead to only give out the bare-minimum that the
user allows. Since the user is not the DaaS' direct customer, this might be too
much to ask. Nevertheless, it's a use-case which is worth thinking about.

## Dapps

The idea of decentralized applications, or dapps, has begun to gain traction.
While not mainstream yet, I think they have potential, and it's necessary to
discuss how a DaaS would operate in a world where the internet is no longer
hosted in central datacenters.

Consider an Ethereum-based dapp. If a user were to register one ethereum address
(which are really public keys) with their DaaS account, the following use-case
could be implemented:

* A charity dapp has an ethereum contract, which receives a call from an
  ethereum address asking for money. The dapp wants to ensure every person it
  sends money to hasn't received any that day.

* The DaaS has a separate ethereum contract it manages, where it stores all
  addresses which have been registered to a user. There is no need to keep any
  other user information in the contract.

* The charity dapp's contract calls the DaaS' contract, asking it if the address
  is one of its addresses. If so, and if the charity contract hasn't given to
  that address yet today, it can send money to that address.

There would perhaps need to be some mechanism by which a user could change their
address, which would be complex since that address might be in use by a dapp
already, but it's likely a solvable problem.

A charity dapp is a bit of a silly example; ideally with a charity dapp there'd
also be some mechanism to ensure a person actually _needs_ the money. But
there's other dapp ideas which would become feasible, due to the inability of a
person to impersonate many people, if DaaS use becomes normal.

## Why Did I Write This?

Perhaps you've gotten this far and are asking: "Clearly you've thought about
this a lot, why don't you make this yourself and make some phat stacks of cash
with a startup?" The answer is that this project would need to be started and
run by serious people, who can be dedicated and thorough and responsible. I'm
not sure I'm one of those people; I get distracted easily. But I would like to
see this idea tried, and so I've written this up thinking maybe someone else
would take the reins.

I'm not asking for equity or anything, if you want to try; it's a free idea for
the taking. But if it turns out to be a bazillion dollar Good Idea™, I won't say
no to a donation...
