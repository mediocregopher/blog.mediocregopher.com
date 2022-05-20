// Compiled by ClojureScript 1.9.473 {}
goog.provide('viz.core');
goog.require('cljs.core');
goog.require('quil.core');
goog.require('quil.middleware');
goog.require('viz.forest');
goog.require('viz.grid');
goog.require('viz.ghost');
goog.require('goog.string');
goog.require('goog.string.format');
viz.core.debug = (function viz$core$debug(var_args){
var args__7934__auto__ = [];
var len__7927__auto___18287 = arguments.length;
var i__7928__auto___18288 = (0);
while(true){
if((i__7928__auto___18288 < len__7927__auto___18287)){
args__7934__auto__.push((arguments[i__7928__auto___18288]));

var G__18289 = (i__7928__auto___18288 + (1));
i__7928__auto___18288 = G__18289;
continue;
} else {
}
break;
}

var argseq__7935__auto__ = ((((0) < args__7934__auto__.length))?(new cljs.core.IndexedSeq(args__7934__auto__.slice((0)),(0),null)):null);
return viz.core.debug.cljs$core$IFn$_invoke$arity$variadic(argseq__7935__auto__);
});

viz.core.debug.cljs$core$IFn$_invoke$arity$variadic = (function (args){
return console.log(clojure.string.join.call(null," ",cljs.core.map.call(null,cljs.core.str,args)));
});

viz.core.debug.cljs$lang$maxFixedArity = (0);

viz.core.debug.cljs$lang$applyTo = (function (seq18286){
return viz.core.debug.cljs$core$IFn$_invoke$arity$variadic(cljs.core.seq.call(null,seq18286));
});

