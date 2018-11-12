// Compiled by ClojureScript 1.9.473 {}
goog.provide('viz.core');
goog.require('cljs.core');
goog.require('quil.core');
goog.require('quil.middleware');
goog.require('viz.forest');
goog.require('viz.grid');
goog.require('viz.ghost');
goog.require('viz.dial');
goog.require('goog.string');
goog.require('goog.string.format');
viz.core.debug = (function viz$core$debug(var_args){
var args__7934__auto__ = [];
var len__7927__auto___15474 = arguments.length;
var i__7928__auto___15475 = (0);
while(true){
if((i__7928__auto___15475 < len__7927__auto___15474)){
args__7934__auto__.push((arguments[i__7928__auto___15475]));

var G__15476 = (i__7928__auto___15475 + (1));
i__7928__auto___15475 = G__15476;
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

viz.core.debug.cljs$lang$applyTo = (function (seq15473){
return viz.core.debug.cljs$core$IFn$_invoke$arity$variadic(cljs.core.seq.call(null,seq15473));
});

viz.core.positive = (function viz$core$positive(n){
if(((0) > n)){
return (- n);
} else {
return n;
}
});
viz.core.window_partial = (function viz$core$window_partial(k){
return (((document["documentElement"][k]) * 0.95) | (0));
});
viz.core.window_size = new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [(function (){var x__7157__auto__ = (1025);
var y__7158__auto__ = viz.core.window_partial.call(null,"clientWidth");
return ((x__7157__auto__ < y__7158__auto__) ? x__7157__auto__ : y__7158__auto__);
})(),((viz.core.window_partial.call(null,"clientHeight") * 0.75) | (0))], null);
viz.core.window_half_size = cljs.core.apply.call(null,cljs.core.vector,cljs.core.map.call(null,(function (p1__15477_SHARP_){
return (p1__15477_SHARP_ / (2));
}),viz.core.window_size));
viz.core.frame_rate = (15);
viz.core.new_state = (function viz$core$new_state(){
return cljs.core.PersistentHashMap.fromArrays([new cljs.core.Keyword(null,"grid-width","grid-width",837583106),new cljs.core.Keyword(null,"frame-rate","frame-rate",-994918942),new cljs.core.Keyword(null,"color-cycle-period","color-cycle-period",1656886882),new cljs.core.Keyword(null,"frame","frame",-1711082588),new cljs.core.Keyword(null,"heartbeat-plot","heartbeat-plot",-98284443),new cljs.core.Keyword(null,"exit-wait-frames","exit-wait-frames",1417213098),new cljs.core.Keyword(null,"gif-seconds","gif-seconds",1861397548),new cljs.core.Keyword(null,"tail-length","tail-length",-2007115089),new cljs.core.Keyword(null,"forest","forest",278860306),new cljs.core.Keyword(null,"dial","dial",1238392184),new cljs.core.Keyword(null,"ghost","ghost",-1531157576)],[(35),viz.core.frame_rate,(2),(0),viz.dial.new_plot.call(null,viz.core.frame_rate,0.86,new cljs.core.PersistentVector(null, 3, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [0.5,0.5], null),new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [0.7,(0)], null),new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [0.8,(1)], null)], null)),(40),(0),(7),viz.forest.new_forest.call(null,viz.grid.isometric),viz.dial.new_dial.call(null),cljs.core.assoc.call(null,viz.ghost.new_ghost.call(null),new cljs.core.Keyword(null,"color","color",1011675173),quil.core.color.call(null,(0),(1),(1)))]);
});
viz.core.new_active_node = (function viz$core$new_active_node(state,pos){
var vec__15481 = viz.forest.add_node.call(null,new cljs.core.Keyword(null,"forest","forest",278860306).cljs$core$IFn$_invoke$arity$1(state),pos);
var forest = cljs.core.nth.call(null,vec__15481,(0),null);
var id = cljs.core.nth.call(null,vec__15481,(1),null);
var ghost = viz.ghost.add_active_node.call(null,new cljs.core.Keyword(null,"ghost","ghost",-1531157576).cljs$core$IFn$_invoke$arity$1(state),id);
return cljs.core.assoc.call(null,state,new cljs.core.Keyword(null,"ghost","ghost",-1531157576),ghost,new cljs.core.Keyword(null,"forest","forest",278860306),forest);
});
viz.core.frames_per_color_cycle = (function viz$core$frames_per_color_cycle(state){
return (new cljs.core.Keyword(null,"color-cycle-period","color-cycle-period",1656886882).cljs$core$IFn$_invoke$arity$1(state) * new cljs.core.Keyword(null,"frame-rate","frame-rate",-994918942).cljs$core$IFn$_invoke$arity$1(state));
});
viz.core.setup = (function viz$core$setup(){
quil.core.color_mode.call(null,new cljs.core.Keyword(null,"hsb","hsb",-753472031),(10),(1),(1));

var state = viz.core.new_active_node.call(null,viz.core.new_state.call(null),new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [(10),(10)], null));
quil.core.frame_rate.call(null,new cljs.core.Keyword(null,"frame-rate","frame-rate",-994918942).cljs$core$IFn$_invoke$arity$1(state));

quil.core.color_mode.call(null,new cljs.core.Keyword(null,"hsb","hsb",-753472031),viz.core.frames_per_color_cycle.call(null,state),(1),(1));

return state;
});
viz.core.curr_second = (function viz$core$curr_second(state){
return (new cljs.core.Keyword(null,"frame","frame",-1711082588).cljs$core$IFn$_invoke$arity$1(state) / new cljs.core.Keyword(null,"frame-rate","frame-rate",-994918942).cljs$core$IFn$_invoke$arity$1(state));
});
viz.core.grid_size = (function viz$core$grid_size(state){
var h = ((viz.core.window_size.call(null,(1)) * (new cljs.core.Keyword(null,"grid-width","grid-width",837583106).cljs$core$IFn$_invoke$arity$1(state) / viz.core.window_size.call(null,(0)))) | (0));
return new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"grid-width","grid-width",837583106).cljs$core$IFn$_invoke$arity$1(state),h], null);
});
viz.core.scale = (function viz$core$scale(state,xy){
return cljs.core.map_indexed.call(null,(function (p1__15485_SHARP_,p2__15484_SHARP_){
return (p2__15484_SHARP_ * (viz.core.window_half_size.call(null,p1__15485_SHARP_) / viz.core.grid_size.call(null,state).call(null,p1__15485_SHARP_)));
}),xy);
});
viz.core.bounds_buffer = (1);
viz.core.in_bounds_QMARK_ = (function viz$core$in_bounds_QMARK_(state,pos){
var vec__15492 = cljs.core.apply.call(null,cljs.core.vector,cljs.core.map.call(null,(function (p1__15486_SHARP_){
return (p1__15486_SHARP_ - viz.core.bounds_buffer);
}),viz.core.grid_size.call(null,state)));
var w = cljs.core.nth.call(null,vec__15492,(0),null);
var h = cljs.core.nth.call(null,vec__15492,(1),null);
var min_bound = new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [(- w),(- h)], null);
var max_bound = new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [w,h], null);
var pos_k = cljs.core.keep_indexed.call(null,((function (vec__15492,w,h,min_bound,max_bound){
return (function (p1__15487_SHARP_,p2__15488_SHARP_){
var mini = min_bound.call(null,p1__15487_SHARP_);
var maxi = max_bound.call(null,p1__15487_SHARP_);
if(((p2__15488_SHARP_ >= mini)) && ((p2__15488_SHARP_ <= maxi))){
return p2__15488_SHARP_;
} else {
return null;
}
});})(vec__15492,w,h,min_bound,max_bound))
,pos);
return cljs.core._EQ_.call(null,cljs.core.count.call(null,pos),cljs.core.count.call(null,pos_k));
});
viz.core.ceil_one = (function viz$core$ceil_one(x){
if((x > (0))){
return (1);
} else {
return (0);
}
});
viz.core.set_dial = (function viz$core$set_dial(state){
return cljs.core.update_in.call(null,state,new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"dial","dial",1238392184)], null),viz.dial.by_plot,new cljs.core.Keyword(null,"heartbeat-plot","heartbeat-plot",-98284443).cljs$core$IFn$_invoke$arity$1(state),new cljs.core.Keyword(null,"frame","frame",-1711082588).cljs$core$IFn$_invoke$arity$1(state));
});
viz.core.dist_from_sqr = (function viz$core$dist_from_sqr(pos1,pos2){
return cljs.core.reduce.call(null,cljs.core._PLUS_,cljs.core.map.call(null,(function (p1__15495_SHARP_){
return (p1__15495_SHARP_ * p1__15495_SHARP_);
}),cljs.core.map.call(null,cljs.core._,pos1,pos2)));
});
viz.core.dist_from = (function viz$core$dist_from(pos1,pos2){
return quil.core.sqrt.call(null,viz.core.dist_from_sqr.call(null,pos1,pos2));
});
viz.core.order_adj_poss_fns = new cljs.core.PersistentArrayMap(null, 2, [new cljs.core.Keyword(null,"random","random",-557811113),(function (state){
return (function (pos,adj_poss){
return cljs.core.shuffle.call(null,adj_poss);
});
}),new cljs.core.Keyword(null,"centered","centered",-515171141),(function (state){
return (function (pos,adj_poss){
return cljs.core.sort_by.call(null,(function (p1__15496_SHARP_){
return viz.core.dist_from_sqr.call(null,p1__15496_SHARP_,new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [(0),(0)], null));
}),adj_poss);
});
})], null);
viz.core.mk_order_adj_poss_fn = (function viz$core$mk_order_adj_poss_fn(var_args){
var args__7934__auto__ = [];
var len__7927__auto___15499 = arguments.length;
var i__7928__auto___15500 = (0);
while(true){
if((i__7928__auto___15500 < len__7927__auto___15499)){
args__7934__auto__.push((arguments[i__7928__auto___15500]));

var G__15501 = (i__7928__auto___15500 + (1));
i__7928__auto___15500 = G__15501;
continue;
} else {
}
break;
}

var argseq__7935__auto__ = ((((0) < args__7934__auto__.length))?(new cljs.core.IndexedSeq(args__7934__auto__.slice((0)),(0),null)):null);
return viz.core.mk_order_adj_poss_fn.cljs$core$IFn$_invoke$arity$variadic(argseq__7935__auto__);
});

