---
title: >-
    How to Secure a Webapp
description: >-
    Get ready to jump through some hoops.
---

In this post I will be documenting all security hoops that one must jump through
in order to consider their webapp secure. This list should not be considered
comprehensive, as I might have forgotten something or some new hoop might have
appeared since writing.

For the context of this post a "webapp" will be considered to be an HTML/CSS/JS
website, loaded in a browser, with which users create and access accounts using
some set of credentials (probably username and password). In other words, most
popular websites today. This post will only cover those concerns which apply to
_all_ webapps of this nature, and so won't dive into any which might be incurred
by using one particular technology or another.

Some of these hoops might seem redundant or optional. That may be the case. But
if you are building a website and are beholden to passing some third-party
security audit for any reason you'll likely find yourself being forced to
implement most, if not all, of these measures anyway.

So without further ado, let's get started!

## HTTPS

At this point you have to use HTTPS, there's not excuse for not doing so. All
attempts to hit an HTTP endpoint should redirect to the equivalent HTTPS
endpoint, and you should be using [HSTS][hsts] to ensure that a browser is never
tricked into falling back to HTTP via some compromised DNS server.

[hsts]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Strict-Transport-Security

## Cookies

Cookies are an old web technology, and have always been essentially broken. Each
cookie can have certain flags set on it which change their behavior, and some of
these flags are required at this point.

### Secure

If you're storing anything sensitive in a cookie (spoiler alert: you will be)
then you need to have the Secure flag set on it. This prevents the cookie from
being sent in a non-HTTPS request.

### HTTPOnly

The HTTPOnly flag protects a cookie from XSS attacks by preventing it from being
accessible from javascript. Any cookie which is storing sensitive information
_must_ have this flag set. In the **Authentication** section we will cover the
storage of session tokens, but the TLDR is that they have to be stored in an
HTTPOnly cookie.

Practically, this means that your sessions architecture _must_ account for the
fact that the webapp itself will not have direct access to its persistent
session token(s), and therefore must have some other way of knowing that it's
logged in (e.g. a secondary, non-HTTPOnly cookie which contains no secrets but
only signals that the browser is logged in).

### SameSite

The SameSite attribute can be set to `Strict`, `Lax`, or `None`. `Lax` is the
default in modern browsers and is sufficient for most security concerns, but if
you can go with `Strict` that would be better. The downside of `Strict` is that
cookies won't be sent on initial page-load of a site.

In any case, even though `Lax` is the default you should still set this
attribute manually (or your auditor might get to add another bullet point to
their report).

## Authentication

Authentication is obviously one of the juiciest targets for an attacker. It's
one thing to be able to trick a user into performing this or that action, but if
one can just log in _as_ the user then they essentially have free-reign over all
their information.

### Password History

Most websites use a username/password system as the first step of login. This
is.... fine. We've accepted it, at any rate. But there's a couple of hoops which
must be jumped through as a result of it, and the first is password history.

I hope it goes without saying that one should be using a hashing algorithm like
bcrypt to store user passwords. But what is often not said is that, for each
user, you need to store the hashes of their last N passwords (where N is
something like 8). This way if they attempt to re-use an old password they are
not able to do so. The users must be protected from themselves, afterall.

### Credential Stuffing/Account Enumeration

A credential stuffing attack is one where credentials are stolen from one
website and then attempted to be used on another, in the hope that users have
re-used their username/password across multiple sites. When they occur it'll
often look like a botnet spamming the authentication endpoint with tons of
different credentials.

Account enumeration is a similar attack: it's where an attacker finds a way to
get the webapp to tell them whether or not an account email/username exists in
the system, without needing to have the right password. This is often done by
analyzing the error messages returned from login or a similar endpoint (e.g.
"Sorry this username is taken"). They then run through all possible values for
that endpoint to try and enumerate which users actually exist in the system.

Account enumeration is tricky because often those errors are extremely helpful,
and we'd _like_ to keep them if we can.