viz.core.window_partial = (function viz$core$window_partial(k){
return (((document["documentElement"][k]) * 0.95) | (0));
});
viz.core.window_size = new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [(function (){var x__7157__auto__ = (1025);
var y__7158__auto__ = viz.core.window_partial.call(null,"clientWidth");
return ((x__7157__auto__ < y__7158__auto__) ? x__7157__auto__ : y__7158__auto__);
})(),((viz.core.window_partial.call(null,"clientHeight") * 0.75) | (0))], null);
viz.core.window_half_size = cljs.core.apply.call(null,cljs.core.vector,cljs.core.map.call(null,(function (p1__18290_SHARP_){
return (p1__18290_SHARP_ / (2));
}),viz.core.window_size));
viz.core.new_state = (function viz$core$new_state(){
return new cljs.core.PersistentArrayMap(null, 7, [new cljs.core.Keyword(null,"frame-rate","frame-rate",-994918942),(15),new cljs.core.Keyword(null,"exit-wait-frames","exit-wait-frames",1417213098),(40),new cljs.core.Keyword(null,"tail-length","tail-length",-2007115089),(15),new cljs.core.Keyword(null,"frame","frame",-1711082588),(0),new cljs.core.Keyword(null,"gif-seconds","gif-seconds",1861397548),(0),new cljs.core.Keyword(null,"grid-width","grid-width",837583106),(30),new cljs.core.Keyword(null,"ghost","ghost",-1531157576),viz.ghost.new_active_node.call(null,viz.ghost.new_ghost.call(null,viz.grid.euclidean),new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [(0),(0)], null))], null);
});
viz.core.curr_second = (function viz$core$curr_second(state){
return (new cljs.core.Keyword(null,"frame","frame",-1711082588).cljs$core$IFn$_invoke$arity$1(state) / new cljs.core.Keyword(null,"frame-rate","frame-rate",-994918942).cljs$core$IFn$_invoke$arity$1(state));
});
viz.core.grid_size = (function viz$core$grid_size(state){
var h = ((viz.core.window_size.call(null,(1)) * (new cljs.core.Keyword(null,"grid-width","grid-width",837583106).cljs$core$IFn$_invoke$arity$1(state) / viz.core.window_size.call(null,(0)))) | (0));
return new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"grid-width","grid-width",837583106).cljs$core$IFn$_invoke$arity$1(state),h], null);
});
viz.core.positive = (function viz$core$positive(n){
if(((0) > n)){
return (- n);
} else {
return n;
}
});
viz.core.spawn_chance = (function viz$core$spawn_chance(state){
var period_seconds = (1);
var period_frames = (new cljs.core.Keyword(null,"frame-rate","frame-rate",-994918942).cljs$core$IFn$_invoke$arity$1(state) * period_seconds);
if((cljs.core.rem.call(null,new cljs.core.Keyword(null,"frame","frame",-1711082588).cljs$core$IFn$_invoke$arity$1(state),period_frames) === (0))){
return (1);
} else {
return (100);
}
});
viz.core.mk_poss_fn = (function viz$core$mk_poss_fn(state){
return (function (pos,adj_poss){
return cljs.core.take.call(null,(2),cljs.core.random_sample.call(null,0.6,adj_poss));
});
});
viz.core.setup = (function viz$core$setup(){
var state = viz.core.new_state.call(null);
quil.core.frame_rate.call(null,new cljs.core.Keyword(null,"frame-rate","frame-rate",-994918942).cljs$core$IFn$_invoke$arity$1(state));

return state;
});
viz.core.scale = (function viz$core$scale(state,xy){
return cljs.core.map_indexed.call(null,(function (p1__18292_SHARP_,p2__18291_SHARP_){
return (p2__18291_SHARP_ * (viz.core.window_half_size.call(null,p1__18292_SHARP_) / viz.core.grid_size.call(null,state).call(null,p1__18292_SHARP_)));
}),xy);
});
viz.core.in_bounds_QMARK_ = (function viz$core$in_bounds_QMARK_(min_bound,max_bound,pos){
var pos_k = cljs.core.keep_indexed.call(null,(function (p1__18293_SHARP_,p2__18294_SHARP_){
var mini = min_bound.call(null,p1__18293_SHARP_);
var maxi = max_bound.call(null,p1__18293_SHARP_);
if(((p2__18294_SHARP_ >= mini)) && ((p2__18294_SHARP_ <= maxi))){
return p2__18294_SHARP_;
} else {
return null;
}
}),pos);
return cljs.core._EQ_.call(null,cljs.core.count.call(null,pos),cljs.core.count.call(null,pos_k));
});
viz.core.quil_bounds = (function viz$core$quil_bounds(state,buffer){
var vec__18299 = cljs.core.apply.call(null,cljs.core.vector,cljs.core.map.call(null,(function (p1__18295_SHARP_){
return (p1__18295_SHARP_ - buffer);
}),viz.core.grid_size.call(null,state)));
var w = cljs.core.nth.call(null,vec__18299,(0),null);
var h = cljs.core.nth.call(null,vec__18299,(1),null);
return new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [(- w),(- h)], null),new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [w,h], null)], null);
});
viz.core.ghost_incr = (function viz$core$ghost_incr(state){
return cljs.core.assoc.call(null,state,new cljs.core.Keyword(null,"ghost","ghost",-1531157576),viz.ghost.filter_active_nodes.call(null,viz.ghost.incr.call(null,new cljs.core.Keyword(null,"ghost","ghost",-1531157576).cljs$core$IFn$_invoke$arity$1(state),viz.core.mk_poss_fn.call(null,state)),(function (p1__18302_SHARP_){
var vec__18306 = viz.core.quil_bounds.call(null,state,(2));
var minb = cljs.core.nth.call(null,vec__18306,(0),null);
var maxb = cljs.core.nth.call(null,vec__18306,(1),null);
return viz.core.in_bounds_QMARK_.call(null,minb,maxb,new cljs.core.Keyword(null,"pos","pos",-864607220).cljs$core$IFn$_invoke$arity$1(p1__18302_SHARP_));
})));
});
viz.core.ghost_expire_roots = (function viz$core$ghost_expire_roots(state){
if(!((new cljs.core.Keyword(null,"tail-length","tail-length",-2007115089).cljs$core$IFn$_invoke$arity$1(state) < new cljs.core.Keyword(null,"frame","frame",-1711082588).cljs$core$IFn$_invoke$arity$1(state)))){
return state;
} else {
return cljs.core.update_in.call(null,state,new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"ghost","ghost",-1531157576)], null),viz.ghost.remove_roots);
}
});
viz.core.maybe_exit = (function viz$core$maybe_exit(state){
if(cljs.core.empty_QMARK_.call(null,cljs.core.get_in.call(null,state,new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"ghost","ghost",-1531157576),new cljs.core.Keyword(null,"active-node-ids","active-node-ids",-398210751)], null)))){
if((new cljs.core.Keyword(null,"exit-wait-frames","exit-wait-frames",1417213098).cljs$core$IFn$_invoke$arity$1(state) === (0))){
return viz.core.new_state.call(null);
} else {
return cljs.core.update_in.call(null,state,new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"exit-wait-frames","exit-wait-frames",1417213098)], null),cljs.core.dec);
}
} else {
return state;
}
});
viz.core.update_state = (function viz$core$update_state(state){
return viz.core.maybe_exit.call(null,cljs.core.update_in.call(null,viz.core.ghost_expire_roots.call(null,viz.core.ghost_incr.call(null,state)),new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"frame","frame",-1711082588)], null),cljs.core.inc));
});
viz.core.draw_ellipse = (function viz$core$draw_ellipse(state,pos,size){
var scaled_pos = viz.core.scale.call(null,state,pos);
var scaled_size = cljs.core.map.call(null,cljs.core.int$,viz.core.scale.call(null,state,size));
return cljs.core.apply.call(null,quil.core.ellipse,cljs.core.concat.call(null,scaled_pos,scaled_size));
});
viz.core.in_line_QMARK_ = (function viz$core$in_line_QMARK_(var_args){
var args__7934__auto__ = [];
var len__7927__auto___18311 = arguments.length;
var i__7928__auto___18312 = (0);
while(true){
if((i__7928__auto___18312 < len__7927__auto___18311)){
args__7934__auto__.push((arguments[i__7928__auto___18312]));

var G__18313 = (i__7928__auto___18312 + (1));
i__7928__auto___18312 = G__18313;
continue;
} else {
}
break;
}

var argseq__7935__auto__ = ((((0) < args__7934__auto__.length))?(new cljs.core.IndexedSeq(args__7934__auto__.slice((0)),(0),null)):null);
return viz.core.in_line_QMARK_.cljs$core$IFn$_invoke$arity$variadic(argseq__7935__auto__);
});

