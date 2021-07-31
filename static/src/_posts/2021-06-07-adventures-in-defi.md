---
title: >-
    Adventures In DeFi
description: >-
    There and Back Again, a Yield Farmer's Tale.
---

It's difficult to be remotely interested in crypto and avoid the world of
decentralized finance (DeFi). Somewhere between the explosion of new projects,
implausible APY percents, complex tokens schemes, new phrases like "yield
farming" and "impermanent loss", rug pulls, hacks, and astronomical ethereum
fees, you simply _must_ have heard of it, even in passing.

In late November of 2020 I decided to jump in and see what would happen. I read
everything I could find, got as educated as I could, did some (but probably not
enough) math, and got to work. Almost immediately afterwards a giant bull
market hit, fees on ethereum shot up to the moon, and my little yield farming
DeFi ship was effectively out to sea.

For the past 200 days I haven't been able to tweak or withdraw any of the DeFi
positions I made, for fear of incurring so many ethereum fees that any gains I
made would be essentially wiped out. But the bull market is finally at a rest,
fees are down, and I'm interested in what the results of my involuntary
long-term experiment were. Before getting to the results though, let's start at
the beginning. I'm going to walk you through all the steps I took, as well as my
decision making process (as flawed as it surely was) and risk assessments.

## Step 1: The Base Positions

My first step was to set aside some ETH and BTC for this experiment. I was (and
remain) confident that these assets would acrue in value, and so wanted to hold
onto them for a long period of time. But while holding onto those assets, why
not make a little interest on them by putting them to use? That's where DeFi
comes in.

I started with 2.04 ETH and 0.04 BTC. The ETH existed as normal ETH on the
ethereum blockchain, while the 0.04 BTC I had to first convert to [renBTC][ren].

renBTC is an ethereum token whose value is pinned to the value of BTC. This is
accomplished via a decentralized locking mechanism, wherein real BTC is
transferred to a decentralized network of ren nodes, and they lock it such that
no individual node has access to the wallet holding the BTC. At the same time
that the BTC is locked, the nodes print and transfer a corresponding amount of
renBTC to a wallet specified in the BTC transaction. It's a very interesting
project, though the exact locking mechanism used was closed-source at the time I
used it, which concerned me somewhat.

[ren]: https://renproject.io/

### Step 1.5: Collateralization

In Step 2 I deposit my assets into liquidity pools. For my renBTC this was no
problem, but for my ETH it wasn't so simple. I'll explain what a liquidity pool
is in the next section, but for now all that needs to be known is that there are
no worthwhile liquidity pools between ETH and anything ostensibly pinned to ETH
(e.g.  WETH). So I needed to first convert my ETH into an asset for which there
are worthwhile liquidity pools, while also not losing my ETH position.

Enter [MakerDAO][makerdao]. MakerDAO runs a decentralized collateralization app,
wheren a user deposits assets into a contract and is granted an amount of DAI
tokens relative to the value of the deposited assets. The value of DAI tokens
are carefully managed via the variable fee structure of the MakerDAO app, such
that 1 DAI is, generally, equal to 1 USD. If the value of the collateralized
assets drops below a certain threshold the position is liquidated, meaning the
user keeps the DAI and MakerDAO keeps the assets. It's not dissimilar to taking
a loan out, using one's house as collateral, except that the collateral is ETH
and not a house.

MakerDAO allows you to choose, within some bounds, how much DAI you withdraw on
your deposited collateral. The more DAI you withdraw, the higher your
liquidation threshold, and if your assets fall in value and hit that threshold
you lose them, so a higher threshold entails more risk. In this way the user has
some say over how risky of a position they want to take out.

In my case I took out a total of 500 DAI on my 2.04 ETH. Even at the time this
was somewhat conservative, but now that the price of ETH has 5x'd it's almost
comical. In any case, I now had 500 DAI to work with, and could move on to the
next step.

[makerdao]: https://makerdao.com/

## Step 2: Liquidity Pools

