---
layout: post
title: Namecoin, A Replacement For SSL
---

At [cryptic.io][cryptic] we are creating a client-side, in-browser encryption
system where a user can upload their already encrypted content to our storage
system and be 100% confident that their data can never be decrypted by anyone
but them.

On of the main problems with this approach is that the client has to be sure
that the code that's being run in their browser is the correct code; that is,
that they aren't the subject of a man-in-the-middle attack where an attacker is
turning our strong encryption into weak encryption that they could later break.

A component of our current solution is to deliver the site's javascript (and all
other assets, for that matter) using SSL encryption. This protects the files
from tampering in-between leaving our servers and being received by the client.
Unfortunately, SSL isn't 100% foolproof. This post aims to show why SSL is
faulty, and propose a solution.

# SSL

SSL is the mechanism by which web-browsers establish an encrypted connection to
web-servers. The goal of this connection is that only the destination
web-browser and the server know what data is passing between them. Anyone spying
on the connection would only see gibberish. To do this a secret key is first
established between the client and the server, and used to encrypt/decrypt all
data. As long as no-one but those parties knows that key, that data will never
be decrypted by anyone else.

SSL is what's used to establish that secret key on a per-session basis, so that
a key isn't ever re-used and so only the client and the server know it.

## Public-Private Key Cryptography

SSL is based around public-private key cryptography. In a public-private key
system, you have both a public key which is generated from a private key. The
public key can be given to anyone, but the private key must remain hidden. There
are two main uses for these two keys:

* Someone can encrypt a message with your public key, and only you (with the
  private key) can decrypt it.

* You can sign a message with your private key, and anyone with your public key
  can verify that it was you and not someone else who signed it.

These are both extremely useful functions, not just for internet traffic but for
any kind of communication form. Unfortunately, there remains a fundamental flaw.
At some point you must give your public key to the other person in an insecure
way. If an attacker was to intercept your message containing your public key and
swap it for their own, then all future communications could be compromised. That
attacker could create messages the other person would think are from you, and
the other person would encrypt messages meant for you but which would be
decrypt-able by the attacker.

## How does SSL work?

SSL is at its heart a public-private key system, but its aim is to be more
secure against the attack described above.

SSL uses a trust-chain to verify that a public key is the intended one. Your web
browser has a built-in set of public keys, called the root certificates, that it
implicitly trusts. These root certificates are managed by a small number of
companies designated by some agency who decides on these things.

When you receive a server's SSL certificate (its public key) that certificate
will be signed by a root certificate. You can verify that signature since you
have the root certificate's public key built into your browser. If the signature
checks out then you know a certificate authority trusts the public key the site
gave you, which means you can trust it too.

There's a bit (a lot!) more to SSL than this, but this is enough to understand
the fundamental problems with it.

## How SSL doesn't work

SSL has a few glaring problems. One, it implies we trust the companies holding
the root certificates to not be compromised. If some malicious agency was to get
ahold of a root certificate they could listen in on any connection on the
internet by swapping a site's real certificate with one they generate on the
fly. They could trivially steal any data we send on the internet.

The second problem is that it's expensive. Really expensive. If you're running a
business you'll have to shell out about $200 a year to keep your SSL certificate
signed (those signatures have an expiration date attached). Since there's very
few root authorities there's an effective monopoly on signatures, and there's
nothing we can do about it. For 200 bucks I know most people simply say "no
thanks" and go unencrypted. The solution is creating a bigger problem.

# Bitcoins

Time to switch gears, and propose a solution to the above issues: namecoins. I'm
going to first talk about what namecoins are, how they work, and why we need
them. To start with, namecoins are based on bitcoins.

If you haven't yet checked out bitcoins, [I highly encourage you to do
so][bitcoins]. They're awesome, and I think they have a chance of really
changing the way we think of and use money in the future. At the moment they're
still a bit of a novelty in the tech realm, but they're growing in popularity.

The rest of this post assumes you know more or less what bitcoins are, and how
they work.

# Namecoins

Few people actually know about bitcoins. Even fewer know that there's other
crypto-currencies besides bitcoins. Basically, developers of these alternative
currencies (altcoins, in the parlance of our times) took the original bitcoin
source code and modified it to produce a new, separate blockchain from the
original bitcoin one. The altcoins are based on the same idea as bitcoins
(namely, a chain of blocks representing all the transactions ever made), but
have slightly different characterstics.

One of these altcoins is called namecoin. Where other altcoins aim to be digital
currencies, and used as such (like bitcoins), namecoin has a different goal. The
point of namecoin is to create a global, distributed, secure key-value store.
You spend namecoins to claim arbitrary keys (once you've claimed it, you own it
for a set period of time) and to give those keys arbitrary values. Anyone else
with namecoind running can see these values.

## Why use it?

A blockchain based on a digital currency seems like a weird idea at first. I
know when I first read about it I was less than thrilled. How is this better
than a DHT? It's a key-value store, why is there a currency involved?

### DHT

DHT stands for Distributed Hash-Table. I'm not going to go too into how they
work, but suffice it to say that they are essentially a distributed key-value
store. Like namecoin. The difference is in the operation. DHTs operate by
spreading and replicating keys and their values across nodes in a P2P mesh. They
have [lots of issues][dht] as far as security goes, the main one being that it's
fairly easy for an attacker to forge the value for a given key, and very
difficult to stop them from doing so or even to detect that it's happened.

