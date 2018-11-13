// Compiled by ClojureScript 1.10.439 {}
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
var args__4647__auto__ = [];
var len__4641__auto___25578 = arguments.length;
var i__4642__auto___25579 = (0);
while(true){
if((i__4642__auto___25579 < len__4641__auto___25578)){
args__4647__auto__.push((arguments[i__4642__auto___25579]));

var G__25580 = (i__4642__auto___25579 + (1));
i__4642__auto___25579 = G__25580;
continue;
} else {
}
break;
}

var argseq__4648__auto__ = ((((0) < args__4647__auto__.length))?(new cljs.core.IndexedSeq(args__4647__auto__.slice((0)),(0),null)):null);
return viz.core.debug.cljs$core$IFn$_invoke$arity$variadic(argseq__4648__auto__);
});

viz.core.debug.cljs$core$IFn$_invoke$arity$variadic = (function (args){
return console.log(clojure.string.join.call(null," ",cljs.core.map.call(null,cljs.core.str,args)));
});

viz.core.debug.cljs$lang$maxFixedArity = (0);

/** @this {Function} */
viz.core.debug.cljs$lang$applyTo = (function (seq25577){
var self__4629__auto__ = this;
return self__4629__auto__.cljs$core$IFn$_invoke$arity$variadic(cljs.core.seq.call(null,seq25577));
});