I've bucketed both of these attacks in the same section because they have a
similar solution: proof-of-work. The idea is that, for each request to some
sensitive endpoint, the client must send some proof that they've done an
intensive CPU computation.

Compared to IP-based rate-limiting, PoW is much more effective against botnets
(which have a limitless set of IPs from which to spam you), while also being
much less intrusive on your real users than a captcha.

PoW stymies botnets because they are generally being hosted by low-power,
compromised machines. In addition the systems that run these botnets are pretty
shallow in capability, because it's more lucrative to rent the botnet out then
to actually use it yourself, so it's rare for a botnet operator to go to the
trouble of implementing your PoW algorithm in the first place.

So stick a PoW requirement on any login or account creation endpoint, or any
other endpoint which might be used to enumerate accounts in the system. You can
even make the PoW difficulty rise in relation to number of recent attempts on
these endpoints, if you're feeling spry.

### MFA

All the PoW checks in the world won't help your poor user who isn't using a
different username/password for each website, and who got unlucky enough to have
those credentials leaked in a hack of a completely separate site than your own.
They also won't help your user if they _are_ using different username/passwords
for everything, but their machine gets straight up stolen IRL and the attacker
gets access to their credential storage.

What _will_ help them in these cases, however, is if your site supports
multi-factor authentication, such as [TOTP][totp]. If it does then your user
will have a further line of defense in the form of another password which
changes every 30 seconds, and which can only be accessed from a secondary device
(like their phone). If your site claims to care about the security of your
user's account then MFA is an absolute requirement.

It should be noted, however, that not all MFA is created equal. A TOTP system
is great, but a one-time code being sent over SMS or email is totally different
and not nearly as great. SMS is vulnerable to [SIM jacking][sim], which can be
easily used in a targeted attack against one of your users. One-time codes over
email are pointless for MFA, as most people have their email logged in on their
machine all the time, so if someone steals your user's machine they're still
screwed.

In summary: MFA is essentially required, _especially_ if the user's account is
linked to anything valuable, and must be done with real MFA systems like TOTP,
not SMS or email.

[totp]: https://www.twilio.com/docs/glossary/totp
[sim]: https://www.vice.com/en/article/3kx4ej/sim-jacking-mobile-phone-fraud

### Login Notifications

Whenever a user successfully logs into their account you should send them email
(or some other notification) letting them know it happened. This way if it
wasn't actually them who did so, but an attacker, they can perhaps act quickly
to lock down their account and prevent any further harm. The login notification
email should have some kind of link in it which can be used to immediately lock
the account.

### Token Storage

Once your user has logged into your webapp, it's up to you, the developer, to
store their session token(s) somewhere. The question is... where? Well this
one's easy, because there's only one right answer: HTTPOnly cookies (as alluded
to earlier).

When storing session tokens you want to guard against XSS attacks which might
grab the tokens and send them to an attacker, allowing that attacker to hijack
the session and pose as the user. This means the following are not suitable
places to store the tokens:

* Local storage.
* `window`, or anything which can be accessed via `window`.
* Non-HTTPOnly cookies.

Any of these are trivial to find by a script running in the browser. If a
session token is ephemeral then it may be stored in a "normal" javascript
variable somewhere _as long as_ that variable isn't accessible from a global
context. But for any tokens which need to be persisted across browser restarts
an HTTPOnly cookie is your only option.

## Cross-Site

Speaking of XSS attacks, we have some more mitigation coming up...

### CSP

Setting a [CSP][csp] for your website is key to preventing XSS. A CSP allows you
to more tightly control the allowed origins of the various entities on your site
(be they scripts, styles, images, etc...). If an entity of unexpected origin
shows up it is disallowed.

Be sure to avoid any usages of the policies labeled "unsafe" (go figure),
otherwise the CSP is rendered somewhat pointless. Also, when using hostname
based allowlisting try to be as narrow as you can in your allowlist, and
especially only include https hosts. If you can you should opt for the `nonce`
or `sha` policies.