Namecoins don't have this problem. To forge a particular key an attacker would
essentially have to create a new blockchain from a certain point in the existing
chain, and then replicate all the work put into the existing chain into that new
compromised one so that the new one is longer and other clients in the network
will except it. This is extremely non-trivial.

### Why a currency?

To answer why a currency needs to be involved, we need to first look at how
bitcoin/namecoin work. When you take an action (send someone money, set a value
to a key) that action gets broadcast to the network. Nodes on the network
collect these actions into a block, which is just a collection of multiple
actions. Their goal is to find a hash of this new block, combined with some data
from the top-most block in the existing chain, combined with some arbitrary
data, such that the first n characters in the resulting hash are zeros (with n
constantly increasing). When they find one they broadcast it out on the network.
Assuming the block is legitimate they receive some number of coins as
compensation.

That compensation is what keeps a blockchain based currency going.  If there
were no compensation there would be no reason to mine except out of goodwill, so
far fewer people would do it. Since the chain can be compromised if a malicious
group has more computing power than all legitimate miners combined, having few
legitimate miners is a serious problem.

In the case of namecoins, there's even more reason to involve a currency. Since
you have to spend money to make changes to the chain there's a disincentive for
attackers (read: idiots) to spam the chain with frivolous changes to keys.

### Why a *new* currency?

I'll admit, it's a bit annoying to see all these altcoins popping up. I'm sure
many of them have some solid ideas backing them, but it also makes things
confusing for newcomers and dilutes the "market" of cryptocoin users; the more
users a particular chain has, the stronger it is. If we have many chains, all we
have are a bunch of weak chains.

The exception to this gripe, for me, is namecoin. When I was first thinking
about this problem my instinct was to just use the existing bitcoin blockchain
as a key-value storage. However, the maintainers of the bitcoin clients
(who are, in effect, the maintainers of the chain) don't want the bitcoin
blockchain polluted with non-commerce related data. At first I disagreed; it's a
P2P network, no-one gets to say what I can or can't use the chain for!  And
that's true. But things work out better for everyone involved if there's two
chains.

Bitcoin is a currency. Namecoin is a key-value store (with a currency as its
driving force). Those are two completely different use-cases, with two
completely difference usage characteristics. And we don't know yet what those
characteristics are, or if they'll change. If the chain-maintainers have to deal
with a mingled chain we could very well be tying their hands with regards to
what they can or can't change with regards to the behavior of the chain, since
improving performance for one use-case may hurt the performance of the other.
With two separate chains the maintainers of each are free to do what they see
fit to keep their respective chains operating as smoothly as possible.
Additionally, if for some reason bitcoins fall by the wayside, namecoin will
still have a shot at continuing operation since it isn't tied to the former.
Tldr: separation of concerns.

# Namecoin as an alternative to SSL

And now to tie it all together.

There are already a number of proposed formats for standardizing how we store
data on the namecoin chain so that we can start building tools around it. I'm
not hugely concerned with the particulars of those standards, only that we can,
in some way, standardize on attaching a public key (or a fingerprint of one) to
some key on the namecoin blockchain. When you visit a website, the server
would then send both its public key and the namecoin chain key to be checked
against to the browser, and the browser would validate that the public key it
received is the same as the one on the namecoin chain.

The main issue with this is that it requires another round-trip when visiting a
website: One for DNS, and one to check the namecoin chain. And where would this
chain even be hosted?

My proposition is there would exist a number of publicly available servers
hosting a namecoind process that anyone in the world could send requests for
values on the chain. Browsers could then be made with a couple of these
hardwired in. ISPs could also run their own copies at various points in their
network to improve response-rates and decrease load on the globally public
servers. Furthermore, the paranoid could host their own and be absolutely sure
that the data they're receiving is valid.

If the above scheme sounds a lot like what we currently use for DNS, that's
because it is. In fact, one of namecoin's major goals is that it be used as a
replacement for DNS, and most of the talk around it is focused on this subject.
DNS has many of the same problems as SSL, namely single-point-of-failure and
that it's run by a centralized agency that we have to pay arbitrarily high fees
to. By switching our DNS and SSL infrastructure to use namecoin we could kill
two horribly annoying, monopolized, expensive birds with a single stone.

That's it. If we use the namecoin chain as a DNS service we get security almost
for free, along with lots of other benefits. To make this happen we need
cooperation from browser makers, and to standardize on a simple way of
retrieving DNS information from the chain that the browsers can use. The
protocol doesn't need to be very complex, I think HTTP/REST should suffice,
since the meat of the data will be embedded in the JSON value on the namecoin
chain.

If you want to contribute or learn more please check out [namecoin][nmc] and
specifically the [d namespace proposal][dns] for it.

[cryptic]: http://cryptic.io
[bitcoins]: http://vimeo.com/63502573
[dht]: http://www.globule.org/publi/SDST_acmcs2009.pdf
[nsa]: https://www.schneier.com/blog/archives/2013/09/new_nsa_leak_sh.html
[nmc]: http://dot-bit.org/Main_Page
[dns]: http://dot-bit.org/Namespace:Domain_names_v2.0
