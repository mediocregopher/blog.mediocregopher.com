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

Many (most?) tutorials advocate using `{active,once}` in your application [0][1][2]. This has to do with usability and
security. When in `{active,true}` it's possible for a client to flood the connection faster than the receiving process
will process those messages, potentially eating up a lot of memory in the VM. However, if you want to be able to receive
both tcp data messages as well as other messages from other erlang processes at the same time you can't use `{active,false}`.
So `{active,once}` is generally preferred because it deals with both of these problems quite well.

# Why not to use `{active,once}`

Here's what your classic `{active,once}` enabled tcp socket implementation will probably look like:

```erlang
-module(tcp_test).
-compile(export_all).

-define(TCP_OPTS, [binary, {packet, raw}, {nodelay,true}, {active, false}, {reuseaddr, true}, {keepalive,true}, {backlog,500}]).

%Start listening
listen(Port) ->
    {ok, L} = gen_tcp:listen(Port, ?TCP_OPTS),
    ?MODULE:accept(L).

%Accept a connection
accept(L) ->
    {ok, Socket} = gen_tcp:accept(L),
    ?MODULE:read_loop(Socket),
    io:fwrite("Done reading, connection was closed\n"),
    ?MODULE:accept(L).

%Read everything it sends us
read_loop(Socket) ->
    inet:setopts(Socket, [{active, once}]),
    receive
    {tcp, _, _} ->
        do_stuff_here,
        ?MODULE:read_loop(Socket);
    {tcp_closed, _}-> donezo;
    {tcp_error, _, _} -> donezo
    end.
```

This code isn't actually usable for a production system; it doesn't even spawn a new process for the new socket. But that's not
the point I'm making. If I run it with `tcp_test:listen(8000)`, and in other window do:

```bash
while [ 1 ]; do echo "aloha"; done | nc localhost 8000
```

We'll be flooding the the server with data pretty well. Using [eprof](http://www.erlang.org/doc/man/eprof.html) we can get an idea
of how our code performs, and where the hang-ups are:

```erlang
1> eprof:start().
{ok,<0.34.0>}

2> P = spawn(tcp_test,listen,[8000]).
<0.36.0>

3> eprof:start_profiling([P]).
profiling

4> running_the_while_loop.
running_the_while_loop

5> eprof:stop_profiling().
profiling_stopped

6> eprof:analyze(procs,[{sort,time}]).

****** Process <0.36.0>    -- 100.00 % of profiled time *** 
FUNCTION                           CALLS      %      TIME  [uS / CALLS]
--------                           -----    ---      ----  [----------]
prim_inet:type_value_2/2               6   0.00         0  [      0.00]

....snip....

prim_inet:enc_opts/2                   6   0.00         8  [      1.33]
prim_inet:setopts/2             12303599   1.85   1466319  [      0.12]
tcp_test:read_loop/1            12303598   2.22   1761775  [      0.14]
prim_inet:encode_opt_val/1      12303599   3.50   2769285  [      0.23]
prim_inet:ctl_cmd/3             12303600   4.29   3399333  [      0.28]
prim_inet:enc_opt_val/2         24607203   5.28   4184818  [      0.17]
inet:setopts/2                  12303598   5.72   4533863  [      0.37]
erlang:port_control/3           12303600  77.13  61085040  [      4.96]
```

eprof shows us where our process is spending the majority of its time. The `%` column indicates percentage of time the process spent
during profiling inside any function. We can pretty clearly see that the vast majority of time was spent inside `erlang:port_control/3`,
the BIF that `inet:setopts/2` uses to switch the socket to `{active,once}` mode. Amongst the calls which were called on every loop,
it takes up by far the most amount of time. In addition all of those other calls are also related to `inet:setopts/2`.

I'm gonna rewrite our little listen server to use `{active,true}`, and we'll do it all again:

```erlang
-module(tcp_test).
-compile(export_all).

-define(TCP_OPTS, [binary, {packet, raw}, {nodelay,true}, {active, false}, {reuseaddr, true}, {keepalive,true}, {backlog,500}]).

%Start listening
listen(Port) ->
    {ok, L} = gen_tcp:listen(Port, ?TCP_OPTS),
    ?MODULE:accept(L).

%Accept a connection
accept(L) ->
    {ok, Socket} = gen_tcp:accept(L),
    inet:setopts(Socket, [{active, true}]), %Well this is new
    ?MODULE:read_loop(Socket),
    io:fwrite("Done reading, connection was closed\n"),
    ?MODULE:accept(L).

%Read everything it sends us
read_loop(Socket) ->
    %inet:setopts(Socket, [{active, once}]),
    receive
    {tcp, _, _} ->
        do_stuff_here,
        ?MODULE:read_loop(Socket);
    {tcp_closed, _}-> donezo;
    {tcp_error, _, _} -> donezo
    end.
```

And the profiling results:

```erlang
1> eprof:start().
{ok,<0.34.0>}

2>  P = spawn(tcp_test,listen,[8000]).
<0.36.0>

3> eprof:start_profiling([P]).
profiling

4>  running_the_while_loop.
running_the_while_loop

5> eprof:stop_profiling().
profiling_stopped

6> eprof:analyze(procs,[{sort,time}]).

****** Process <0.36.0>    -- 100.00 % of profiled time *** 
FUNCTION                           CALLS       %      TIME  [uS / CALLS]
--------                           -----     ---      ----  [----------]
prim_inet:enc_value_1/3                7    0.00         1  [      0.14]
prim_inet:decode_opt_val/1             1    0.00         1  [      1.00]
inet:setopts/2                         1    0.00         2  [      2.00]
prim_inet:setopts/2                    2    0.00         2  [      1.00]
prim_inet:enum_name/2                  1    0.00         2  [      2.00]
erlang:port_set_data/2                 1    0.00         2  [      2.00]
inet_db:register_socket/2              1    0.00         3  [      3.00]
prim_inet:type_value_1/3               7    0.00         3  [      0.43]

.... snip ....

prim_inet:type_opt_1/1                19    0.00         7  [      0.37]
prim_inet:enc_value/3                  7    0.00         7  [      1.00]
prim_inet:enum_val/2                   6    0.00         7  [      1.17]
prim_inet:dec_opt_val/1                7    0.00         7  [      1.00]
prim_inet:dec_value/2                  6    0.00        10  [      1.67]
prim_inet:enc_opt/1                   13    0.00        12  [      0.92]
prim_inet:type_opt/2                  19    0.00        33  [      1.74]
erlang:port_control/3                  3    0.00        59  [     19.67]
tcp_test:read_loop/1            20716370  100.00  12187488  [      0.59]
```

This time our process spent almost no time at all (according to eprof, 0%) fiddling with the socket opts.
Instead it spent all of its time in the read_loop doing the work we actually want to be doing.

# So what does this mean?

I'm by no means advocating never using `{active,once}`. The security concern is still a completely valid concern and one
that `{active,once}` mitigates quite well. I'm simply pointing out that this mitigation has some fairly serious performance
implications which have the potential to bite you if you're not careful, especially in cases where a socket is going to be
receiving a large amount of traffic.

# Meta

These tests were done using R15B03, but I've done similar ones in R14 and found similar results. I have not tested R16.

[0] http://learnyousomeerlang.com/buckets-of-sockets
[1] http://www.erlang.org/doc/man/gen_tcp.html#examples
[2] http://erlycoder.com/25/erlang-tcp-server-tcp-client-sockets-with-gen_tcp