viz.core.in_line_QMARK_.cljs$core$IFn$_invoke$arity$variadic = (function (nodes){
return cljs.core.apply.call(null,cljs.core._EQ_,cljs.core.map.call(null,(function (p1__18309_SHARP_){
return cljs.core.apply.call(null,cljs.core.map,cljs.core._,p1__18309_SHARP_);
}),cljs.core.partition.call(null,(2),(1),cljs.core.map.call(null,new cljs.core.Keyword(null,"pos","pos",-864607220),nodes))));
});

viz.core.in_line_QMARK_.cljs$lang$maxFixedArity = (0);

viz.core.in_line_QMARK_.cljs$lang$applyTo = (function (seq18310){
return viz.core.in_line_QMARK_.cljs$core$IFn$_invoke$arity$variadic(cljs.core.seq.call(null,seq18310));
});

viz.core.draw_lines = (function viz$core$draw_lines(state,forest,parent,node){

quil.core.stroke.call(null,(4278190080));

quil.core.fill.call(null,(4294967295));

var children = cljs.core.map.call(null,(function (p1__18314_SHARP_){
return viz.forest.get_node.call(null,forest,p1__18314_SHARP_);
}),new cljs.core.Keyword(null,"child-ids","child-ids",-604525861).cljs$core$IFn$_invoke$arity$1(node));
if(cljs.core.not.call(null,parent)){
var seq__18325_18333 = cljs.core.seq.call(null,children);
var chunk__18326_18334 = null;
var count__18327_18335 = (0);
var i__18328_18336 = (0);
while(true){
if((i__18328_18336 < count__18327_18335)){
var child_18337 = cljs.core._nth.call(null,chunk__18326_18334,i__18328_18336);
viz.core.draw_lines.call(null,state,forest,node,child_18337);

var G__18338 = seq__18325_18333;
var G__18339 = chunk__18326_18334;
var G__18340 = count__18327_18335;
var G__18341 = (i__18328_18336 + (1));
seq__18325_18333 = G__18338;
chunk__18326_18334 = G__18339;
count__18327_18335 = G__18340;
i__18328_18336 = G__18341;
continue;
} else {
var temp__4657__auto___18342 = cljs.core.seq.call(null,seq__18325_18333);
if(temp__4657__auto___18342){
var seq__18325_18343__$1 = temp__4657__auto___18342;
if(cljs.core.chunked_seq_QMARK_.call(null,seq__18325_18343__$1)){
var c__7633__auto___18344 = cljs.core.chunk_first.call(null,seq__18325_18343__$1);
var G__18345 = cljs.core.chunk_rest.call(null,seq__18325_18343__$1);
var G__18346 = c__7633__auto___18344;
var G__18347 = cljs.core.count.call(null,c__7633__auto___18344);
var G__18348 = (0);
seq__18325_18333 = G__18345;
chunk__18326_18334 = G__18346;
count__18327_18335 = G__18347;
i__18328_18336 = G__18348;
continue;
} else {
var child_18349 = cljs.core.first.call(null,seq__18325_18343__$1);
viz.core.draw_lines.call(null,state,forest,node,child_18349);

var G__18350 = cljs.core.next.call(null,seq__18325_18343__$1);
var G__18351 = null;
var G__18352 = (0);
var G__18353 = (0);
seq__18325_18333 = G__18350;
chunk__18326_18334 = G__18351;
count__18327_18335 = G__18352;
i__18328_18336 = G__18353;
continue;
}
} else {
}
}
break;
}
} else {
var in_line_child_18354 = cljs.core.some.call(null,((function (children){
return (function (p1__18315_SHARP_){
if(cljs.core.truth_(viz.core.in_line_QMARK_.call(null,parent,node,p1__18315_SHARP_))){
return p1__18315_SHARP_;
} else {
return null;
}
});})(children))
,children);
var seq__18329_18355 = cljs.core.seq.call(null,children);
var chunk__18330_18356 = null;
var count__18331_18357 = (0);
var i__18332_18358 = (0);
while(true){
if((i__18332_18358 < count__18331_18357)){
var child_18359 = cljs.core._nth.call(null,chunk__18330_18356,i__18332_18358);
if(cljs.core.truth_((function (){var and__6802__auto__ = in_line_child_18354;
if(cljs.core.truth_(and__6802__auto__)){
return cljs.core._EQ_.call(null,in_line_child_18354,child_18359);
} else {
return and__6802__auto__;
}
})())){
viz.core.draw_lines.call(null,state,forest,parent,child_18359);
} else {
viz.core.draw_lines.call(null,state,forest,node,child_18359);
}

var G__18360 = seq__18329_18355;
var G__18361 = chunk__18330_18356;
var G__18362 = count__18331_18357;
var G__18363 = (i__18332_18358 + (1));
seq__18329_18355 = G__18360;
chunk__18330_18356 = G__18361;
count__18331_18357 = G__18362;
i__18332_18358 = G__18363;
continue;
} else {
var temp__4657__auto___18364 = cljs.core.seq.call(null,seq__18329_18355);
if(temp__4657__auto___18364){
var seq__18329_18365__$1 = temp__4657__auto___18364;
if(cljs.core.chunked_seq_QMARK_.call(null,seq__18329_18365__$1)){
var c__7633__auto___18366 = cljs.core.chunk_first.call(null,seq__18329_18365__$1);
var G__18367 = cljs.core.chunk_rest.call(null,seq__18329_18365__$1);
var G__18368 = c__7633__auto___18366;
var G__18369 = cljs.core.count.call(null,c__7633__auto___18366);
var G__18370 = (0);
seq__18329_18355 = G__18367;
chunk__18330_18356 = G__18368;
count__18331_18357 = G__18369;
i__18332_18358 = G__18370;
continue;
} else {
var child_18371 = cljs.core.first.call(null,seq__18329_18365__$1);
if(cljs.core.truth_((function (){var and__6802__auto__ = in_line_child_18354;
if(cljs.core.truth_(and__6802__auto__)){
return cljs.core._EQ_.call(null,in_line_child_18354,child_18371);
} else {
return and__6802__auto__;
}
})())){
viz.core.draw_lines.call(null,state,forest,parent,child_18371);
} else {
viz.core.draw_lines.call(null,state,forest,node,child_18371);
}

var G__18372 = cljs.core.next.call(null,seq__18329_18365__$1);
var G__18373 = null;
var G__18374 = (0);
var G__18375 = (0);
seq__18329_18355 = G__18372;
chunk__18330_18356 = G__18373;
count__18331_18357 = G__18374;
i__18332_18358 = G__18375;
continue;
}
} else {
}
}
break;
}

if(cljs.core.truth_(in_line_child_18354)){
} else {
cljs.core.apply.call(null,quil.core.line,cljs.core.apply.call(null,cljs.core.concat,cljs.core.map.call(null,((function (in_line_child_18354,children){
return (function (p1__18316_SHARP_){
return viz.core.scale.call(null,state,p1__18316_SHARP_);
});})(in_line_child_18354,children))
,cljs.core.map.call(null,new cljs.core.Keyword(null,"pos","pos",-864607220),(function (){var x__7656__auto__ = parent;
return cljs.core._conj.call(null,(function (){var x__7656__auto____$1 = node;
return cljs.core._conj.call(null,cljs.core.List.EMPTY,x__7656__auto____$1);
})(),x__7656__auto__);
})()))));
}
}

if(cljs.core.empty_QMARK_.call(null,children)){
return viz.core.draw_ellipse.call(null,state,new cljs.core.Keyword(null,"pos","pos",-864607220).cljs$core$IFn$_invoke$arity$1(node),new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [0.3,0.3], null));
} else {
return null;
}
});
viz.core.draw_state = (function viz$core$draw_state(state){
quil.core.background.call(null,(4294967295));

var tr__8398__auto__ = new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [(viz.core.window_size.call(null,(0)) / (2)),(viz.core.window_size.call(null,(1)) / (2))], null);
quil.core.push_matrix.call(null);

try{quil.core.translate.call(null,tr__8398__auto__);

var lines = viz.forest.lines.call(null,cljs.core.get_in.call(null,state,new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"ghost","ghost",-1531157576),new cljs.core.Keyword(null,"forest","forest",278860306)], null)));
var leaves = viz.forest.leaves.call(null,cljs.core.get_in.call(null,state,new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"ghost","ghost",-1531157576),new cljs.core.Keyword(null,"forest","forest",278860306)], null)));
var active = viz.ghost.active_nodes.call(null,new cljs.core.Keyword(null,"ghost","ghost",-1531157576).cljs$core$IFn$_invoke$arity$1(state));
var roots = viz.forest.roots.call(null,cljs.core.get_in.call(null,state,new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"ghost","ghost",-1531157576),new cljs.core.Keyword(null,"forest","forest",278860306)], null)));
quil.core.stroke.call(null,(4278190080));

