// Compiled by ClojureScript 1.9.473 {}
goog.provide('viz.ghost');
goog.require('cljs.core');
goog.require('viz.forest');
goog.require('viz.grid');
goog.require('clojure.set');
viz.ghost.new_ghost = (function viz$ghost$new_ghost(){
return new cljs.core.PersistentArrayMap(null, 2, [new cljs.core.Keyword(null,"active-node-ids","active-node-ids",-398210751),cljs.core.PersistentHashSet.EMPTY,new cljs.core.Keyword(null,"color","color",1011675173),(4278190080)], null);
});
viz.ghost.add_active_node = (function viz$ghost$add_active_node(ghost,id){
return cljs.core.update_in.call(null,ghost,new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"active-node-ids","active-node-ids",-398210751)], null),cljs.core.conj,id);
});
viz.ghost.rm_active_node = (function viz$ghost$rm_active_node(ghost,id){
return cljs.core.update_in.call(null,ghost,new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"active-node-ids","active-node-ids",-398210751)], null),cljs.core.disj,id);
});
viz.ghost.gen_new_poss = (function viz$ghost$gen_new_poss(forest,poss_fn,id){

var pos = new cljs.core.Keyword(null,"pos","pos",-864607220).cljs$core$IFn$_invoke$arity$1(viz.forest.get_node.call(null,forest,id));
var adj_poss = viz.forest.empty_adjacent_points.call(null,forest,pos);
var new_poss = poss_fn.call(null,pos,adj_poss);
return new_poss;
});
viz.ghost.spawn_children = (function viz$ghost$spawn_children(forest,poss_fn,id){
return cljs.core.reduce.call(null,(function (p__13494,pos){
var vec__13495 = p__13494;
var forest__$1 = cljs.core.nth.call(null,vec__13495,(0),null);
var new_ids = cljs.core.nth.call(null,vec__13495,(1),null);
var vec__13498 = viz.forest.spawn_child.call(null,forest__$1,id,pos);
var forest__$2 = cljs.core.nth.call(null,vec__13498,(0),null);
var new_id = cljs.core.nth.call(null,vec__13498,(1),null);
return new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [forest__$2,cljs.core.conj.call(null,new_ids,new_id)], null);
}),new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [forest,cljs.core.PersistentHashSet.EMPTY], null),viz.ghost.gen_new_poss.call(null,forest,poss_fn,id));
});
viz.ghost.spawn_children_multi = (function viz$ghost$spawn_children_multi(forest,poss_fn,ids){
return cljs.core.reduce.call(null,(function (p__13508,id){
var vec__13509 = p__13508;
var forest__$1 = cljs.core.nth.call(null,vec__13509,(0),null);
var new_ids = cljs.core.nth.call(null,vec__13509,(1),null);
var vec__13512 = viz.ghost.spawn_children.call(null,forest__$1,poss_fn,id);
var forest__$2 = cljs.core.nth.call(null,vec__13512,(0),null);
var this_new_ids = cljs.core.nth.call(null,vec__13512,(1),null);
return new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [forest__$2,clojure.set.union.call(null,new_ids,this_new_ids)], null);
}),new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [forest,cljs.core.PersistentHashSet.EMPTY], null),ids);
});
viz.ghost.incr = (function viz$ghost$incr(ghost,forest,poss_fn){
var vec__13518 = viz.ghost.spawn_children_multi.call(null,forest,poss_fn,new cljs.core.Keyword(null,"active-node-ids","active-node-ids",-398210751).cljs$core$IFn$_invoke$arity$1(ghost));
var forest__$1 = cljs.core.nth.call(null,vec__13518,(0),null);
var new_ids = cljs.core.nth.call(null,vec__13518,(1),null);
return new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [cljs.core.assoc.call(null,ghost,new cljs.core.Keyword(null,"active-node-ids","active-node-ids",-398210751),new_ids),forest__$1], null);
});
viz.ghost.active_nodes = (function viz$ghost$active_nodes(ghost,forest){
return cljs.core.map.call(null,(function (p1__13521_SHARP_){
return viz.forest.get_node.call(null,forest,p1__13521_SHARP_);
}),new cljs.core.Keyword(null,"active-node-ids","active-node-ids",-398210751).cljs$core$IFn$_invoke$arity$1(ghost));
});
viz.ghost.filter_active_nodes = (function viz$ghost$filter_active_nodes(ghost,pred){
return cljs.core.assoc.call(null,ghost,new cljs.core.Keyword(null,"active-node-ids","active-node-ids",-398210751),cljs.core.reduce.call(null,(function (p1__13523_SHARP_,p2__13522_SHARP_){
if(cljs.core.truth_(pred.call(null,p2__13522_SHARP_))){
return cljs.core.conj.call(null,p1__13523_SHARP_,new cljs.core.Keyword(null,"id","id",-1388402092).cljs$core$IFn$_invoke$arity$1(p2__13522_SHARP_));
} else {
return p1__13523_SHARP_;
}
}),cljs.core.PersistentHashSet.EMPTY,viz.ghost.active_nodes.call(null,ghost)));
});
viz.ghost.eg_poss_fn = (function viz$ghost$eg_poss_fn(pos,adj_poss){
return cljs.core.take.call(null,(2),cljs.core.random_sample.call(null,0.6,adj_poss));
});

//# sourceMappingURL=ghost.js.map