viz.core.mk_order_adj_poss_fn.cljs$core$IFn$_invoke$arity$variadic = (function (ks){
return (function (state){
var fns = cljs.core.map.call(null,(function (p1__15497_SHARP_){
return p1__15497_SHARP_.call(null,state);
}),cljs.core.map.call(null,viz.core.order_adj_poss_fns,ks));
return ((function (fns){
return (function (pos,adj_poss){
return cljs.core.reduce.call(null,((function (fns){
return (function (inner_adj_poss,next_fn){
return next_fn.call(null,pos,inner_adj_poss);
});})(fns))
,adj_poss,fns);
});
;})(fns))
});
});

viz.core.mk_order_adj_poss_fn.cljs$lang$maxFixedArity = (0);

viz.core.mk_order_adj_poss_fn.cljs$lang$applyTo = (function (seq15498){
return viz.core.mk_order_adj_poss_fn.cljs$core$IFn$_invoke$arity$variadic(cljs.core.seq.call(null,seq15498));
});

viz.core.take_adj_poss_fns = new cljs.core.PersistentArrayMap(null, 3, [new cljs.core.Keyword(null,"random","random",-557811113),(function (state){
return (function (pos,adj_poss){
return quil.core.map_range.call(null,cljs.core.rand.call(null),(0),(1),0.75,(1));
});
}),new cljs.core.Keyword(null,"dial","dial",1238392184),(function (state){
return (function (pos,adj_poss){
return new cljs.core.Keyword(null,"val","val",128701612).cljs$core$IFn$_invoke$arity$1(viz.dial.scaled.call(null,new cljs.core.Keyword(null,"dial","dial",1238392184).cljs$core$IFn$_invoke$arity$1(state),-0.25,1.75));
});
}),new cljs.core.Keyword(null,"centered","centered",-515171141),(function (state){
return (function (pos,adj_poss){
var d = viz.core.dist_from.call(null,new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [(0),(0)], null),pos);
var max_d = state.call(null,new cljs.core.Keyword(null,"grid-width","grid-width",837583106));
var norm_d = (d / max_d);
return ((1) - norm_d);
});
})], null);
viz.core.mk_take_adj_poss_fn = (function viz$core$mk_take_adj_poss_fn(var_args){
var args__7934__auto__ = [];
var len__7927__auto___15505 = arguments.length;
var i__7928__auto___15506 = (0);
while(true){
if((i__7928__auto___15506 < len__7927__auto___15505)){
args__7934__auto__.push((arguments[i__7928__auto___15506]));

var G__15507 = (i__7928__auto___15506 + (1));
i__7928__auto___15506 = G__15507;
continue;
} else {
}
break;
}

var argseq__7935__auto__ = ((((0) < args__7934__auto__.length))?(new cljs.core.IndexedSeq(args__7934__auto__.slice((0)),(0),null)):null);
return viz.core.mk_take_adj_poss_fn.cljs$core$IFn$_invoke$arity$variadic(argseq__7935__auto__);
});