var seq__18384_18392 = cljs.core.seq.call(null,roots);
var chunk__18385_18393 = null;
var count__18386_18394 = (0);
var i__18387_18395 = (0);
while(true){
if((i__18387_18395 < count__18386_18394)){
var root_18396 = cljs.core._nth.call(null,chunk__18385_18393,i__18387_18395);
viz.core.draw_lines.call(null,state,cljs.core.get_in.call(null,state,new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"ghost","ghost",-1531157576),new cljs.core.Keyword(null,"forest","forest",278860306)], null)),null,root_18396);

var G__18397 = seq__18384_18392;
var G__18398 = chunk__18385_18393;
var G__18399 = count__18386_18394;
var G__18400 = (i__18387_18395 + (1));
seq__18384_18392 = G__18397;
chunk__18385_18393 = G__18398;
count__18386_18394 = G__18399;
i__18387_18395 = G__18400;
continue;
} else {
var temp__4657__auto___18401 = cljs.core.seq.call(null,seq__18384_18392);
if(temp__4657__auto___18401){
var seq__18384_18402__$1 = temp__4657__auto___18401;
if(cljs.core.chunked_seq_QMARK_.call(null,seq__18384_18402__$1)){
var c__7633__auto___18403 = cljs.core.chunk_first.call(null,seq__18384_18402__$1);
var G__18404 = cljs.core.chunk_rest.call(null,seq__18384_18402__$1);
var G__18405 = c__7633__auto___18403;
var G__18406 = cljs.core.count.call(null,c__7633__auto___18403);
var G__18407 = (0);
seq__18384_18392 = G__18404;
chunk__18385_18393 = G__18405;
count__18386_18394 = G__18406;
i__18387_18395 = G__18407;
continue;
} else {
var root_18408 = cljs.core.first.call(null,seq__18384_18402__$1);
viz.core.draw_lines.call(null,state,cljs.core.get_in.call(null,state,new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"ghost","ghost",-1531157576),new cljs.core.Keyword(null,"forest","forest",278860306)], null)),null,root_18408);

