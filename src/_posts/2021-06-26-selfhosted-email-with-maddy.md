---
title: >-
    Self-Hosted Email With maddy: A Naive First Attempt
description: >-
    How hard could it be?
tags: tech
series: selfhost
---

For a _long_ time now I've wanted to get off gmail and host my own email
domains. I've looked into it a few times, but have been discouraged on multiple
fronts:

* Understanding the protocols underlying email isn't straightforward; it's an
  old system, there's a lot of cruft, lots of auxiliary protocols that are now
  essentially required, and a lot of different services required to tape it all
  together.

* The services which are required are themselves old, and use operational
  patterns that maybe used to make sense but are now pretty freaking cumbersome.
  For example, postfix requires something like 3 different system accounts.

* Deviating from the non-standard route and using something like
  [Mail-in-a-box][miab] involves running docker, which I'm trying to avoid.

So up till now I had let the idea sit, waiting for something better to come
along.

[maddy][maddy] is, I think, something better. According to the homepage
"\[maddy\] replaces Postfix, Dovecot, OpenDKIM, OpenSPF, OpenDMARC and more with
one daemon with uniform configuration and minimal maintenance cost." Sounds
perfect! The homepage is clean and to the point, it's written in go, and the
docs appear to be reasonably well written. And, to top it all off, it's already
been added to [nixpkgs][nixpkgs]!

So in this post (and subsequent posts) I'll be documenting my journey into
getting a maddy server running to see how well it works out.

## Just Do It

I'm almost 100% sure this won't work, but to start with I'm going to simply get
maddy up and running on my home media server as per the tutorial on its site,
and go from there.

First there's some global system configuration I need to perform. Ideally maddy
could be completely packaged up and not pollute the rest of the system at all,
and if I was using NixOS I think that would be possible, but as it is I need to
create a user for maddy and ensure it's able to read the TLS certificates that I
manage via [LetsEncrypt][le].

```bash
sudo useradd -mrU -s /sbin/nologin -d /var/lib/maddy -c "maddy mail server" maddy
sudo setfacl -R -m u:maddy:rX /etc/letsencrypt/{live,archive}
```

The next step is to set up the nix build of the systemd service file. This is a
strategy I've been using recently to nix-ify my services without needing to deal
with nix profiles. The idea is to encode the nix store path to everything
directly into the systemd service file, and install that file normally. In this
case this looks something like:

```
pkgs.writeTextFile {
    name = "mediocregopher-maddy-service";
    text = ''
        [Unit]
        Description=mediocregopher maddy
        Documentation=man:maddy(1)
        Documentation=man:maddy.conf(5)
        Documentation=https://maddy.email
        After=network.target

        [Service]
        Type=notify
        NotifyAccess=main
        Restart=always
        RestartSec=1s

        User=maddy
        Group=maddy

        # cd to state directory to make sure any relative paths
        # in config will be relative to it unless handled specially.
        WorkingDirectory=/mnt/vol1/maddy
        ReadWritePaths=/mnt/vol1/maddy

        # ... lots of directives from
        # https://github.com/foxcpp/maddy/blob/master/dist/systemd/maddy.service
        # that we'll elide here ...

        ExecStart=${pkgs.maddy}/bin/maddy -config ${./maddy.conf}

        ExecReload=/bin/kill -USR1 $MAINPID
        ExecReload=/bin/kill -USR2 $MAINPID

        [Install]
        WantedBy=multi-user.target
    '';
}
```

With the service now testable, it falls on me to actually go through the setup
steps described in the [tutorial][tutorial].

## Following The Tutorial

The first step in the tutorial is setting up of domain names, which I first
perform in cloudflare (where my DNS is hosted) and then reflect into the conf
file. Then I point the `tls file` configuration line at my LetsEncrypt
directory by changing the line to:

```
tls file /etc/letsencrypt/live/$(hostname)/fullchain.pem /etc/letsencrypt/live/$(hostname)/privkey.pem
```


maddy can access these files thanks to the `setfacl` command I performed
earlier.

At this point the server should be effectively configured. However, starting it
via systemd results in this error:

```
failed to load /etc/letsencrypt/live/mx.mydomain.com/fullchain.pem and /etc/letsencrypt/live/mx.mydomain.com/privkey.pem
```