viz.core.mk_take_adj_poss_fn.cljs$core$IFn$_invoke$arity$variadic = (function (ks){
return (function (state){
var fns = cljs.core.map.call(null,(function (p1__15502_SHARP_){
return p1__15502_SHARP_.call(null,state);
}),cljs.core.map.call(null,viz.core.take_adj_poss_fns,ks));
return ((function (fns){
return (function (pos,adj_poss){
var mults = cljs.core.map.call(null,((function (fns){
return (function (p1__15503_SHARP_){
return p1__15503_SHARP_.call(null,pos,adj_poss);
});})(fns))
,fns);
var mult = cljs.core.reduce.call(null,cljs.core._STAR_,(1),mults);
var to_take = ((mult * cljs.core.count.call(null,adj_poss)) | (0));
return cljs.core.take.call(null,to_take,adj_poss);
});
;})(fns))
});
});

viz.core.mk_take_adj_poss_fn.cljs$lang$maxFixedArity = (0);

viz.core.mk_take_adj_poss_fn.cljs$lang$applyTo = (function (seq15504){
return viz.core.mk_take_adj_poss_fn.cljs$core$IFn$_invoke$arity$variadic(cljs.core.seq.call(null,seq15504));
});

viz.core.order_adj_poss_fn = viz.core.mk_order_adj_poss_fn.call(null,new cljs.core.Keyword(null,"centered","centered",-515171141));
viz.core.take_adj_poss_fn = viz.core.mk_take_adj_poss_fn.call(null,new cljs.core.Keyword(null,"centered","centered",-515171141),new cljs.core.Keyword(null,"random","random",-557811113));
viz.core.mk_poss_fn = (function viz$core$mk_poss_fn(state){
var order_inner_fn = viz.core.order_adj_poss_fn.call(null,state);
var take_inner_fn = viz.core.take_adj_poss_fn.call(null,state);
return ((function (order_inner_fn,take_inner_fn){
return (function (pos,adj_poss){
var adj_poss__$1 = cljs.core.filter.call(null,((function (order_inner_fn,take_inner_fn){
return (function (p1__15508_SHARP_){
return viz.core.in_bounds_QMARK_.call(null,state,p1__15508_SHARP_);
});})(order_inner_fn,take_inner_fn))
,adj_poss);
var adj_poss_ordered = order_inner_fn.call(null,pos,adj_poss__$1);
var to_take = take_inner_fn.call(null,pos,adj_poss__$1);
return take_inner_fn.call(null,pos,adj_poss_ordered);
});
;})(order_inner_fn,take_inner_fn))
});
viz.core.ghost_incr = (function viz$core$ghost_incr(state,poss_fn){
var vec__15512 = viz.ghost.incr.call(null,new cljs.core.Keyword(null,"ghost","ghost",-1531157576).cljs$core$IFn$_invoke$arity$1(state),new cljs.core.Keyword(null,"forest","forest",278860306).cljs$core$IFn$_invoke$arity$1(state),poss_fn);
var ghost = cljs.core.nth.call(null,vec__15512,(0),null);
var forest = cljs.core.nth.call(null,vec__15512,(1),null);
return cljs.core.assoc.call(null,state,new cljs.core.Keyword(null,"ghost","ghost",-1531157576),ghost,new cljs.core.Keyword(null,"forest","forest",278860306),forest);
});
viz.core.maybe_remove_roots = (function viz$core$maybe_remove_roots(state){
if((new cljs.core.Keyword(null,"tail-length","tail-length",-2007115089).cljs$core$IFn$_invoke$arity$1(state) >= new cljs.core.Keyword(null,"frame","frame",-1711082588).cljs$core$IFn$_invoke$arity$1(state))){
return state;
} else {
var roots = viz.forest.roots.call(null,new cljs.core.Keyword(null,"forest","forest",278860306).cljs$core$IFn$_invoke$arity$1(state));
var root_ids = cljs.core.map.call(null,new cljs.core.Keyword(null,"id","id",-1388402092),roots);
return cljs.core.update_in.call(null,cljs.core.update_in.call(null,state,new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"ghost","ghost",-1531157576)], null),((function (roots,root_ids){
return (function (p1__15515_SHARP_){
return cljs.core.reduce.call(null,viz.ghost.rm_active_node,p1__15515_SHARP_,root_ids);
});})(roots,root_ids))
),new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"forest","forest",278860306)], null),((function (roots,root_ids){
return (function (p1__15516_SHARP_){
return cljs.core.reduce.call(null,viz.forest.remove_node,p1__15516_SHARP_,root_ids);
});})(roots,root_ids))
);
}
});
viz.core.update_node_meta = (function viz$core$update_node_meta(state,id,f){
return cljs.core.update_in.call(null,state,new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"forest","forest",278860306)], null),viz.forest.update_node_meta,id,f);
});
viz.core.ghost_set_active_nodes_color = (function viz$core$ghost_set_active_nodes_color(state){
var color = quil.core.color.call(null,cljs.core.mod.call(null,new cljs.core.Keyword(null,"frame","frame",-1711082588).cljs$core$IFn$_invoke$arity$1(state),viz.core.frames_per_color_cycle.call(null,state)),(1),(1));
return cljs.core.reduce.call(null,((function (color){
return (function (state__$1,id){
return viz.core.update_node_meta.call(null,state__$1,id,((function (color){
return (function (p1__15517_SHARP_){
return cljs.core.assoc.call(null,p1__15517_SHARP_,new cljs.core.Keyword(null,"color","color",1011675173),color);
});})(color))
);
});})(color))
,state,cljs.core.get_in.call(null,state,new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"ghost","ghost",-1531157576),new cljs.core.Keyword(null,"active-node-ids","active-node-ids",-398210751)], null)));
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
var poss_fn = viz.core.mk_poss_fn.call(null,state);
return viz.core.maybe_exit.call(null,cljs.core.update_in.call(null,viz.core.maybe_remove_roots.call(null,viz.core.ghost_set_active_nodes_color.call(null,viz.core.ghost_incr.call(null,state,poss_fn))),new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"frame","frame",-1711082588)], null),cljs.core.inc));
});
viz.core.draw_ellipse = (function viz$core$draw_ellipse(state,pos,size){
var scaled_pos = viz.core.scale.call(null,state,pos);
var scaled_size = cljs.core.map.call(null,cljs.core.int$,viz.core.scale.call(null,state,size));
return cljs.core.apply.call(null,quil.core.ellipse,cljs.core.concat.call(null,scaled_pos,scaled_size));
});
viz.core.in_line_QMARK_ = (function viz$core$in_line_QMARK_(var_args){
var args__7934__auto__ = [];
var len__7927__auto___15520 = arguments.length;
var i__7928__auto___15521 = (0);
while(true){
if((i__7928__auto___15521 < len__7927__auto___15520)){
args__7934__auto__.push((arguments[i__7928__auto___15521]));

var G__15522 = (i__7928__auto___15521 + (1));
i__7928__auto___15521 = G__15522;
continue;
} else {
}
break;
}

var argseq__7935__auto__ = ((((0) < args__7934__auto__.length))?(new cljs.core.IndexedSeq(args__7934__auto__.slice((0)),(0),null)):null);
return viz.core.in_line_QMARK_.cljs$core$IFn$_invoke$arity$variadic(argseq__7935__auto__);
});

