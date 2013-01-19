# Erlang, tcp sockets, and active true

If you don't know erlang then [you're missing out](http://learnyousomeerlang.com/content).
If you do know erlang, you've probably at some point done something with tcp sockets. Erlang's
highly concurrent model of execution lends itself well to server programs where a high number
of active connections is desired. Each thread can autonomously handle its single client,
greatly simplifying the logic of the whole application while still retaining
[great performance characteristics](http://www.metabrew.com/article/a-million-user-comet-application-with-mochiweb-part-1).

# Background

For an erlang thread which owns a single socket there are three different ways to receive data
off of that socket. These all revolve around the `active` [setopts](http://www.erlang.org/doc/man/inet.html#setopts-2)
flag. A socket can be set to one of:

* `{active,false}` - All data must be obtained through [recv/2](http://www.erlang.org/doc/man/gen_tcp.html#recv-2)
                     calls. This amounts to syncronous socket reading.
* `{active,true}`  - All data on the socket gets sent to the controlling thread as a normal erlang
                     message. It is the thread's responsibility to keep up with the buffered data
                     in the message queue. This amounts to asyncronous socket reading.
* `{active,once}`  - When set the socket is placed in `{active,true}` for a single packet. That
                     is, once set the thread can expect a single message to be sent to when data
                     comes in. To receive any more data off of the socket the socket must either
                     be read from using [recv/2](http://www.erlang.org/doc/man/gen_tcp.html#recv-2)
                     or be put in `{active,once}` or `{active,true}`.

# Which to use?

<Explanation of how other sources claim you should use active,once, and why>

# Why not to use it