(For my own security I'm not going to be using the actual email domain in this
post, I'll use `mydomain.com` instead.)

This makes sense... I use a wildcard domain with LetsEncrypt, so certs for the
`mx` sub-domain specifically won't exist. I need to figure out how to tell maddy
to use the wildcard, or actually create a separate certificate for the `mx`
sub-domain. I'd rather the former, obviously, as it's far less work.

Luckily, making it use the wildcard isn't too hard, all that is needed is to
change the `tls file` line to:

```
tls file /etc/letsencrypt/live/$(primary_domain)/fullchain.pem /etc/letsencrypt/live/$(primary_domain)/privkey.pem
```

This works because my `primary_domain` domain is set to the top-level
(`mydomain.com`), which is what the wildcard cert is issued for.

At this point maddy is up and running, but there's still a slight problem. maddy
appears to be placing all of its state files in `/var/lib/maddy`, even though
I'd like to place them in `/mnt/vol1/maddy`. I had set the `WorkingDirectory` in
the systemd service file to this, but apparently that's not enough. After
digging through the codebase I discover an undocumented directive which can be
added to the conf file:

```
state_dir /mnt/vol1/maddy
```

Kind of annoying, but at least it works.

The next step is to fiddle with DNS records some more. I add the SPF, DMARC and
DKIM records to cloudflare as described by the tutorial (what do these do? I
have no fuckin clue).

I also need to set up MTA-STS (again, not really knowing what that is). The
tutorial says I need to make a file with certain contents available at the URL
`https://mta-sts.mydomain.com/.well-known/mta-sts.txt`. I love it when protocol
has to give up and resort to another one in order to keep itself afloat, it
really inspires confidence.

Anyway, I set that subdomain up in cloudflare, and add the following to my nginx
configuration:

```
server {
    listen      80;
    server_name mta-sts.mydomain.com;
    include     include/public_whitelist.conf;

    location / {
        return 404;
    }

    location /.well-known/mta-sts.txt {

        # Check out openresty if you want to get super useful nginx plugins, like
        # the echo module, out-of-the-box.
        echo 'mode: enforce';
        echo 'max_age: 604800';
        echo 'mx: mx.mydomain.com';
    }
}
```

(Note: my `public_whitelist.conf` only allows cloudflare IPs to access this
sub-domain, which is something I do for all sub-domains which I can put through
cloudflare.)

Finally, I need to create some actual credentials in maddy with which to send my
email. I do this via the `maddyctl` command-line utility:

```
> sudo maddyctl --config maddy.conf creds create 'me@mydomain.com'
Enter password for new user:
> sudo maddyctl --config maddy.conf imap-acct create 'me@mydomain.com'
```

## Send It!

At this point I'm ready to actually test the email sending. I'm going to use
[S-nail][snail] to do so, and after reading through the docs there I put the
following in my `~/.mailrc`:

```
set v15-compat
set mta=smtp://me%40mydomain.com:password@localhost:587 smtp-use-starttls
```

And attempt the following `mailx` command to send an email from my new mail
server:

```
> echo 'Hello! This is a cool email' | mailx -s 'Subject' -r 'Me <me@mydomain.com>' 'test.email@gmail.com'
reproducible_build: TLS certificate does not match: localhost:587
/home/mediocregopher/dead.letter 10/313
reproducible_build: ... message not sent
```

Damn. TLS is failing because I'm connecting over `localhost`, but maddy is
serving the TLS certs for `mydomain.com`. Since I haven't gone through the steps
of exposing maddy publicly yet (which would require port forwarding in my
router, as well as opening a port in iptables) I can't properly test this with
TLS not being required. _It's very important that I remember to re-require TLS
before putting anything public._

In the meantime I remove the `smtp-use-starttls` entry from my `~/.mailrc`, and
retry the `mailx` command. This time I get a different error:

```
reproducible_build: SMTP server: 523 5.7.10 TLS is required
```

It turns out there's a further configuration directive I need to add, this time
in `maddy.conf`. Within my `submission` configuration block I add the following
line:

```
insecure_auth true
```

This allows plaintext auth over non-TLS connections. Kind of sketchy, but again
I'll undo this before putting anything public.

Finally, I try the `mailx` command one more time, and it successfully returns!

Unfortunately, no email is ever received in my gmail :( I check the maddy logs
and see what I feared most all along:

```
Jun 29 08:44:58 maddy[127396]: remote: cannot use MX        {"domain":"gmail.com","io_op":"dial","msg_id":"5c23d76a-60db30e7","reason":"dial tcp 142.250.152.26:25: connect: connection timed out","remote_addr":"142.250.152.
26:25","remote_server":"alt1.gmail-smtp-in.l.google.com.","smtp_code":450,"smtp_enchcode":"4.4.2","smtp_msg":"Network I/O error"}
```

My ISP is blocking outbound connections on port 25. This is classic email
bullshit; ISPs essentially can't allow outbound SMTP connections, as email is so
easily abusable it would drastically increase the amount of spam being sent from
their networks.

## Lessons Learned

The next attempt will involve an external VPS which allows SMTP, and a lot more
interesting configuration. But for now I'm forced to turn off maddy and let this
dream sit for a little while longer.

[miab]: https://mailinabox.email/
[maddy]: https://maddy.email
[nixpkgs]: https://search.nixos.org/packages?channel=21.05&from=0&size=50&sort=relevance&query=maddy
[tutorial]: https://maddy.email/tutorials/setting-up/
[le]: https://letsencrypt.org/
[snail]: https://wiki.archlinux.org/title/S-nail