viz.core.in_line_QMARK_.cljs$core$IFn$_invoke$arity$variadic = (function (nodes){
return cljs.core.apply.call(null,cljs.core._EQ_,cljs.core.map.call(null,(function (p1__15518_SHARP_){
return cljs.core.apply.call(null,cljs.core.map,cljs.core._,p1__15518_SHARP_);
}),cljs.core.partition.call(null,(2),(1),cljs.core.map.call(null,new cljs.core.Keyword(null,"pos","pos",-864607220),nodes))));
});

viz.core.in_line_QMARK_.cljs$lang$maxFixedArity = (0);

viz.core.in_line_QMARK_.cljs$lang$applyTo = (function (seq15519){
return viz.core.in_line_QMARK_.cljs$core$IFn$_invoke$arity$variadic(cljs.core.seq.call(null,seq15519));
});

viz.core.draw_node = (function viz$core$draw_node(state,node,active_QMARK_){
var pos = new cljs.core.Keyword(null,"pos","pos",-864607220).cljs$core$IFn$_invoke$arity$1(node);
var stroke = cljs.core.get_in.call(null,node,new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"meta","meta",1499536964),new cljs.core.Keyword(null,"color","color",1011675173)], null));
var fill = (cljs.core.truth_(active_QMARK_)?stroke:(4294967295));
var size = new cljs.core.Keyword(null,"val","val",128701612).cljs$core$IFn$_invoke$arity$1(viz.dial.scaled.call(null,new cljs.core.Keyword(null,"dial","dial",1238392184).cljs$core$IFn$_invoke$arity$1(state),0.25,0.45));
quil.core.stroke.call(null,stroke);