My assets were ready to get put to work, and the work they got put to was in
liquidity pools (LPs). The function of an LP is to facilitate the exchange of
one asset for another between users. They play the same role as a centralized
exchange like Kraken or Binance, but are able to operate on decentralized chains
by using a different exchange mechanism.

I won't go into the details of how LPs work here, as it's not super pertinent.
There's great explainers, like [this one][lp], that are easy to find. Suffice it
to say that each LP operates on a set of assets that it allows users to convert
between, and LP providers can deposit one or more of those assets into the pool
in order to earn fees on each conversion.

When you deposit an asset into an LP you receive back a corresponding amount of
tokens representing your position in that LP. Each LP has its own token, and
each token represents a share of of the pool that the provider owns. The value
of each token goes up over time as fees are collected, and so acts as the
mechanism by which the provider ultimately collects their yield.

In addition to the yield one gets from users making conversions via the LP, LP
providers are often also further incentivized by being granted governance tokens
in the LPs they provide for, which they can then turn around and sell directly
or hold onto as an investment. These are usually granted via a staking
mechanism, where the LP provider stakes (or "locks") their LP tokens into the
platform, and is able to withdraw the incentive token based on how long and how
much they've staked.

Some LP projects, such as [Sushi][sushi], have gone further and completely
gamified the whole experience, and are the cause of the multi thousand percent
APYs that DeFi has become somewhat famous for. These projects are flashy, but I
couldn't find myself placing any trust in them.

There is a risk in being an LP provider, and it's called ["impermanent
loss"][il]. This is another area where it's not worth going into super detail,
so I'll just say that impermanent loss occurs when the relative value of the
assets in the pool diverges significantly. For example, if you are a provider in
a BTC/USDC pool, and the value of BTC relative to USD either tanks or
skyrockets, you will have ended up losing money.

I wanted to avoid impermanent loss, and so focused on pools where the assets
have little chance of diverging. These would be pools where the assets are
ostensibly pinned in value, for example a pool between DAI and USDC, or between
renBTC and WBTC. These are called stable pools. By choosing such pools my only
risk was in one of the pooled assets suddenly losing all of its value due to a
flaw in its mechanism, for example if MakerDAO's smart contract were to be
hacked. Unfortunately, stable pools don't have as great yields as their volatile
counterparts, but given that this was all gravy on top of the appreciation of
the underlying ETH and BTC I didn't mind this as much.

I chose the [Curve][curve] project as my LP project of choice. Curve focuses
mainly on stable pools, and provides decent yield percents in that area while
also being a relatively trusted and actively developed project.

I made the following deposits into Curve:

* 200 DAI into the [Y Pool][ypool], receiving back 188 LP tokens.
* 300 DAI into the [USDN Pool][usdnpool], receiving back 299 LP tokens.
* 0.04 renBTC into the [tBTC Pool][tbtcpool], receiving back 0.039 LP tokens.

[lp]: https://finematics.com/liquidity-pools-explained/
[il]: https://finematics.com/impermanent-loss-explained/
[sushi]: https://www.sushi.com/
[curve]: https://curve.fi
[ypool]: https://curve.fi/iearn
[usdnpool]: https://curve.fi/usdn
[tbtcpool]: https://curve.fi/tbtc

## Step 3: Yield Farming

At this point I could have taken the next step of staking my LP tokens into the
Curve platform, and periodically going in and reaping the incentive tokens that
doing so would earn me. I could then sell these tokens and re-invest the profits
back into the LP, and then stake the resulting LP tokens back into Curve,
resulting in a higher yield the next time I reap the incentives, ad neaseaum
forever.

This is a fine strategy, but it has two major drawbacks:

* I don't have the time, nor the patience, to implement it.
* ETH transaction fees would make it completely impractical.

Luckily, yield farming platforms exist. Rather than staking your LP tokens
yourself, you instead deposit them into a yield farming platform. The platform
aggregates everyone's LP tokens, stakes them, and automatically collects and
re-invests incentives in large batches. By using a yield farming platform,
small, humble yield farmers like myself can pool our resources together to take
advantage of scale we wouldn't normally have.