viz.core.observe = (function viz$core$observe(v){
viz.core.debug.call(null,v);

return v;
});
viz.core.positive = (function viz$core$positive(n){
if(((0) > n)){
return (- n);
} else {
return n;
}
});
viz.core.window_partial = (function viz$core$window_partial(k){
return ((document["documentElement"][k]) | (0));
});
viz.core.window_size = (function (){var w = ((function (){var x__4138__auto__ = (1024);
var y__4139__auto__ = viz.core.window_partial.call(null,"clientWidth");
return ((x__4138__auto__ < y__4139__auto__) ? x__4138__auto__ : y__4139__auto__);
})() | (0));
return new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [w,((function (){var x__4138__auto__ = (w * 0.75);
var y__4139__auto__ = viz.core.window_partial.call(null,"clientHeight");
return ((x__4138__auto__ < y__4139__auto__) ? x__4138__auto__ : y__4139__auto__);
})() | (0))], null);
})();
viz.core.window_half_size = cljs.core.apply.call(null,cljs.core.vector,cljs.core.map.call(null,(function (p1__25581_SHARP_){
return (p1__25581_SHARP_ / (2));
}),viz.core.window_size));
viz.core.set_grid_size = (function viz$core$set_grid_size(state){
var h = ((viz.core.window_size.call(null,(1)) * (new cljs.core.Keyword(null,"grid-width","grid-width",837583106).cljs$core$IFn$_invoke$arity$1(state) / viz.core.window_size.call(null,(0)))) | (0));
return cljs.core.assoc.call(null,state,new cljs.core.Keyword(null,"grid-size","grid-size",2138480144),new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"grid-width","grid-width",837583106).cljs$core$IFn$_invoke$arity$1(state),h], null));
});
viz.core.add_ghost = (function viz$core$add_ghost(state,ghost_def){
var vec__25582 = viz.forest.add_node.call(null,new cljs.core.Keyword(null,"forest","forest",278860306).cljs$core$IFn$_invoke$arity$1(state),new cljs.core.Keyword(null,"start-pos","start-pos",668789086).cljs$core$IFn$_invoke$arity$1(ghost_def));
var forest = cljs.core.nth.call(null,vec__25582,(0),null);
var id = cljs.core.nth.call(null,vec__25582,(1),null);
var ghost = cljs.core.assoc.call(null,viz.ghost.add_active_node.call(null,viz.ghost.new_ghost.call(null),id),new cljs.core.Keyword(null,"ghost-def","ghost-def",1211539367),ghost_def);
return cljs.core.assoc.call(null,state,new cljs.core.Keyword(null,"forest","forest",278860306),forest,new cljs.core.Keyword(null,"ghosts","ghosts",665819293),cljs.core.cons.call(null,ghost,new cljs.core.Keyword(null,"ghosts","ghosts",665819293).cljs$core$IFn$_invoke$arity$1(state)));
});
viz.core.new_state = (function viz$core$new_state(){
return viz.core.add_ghost.call(null,viz.core.set_grid_size.call(null,new cljs.core.PersistentArrayMap(null, 6, [new cljs.core.Keyword(null,"frame-rate","frame-rate",-994918942),(15),new cljs.core.Keyword(null,"color-cycle-period","color-cycle-period",1656886882),(8),new cljs.core.Keyword(null,"tail-length","tail-length",-2007115089),(7),new cljs.core.Keyword(null,"frame","frame",-1711082588),(0),new cljs.core.Keyword(null,"grid-width","grid-width",837583106),(45),new cljs.core.Keyword(null,"forest","forest",278860306),viz.forest.new_forest.call(null,viz.grid.isometric)], null)),new cljs.core.PersistentArrayMap(null, 2, [new cljs.core.Keyword(null,"start-pos","start-pos",668789086),new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [(-10),(-10)], null),new cljs.core.Keyword(null,"color-fn","color-fn",1518098073),(function (state){
var frames_per_color_cycle = (new cljs.core.Keyword(null,"color-cycle-period","color-cycle-period",1656886882).cljs$core$IFn$_invoke$arity$1(state) * new cljs.core.Keyword(null,"frame-rate","frame-rate",-994918942).cljs$core$IFn$_invoke$arity$1(state));
return quil.core.color.call(null,(cljs.core.mod.call(null,new cljs.core.Keyword(null,"frame","frame",-1711082588).cljs$core$IFn$_invoke$arity$1(state),frames_per_color_cycle) / frames_per_color_cycle),(1),(1));
})], null));
});
viz.core.setup = (function viz$core$setup(){
quil.core.color_mode.call(null,new cljs.core.Keyword(null,"hsb","hsb",-753472031),(1),(1),(1));

return viz.core.new_state.call(null);
});
viz.core.curr_second = (function viz$core$curr_second(state){
return (new cljs.core.Keyword(null,"frame","frame",-1711082588).cljs$core$IFn$_invoke$arity$1(state) / new cljs.core.Keyword(null,"frame-rate","frame-rate",-994918942).cljs$core$IFn$_invoke$arity$1(state));
});
viz.core.scale = (function viz$core$scale(grid_size,xy){
return cljs.core.map_indexed.call(null,(function (p1__25586_SHARP_,p2__25585_SHARP_){
return (p2__25585_SHARP_ * (viz.core.window_half_size.call(null,p1__25586_SHARP_) / grid_size.call(null,p1__25586_SHARP_)));
}),xy);
});
viz.core.bounds_buffer = (1);
viz.core.in_bounds_QMARK_ = (function viz$core$in_bounds_QMARK_(grid_size,pos){
var vec__25589 = cljs.core.apply.call(null,cljs.core.vector,cljs.core.map.call(null,(function (p1__25587_SHARP_){
return (p1__25587_SHARP_ - viz.core.bounds_buffer);
}),grid_size));
var w = cljs.core.nth.call(null,vec__25589,(0),null);
var h = cljs.core.nth.call(null,vec__25589,(1),null);
return cljs.core.every_QMARK_.call(null,((function (vec__25589,w,h){
return (function (p1__25588_SHARP_){
return (((p1__25588_SHARP_.call(null,(1)) >= (- p1__25588_SHARP_.call(null,(0))))) && ((p1__25588_SHARP_.call(null,(1)) <= p1__25588_SHARP_.call(null,(0)))));
});})(vec__25589,w,h))
,cljs.core.map.call(null,cljs.core.vector,new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [w,h], null),pos));
});
viz.core.dist_from_sqr = (function viz$core$dist_from_sqr(pos1,pos2){
return cljs.core.reduce.call(null,cljs.core._PLUS_,cljs.core.map.call(null,(function (p1__25592_SHARP_){
return (p1__25592_SHARP_ * p1__25592_SHARP_);
}),cljs.core.map.call(null,cljs.core._,pos1,pos2)));
});
viz.core.dist_from = (function viz$core$dist_from(pos1,pos2){
return quil.core.sqrt.call(null,viz.core.dist_from_sqr.call(null,pos1,pos2));
});
viz.core.take_adj_poss = (function viz$core$take_adj_poss(grid_width,pos,adj_poss){
var dist_from_center = viz.core.dist_from.call(null,new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [(0),(0)], null),pos);
var width = grid_width;
var dist_ratio = ((width - dist_from_center) / width);
return cljs.core.take.call(null,(((quil.core.map_range.call(null,cljs.core.rand.call(null),(0),(1),0.75,(1)) * dist_ratio) * cljs.core.count.call(null,adj_poss)) | (0)),adj_poss);
});
viz.core.mk_poss_fn = (function viz$core$mk_poss_fn(state){
var grid_size = new cljs.core.Keyword(null,"grid-size","grid-size",2138480144).cljs$core$IFn$_invoke$arity$1(state);
return ((function (grid_size){
return (function (pos,adj_poss){
return viz.core.take_adj_poss.call(null,grid_size.call(null,(0)),pos,cljs.core.sort_by.call(null,((function (grid_size){
return (function (p1__25594_SHARP_){
return viz.core.dist_from_sqr.call(null,p1__25594_SHARP_,new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [(0),(0)], null));
});})(grid_size))
,cljs.core.filter.call(null,((function (grid_size){
return (function (p1__25593_SHARP_){
return viz.core.in_bounds_QMARK_.call(null,grid_size,p1__25593_SHARP_);
});})(grid_size))
,adj_poss)));
});
;})(grid_size))
});
viz.core.update_ghost_forest = (function viz$core$update_ghost_forest(state,update_fn){
var vec__25595 = cljs.core.reduce.call(null,(function (p__25598,ghost){
var vec__25599 = p__25598;
var ghosts = cljs.core.nth.call(null,vec__25599,(0),null);
var forest = cljs.core.nth.call(null,vec__25599,(1),null);
var vec__25602 = update_fn.call(null,ghost,forest);
var ghost__$1 = cljs.core.nth.call(null,vec__25602,(0),null);
var forest__$1 = cljs.core.nth.call(null,vec__25602,(1),null);
return new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [cljs.core.cons.call(null,ghost__$1,ghosts),forest__$1], null);
}),new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [null,new cljs.core.Keyword(null,"forest","forest",278860306).cljs$core$IFn$_invoke$arity$1(state)], null),new cljs.core.Keyword(null,"ghosts","ghosts",665819293).cljs$core$IFn$_invoke$arity$1(state));
var ghosts = cljs.core.nth.call(null,vec__25595,(0),null);
var forest = cljs.core.nth.call(null,vec__25595,(1),null);
return cljs.core.assoc.call(null,state,new cljs.core.Keyword(null,"ghosts","ghosts",665819293),cljs.core.reverse.call(null,ghosts),new cljs.core.Keyword(null,"forest","forest",278860306),forest);
});
viz.core.ghost_incr = (function viz$core$ghost_incr(state,poss_fn){
return viz.core.update_ghost_forest.call(null,state,(function (p1__25605_SHARP_,p2__25606_SHARP_){
return viz.ghost.incr.call(null,p1__25605_SHARP_,p2__25606_SHARP_,poss_fn);
}));
});
viz.core.rm_nodes = (function viz$core$rm_nodes(state,node_ids){
return viz.core.update_ghost_forest.call(null,state,(function (ghost,forest){
return new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [cljs.core.reduce.call(null,viz.ghost.rm_active_node,ghost,node_ids),cljs.core.reduce.call(null,viz.forest.remove_node,forest,node_ids)], null);
}));
});
viz.core.maybe_remove_roots = (function viz$core$maybe_remove_roots(state){
if((new cljs.core.Keyword(null,"tail-length","tail-length",-2007115089).cljs$core$IFn$_invoke$arity$1(state) >= new cljs.core.Keyword(null,"frame","frame",-1711082588).cljs$core$IFn$_invoke$arity$1(state))){
return state;
} else {
return viz.core.rm_nodes.call(null,state,cljs.core.map.call(null,new cljs.core.Keyword(null,"id","id",-1388402092),viz.forest.roots.call(null,new cljs.core.Keyword(null,"forest","forest",278860306).cljs$core$IFn$_invoke$arity$1(state))));
}
});
viz.core.ghost_set_color = (function viz$core$ghost_set_color(state){
return viz.core.update_ghost_forest.call(null,state,(function (ghost,forest){
var color = cljs.core.get_in.call(null,ghost,new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"ghost-def","ghost-def",1211539367),new cljs.core.Keyword(null,"color-fn","color-fn",1518098073)], null)).call(null,state);
return new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [cljs.core.assoc.call(null,ghost,new cljs.core.Keyword(null,"color","color",1011675173),color),forest], null);
}));
});
viz.core.update_state = (function viz$core$update_state(state){
var poss_fn = viz.core.mk_poss_fn.call(null,state);
return cljs.core.update_in.call(null,viz.core.maybe_remove_roots.call(null,viz.core.ghost_incr.call(null,viz.core.ghost_set_color.call(null,state),poss_fn)),new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"frame","frame",-1711082588)], null),cljs.core.inc);
});
viz.core.draw_ellipse = (function viz$core$draw_ellipse(pos,size,scale_fn){
var scaled_pos = scale_fn.call(null,pos);
var scaled_size = cljs.core.map.call(null,cljs.core.int$,scale_fn.call(null,size));
return cljs.core.apply.call(null,quil.core.ellipse,cljs.core.concat.call(null,scaled_pos,scaled_size));
});
viz.core.in_line_QMARK_ = (function viz$core$in_line_QMARK_(var_args){
var args__4647__auto__ = [];
var len__4641__auto___25609 = arguments.length;
var i__4642__auto___25610 = (0);
while(true){
if((i__4642__auto___25610 < len__4641__auto___25609)){
args__4647__auto__.push((arguments[i__4642__auto___25610]));

var G__25611 = (i__4642__auto___25610 + (1));
i__4642__auto___25610 = G__25611;
continue;
} else {
}
break;
}

var argseq__4648__auto__ = ((((0) < args__4647__auto__.length))?(new cljs.core.IndexedSeq(args__4647__auto__.slice((0)),(0),null)):null);
return viz.core.in_line_QMARK_.cljs$core$IFn$_invoke$arity$variadic(argseq__4648__auto__);
});