quil.core.fill.call(null,fill);

return viz.core.draw_ellipse.call(null,state,pos,new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [size,size], null));
});
viz.core.draw_line = (function viz$core$draw_line(state,node,parent){
var node_color = cljs.core.get_in.call(null,node,new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"meta","meta",1499536964),new cljs.core.Keyword(null,"color","color",1011675173)], null));
var parent_color = cljs.core.get_in.call(null,node,new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"meta","meta",1499536964),new cljs.core.Keyword(null,"color","color",1011675173)], null));
var color = quil.core.lerp_color.call(null,node_color,parent_color,0.5);
var weight = new cljs.core.Keyword(null,"val","val",128701612).cljs$core$IFn$_invoke$arity$1(viz.dial.scaled.call(null,new cljs.core.Keyword(null,"dial","dial",1238392184).cljs$core$IFn$_invoke$arity$1(state),(-1),(3)));
quil.core.stroke.call(null,color);

quil.core.stroke_weight.call(null,weight);

return cljs.core.apply.call(null,quil.core.line,cljs.core.map.call(null,((function (node_color,parent_color,color,weight){
return (function (p1__15523_SHARP_){
return viz.core.scale.call(null,state,p1__15523_SHARP_);
});})(node_color,parent_color,color,weight))
,cljs.core.map.call(null,new cljs.core.Keyword(null,"pos","pos",-864607220),(function (){var x__7656__auto__ = parent;
return cljs.core._conj.call(null,(function (){var x__7656__auto____$1 = node;
return cljs.core._conj.call(null,cljs.core.List.EMPTY,x__7656__auto____$1);
})(),x__7656__auto__);
})())));
});
viz.core.draw_lines = (function viz$core$draw_lines(state,forest,parent,node){

var children = cljs.core.map.call(null,(function (p1__15524_SHARP_){
return viz.forest.get_node.call(null,forest,p1__15524_SHARP_);
}),new cljs.core.Keyword(null,"child-ids","child-ids",-604525861).cljs$core$IFn$_invoke$arity$1(node));
if(cljs.core.not.call(null,parent)){
var seq__15534_15542 = cljs.core.seq.call(null,children);
var chunk__15535_15543 = null;
var count__15536_15544 = (0);
var i__15537_15545 = (0);
while(true){
if((i__15537_15545 < count__15536_15544)){
var child_15546 = cljs.core._nth.call(null,chunk__15535_15543,i__15537_15545);
viz.core.draw_lines.call(null,state,forest,node,child_15546);

var G__15547 = seq__15534_15542;
var G__15548 = chunk__15535_15543;
var G__15549 = count__15536_15544;
var G__15550 = (i__15537_15545 + (1));
seq__15534_15542 = G__15547;
chunk__15535_15543 = G__15548;
count__15536_15544 = G__15549;
i__15537_15545 = G__15550;
continue;
} else {
var temp__4657__auto___15551 = cljs.core.seq.call(null,seq__15534_15542);
if(temp__4657__auto___15551){
var seq__15534_15552__$1 = temp__4657__auto___15551;
if(cljs.core.chunked_seq_QMARK_.call(null,seq__15534_15552__$1)){
var c__7633__auto___15553 = cljs.core.chunk_first.call(null,seq__15534_15552__$1);
var G__15554 = cljs.core.chunk_rest.call(null,seq__15534_15552__$1);
var G__15555 = c__7633__auto___15553;
var G__15556 = cljs.core.count.call(null,c__7633__auto___15553);
var G__15557 = (0);
seq__15534_15542 = G__15554;
chunk__15535_15543 = G__15555;
count__15536_15544 = G__15556;
i__15537_15545 = G__15557;
continue;
} else {
var child_15558 = cljs.core.first.call(null,seq__15534_15552__$1);
viz.core.draw_lines.call(null,state,forest,node,child_15558);

var G__15559 = cljs.core.next.call(null,seq__15534_15552__$1);
var G__15560 = null;
var G__15561 = (0);
var G__15562 = (0);
seq__15534_15542 = G__15559;
chunk__15535_15543 = G__15560;
count__15536_15544 = G__15561;
i__15537_15545 = G__15562;
continue;
}
} else {
}
}
break;
}
} else {
var in_line_child_15563 = cljs.core.some.call(null,((function (children){
return (function (p1__15525_SHARP_){
if(cljs.core.truth_(viz.core.in_line_QMARK_.call(null,parent,node,p1__15525_SHARP_))){
return p1__15525_SHARP_;
} else {
return null;
}
});})(children))
,children);
var seq__15538_15564 = cljs.core.seq.call(null,children);
var chunk__15539_15565 = null;
var count__15540_15566 = (0);
var i__15541_15567 = (0);
while(true){
if((i__15541_15567 < count__15540_15566)){
var child_15568 = cljs.core._nth.call(null,chunk__15539_15565,i__15541_15567);
if(cljs.core.truth_((function (){var and__6802__auto__ = in_line_child_15563;
if(cljs.core.truth_(and__6802__auto__)){
return cljs.core._EQ_.call(null,in_line_child_15563,child_15568);
} else {
return and__6802__auto__;
}
})())){
viz.core.draw_lines.call(null,state,forest,parent,child_15568);
} else {
viz.core.draw_lines.call(null,state,forest,node,child_15568);
}

var G__15569 = seq__15538_15564;
var G__15570 = chunk__15539_15565;
var G__15571 = count__15540_15566;
var G__15572 = (i__15541_15567 + (1));
seq__15538_15564 = G__15569;
chunk__15539_15565 = G__15570;
count__15540_15566 = G__15571;
i__15541_15567 = G__15572;
continue;
} else {
var temp__4657__auto___15573 = cljs.core.seq.call(null,seq__15538_15564);
if(temp__4657__auto___15573){
var seq__15538_15574__$1 = temp__4657__auto___15573;
if(cljs.core.chunked_seq_QMARK_.call(null,seq__15538_15574__$1)){
var c__7633__auto___15575 = cljs.core.chunk_first.call(null,seq__15538_15574__$1);
var G__15576 = cljs.core.chunk_rest.call(null,seq__15538_15574__$1);
var G__15577 = c__7633__auto___15575;
var G__15578 = cljs.core.count.call(null,c__7633__auto___15575);
var G__15579 = (0);
seq__15538_15564 = G__15576;
chunk__15539_15565 = G__15577;
count__15540_15566 = G__15578;
i__15541_15567 = G__15579;
continue;
} else {
var child_15580 = cljs.core.first.call(null,seq__15538_15574__$1);
if(cljs.core.truth_((function (){var and__6802__auto__ = in_line_child_15563;
if(cljs.core.truth_(and__6802__auto__)){
return cljs.core._EQ_.call(null,in_line_child_15563,child_15580);
} else {
return and__6802__auto__;
}
})())){
viz.core.draw_lines.call(null,state,forest,parent,child_15580);
} else {
viz.core.draw_lines.call(null,state,forest,node,child_15580);
}

var G__15581 = cljs.core.next.call(null,seq__15538_15574__$1);
var G__15582 = null;
var G__15583 = (0);
var G__15584 = (0);
seq__15538_15564 = G__15581;
chunk__15539_15565 = G__15582;
count__15540_15566 = G__15583;
i__15541_15567 = G__15584;
continue;
}
} else {
}
}
break;
}

if(cljs.core.truth_(in_line_child_15563)){
} else {
viz.core.draw_line.call(null,state,node,parent);
}
}

if(cljs.core.empty_QMARK_.call(null,children)){
return viz.core.draw_node.call(null,state,node,false);
} else {
return null;
}
});
viz.core.draw_dial = (function viz$core$draw_dial(state,dial,posL,posR){
var dial_norm = quil.core.norm.call(null,new cljs.core.Keyword(null,"val","val",128701612).cljs$core$IFn$_invoke$arity$1(dial),new cljs.core.Keyword(null,"min","min",444991522).cljs$core$IFn$_invoke$arity$1(dial),new cljs.core.Keyword(null,"max","max",61366548).cljs$core$IFn$_invoke$arity$1(dial));
var dial_pos = cljs.core.map.call(null,((function (dial_norm){
return (function (p1__15585_SHARP_,p2__15586_SHARP_){
return quil.core.lerp.call(null,p1__15585_SHARP_,p2__15586_SHARP_,dial_norm);
});})(dial_norm))
,posL,posR);
quil.core.stroke.call(null,(4278190080));

quil.core.stroke_weight.call(null,(1));

quil.core.fill.call(null,(4278190080));

cljs.core.apply.call(null,quil.core.line,cljs.core.concat.call(null,posL,posR));

return cljs.core.apply.call(null,quil.core.ellipse,cljs.core.concat.call(null,dial_pos,new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [(5),(5)], null)));
});
viz.core.draw_state = (function viz$core$draw_state(state){
quil.core.background.call(null,(4294967295));

var tr__8398__auto__ = new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [(viz.core.window_size.call(null,(0)) / (2)),(viz.core.window_size.call(null,(1)) / (2))], null);
quil.core.push_matrix.call(null);

try{quil.core.translate.call(null,tr__8398__auto__);

var lines = viz.forest.lines.call(null,new cljs.core.Keyword(null,"forest","forest",278860306).cljs$core$IFn$_invoke$arity$1(state));
var leaves = viz.forest.leaves.call(null,new cljs.core.Keyword(null,"forest","forest",278860306).cljs$core$IFn$_invoke$arity$1(state));
var active = viz.ghost.active_nodes.call(null,new cljs.core.Keyword(null,"ghost","ghost",-1531157576).cljs$core$IFn$_invoke$arity$1(state),new cljs.core.Keyword(null,"forest","forest",278860306).cljs$core$IFn$_invoke$arity$1(state));
var roots = viz.forest.roots.call(null,new cljs.core.Keyword(null,"forest","forest",278860306).cljs$core$IFn$_invoke$arity$1(state));
var seq__15595_15603 = cljs.core.seq.call(null,roots);
var chunk__15596_15604 = null;
var count__15597_15605 = (0);
var i__15598_15606 = (0);
while(true){
if((i__15598_15606 < count__15597_15605)){
var root_15607 = cljs.core._nth.call(null,chunk__15596_15604,i__15598_15606);
viz.core.draw_lines.call(null,state,new cljs.core.Keyword(null,"forest","forest",278860306).cljs$core$IFn$_invoke$arity$1(state),null,root_15607);

var G__15608 = seq__15595_15603;
var G__15609 = chunk__15596_15604;
var G__15610 = count__15597_15605;
var G__15611 = (i__15598_15606 + (1));
seq__15595_15603 = G__15608;
chunk__15596_15604 = G__15609;
count__15597_15605 = G__15610;
i__15598_15606 = G__15611;
continue;
} else {
var temp__4657__auto___15612 = cljs.core.seq.call(null,seq__15595_15603);
if(temp__4657__auto___15612){
var seq__15595_15613__$1 = temp__4657__auto___15612;
if(cljs.core.chunked_seq_QMARK_.call(null,seq__15595_15613__$1)){
var c__7633__auto___15614 = cljs.core.chunk_first.call(null,seq__15595_15613__$1);
var G__15615 = cljs.core.chunk_rest.call(null,seq__15595_15613__$1);
var G__15616 = c__7633__auto___15614;
var G__15617 = cljs.core.count.call(null,c__7633__auto___15614);
var G__15618 = (0);
seq__15595_15603 = G__15615;
chunk__15596_15604 = G__15616;
count__15597_15605 = G__15617;
i__15598_15606 = G__15618;
continue;
} else {
var root_15619 = cljs.core.first.call(null,seq__15595_15613__$1);
viz.core.draw_lines.call(null,state,new cljs.core.Keyword(null,"forest","forest",278860306).cljs$core$IFn$_invoke$arity$1(state),null,root_15619);

var G__15620 = cljs.core.next.call(null,seq__15595_15613__$1);
var G__15621 = null;
var G__15622 = (0);
var G__15623 = (0);
seq__15595_15603 = G__15620;
chunk__15596_15604 = G__15621;
count__15597_15605 = G__15622;
i__15598_15606 = G__15623;
continue;
}
} else {
}
}
break;
}

var seq__15599 = cljs.core.seq.call(null,active);
var chunk__15600 = null;
var count__15601 = (0);
var i__15602 = (0);
while(true){
if((i__15602 < count__15601)){
var active_node = cljs.core._nth.call(null,chunk__15600,i__15602);
viz.core.draw_node.call(null,state,active_node,true);

var G__15624 = seq__15599;
var G__15625 = chunk__15600;
var G__15626 = count__15601;
var G__15627 = (i__15602 + (1));
seq__15599 = G__15624;
chunk__15600 = G__15625;
count__15601 = G__15626;
i__15602 = G__15627;
continue;
} else {
var temp__4657__auto__ = cljs.core.seq.call(null,seq__15599);
if(temp__4657__auto__){
var seq__15599__$1 = temp__4657__auto__;
if(cljs.core.chunked_seq_QMARK_.call(null,seq__15599__$1)){
var c__7633__auto__ = cljs.core.chunk_first.call(null,seq__15599__$1);
var G__15628 = cljs.core.chunk_rest.call(null,seq__15599__$1);
var G__15629 = c__7633__auto__;
var G__15630 = cljs.core.count.call(null,c__7633__auto__);
var G__15631 = (0);
seq__15599 = G__15628;
chunk__15600 = G__15629;
count__15601 = G__15630;
i__15602 = G__15631;
continue;
} else {
var active_node = cljs.core.first.call(null,seq__15599__$1);
viz.core.draw_node.call(null,state,active_node,true);

var G__15632 = cljs.core.next.call(null,seq__15599__$1);
var G__15633 = null;
var G__15634 = (0);
var G__15635 = (0);
seq__15599 = G__15632;
chunk__15600 = G__15633;
count__15601 = G__15634;
i__15602 = G__15635;
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
var G__15636__delegate = function (args){
return cljs.core.apply.call(null,viz.core.update_state,args);
};
var G__15636 = function (var_args){
var args = null;
if (arguments.length > 0) {
var G__15637__i = 0, G__15637__a = new Array(arguments.length -  0);
while (G__15637__i < G__15637__a.length) {G__15637__a[G__15637__i] = arguments[G__15637__i + 0]; ++G__15637__i;}
  args = new cljs.core.IndexedSeq(G__15637__a,0);
} 
return G__15636__delegate.call(this,args);};
G__15636.cljs$lang$maxFixedArity = 0;
G__15636.cljs$lang$applyTo = (function (arglist__15638){
var args = cljs.core.seq(arglist__15638);
return G__15636__delegate(args);
});
G__15636.cljs$core$IFn$_invoke$arity$variadic = G__15636__delegate;
return G__15636;
})()
:viz.core.update_state),new cljs.core.Keyword(null,"size","size",1098693007),((cljs.core.fn_QMARK_.call(null,viz.core.window_size))?(function() { 
var G__15639__delegate = function (args){
return cljs.core.apply.call(null,viz.core.window_size,args);
};
var G__15639 = function (var_args){
var args = null;
if (arguments.length > 0) {
var G__15640__i = 0, G__15640__a = new Array(arguments.length -  0);
while (G__15640__i < G__15640__a.length) {G__15640__a[G__15640__i] = arguments[G__15640__i + 0]; ++G__15640__i;}
  args = new cljs.core.IndexedSeq(G__15640__a,0);
} 
return G__15639__delegate.call(this,args);};
G__15639.cljs$lang$maxFixedArity = 0;
G__15639.cljs$lang$applyTo = (function (arglist__15641){
var args = cljs.core.seq(arglist__15641);
return G__15639__delegate(args);
});
G__15639.cljs$core$IFn$_invoke$arity$variadic = G__15639__delegate;
return G__15639;
})()
:viz.core.window_size),new cljs.core.Keyword(null,"title","title",636505583),"",new cljs.core.Keyword(null,"setup","setup",1987730512),((cljs.core.fn_QMARK_.call(null,viz.core.setup))?(function() { 
var G__15642__delegate = function (args){
return cljs.core.apply.call(null,viz.core.setup,args);
};
var G__15642 = function (var_args){
var args = null;
if (arguments.length > 0) {
var G__15643__i = 0, G__15643__a = new Array(arguments.length -  0);
while (G__15643__i < G__15643__a.length) {G__15643__a[G__15643__i] = arguments[G__15643__i + 0]; ++G__15643__i;}
  args = new cljs.core.IndexedSeq(G__15643__a,0);
} 
return G__15642__delegate.call(this,args);};
G__15642.cljs$lang$maxFixedArity = 0;
G__15642.cljs$lang$applyTo = (function (arglist__15644){
var args = cljs.core.seq(arglist__15644);
return G__15642__delegate(args);
});
G__15642.cljs$core$IFn$_invoke$arity$variadic = G__15642__delegate;
return G__15642;
})()
:viz.core.setup),new cljs.core.Keyword(null,"middleware","middleware",1462115504),new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [quil.middleware.fun_mode], null),new cljs.core.Keyword(null,"draw","draw",1358331674),((cljs.core.fn_QMARK_.call(null,viz.core.draw_state))?(function() { 
var G__15645__delegate = function (args){
return cljs.core.apply.call(null,viz.core.draw_state,args);
};
var G__15645 = function (var_args){
var args = null;
if (arguments.length > 0) {
var G__15646__i = 0, G__15646__a = new Array(arguments.length -  0);
while (G__15646__i < G__15646__a.length) {G__15646__a[G__15646__i] = arguments[G__15646__i + 0]; ++G__15646__i;}
  args = new cljs.core.IndexedSeq(G__15646__a,0);
} 
return G__15645__delegate.call(this,args);};
G__15645.cljs$lang$maxFixedArity = 0;
G__15645.cljs$lang$applyTo = (function (arglist__15647){
var args = cljs.core.seq(arglist__15647);
return G__15645__delegate(args);
});
G__15645.cljs$core$IFn$_invoke$arity$variadic = G__15645__delegate;
return G__15645;
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