var G__18409 = cljs.core.next.call(null,seq__18384_18402__$1);
var G__18410 = null;
var G__18411 = (0);
var G__18412 = (0);
seq__18384_18392 = G__18409;
chunk__18385_18393 = G__18410;
count__18386_18394 = G__18411;
i__18387_18395 = G__18412;
continue;
}
} else {
}
}
break;
}

quil.core.stroke.call(null,(4278190080));

quil.core.fill.call(null,(4278190080));

var seq__18388 = cljs.core.seq.call(null,active);
var chunk__18389 = null;
var count__18390 = (0);
var i__18391 = (0);
while(true){
if((i__18391 < count__18390)){
var active_node = cljs.core._nth.call(null,chunk__18389,i__18391);
var pos_18413 = new cljs.core.Keyword(null,"pos","pos",-864607220).cljs$core$IFn$_invoke$arity$1(active_node);
viz.core.draw_ellipse.call(null,state,pos_18413,new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [0.35,0.35], null));

var G__18414 = seq__18388;
var G__18415 = chunk__18389;
var G__18416 = count__18390;
var G__18417 = (i__18391 + (1));
seq__18388 = G__18414;
chunk__18389 = G__18415;
count__18390 = G__18416;
i__18391 = G__18417;
continue;
} else {
var temp__4657__auto__ = cljs.core.seq.call(null,seq__18388);
if(temp__4657__auto__){
var seq__18388__$1 = temp__4657__auto__;
if(cljs.core.chunked_seq_QMARK_.call(null,seq__18388__$1)){
var c__7633__auto__ = cljs.core.chunk_first.call(null,seq__18388__$1);
var G__18418 = cljs.core.chunk_rest.call(null,seq__18388__$1);
var G__18419 = c__7633__auto__;
var G__18420 = cljs.core.count.call(null,c__7633__auto__);
var G__18421 = (0);
seq__18388 = G__18418;
chunk__18389 = G__18419;
count__18390 = G__18420;
i__18391 = G__18421;
continue;
} else {
var active_node = cljs.core.first.call(null,seq__18388__$1);
var pos_18422 = new cljs.core.Keyword(null,"pos","pos",-864607220).cljs$core$IFn$_invoke$arity$1(active_node);
viz.core.draw_ellipse.call(null,state,pos_18422,new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [0.35,0.35], null));

var G__18423 = cljs.core.next.call(null,seq__18388__$1);
var G__18424 = null;
var G__18425 = (0);
var G__18426 = (0);
seq__18388 = G__18423;
chunk__18389 = G__18424;
count__18390 = G__18425;
i__18391 = G__18426;
continue;
}
} else {
return null;
}
}
break;
}
}finally {quil.core.pop_matrix.call(null);
}});
viz.core.viz = (function viz$core$viz(){
return quil.sketch.sketch.call(null,new cljs.core.Keyword(null,"host","host",-1558485167),"viz",new cljs.core.Keyword(null,"features","features",-1146962336),new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"keep-on-top","keep-on-top",-970284267)], null),new cljs.core.Keyword(null,"update","update",1045576396),((cljs.core.fn_QMARK_.call(null,viz.core.update_state))?(function() { 
var G__18427__delegate = function (args){
return cljs.core.apply.call(null,viz.core.update_state,args);
};
var G__18427 = function (var_args){
var args = null;
if (arguments.length > 0) {
var G__18428__i = 0, G__18428__a = new Array(arguments.length -  0);
while (G__18428__i < G__18428__a.length) {G__18428__a[G__18428__i] = arguments[G__18428__i + 0]; ++G__18428__i;}
  args = new cljs.core.IndexedSeq(G__18428__a,0);
} 
return G__18427__delegate.call(this,args);};
G__18427.cljs$lang$maxFixedArity = 0;
G__18427.cljs$lang$applyTo = (function (arglist__18429){
var args = cljs.core.seq(arglist__18429);
return G__18427__delegate(args);
});
G__18427.cljs$core$IFn$_invoke$arity$variadic = G__18427__delegate;
return G__18427;
})()
:viz.core.update_state),new cljs.core.Keyword(null,"size","size",1098693007),((cljs.core.fn_QMARK_.call(null,viz.core.window_size))?(function() { 
var G__18430__delegate = function (args){
return cljs.core.apply.call(null,viz.core.window_size,args);
};
var G__18430 = function (var_args){
var args = null;
if (arguments.length > 0) {
var G__18431__i = 0, G__18431__a = new Array(arguments.length -  0);
while (G__18431__i < G__18431__a.length) {G__18431__a[G__18431__i] = arguments[G__18431__i + 0]; ++G__18431__i;}
  args = new cljs.core.IndexedSeq(G__18431__a,0);
} 
return G__18430__delegate.call(this,args);};
G__18430.cljs$lang$maxFixedArity = 0;
G__18430.cljs$lang$applyTo = (function (arglist__18432){
var args = cljs.core.seq(arglist__18432);
return G__18430__delegate(args);
});
G__18430.cljs$core$IFn$_invoke$arity$variadic = G__18430__delegate;
return G__18430;
})()
:viz.core.window_size),new cljs.core.Keyword(null,"title","title",636505583),"",new cljs.core.Keyword(null,"setup","setup",1987730512),((cljs.core.fn_QMARK_.call(null,viz.core.setup))?(function() { 
var G__18433__delegate = function (args){
return cljs.core.apply.call(null,viz.core.setup,args);
};
var G__18433 = function (var_args){
var args = null;
if (arguments.length > 0) {
var G__18434__i = 0, G__18434__a = new Array(arguments.length -  0);
while (G__18434__i < G__18434__a.length) {G__18434__a[G__18434__i] = arguments[G__18434__i + 0]; ++G__18434__i;}
  args = new cljs.core.IndexedSeq(G__18434__a,0);
} 
return G__18433__delegate.call(this,args);};
G__18433.cljs$lang$maxFixedArity = 0;
G__18433.cljs$lang$applyTo = (function (arglist__18435){
var args = cljs.core.seq(arglist__18435);
return G__18433__delegate(args);
});
G__18433.cljs$core$IFn$_invoke$arity$variadic = G__18433__delegate;
return G__18433;
})()
:viz.core.setup),new cljs.core.Keyword(null,"middleware","middleware",1462115504),new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [quil.middleware.fun_mode], null),new cljs.core.Keyword(null,"draw","draw",1358331674),((cljs.core.fn_QMARK_.call(null,viz.core.draw_state))?(function() { 
var G__18436__delegate = function (args){
return cljs.core.apply.call(null,viz.core.draw_state,args);
};
var G__18436 = function (var_args){
var args = null;
if (arguments.length > 0) {
var G__18437__i = 0, G__18437__a = new Array(arguments.length -  0);
while (G__18437__i < G__18437__a.length) {G__18437__a[G__18437__i] = arguments[G__18437__i + 0]; ++G__18437__i;}
  args = new cljs.core.IndexedSeq(G__18437__a,0);
} 
return G__18436__delegate.call(this,args);};
G__18436.cljs$lang$maxFixedArity = 0;
G__18436.cljs$lang$applyTo = (function (arglist__18438){
var args = cljs.core.seq(arglist__18438);
return G__18436__delegate(args);
});
G__18436.cljs$core$IFn$_invoke$arity$variadic = G__18436__delegate;
return G__18436;
})()
:viz.core.draw_state));
});
goog.exportSymbol('viz.core.viz', viz.core.viz);

if(cljs.core.truth_(cljs.core.some.call(null,(function (p1__8011__8012__auto__){
return cljs.core._EQ_.call(null,new cljs.core.Keyword(null,"no-start","no-start",1381488856),p1__8011__8012__auto__);
}),new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"keep-on-top","keep-on-top",-970284267)], null)))){
} else {
quil.sketch.add_sketch_to_init_list.call(null,new cljs.core.PersistentArrayMap(null, 2, [new cljs.core.Keyword(null,"fn","fn",-1175266204),viz.core.viz,new cljs.core.Keyword(null,"host-id","host-id",742376279),"viz"], null));
}

//# sourceMappingURL=core.js.map