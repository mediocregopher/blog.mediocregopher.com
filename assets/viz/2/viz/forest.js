// Compiled by ClojureScript 1.10.439 {}
goog.provide('viz.forest');
goog.require('cljs.core');
goog.require('viz.grid');
viz.forest.new_forest = (function viz$forest$new_forest(grid_def){
return new cljs.core.PersistentArrayMap(null, 5, [new cljs.core.Keyword(null,"grid","grid",402978600),viz.grid.new_grid.call(null,grid_def),new cljs.core.Keyword(null,"nodes","nodes",-2099585805),cljs.core.PersistentArrayMap.EMPTY,new cljs.core.Keyword(null,"roots","roots",-1088919250),cljs.core.PersistentHashSet.EMPTY,new cljs.core.Keyword(null,"leaves","leaves",-2143630574),cljs.core.PersistentHashSet.EMPTY,new cljs.core.Keyword(null,"next-id","next-id",-224240762),(0)], null);
});
viz.forest.new_id = (function viz$forest$new_id(forest){
var id = new cljs.core.Keyword(null,"next-id","next-id",-224240762).cljs$core$IFn$_invoke$arity$1(forest);
return new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [cljs.core.assoc.call(null,forest,new cljs.core.Keyword(null,"next-id","next-id",-224240762),(id + (1))),id], null);
});
viz.forest.unset_parent = (function viz$forest$unset_parent(forest,id,parent_id){
return cljs.core.update_in.call(null,cljs.core.update_in.call(null,cljs.core.update_in.call(null,cljs.core.update_in.call(null,forest,new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"nodes","nodes",-2099585805),id], null),cljs.core.dissoc,new cljs.core.Keyword(null,"parent-id","parent-id",-1400729131),new cljs.core.Keyword(null,"parent-pos","parent-pos",-282368566)),new cljs.core.PersistentVector(null, 3, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"nodes","nodes",-2099585805),parent_id,new cljs.core.Keyword(null,"child-ids","child-ids",-604525861)], null),cljs.core.disj,id),new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"roots","roots",-1088919250)], null),cljs.core.conj,id),new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"leaves","leaves",-2143630574)], null),cljs.core.conj,parent_id);
});
viz.forest.set_parent = (function viz$forest$set_parent(forest,id,parent_id){
var parent_pos = cljs.core.get_in.call(null,forest,new cljs.core.PersistentVector(null, 3, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"nodes","nodes",-2099585805),parent_id,new cljs.core.Keyword(null,"pos","pos",-864607220)], null));
var prev_parent_id = cljs.core.get_in.call(null,forest,new cljs.core.PersistentVector(null, 3, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"nodes","nodes",-2099585805),id,new cljs.core.Keyword(null,"parent-id","parent-id",-1400729131)], null));
return ((function (parent_pos,prev_parent_id){
return (function (p1__10848_SHARP_){
if(cljs.core.truth_(prev_parent_id)){
return viz.forest.unset_parent.call(null,p1__10848_SHARP_,id,prev_parent_id);
} else {
return p1__10848_SHARP_;
}
});})(parent_pos,prev_parent_id))
.call(null,cljs.core.update_in.call(null,cljs.core.update_in.call(null,cljs.core.update_in.call(null,cljs.core.assoc_in.call(null,cljs.core.assoc_in.call(null,forest,new cljs.core.PersistentVector(null, 3, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"nodes","nodes",-2099585805),id,new cljs.core.Keyword(null,"parent-id","parent-id",-1400729131)], null),parent_id),new cljs.core.PersistentVector(null, 3, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"nodes","nodes",-2099585805),id,new cljs.core.Keyword(null,"parent-pos","parent-pos",-282368566)], null),parent_pos),new cljs.core.PersistentVector(null, 3, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"nodes","nodes",-2099585805),parent_id,new cljs.core.Keyword(null,"child-ids","child-ids",-604525861)], null),((function (parent_pos,prev_parent_id){
return (function (p1__10847_SHARP_){
if(cljs.core.truth_(p1__10847_SHARP_)){
return cljs.core.conj.call(null,p1__10847_SHARP_,id);
} else {
return cljs.core.PersistentHashSet.createAsIfByAssoc([id]);
}
});})(parent_pos,prev_parent_id))
),new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"roots","roots",-1088919250)], null),cljs.core.disj,id),new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"leaves","leaves",-2143630574)], null),cljs.core.disj,parent_id));
});
viz.forest.node_at_pos_QMARK_ = (function viz$forest$node_at_pos_QMARK_(forest,pos){
return cljs.core.boolean$.call(null,cljs.core.some.call(null,(function (p1__10849_SHARP_){
return cljs.core._EQ_.call(null,pos,new cljs.core.Keyword(null,"pos","pos",-864607220).cljs$core$IFn$_invoke$arity$1(p1__10849_SHARP_));
}),cljs.core.vals.call(null,new cljs.core.Keyword(null,"nodes","nodes",-2099585805).cljs$core$IFn$_invoke$arity$1(forest))));
});
viz.forest.empty_adjacent_points = (function viz$forest$empty_adjacent_points(forest,pos){
return viz.grid.empty_adjacent_points.call(null,new cljs.core.Keyword(null,"grid","grid",402978600).cljs$core$IFn$_invoke$arity$1(forest),pos);
});
viz.forest.add_node = (function viz$forest$add_node(forest,pos){
var vec__10850 = viz.forest.new_id.call(null,forest);
var forest__$1 = cljs.core.nth.call(null,vec__10850,(0),null);
var id = cljs.core.nth.call(null,vec__10850,(1),null);
var forest__$2 = cljs.core.update_in.call(null,cljs.core.update_in.call(null,cljs.core.assoc_in.call(null,cljs.core.update_in.call(null,forest__$1,new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"grid","grid",402978600)], null),viz.grid.add_point,pos),new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"nodes","nodes",-2099585805),id], null),new cljs.core.PersistentArrayMap(null, 2, [new cljs.core.Keyword(null,"id","id",-1388402092),id,new cljs.core.Keyword(null,"pos","pos",-864607220),pos], null)),new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"roots","roots",-1088919250)], null),cljs.core.conj,id),new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"leaves","leaves",-2143630574)], null),cljs.core.conj,id);
return new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [forest__$2,id], null);
});
viz.forest.remove_node = (function viz$forest$remove_node(forest,id){
var node = cljs.core.get_in.call(null,forest,new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"nodes","nodes",-2099585805),id], null));
var child_ids = new cljs.core.Keyword(null,"child-ids","child-ids",-604525861).cljs$core$IFn$_invoke$arity$1(node);
var parent_id = new cljs.core.Keyword(null,"parent-id","parent-id",-1400729131).cljs$core$IFn$_invoke$arity$1(node);
return cljs.core.update_in.call(null,cljs.core.update_in.call(null,cljs.core.update_in.call(null,((function (node,child_ids,parent_id){
return (function (forest__$1){
return cljs.core.reduce.call(null,((function (node,child_ids,parent_id){
return (function (p1__10854_SHARP_,p2__10855_SHARP_){
return viz.forest.unset_parent.call(null,p1__10854_SHARP_,p2__10855_SHARP_,id);
});})(node,child_ids,parent_id))
,forest__$1,child_ids);
});})(node,child_ids,parent_id))
.call(null,((function (node,child_ids,parent_id){
return (function (p1__10853_SHARP_){
if(cljs.core.truth_(parent_id)){
return viz.forest.unset_parent.call(null,p1__10853_SHARP_,id,parent_id);
} else {
return p1__10853_SHARP_;
}
});})(node,child_ids,parent_id))
.call(null,cljs.core.update_in.call(null,forest,new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"grid","grid",402978600)], null),viz.grid.rm_point,new cljs.core.Keyword(null,"pos","pos",-864607220).cljs$core$IFn$_invoke$arity$1(node)))),new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"nodes","nodes",-2099585805)], null),cljs.core.dissoc,id),new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"roots","roots",-1088919250)], null),cljs.core.disj,id),new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"leaves","leaves",-2143630574)], null),cljs.core.disj,id);
});
viz.forest.update_node_meta = (function viz$forest$update_node_meta(forest,id,f){
return cljs.core.update_in.call(null,forest,new cljs.core.PersistentVector(null, 3, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"nodes","nodes",-2099585805),id,new cljs.core.Keyword(null,"meta","meta",1499536964)], null),f);
});
viz.forest.get_node_meta = (function viz$forest$get_node_meta(forest,id){
return cljs.core.get_in.call(null,forest,new cljs.core.PersistentVector(null, 3, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"nodes","nodes",-2099585805),id,new cljs.core.Keyword(null,"meta","meta",1499536964)], null));
});
viz.forest.get_node = (function viz$forest$get_node(forest,id){
return cljs.core.get_in.call(null,forest,new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"nodes","nodes",-2099585805),id], null));
});
viz.forest.spawn_child = (function viz$forest$spawn_child(forest,parent_id,pos){
var vec__10856 = viz.forest.add_node.call(null,forest,pos);
var forest__$1 = cljs.core.nth.call(null,vec__10856,(0),null);
var id = cljs.core.nth.call(null,vec__10856,(1),null);
var forest__$2 = viz.forest.set_parent.call(null,forest__$1,id,parent_id);
return new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [forest__$2,id], null);
});
viz.forest.roots = (function viz$forest$roots(forest){
return cljs.core.vals.call(null,cljs.core.select_keys.call(null,new cljs.core.Keyword(null,"nodes","nodes",-2099585805).cljs$core$IFn$_invoke$arity$1(forest),new cljs.core.Keyword(null,"roots","roots",-1088919250).cljs$core$IFn$_invoke$arity$1(forest)));
});
viz.forest.root_QMARK_ = (function viz$forest$root_QMARK_(node){
return (!(cljs.core.boolean$.call(null,new cljs.core.Keyword(null,"parent-id","parent-id",-1400729131).cljs$core$IFn$_invoke$arity$1(node))));
});
viz.forest.leaves = (function viz$forest$leaves(forest){
return cljs.core.vals.call(null,cljs.core.select_keys.call(null,new cljs.core.Keyword(null,"nodes","nodes",-2099585805).cljs$core$IFn$_invoke$arity$1(forest),new cljs.core.Keyword(null,"leaves","leaves",-2143630574).cljs$core$IFn$_invoke$arity$1(forest)));
});
viz.forest.leaf_QMARK_ = (function viz$forest$leaf_QMARK_(node){
return cljs.core.empty_QMARK_.call(null,new cljs.core.Keyword(null,"child-ids","child-ids",-604525861).cljs$core$IFn$_invoke$arity$1(node));
});
viz.forest.lines = (function viz$forest$lines(forest){
return cljs.core.map.call(null,(function (p1__10860_SHARP_){
return (new cljs.core.PersistentVector(null,2,(5),cljs.core.PersistentVector.EMPTY_NODE,[new cljs.core.Keyword(null,"pos","pos",-864607220).cljs$core$IFn$_invoke$arity$1(p1__10860_SHARP_),new cljs.core.Keyword(null,"parent-pos","parent-pos",-282368566).cljs$core$IFn$_invoke$arity$1(p1__10860_SHARP_)],null));
}),cljs.core.remove.call(null,(function (p1__10859_SHARP_){
return cljs.core.empty_QMARK_.call(null,new cljs.core.Keyword(null,"parent-pos","parent-pos",-282368566).cljs$core$IFn$_invoke$arity$1(p1__10859_SHARP_));
}),cljs.core.vals.call(null,new cljs.core.Keyword(null,"nodes","nodes",-2099585805).cljs$core$IFn$_invoke$arity$1(forest))));
});

//# sourceMappingURL=forest.js.map