viz.core.in_line_QMARK_.cljs$core$IFn$_invoke$arity$variadic = (function (nodes){
return cljs.core.apply.call(null,cljs.core._EQ_,cljs.core.map.call(null,(function (p1__25607_SHARP_){
return cljs.core.apply.call(null,cljs.core.map,cljs.core._,p1__25607_SHARP_);
}),cljs.core.partition.call(null,(2),(1),cljs.core.map.call(null,new cljs.core.Keyword(null,"pos","pos",-864607220),nodes))));
});

viz.core.in_line_QMARK_.cljs$lang$maxFixedArity = (0);

/** @this {Function} */
viz.core.in_line_QMARK_.cljs$lang$applyTo = (function (seq25608){
var self__4629__auto__ = this;
return self__4629__auto__.cljs$core$IFn$_invoke$arity$variadic(cljs.core.seq.call(null,seq25608));
});

viz.core.draw_node = (function viz$core$draw_node(node,active_QMARK_,scale_fn){
var pos = new cljs.core.Keyword(null,"pos","pos",-864607220).cljs$core$IFn$_invoke$arity$1(node);
var stroke = cljs.core.get_in.call(null,node,new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"meta","meta",1499536964),new cljs.core.Keyword(null,"color","color",1011675173)], null));
var fill = (cljs.core.truth_(active_QMARK_)?stroke:(4294967295));
quil.core.stroke.call(null,stroke);

