// Compiled by ClojureScript 1.9.473 {}
goog.provide('viz.forest');
goog.require('cljs.core');
viz.forest.new_forest = (function viz$forest$new_forest(){
return new cljs.core.PersistentArrayMap(null, 4, [new cljs.core.Keyword(null,"nodes","nodes",-2099585805),cljs.core.PersistentArrayMap.EMPTY,new cljs.core.Keyword(null,"roots","roots",-1088919250),cljs.core.PersistentHashSet.EMPTY,new cljs.core.Keyword(null,"leaves","leaves",-2143630574),cljs.core.PersistentHashSet.EMPTY,new cljs.core.Keyword(null,"next-id","next-id",-224240762),(0)], null);
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
return (function (p1__7982_SHARP_){
if(cljs.core.truth_(prev_parent_id)){
return viz.forest.unset_parent.call(null,p1__7982_SHARP_,id,prev_parent_id);
} else {
return p1__7982_SHARP_;
}
});})(parent_pos,prev_parent_id))
.call(null,cljs.core.update_in.call(null,cljs.core.update_in.call(null,cljs.core.update_in.call(null,cljs.core.assoc_in.call(null,cljs.core.assoc_in.call(null,forest,new cljs.core.PersistentVector(null, 3, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"nodes","nodes",-2099585805),id,new cljs.core.Keyword(null,"parent-id","parent-id",-1400729131)], null),parent_id),new cljs.core.PersistentVector(null, 3, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"nodes","nodes",-2099585805),id,new cljs.core.Keyword(null,"parent-pos","parent-pos",-282368566)], null),parent_pos),new cljs.core.PersistentVector(null, 3, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"nodes","nodes",-2099585805),parent_id,new cljs.core.Keyword(null,"child-ids","child-ids",-604525861)], null),((function (parent_pos,prev_parent_id){
return (function (p1__7981_SHARP_){
if(cljs.core.truth_(p1__7981_SHARP_)){
return cljs.core.conj.call(null,p1__7981_SHARP_,id);
} else {
return cljs.core.PersistentHashSet.createAsIfByAssoc([id], true);
}
});})(parent_pos,prev_parent_id))
),new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"roots","roots",-1088919250)], null),cljs.core.disj,id),new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"leaves","leaves",-2143630574)], null),cljs.core.disj,parent_id));
});
viz.forest.add_node = (function viz$forest$add_node(forest,pos){
var vec__7986 = viz.forest.new_id.call(null,forest);
var forest__$1 = cljs.core.nth.call(null,vec__7986,(0),null);
var id = cljs.core.nth.call(null,vec__7986,(1),null);
var forest__$2 = cljs.core.update_in.call(null,cljs.core.update_in.call(null,cljs.core.assoc_in.call(null,forest__$1,new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"nodes","nodes",-2099585805),id], null),new cljs.core.PersistentArrayMap(null, 2, [new cljs.core.Keyword(null,"id","id",-1388402092),id,new cljs.core.Keyword(null,"pos","pos",-864607220),pos], null)),new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"roots","roots",-1088919250)], null),cljs.core.conj,id),new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"leaves","leaves",-2143630574)], null),cljs.core.conj,id);
return new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [forest__$2,id], null);
});
viz.forest.remove_node = (function viz$forest$remove_node(forest,id){
var child_ids = cljs.core.get_in.call(null,forest,new cljs.core.PersistentVector(null, 3, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"nodes","nodes",-2099585805),id,new cljs.core.Keyword(null,"child-ids","child-ids",-604525861)], null));
var parent_id = cljs.core.get_in.call(null,forest,new cljs.core.PersistentVector(null, 3, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"nodes","nodes",-2099585805),id,new cljs.core.Keyword(null,"parent-id","parent-id",-1400729131)], null));
return cljs.core.update_in.call(null,cljs.core.update_in.call(null,cljs.core.update_in.call(null,((function (child_ids,parent_id){
return (function (forest__$1){
return cljs.core.reduce.call(null,((function (child_ids,parent_id){
return (function (p1__7990_SHARP_,p2__7991_SHARP_){
return viz.forest.unset_parent.call(null,p1__7990_SHARP_,p2__7991_SHARP_,id);
});})(child_ids,parent_id))
,forest__$1,child_ids);
});})(child_ids,parent_id))
.call(null,((function (child_ids,parent_id){
return (function (p1__7989_SHARP_){
if(cljs.core.truth_(parent_id)){
return viz.forest.unset_parent.call(null,p1__7989_SHARP_,id,parent_id);
} else {
return p1__7989_SHARP_;
}
});})(child_ids,parent_id))
.call(null,forest)),new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"nodes","nodes",-2099585805)], null),cljs.core.dissoc,id),new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"roots","roots",-1088919250)], null),cljs.core.disj,id),new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"leaves","leaves",-2143630574)], null),cljs.core.disj,id);
});
viz.forest.get_node = (function viz$forest$get_node(forest,id){
return cljs.core.get_in.call(null,forest,new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"nodes","nodes",-2099585805),id], null));
});
viz.forest.spawn_child = (function viz$forest$spawn_child(forest,parent_id,pos){
var vec__7995 = viz.forest.add_node.call(null,forest,pos);
var forest__$1 = cljs.core.nth.call(null,vec__7995,(0),null);
var id = cljs.core.nth.call(null,vec__7995,(1),null);
var forest__$2 = viz.forest.set_parent.call(null,forest__$1,id,parent_id);
return new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [forest__$2,id], null);
});
viz.forest.roots = (function viz$forest$roots(forest){
return cljs.core.vals.call(null,cljs.core.select_keys.call(null,new cljs.core.Keyword(null,"nodes","nodes",-2099585805).cljs$core$IFn$_invoke$arity$1(forest),new cljs.core.Keyword(null,"roots","roots",-1088919250).cljs$core$IFn$_invoke$arity$1(forest)));
});
viz.forest.root_QMARK_ = (function viz$forest$root_QMARK_(node){
return !(cljs.core.boolean$.call(null,new cljs.core.Keyword(null,"parent-id","parent-id",-1400729131).cljs$core$IFn$_invoke$arity$1(node)));
});
viz.forest.leaves = (function viz$forest$leaves(forest){
return cljs.core.vals.call(null,cljs.core.select_keys.call(null,new cljs.core.Keyword(null,"nodes","nodes",-2099585805).cljs$core$IFn$_invoke$arity$1(forest),new cljs.core.Keyword(null,"leaves","leaves",-2143630574).cljs$core$IFn$_invoke$arity$1(forest)));
});
viz.forest.leaf_QMARK_ = (function viz$forest$leaf_QMARK_(node){
return cljs.core.empty_QMARK_.call(null,new cljs.core.Keyword(null,"child-ids","child-ids",-604525861).cljs$core$IFn$_invoke$arity$1(node));
});
viz.forest.lines = (function viz$forest$lines(forest){
return cljs.core.map.call(null,(function (p1__7999_SHARP_){
return (new cljs.core.PersistentVector(null,2,(5),cljs.core.PersistentVector.EMPTY_NODE,[new cljs.core.Keyword(null,"pos","pos",-864607220).cljs$core$IFn$_invoke$arity$1(p1__7999_SHARP_),new cljs.core.Keyword(null,"parent-pos","parent-pos",-282368566).cljs$core$IFn$_invoke$arity$1(p1__7999_SHARP_)],null));
}),cljs.core.remove.call(null,(function (p1__7998_SHARP_){
return cljs.core.empty_QMARK_.call(null,new cljs.core.Keyword(null,"parent-pos","parent-pos",-282368566).cljs$core$IFn$_invoke$arity$1(p1__7998_SHARP_));
}),cljs.core.vals.call(null,new cljs.core.Keyword(null,"nodes","nodes",-2099585805).cljs$core$IFn$_invoke$arity$1(forest))));
});
viz.forest.my_forest = (function (){var forest = viz.forest.new_forest.call(null);
var vec__8000 = viz.forest.add_node.call(null,forest,new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [(0),(0)], null));
var forest__$1 = cljs.core.nth.call(null,vec__8000,(0),null);
var id0 = cljs.core.nth.call(null,vec__8000,(1),null);
var vec__8003 = viz.forest.spawn_child.call(null,forest__$1,id0,new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [(1),(1)], null));
var forest__$2 = cljs.core.nth.call(null,vec__8003,(0),null);
var id1 = cljs.core.nth.call(null,vec__8003,(1),null);
var vec__8006 = viz.forest.spawn_child.call(null,forest__$2,id0,new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [(-1),(-1)], null));
var forest__$3 = cljs.core.nth.call(null,vec__8006,(0),null);
var id2 = cljs.core.nth.call(null,vec__8006,(1),null);
var vec__8009 = viz.forest.spawn_child.call(null,forest__$3,id1,new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [(2),(2)], null));
var forest__$4 = cljs.core.nth.call(null,vec__8009,(0),null);
var id3 = cljs.core.nth.call(null,vec__8009,(1),null);
var forest__$5 = viz.forest.remove_node.call(null,forest__$4,id1);
return forest__$5;
})();
cljs.core.identity.call(null,viz.forest.my_forest);
viz.forest.lines.call(null,viz.forest.my_forest);

//# sourceMappingURL=forest.js.map