Of course, yield farming adds yet another gamification layer to the whole
system, and complicates everything. You'll see what I mean in a moment.

The yield farming platform I chose was [Harvest][harvest]. Overall
Harvest had the best advertised APYs (though those can obviously change on a
dime), a large number of farmed pools that gets updated regularly, as well as a
simple interface that I could sort of understand. The project is a _bit_ of a
mess, and there's probably better options now, but it was what I had at the
time.

For each of the 3 kinds of LP tokens I had collected in Step 2 I deposited them
into the corresponding farming pool on Harvest. As with the LPs, for each
farming pool you deposit into you receive back a corresponding amount of farming
pool tokens which you can then stake back into Harvest. Based on how much you
stake into Harvest you can collect a certain amount of FARM tokens periodically,
which you can then sell, yada yada yada. It's farming all the way down. I didn't
bother much with this.

[harvest]: https://harvest.finance

## Step 4: Wait

At this point the market picked up, ethereum transactions shot up from 20 to 200
gwei, and I was no longer able to play with my DeFi money without incurring huge
losses. So I mostly forgot about it, and only now am coming back to it to see
the damage.

## Step 5: Reap What I've Sown

It's 200 days later, fees are down again, and enough time has passed that I
could plausibly evaluate my strategy, I've gone through the trouble of undoing
all my positions in order to arrive back at my base assets, ETC and BTC. While
it's tempting to just keep the DeFi ship floating on, I think I need to redo it
in a way that I won't be paralyzed during the next market turn, and I'd like to
evaluate other chains besides ethereum.

First, I've unrolled my Harvest positions, collecting the original LP tokens
back plus whatever yield the farming was able to generate. The results of that
step are:

* 194 Y Pool tokens (originally 188).
* 336 USDN Pool tokens (originally 299).
* 0.0405 tBTC Pool tokens (originally 0.039).

Second, I've burned those LP tokens to collect back the original assets from the
LPs, resulting in:

* 215.83 DAI from the Y Pool (originally 200).
* 346.45 DAI from the USDN Pool (originally 300).
* 0.0405 renBTC from the tBTC Pool (originally 0.04).

For a total DAI of 562.28.

Finally, I've re-deposited the DAI back into MakerDAO to reclaim my original
ETH. I had originally withdrawn 500 DAI, but due to interest I now owed 511
DAI. So after reclaiming my full 2.04 ETH I have ~51 DAI leftover.

## Insane Profits

Calculating actual APY for the BTC investment is straightforward: it came out to
about 4.20% APY. Not too bad, considering the position is fairly immune to price
movements.

Calculating for ETH is a bit trickier, since in the end I ended up with the same
ETH as I started with (2.04) plus 51 DAI. If I were to purchase ETH with that
DAI now, it would get me ~0.02 further ETH. Not a whole heck of a lot. And that
doesn't even account for ethereum fees! I made 22 ethereum transactions
throughout this whole process, resulting in ~0.098 ETH spent on transaction
fees.

So in the end, I lost 0.078 ETH, but gained 0.0005 BTC. If I were to
convert the BTC gain to ETH now it would give me a net total profit of:

**-0.071 ETH**

A net loss, how fun!

## Conclusions

There were a lot of takeaways from this experiment:

* ETH fees will get ya, even in the good times. I would need to be working with
  at least an order of magnitude higher base position in order for this to work
  out in my favor.

* I should have put all my DAI in the Curve USDN pool, and not bothered with the
  Y pool. It had almost double the percent return in the end.

* Borrowing DAI on my ETH was fun, but it really cuts down on how much of my ETH
  value I'm able to take advantage of. My BTC was able to be fully invested,
  whereas at most half of my ETH value was.

* If I have a large USD position I want to sit on, the USDN pool on its own is
  not the worst place to park it. The APY on it was about 30%!

I _will_ be trying this again, albeit with a bigger budget and more knowledge. I
want to check out other chains besides ethereum, so as to avoid the fees, as
well as other yield mechanisms besides LPs, and other yield farming platforms
besides Harvest.

Until then!