quil.core.fill.call(null,fill);

return viz.core.draw_ellipse.call(null,pos,new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [0.3,0.3], null),scale_fn);
});
viz.core.draw_line = (function viz$core$draw_line(node,parent,scale_fn){
var node_color = cljs.core.get_in.call(null,node,new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"meta","meta",1499536964),new cljs.core.Keyword(null,"color","color",1011675173)], null));
var parent_color = cljs.core.get_in.call(null,node,new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"meta","meta",1499536964),new cljs.core.Keyword(null,"color","color",1011675173)], null));
var color = quil.core.lerp_color.call(null,node_color,parent_color,0.5);
quil.core.stroke.call(null,color);

quil.core.stroke_weight.call(null,(1));

return cljs.core.apply.call(null,quil.core.line,cljs.core.map.call(null,scale_fn,cljs.core.map.call(null,new cljs.core.Keyword(null,"pos","pos",-864607220),(new cljs.core.List(null,parent,(new cljs.core.List(null,node,null,(1),null)),(2),null)))));
});
viz.core.draw_lines = (function viz$core$draw_lines(forest,parent,node,scale_fn){

var children = cljs.core.map.call(null,(function (p1__25612_SHARP_){
return viz.forest.get_node.call(null,forest,p1__25612_SHARP_);
}),new cljs.core.Keyword(null,"child-ids","child-ids",-604525861).cljs$core$IFn$_invoke$arity$1(node));
if(cljs.core.not.call(null,parent)){
var seq__25614_25622 = cljs.core.seq.call(null,children);
var chunk__25615_25623 = null;
var count__25616_25624 = (0);
var i__25617_25625 = (0);
while(true){
if((i__25617_25625 < count__25616_25624)){
var child_25626 = cljs.core._nth.call(null,chunk__25615_25623,i__25617_25625);
viz.core.draw_lines.call(null,forest,node,child_25626,scale_fn);


var G__25627 = seq__25614_25622;
var G__25628 = chunk__25615_25623;
var G__25629 = count__25616_25624;
var G__25630 = (i__25617_25625 + (1));
seq__25614_25622 = G__25627;
chunk__25615_25623 = G__25628;
count__25616_25624 = G__25629;
i__25617_25625 = G__25630;
continue;
} else {
var temp__4657__auto___25631 = cljs.core.seq.call(null,seq__25614_25622);
if(temp__4657__auto___25631){
var seq__25614_25632__$1 = temp__4657__auto___25631;
if(cljs.core.chunked_seq_QMARK_.call(null,seq__25614_25632__$1)){
var c__4461__auto___25633 = cljs.core.chunk_first.call(null,seq__25614_25632__$1);
var G__25634 = cljs.core.chunk_rest.call(null,seq__25614_25632__$1);
var G__25635 = c__4461__auto___25633;
var G__25636 = cljs.core.count.call(null,c__4461__auto___25633);
var G__25637 = (0);
seq__25614_25622 = G__25634;
chunk__25615_25623 = G__25635;
count__25616_25624 = G__25636;
i__25617_25625 = G__25637;
continue;
} else {
var child_25638 = cljs.core.first.call(null,seq__25614_25632__$1);
viz.core.draw_lines.call(null,forest,node,child_25638,scale_fn);


var G__25639 = cljs.core.next.call(null,seq__25614_25632__$1);
var G__25640 = null;
var G__25641 = (0);
var G__25642 = (0);
seq__25614_25622 = G__25639;
chunk__25615_25623 = G__25640;
count__25616_25624 = G__25641;
i__25617_25625 = G__25642;
continue;
}
} else {
}
}
break;
}
} else {
var in_line_child_25643 = cljs.core.some.call(null,((function (children){
return (function (p1__25613_SHARP_){
if(cljs.core.truth_(viz.core.in_line_QMARK_.call(null,parent,node,p1__25613_SHARP_))){
return p1__25613_SHARP_;
} else {
return null;
}
});})(children))
,children);
var seq__25618_25644 = cljs.core.seq.call(null,children);
var chunk__25619_25645 = null;
var count__25620_25646 = (0);
var i__25621_25647 = (0);
while(true){
if((i__25621_25647 < count__25620_25646)){
var child_25648 = cljs.core._nth.call(null,chunk__25619_25645,i__25621_25647);
if(cljs.core.truth_((function (){var and__4036__auto__ = in_line_child_25643;
if(cljs.core.truth_(and__4036__auto__)){
return cljs.core._EQ_.call(null,in_line_child_25643,child_25648);
} else {
return and__4036__auto__;
}
})())){
viz.core.draw_lines.call(null,forest,parent,child_25648,scale_fn);
} else {
viz.core.draw_lines.call(null,forest,node,child_25648,scale_fn);
}


var G__25649 = seq__25618_25644;
var G__25650 = chunk__25619_25645;
var G__25651 = count__25620_25646;
var G__25652 = (i__25621_25647 + (1));
seq__25618_25644 = G__25649;
chunk__25619_25645 = G__25650;
count__25620_25646 = G__25651;
i__25621_25647 = G__25652;
continue;
} else {
var temp__4657__auto___25653 = cljs.core.seq.call(null,seq__25618_25644);
if(temp__4657__auto___25653){
var seq__25618_25654__$1 = temp__4657__auto___25653;
if(cljs.core.chunked_seq_QMARK_.call(null,seq__25618_25654__$1)){
var c__4461__auto___25655 = cljs.core.chunk_first.call(null,seq__25618_25654__$1);
var G__25656 = cljs.core.chunk_rest.call(null,seq__25618_25654__$1);
var G__25657 = c__4461__auto___25655;
var G__25658 = cljs.core.count.call(null,c__4461__auto___25655);
var G__25659 = (0);
seq__25618_25644 = G__25656;
chunk__25619_25645 = G__25657;
count__25620_25646 = G__25658;
i__25621_25647 = G__25659;
continue;
} else {
var child_25660 = cljs.core.first.call(null,seq__25618_25654__$1);
if(cljs.core.truth_((function (){var and__4036__auto__ = in_line_child_25643;
if(cljs.core.truth_(and__4036__auto__)){
return cljs.core._EQ_.call(null,in_line_child_25643,child_25660);
} else {
return and__4036__auto__;
}
})())){
viz.core.draw_lines.call(null,forest,parent,child_25660,scale_fn);
} else {
viz.core.draw_lines.call(null,forest,node,child_25660,scale_fn);
}


var G__25661 = cljs.core.next.call(null,seq__25618_25654__$1);
var G__25662 = null;
var G__25663 = (0);
var G__25664 = (0);
seq__25618_25644 = G__25661;
chunk__25619_25645 = G__25662;
count__25620_25646 = G__25663;
i__25621_25647 = G__25664;
continue;
}
} else {
}
}
break;
}

if(cljs.core.truth_(in_line_child_25643)){
} else {
viz.core.draw_line.call(null,node,parent,scale_fn);
}
}

if(cljs.core.empty_QMARK_.call(null,children)){
return viz.core.draw_node.call(null,node,false,scale_fn);
} else {
return null;
}
});
viz.core.draw_dial = (function viz$core$draw_dial(state,dial,posL,posR){
var dial_norm = quil.core.norm.call(null,new cljs.core.Keyword(null,"val","val",128701612).cljs$core$IFn$_invoke$arity$1(dial),new cljs.core.Keyword(null,"min","min",444991522).cljs$core$IFn$_invoke$arity$1(dial),new cljs.core.Keyword(null,"max","max",61366548).cljs$core$IFn$_invoke$arity$1(dial));
var dial_pos = cljs.core.map.call(null,((function (dial_norm){
return (function (p1__25665_SHARP_,p2__25666_SHARP_){
return quil.core.lerp.call(null,p1__25665_SHARP_,p2__25666_SHARP_,dial_norm);
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

var tr__1504__auto__ = viz.core.window_half_size;
quil.core.push_matrix.call(null);

try{quil.core.translate.call(null,tr__1504__auto__);

var grid_size = new cljs.core.Keyword(null,"grid-size","grid-size",2138480144).cljs$core$IFn$_invoke$arity$1(state);
var scale_fn = ((function (grid_size,tr__1504__auto__){
return (function (p1__25667_SHARP_){
return viz.core.scale.call(null,grid_size,p1__25667_SHARP_);
});})(grid_size,tr__1504__auto__))
;
var ghost = new cljs.core.Keyword(null,"ghost","ghost",-1531157576).cljs$core$IFn$_invoke$arity$1(state);
var forest = new cljs.core.Keyword(null,"forest","forest",278860306).cljs$core$IFn$_invoke$arity$1(state);
var roots = viz.forest.roots.call(null,forest);
var seq__25669_25685 = cljs.core.seq.call(null,roots);
var chunk__25670_25686 = null;
var count__25671_25687 = (0);
var i__25672_25688 = (0);
while(true){
if((i__25672_25688 < count__25671_25687)){
var root_25689 = cljs.core._nth.call(null,chunk__25670_25686,i__25672_25688);
viz.core.draw_lines.call(null,forest,null,root_25689,scale_fn);


var G__25690 = seq__25669_25685;
var G__25691 = chunk__25670_25686;
var G__25692 = count__25671_25687;
var G__25693 = (i__25672_25688 + (1));
seq__25669_25685 = G__25690;
chunk__25670_25686 = G__25691;
count__25671_25687 = G__25692;
i__25672_25688 = G__25693;
continue;
} else {
var temp__4657__auto___25694 = cljs.core.seq.call(null,seq__25669_25685);
if(temp__4657__auto___25694){
var seq__25669_25695__$1 = temp__4657__auto___25694;
if(cljs.core.chunked_seq_QMARK_.call(null,seq__25669_25695__$1)){
var c__4461__auto___25696 = cljs.core.chunk_first.call(null,seq__25669_25695__$1);
var G__25697 = cljs.core.chunk_rest.call(null,seq__25669_25695__$1);
var G__25698 = c__4461__auto___25696;
var G__25699 = cljs.core.count.call(null,c__4461__auto___25696);
var G__25700 = (0);
seq__25669_25685 = G__25697;
chunk__25670_25686 = G__25698;
count__25671_25687 = G__25699;
i__25672_25688 = G__25700;
continue;
} else {
var root_25701 = cljs.core.first.call(null,seq__25669_25695__$1);
viz.core.draw_lines.call(null,forest,null,root_25701,scale_fn);


var G__25702 = cljs.core.next.call(null,seq__25669_25695__$1);
var G__25703 = null;
var G__25704 = (0);
var G__25705 = (0);
seq__25669_25685 = G__25702;
chunk__25670_25686 = G__25703;
count__25671_25687 = G__25704;
i__25672_25688 = G__25705;
continue;
}
} else {
}
}
break;
}

var seq__25673 = cljs.core.seq.call(null,new cljs.core.Keyword(null,"ghosts","ghosts",665819293).cljs$core$IFn$_invoke$arity$1(state));
var chunk__25674 = null;
var count__25675 = (0);
var i__25676 = (0);
while(true){
if((i__25676 < count__25675)){
var ghost__$1 = cljs.core._nth.call(null,chunk__25674,i__25676);
var seq__25677_25706 = cljs.core.seq.call(null,cljs.core.map.call(null,((function (seq__25673,chunk__25674,count__25675,i__25676,ghost__$1,grid_size,scale_fn,ghost,forest,roots,tr__1504__auto__){
return (function (p1__25668_SHARP_){
return viz.forest.get_node.call(null,forest,p1__25668_SHARP_);
});})(seq__25673,chunk__25674,count__25675,i__25676,ghost__$1,grid_size,scale_fn,ghost,forest,roots,tr__1504__auto__))
,new cljs.core.Keyword(null,"active-node-ids","active-node-ids",-398210751).cljs$core$IFn$_invoke$arity$1(ghost__$1)));
var chunk__25678_25707 = null;
var count__25679_25708 = (0);
var i__25680_25709 = (0);
while(true){
if((i__25680_25709 < count__25679_25708)){
var active_node_25710 = cljs.core._nth.call(null,chunk__25678_25707,i__25680_25709);
viz.core.draw_node.call(null,active_node_25710,true,scale_fn);


var G__25711 = seq__25677_25706;
var G__25712 = chunk__25678_25707;
var G__25713 = count__25679_25708;
var G__25714 = (i__25680_25709 + (1));
seq__25677_25706 = G__25711;
chunk__25678_25707 = G__25712;
count__25679_25708 = G__25713;
i__25680_25709 = G__25714;
continue;
} else {
var temp__4657__auto___25715 = cljs.core.seq.call(null,seq__25677_25706);
if(temp__4657__auto___25715){
var seq__25677_25716__$1 = temp__4657__auto___25715;
if(cljs.core.chunked_seq_QMARK_.call(null,seq__25677_25716__$1)){
var c__4461__auto___25717 = cljs.core.chunk_first.call(null,seq__25677_25716__$1);
var G__25718 = cljs.core.chunk_rest.call(null,seq__25677_25716__$1);
var G__25719 = c__4461__auto___25717;
var G__25720 = cljs.core.count.call(null,c__4461__auto___25717);
var G__25721 = (0);
seq__25677_25706 = G__25718;
chunk__25678_25707 = G__25719;
count__25679_25708 = G__25720;
i__25680_25709 = G__25721;
continue;
} else {
var active_node_25722 = cljs.core.first.call(null,seq__25677_25716__$1);
viz.core.draw_node.call(null,active_node_25722,true,scale_fn);


var G__25723 = cljs.core.next.call(null,seq__25677_25716__$1);
var G__25724 = null;
var G__25725 = (0);
var G__25726 = (0);
seq__25677_25706 = G__25723;
chunk__25678_25707 = G__25724;
count__25679_25708 = G__25725;
i__25680_25709 = G__25726;
continue;
}
} else {
}
}
break;
}


var G__25727 = seq__25673;
var G__25728 = chunk__25674;
var G__25729 = count__25675;
var G__25730 = (i__25676 + (1));
seq__25673 = G__25727;
chunk__25674 = G__25728;
count__25675 = G__25729;
i__25676 = G__25730;
continue;
} else {
var temp__4657__auto__ = cljs.core.seq.call(null,seq__25673);
if(temp__4657__auto__){
var seq__25673__$1 = temp__4657__auto__;
if(cljs.core.chunked_seq_QMARK_.call(null,seq__25673__$1)){
var c__4461__auto__ = cljs.core.chunk_first.call(null,seq__25673__$1);
var G__25731 = cljs.core.chunk_rest.call(null,seq__25673__$1);
var G__25732 = c__4461__auto__;
var G__25733 = cljs.core.count.call(null,c__4461__auto__);
var G__25734 = (0);
seq__25673 = G__25731;
chunk__25674 = G__25732;
count__25675 = G__25733;
i__25676 = G__25734;
continue;
} else {
var ghost__$1 = cljs.core.first.call(null,seq__25673__$1);
var seq__25681_25735 = cljs.core.seq.call(null,cljs.core.map.call(null,((function (seq__25673,chunk__25674,count__25675,i__25676,ghost__$1,seq__25673__$1,temp__4657__auto__,grid_size,scale_fn,ghost,forest,roots,tr__1504__auto__){
return (function (p1__25668_SHARP_){
return viz.forest.get_node.call(null,forest,p1__25668_SHARP_);
});})(seq__25673,chunk__25674,count__25675,i__25676,ghost__$1,seq__25673__$1,temp__4657__auto__,grid_size,scale_fn,ghost,forest,roots,tr__1504__auto__))
,new cljs.core.Keyword(null,"active-node-ids","active-node-ids",-398210751).cljs$core$IFn$_invoke$arity$1(ghost__$1)));
var chunk__25682_25736 = null;
var count__25683_25737 = (0);
var i__25684_25738 = (0);
while(true){
if((i__25684_25738 < count__25683_25737)){
var active_node_25739 = cljs.core._nth.call(null,chunk__25682_25736,i__25684_25738);
viz.core.draw_node.call(null,active_node_25739,true,scale_fn);


var G__25740 = seq__25681_25735;
var G__25741 = chunk__25682_25736;
var G__25742 = count__25683_25737;
var G__25743 = (i__25684_25738 + (1));
seq__25681_25735 = G__25740;
chunk__25682_25736 = G__25741;
count__25683_25737 = G__25742;
i__25684_25738 = G__25743;
continue;
} else {
var temp__4657__auto___25744__$1 = cljs.core.seq.call(null,seq__25681_25735);
if(temp__4657__auto___25744__$1){
var seq__25681_25745__$1 = temp__4657__auto___25744__$1;
if(cljs.core.chunked_seq_QMARK_.call(null,seq__25681_25745__$1)){
var c__4461__auto___25746 = cljs.core.chunk_first.call(null,seq__25681_25745__$1);
var G__25747 = cljs.core.chunk_rest.call(null,seq__25681_25745__$1);
var G__25748 = c__4461__auto___25746;
var G__25749 = cljs.core.count.call(null,c__4461__auto___25746);
var G__25750 = (0);
seq__25681_25735 = G__25747;
chunk__25682_25736 = G__25748;
count__25683_25737 = G__25749;
i__25684_25738 = G__25750;
continue;
} else {
var active_node_25751 = cljs.core.first.call(null,seq__25681_25745__$1);
viz.core.draw_node.call(null,active_node_25751,true,scale_fn);


var G__25752 = cljs.core.next.call(null,seq__25681_25745__$1);
var G__25753 = null;
var G__25754 = (0);
var G__25755 = (0);
seq__25681_25735 = G__25752;
chunk__25682_25736 = G__25753;
count__25683_25737 = G__25754;
i__25684_25738 = G__25755;
continue;
}
} else {
}
}
break;
}


var G__25756 = cljs.core.next.call(null,seq__25673__$1);
var G__25757 = null;
var G__25758 = (0);
var G__25759 = (0);
seq__25673 = G__25756;
chunk__25674 = G__25757;
count__25675 = G__25758;
i__25676 = G__25759;
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
var G__25760__delegate = function (args){
return cljs.core.apply.call(null,viz.core.update_state,args);
};
var G__25760 = function (var_args){
var args = null;
if (arguments.length > 0) {
var G__25761__i = 0, G__25761__a = new Array(arguments.length -  0);
while (G__25761__i < G__25761__a.length) {G__25761__a[G__25761__i] = arguments[G__25761__i + 0]; ++G__25761__i;}
  args = new cljs.core.IndexedSeq(G__25761__a,0,null);
} 
return G__25760__delegate.call(this,args);};
G__25760.cljs$lang$maxFixedArity = 0;
G__25760.cljs$lang$applyTo = (function (arglist__25762){
var args = cljs.core.seq(arglist__25762);
return G__25760__delegate(args);
});
G__25760.cljs$core$IFn$_invoke$arity$variadic = G__25760__delegate;
return G__25760;
})()
:viz.core.update_state),new cljs.core.Keyword(null,"size","size",1098693007),((cljs.core.fn_QMARK_.call(null,viz.core.window_size))?(function() { 
var G__25763__delegate = function (args){
return cljs.core.apply.call(null,viz.core.window_size,args);
};
var G__25763 = function (var_args){
var args = null;
if (arguments.length > 0) {
var G__25764__i = 0, G__25764__a = new Array(arguments.length -  0);
while (G__25764__i < G__25764__a.length) {G__25764__a[G__25764__i] = arguments[G__25764__i + 0]; ++G__25764__i;}
  args = new cljs.core.IndexedSeq(G__25764__a,0,null);
} 
return G__25763__delegate.call(this,args);};
G__25763.cljs$lang$maxFixedArity = 0;
G__25763.cljs$lang$applyTo = (function (arglist__25765){
var args = cljs.core.seq(arglist__25765);
return G__25763__delegate(args);
});
G__25763.cljs$core$IFn$_invoke$arity$variadic = G__25763__delegate;
return G__25763;
})()
:viz.core.window_size),new cljs.core.Keyword(null,"title","title",636505583),"",new cljs.core.Keyword(null,"setup","setup",1987730512),((cljs.core.fn_QMARK_.call(null,viz.core.setup))?(function() { 
var G__25766__delegate = function (args){
return cljs.core.apply.call(null,viz.core.setup,args);
};
var G__25766 = function (var_args){
var args = null;
if (arguments.length > 0) {
var G__25767__i = 0, G__25767__a = new Array(arguments.length -  0);
while (G__25767__i < G__25767__a.length) {G__25767__a[G__25767__i] = arguments[G__25767__i + 0]; ++G__25767__i;}
  args = new cljs.core.IndexedSeq(G__25767__a,0,null);
} 
return G__25766__delegate.call(this,args);};
G__25766.cljs$lang$maxFixedArity = 0;
G__25766.cljs$lang$applyTo = (function (arglist__25768){
var args = cljs.core.seq(arglist__25768);
return G__25766__delegate(args);
});
G__25766.cljs$core$IFn$_invoke$arity$variadic = G__25766__delegate;
return G__25766;
})()
:viz.core.setup),new cljs.core.Keyword(null,"middleware","middleware",1462115504),new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [quil.middleware.fun_mode], null),new cljs.core.Keyword(null,"draw","draw",1358331674),((cljs.core.fn_QMARK_.call(null,viz.core.draw_state))?(function() { 
var G__25769__delegate = function (args){
return cljs.core.apply.call(null,viz.core.draw_state,args);
};
var G__25769 = function (var_args){
var args = null;
if (arguments.length > 0) {
var G__25770__i = 0, G__25770__a = new Array(arguments.length -  0);
while (G__25770__i < G__25770__a.length) {G__25770__a[G__25770__i] = arguments[G__25770__i + 0]; ++G__25770__i;}
  args = new cljs.core.IndexedSeq(G__25770__a,0,null);
} 
return G__25769__delegate.call(this,args);};
G__25769.cljs$lang$maxFixedArity = 0;
G__25769.cljs$lang$applyTo = (function (arglist__25771){
var args = cljs.core.seq(arglist__25771);
return G__25769__delegate(args);
});
G__25769.cljs$core$IFn$_invoke$arity$variadic = G__25769__delegate;
return G__25769;
})()
:viz.core.draw_state));
});
goog.exportSymbol('viz.core.viz', viz.core.viz);

if(cljs.core.truth_(cljs.core.some.call(null,(function (p1__1117__1118__auto__){
return cljs.core._EQ_.call(null,new cljs.core.Keyword(null,"no-start","no-start",1381488856),p1__1117__1118__auto__);
}),new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"keep-on-top","keep-on-top",-970284267)], null)))){
} else {
quil.sketch.add_sketch_to_init_list.call(null,new cljs.core.PersistentArrayMap(null, 2, [new cljs.core.Keyword(null,"fn","fn",-1175266204),viz.core.viz,new cljs.core.Keyword(null,"host-id","host-id",742376279),"viz"], null));
}

//# sourceMappingURL=core.js.map
