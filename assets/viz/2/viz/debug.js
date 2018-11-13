// Compiled by ClojureScript 1.10.439 {}
goog.provide('viz.debug');
goog.require('cljs.core');
viz.debug.log = (function viz$debug$log(var_args){
var args__4647__auto__ = [];
var len__4641__auto___2514 = arguments.length;
var i__4642__auto___2515 = (0);
while(true){
if((i__4642__auto___2515 < len__4641__auto___2514)){
args__4647__auto__.push((arguments[i__4642__auto___2515]));

var G__2516 = (i__4642__auto___2515 + (1));
i__4642__auto___2515 = G__2516;
continue;
} else {
}
break;
}

var argseq__4648__auto__ = ((((0) < args__4647__auto__.length))?(new cljs.core.IndexedSeq(args__4647__auto__.slice((0)),(0),null)):null);
return viz.debug.log.cljs$core$IFn$_invoke$arity$variadic(argseq__4648__auto__);
});

viz.debug.log.cljs$core$IFn$_invoke$arity$variadic = (function (args){
return console.log(clojure.string.join.call(null," ",cljs.core.map.call(null,cljs.core.str,args)));
});

viz.debug.log.cljs$lang$maxFixedArity = (0);

/** @this {Function} */
viz.debug.log.cljs$lang$applyTo = (function (seq2513){
var self__4629__auto__ = this;
return self__4629__auto__.cljs$core$IFn$_invoke$arity$variadic(cljs.core.seq.call(null,seq2513));
});


//# sourceMappingURL=debug.js.map