[csp]: https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP

### SVG

A small but important note: if you're website allows users to upload images,
then be _very_ careful about allowing users to upload SVGs. SVGs are actually
XML documents, and even worse than that they allow `<script>` tags within them!
So you need to be very careful about allowing SVGs to be uploaded. If you can
get away with it, it's better to disallow their use at all.

## CSRF

The web was designed in a time when cross-site requests were a considered
feature. This has proven to be a massive mistake. We have two cross-site request
prevention techniques in this list. The first is CSRF.

CSRF protection will cover you from a variety of attacks, mostly of the kind
where an attacker embeds a `<form>` on their own webpage, with the form set up
to POST to _your_ website in some way. When a user of your website lands on the
attacker's page and triggers the POST, the POST will be sent along with whatever
cookies the user has stored in their browser for _your_ site!

The attacker could, potentially, trick a user into submitting a password-reset
request using a known value, or withdrawing all their money into the attacker's
bank account, or anything else the user might be able to do on their own.

The idea with CSRF is that any HTTP request made against an API should have an
unguessable token as a required parameter, called the CSRF token. The CSRF token
should be given to your webapp in a way where only your webapp could know it.
There are many ways to accomplish this, including a cookie, server-side embedded
value, etc... OWASP has put together an [entire cheatsheet full of CSRF
methods][csrf] which is well worth checking out.

[csrf]: https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html

## CORS

CORS is the other half of cross-site protection. With CSRF in place it's
somewhat redundant, but it's good to have multiple layers of protection in place
(in case you fuck up one of them by accident).

The key thing one must do for CORS protection is to set the
`Access-Control-Allow-Origin` to the origin a request is being sent from _only
if you trust that origin_. If you stick a wildcard in that header then you're
not doing anything.

## Random Headers

The rest of this is random HTTP headers which must be set in various contexts to
protect your users.

### Permissions Policy

The [Permissions-Policy][pp] header is fairly new and not fully standardized
yet, but there is support for it so it's worth using. It allows you to specify
exactly which browser features you expect your webapp to need, and therefore
prevent an attacker from taking advantage of some other feature that you were
never going to use anyway.

[pp]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Feature-Policy

### X-Content-Type-Options

It's important to set `X-Content-Type-Options: nosniff` on virtually all
HTTP responses, in order to (theoretically) prevent a browser from inferring the
MIME of the returned content.

### X-Frame-Options

Set `X-Frame-Options: deny` to prevent your webapp from being rendered in a
frame or iframe on someone else's site, which might then be used to trick one of
your users into doing something stupid.

### X-XSS-Protection

Set `X-XSS-Protection: 1; mode=block` to give older browsers which lack CSP
support some extra defense against XSS attacks. It's not super clear to me what
exactly this actually does, but it's easy enough to set.

### Referrer-Policy

Set the `Referrer-Policy` to inform your users' browsers to not send the
`Referer` header to third-party sites when your users navigate away from your
site. You don't want other websites to be able to see _yours_ in their logs, as
they could then correlate which users of theirs have accounts with you (and so
potentially have some easy targets).

### Cache-Control/Pragma

For all requests which return sensitive information (i.e. any authenticated
requests) it's important to set `Cache-Control: no-store` and `Pragma: no-cache`
on the response. This prevents some middle server or the browser from caching
the response, and potentially returning it later to someone else using your site
from the same location.

## That's It

It's probably not it, actually, these are just what I could think of off the top
of my head. Please email me if I missed any.

If you, like me, find yourself asking "how is anyone supposed to have figured
this out?" then you should A) thank me for writing it all down for you and B)
realize that at least 50% of this list has nothing to do with the web, really,
and everything to do with covering up holes that backwards compatibility has
left open. We can cover these holes, we just need everyone to agree on the path
to doing so, and to allow ourselves to leave some ancient users behind.
