// Compiled by ClojureScript 1.9.473 {}
goog.provide('viz.ghost');
goog.require('cljs.core');
goog.require('viz.forest');
goog.require('viz.grid');
goog.require('clojure.set');
viz.ghost.new_ghost = (function viz$ghost$new_ghost(grid_def){
return new cljs.core.PersistentArrayMap(null, 3, [new cljs.core.Keyword(null,"grid","grid",402978600),viz.grid.new_grid.call(null,grid_def),new cljs.core.Keyword(null,"forest","forest",278860306),viz.forest.new_forest.call(null),new cljs.core.Keyword(null,"active-node-ids","active-node-ids",-398210751),cljs.core.PersistentHashSet.EMPTY], null);
});
viz.ghost.new_active_node = (function viz$ghost$new_active_node(ghost,pos){
var vec__8057 = viz.forest.add_node.call(null,new cljs.core.Keyword(null,"forest","forest",278860306).cljs$core$IFn$_invoke$arity$1(ghost),pos);
var forest = cljs.core.nth.call(null,vec__8057,(0),null);
var id = cljs.core.nth.call(null,vec__8057,(1),null);
var grid = viz.grid.add_point.call(null,new cljs.core.Keyword(null,"grid","grid",402978600).cljs$core$IFn$_invoke$arity$1(ghost),pos);
return cljs.core.update_in.call(null,cljs.core.assoc.call(null,ghost,new cljs.core.Keyword(null,"grid","grid",402978600),grid,new cljs.core.Keyword(null,"forest","forest",278860306),forest),new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"active-node-ids","active-node-ids",-398210751)], null),cljs.core.conj,id);
});
viz.ghost.gen_new_poss = (function viz$ghost$gen_new_poss(ghost,poss_fn,id){

var pos = new cljs.core.Keyword(null,"pos","pos",-864607220).cljs$core$IFn$_invoke$arity$1(viz.forest.get_node.call(null,new cljs.core.Keyword(null,"forest","forest",278860306).cljs$core$IFn$_invoke$arity$1(ghost),id));
var adj_poss = viz.grid.empty_adjacent_points.call(null,new cljs.core.Keyword(null,"grid","grid",402978600).cljs$core$IFn$_invoke$arity$1(ghost),pos);
return poss_fn.call(null,pos,adj_poss);
});
viz.ghost.spawn_children = (function viz$ghost$spawn_children(ghost,poss_fn,id){
return cljs.core.reduce.call(null,(function (p__8067,pos){
var vec__8068 = p__8067;
var ghost__$1 = cljs.core.nth.call(null,vec__8068,(0),null);
var new_ids = cljs.core.nth.call(null,vec__8068,(1),null);
var vec__8071 = viz.forest.spawn_child.call(null,new cljs.core.Keyword(null,"forest","forest",278860306).cljs$core$IFn$_invoke$arity$1(ghost__$1),id,pos);
var forest = cljs.core.nth.call(null,vec__8071,(0),null);
var new_id = cljs.core.nth.call(null,vec__8071,(1),null);
var grid = viz.grid.add_point.call(null,new cljs.core.Keyword(null,"grid","grid",402978600).cljs$core$IFn$_invoke$arity$1(ghost__$1),pos);
return new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [cljs.core.assoc.call(null,ghost__$1,new cljs.core.Keyword(null,"forest","forest",278860306),forest,new cljs.core.Keyword(null,"grid","grid",402978600),grid),cljs.core.conj.call(null,new_ids,new_id)], null);
}),new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [ghost,cljs.core.PersistentHashSet.EMPTY], null),viz.ghost.gen_new_poss.call(null,ghost,poss_fn,id));
});
viz.ghost.spawn_children_multi = (function viz$ghost$spawn_children_multi(ghost,poss_fn,ids){
return cljs.core.reduce.call(null,(function (p__8081,id){
var vec__8082 = p__8081;
var ghost__$1 = cljs.core.nth.call(null,vec__8082,(0),null);
var new_ids = cljs.core.nth.call(null,vec__8082,(1),null);
var vec__8085 = viz.ghost.spawn_children.call(null,ghost__$1,poss_fn,id);
var ghost__$2 = cljs.core.nth.call(null,vec__8085,(0),null);
var this_new_ids = cljs.core.nth.call(null,vec__8085,(1),null);
return new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [ghost__$2,clojure.set.union.call(null,new_ids,this_new_ids)], null);
}),new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [ghost,cljs.core.PersistentHashSet.EMPTY], null),ids);
});
viz.ghost.incr = (function viz$ghost$incr(ghost,poss_fn){
var vec__8091 = viz.ghost.spawn_children_multi.call(null,ghost,poss_fn,new cljs.core.Keyword(null,"active-node-ids","active-node-ids",-398210751).cljs$core$IFn$_invoke$arity$1(ghost));
var ghost__$1 = cljs.core.nth.call(null,vec__8091,(0),null);
var new_ids = cljs.core.nth.call(null,vec__8091,(1),null);
return cljs.core.assoc.call(null,ghost__$1,new cljs.core.Keyword(null,"active-node-ids","active-node-ids",-398210751),new_ids);
});
viz.ghost.active_nodes = (function viz$ghost$active_nodes(ghost){
return cljs.core.map.call(null,(function (p1__8094_SHARP_){
return cljs.core.get_in.call(null,ghost,new cljs.core.PersistentVector(null, 3, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"forest","forest",278860306),new cljs.core.Keyword(null,"nodes","nodes",-2099585805),p1__8094_SHARP_], null));
}),new cljs.core.Keyword(null,"active-node-ids","active-node-ids",-398210751).cljs$core$IFn$_invoke$arity$1(ghost));
});
viz.ghost.filter_active_nodes = (function viz$ghost$filter_active_nodes(ghost,pred){
return cljs.core.assoc.call(null,ghost,new cljs.core.Keyword(null,"active-node-ids","active-node-ids",-398210751),cljs.core.reduce.call(null,(function (p1__8096_SHARP_,p2__8095_SHARP_){
if(cljs.core.truth_(pred.call(null,p2__8095_SHARP_))){
return cljs.core.conj.call(null,p1__8096_SHARP_,new cljs.core.Keyword(null,"id","id",-1388402092).cljs$core$IFn$_invoke$arity$1(p2__8095_SHARP_));
} else {
return p1__8096_SHARP_;
}
}),cljs.core.PersistentHashSet.EMPTY,viz.ghost.active_nodes.call(null,ghost)));
});
viz.ghost.remove_roots = (function viz$ghost$remove_roots(ghost){
var roots = viz.forest.roots.call(null,new cljs.core.Keyword(null,"forest","forest",278860306).cljs$core$IFn$_invoke$arity$1(ghost));
var root_ids = cljs.core.map.call(null,new cljs.core.Keyword(null,"id","id",-1388402092),roots);
var root_poss = cljs.core.map.call(null,new cljs.core.Keyword(null,"pos","pos",-864607220),roots);
return cljs.core.update_in.call(null,cljs.core.update_in.call(null,cljs.core.update_in.call(null,ghost,new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"active-node-ids","active-node-ids",-398210751)], null),((function (roots,root_ids,root_poss){
return (function (p1__8097_SHARP_){
return cljs.core.reduce.call(null,cljs.core.disj,p1__8097_SHARP_,root_ids);
});})(roots,root_ids,root_poss))
),new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"forest","forest",278860306)], null),((function (roots,root_ids,root_poss){
return (function (p1__8098_SHARP_){
return cljs.core.reduce.call(null,viz.forest.remove_node,p1__8098_SHARP_,root_ids);
});})(roots,root_ids,root_poss))
),new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"grid","grid",402978600)], null),((function (roots,root_ids,root_poss){
return (function (p1__8099_SHARP_){
return cljs.core.reduce.call(null,viz.grid.rm_point,p1__8099_SHARP_,root_poss);
});})(roots,root_ids,root_poss))
);
});
viz.ghost.eg_poss_fn = (function viz$ghost$eg_poss_fn(pos,adj_poss){
return cljs.core.take.call(null,(2),cljs.core.random_sample.call(null,0.6,adj_poss));
});
viz.ghost.remove_roots.call(null,viz.ghost.incr.call(null,viz.ghost.incr.call(null,viz.ghost.incr.call(null,viz.ghost.new_active_node.call(null,viz.ghost.new_ghost.call(null,viz.grid.euclidean),new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [(0),(0)], null)),viz.ghost.eg_poss_fn),viz.ghost.eg_poss_fn),viz.ghost.eg_poss_fn));

//# sourceMappingURL=ghost.js.map