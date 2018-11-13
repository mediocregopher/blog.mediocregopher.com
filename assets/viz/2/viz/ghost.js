// Compiled by ClojureScript 1.10.439 {}
goog.provide('viz.ghost');
goog.require('cljs.core');
goog.require('quil.core');
goog.require('quil.middleware');
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
return cljs.core.reduce.call(null,(function (p__10863,pos){
var vec__10864 = p__10863;
var forest__$1 = cljs.core.nth.call(null,vec__10864,(0),null);
var new_ids = cljs.core.nth.call(null,vec__10864,(1),null);
var vec__10867 = viz.forest.spawn_child.call(null,forest__$1,id,pos);
var forest__$2 = cljs.core.nth.call(null,vec__10867,(0),null);
var new_id = cljs.core.nth.call(null,vec__10867,(1),null);
return new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [forest__$2,cljs.core.conj.call(null,new_ids,new_id)], null);
}),new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [forest,cljs.core.PersistentHashSet.EMPTY], null),viz.ghost.gen_new_poss.call(null,forest,poss_fn,id));
});
viz.ghost.spawn_children_multi = (function viz$ghost$spawn_children_multi(forest,poss_fn,ids){
return cljs.core.reduce.call(null,(function (p__10870,id){
var vec__10871 = p__10870;
var forest__$1 = cljs.core.nth.call(null,vec__10871,(0),null);
var new_ids = cljs.core.nth.call(null,vec__10871,(1),null);
var vec__10874 = viz.ghost.spawn_children.call(null,forest__$1,poss_fn,id);
var forest__$2 = cljs.core.nth.call(null,vec__10874,(0),null);
var this_new_ids = cljs.core.nth.call(null,vec__10874,(1),null);
return new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [forest__$2,clojure.set.union.call(null,new_ids,this_new_ids)], null);
}),new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [forest,cljs.core.PersistentHashSet.EMPTY], null),ids);
});
viz.ghost.incr = (function viz$ghost$incr(ghost,forest,poss_fn){
var vec__10877 = viz.ghost.spawn_children_multi.call(null,forest,poss_fn,new cljs.core.Keyword(null,"active-node-ids","active-node-ids",-398210751).cljs$core$IFn$_invoke$arity$1(ghost));
var forest__$1 = cljs.core.nth.call(null,vec__10877,(0),null);
var new_ids = cljs.core.nth.call(null,vec__10877,(1),null);
return new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [cljs.core.assoc.call(null,ghost,new cljs.core.Keyword(null,"active-node-ids","active-node-ids",-398210751),new_ids),cljs.core.reduce.call(null,((function (vec__10877,forest__$1,new_ids){
return (function (forest__$2,id){
return viz.forest.update_node_meta.call(null,forest__$2,id,((function (vec__10877,forest__$1,new_ids){
return (function (m){
return cljs.core.assoc.call(null,m,new cljs.core.Keyword(null,"color","color",1011675173),new cljs.core.Keyword(null,"color","color",1011675173).cljs$core$IFn$_invoke$arity$1(ghost));
});})(vec__10877,forest__$1,new_ids))
);
});})(vec__10877,forest__$1,new_ids))
,forest__$1,new_ids)], null);
});
viz.ghost.eg_poss_fn = (function viz$ghost$eg_poss_fn(pos,adj_poss){
return cljs.core.take.call(null,(2),cljs.core.random_sample.call(null,0.6,adj_poss));
});

//# sourceMappingURL=ghost.js.map
