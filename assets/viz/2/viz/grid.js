// Compiled by ClojureScript 1.9.473 {}
goog.provide('viz.grid');
goog.require('cljs.core');
viz.grid.euclidean = new cljs.core.PersistentVector(null, 4, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [(0),(-1)], null),new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [(-1),(0)], null),new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [(1),(0)], null),new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [(0),(1)], null)], null);
viz.grid.isometric = new cljs.core.PersistentVector(null, 6, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [(-1),(-1)], null),new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [(0),(-2)], null),new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [(1),(-1)], null),new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [(-1),(1)], null),new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [(0),(2)], null),new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [(1),(1)], null)], null);
viz.grid.hexagonal = new cljs.core.PersistentVector(null, 3, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [(0),(-1)], null),new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [(-1),(1)], null),new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [(1),(1)], null)], null);
viz.grid.new_grid = (function viz$grid$new_grid(grid_def){
return new cljs.core.PersistentArrayMap(null, 2, [new cljs.core.Keyword(null,"grid-def","grid-def",-392588768),grid_def,new cljs.core.Keyword(null,"points","points",-1486596883),cljs.core.PersistentHashSet.EMPTY], null);
});
viz.grid.add_point = (function viz$grid$add_point(grid,point){
return cljs.core.update_in.call(null,grid,new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"points","points",-1486596883)], null),cljs.core.conj,point);
});
viz.grid.my_grid = viz.grid.add_point.call(null,viz.grid.new_grid.call(null,viz.grid.euclidean),new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [(0),(1)], null));
viz.grid.rm_point = (function viz$grid$rm_point(grid,point){
return cljs.core.update_in.call(null,grid,new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [new cljs.core.Keyword(null,"points","points",-1486596883)], null),cljs.core.disj,point);
});
viz.grid.adjacent_points = (function viz$grid$adjacent_points(grid,point){
return cljs.core.map.call(null,(function (p1__19767_SHARP_){
return cljs.core.map.call(null,cljs.core._PLUS_,p1__19767_SHARP_,point);
}),new cljs.core.Keyword(null,"grid-def","grid-def",-392588768).cljs$core$IFn$_invoke$arity$1(grid));
});
viz.grid.empty_adjacent_points = (function viz$grid$empty_adjacent_points(grid,point){
return cljs.core.remove.call(null,new cljs.core.Keyword(null,"points","points",-1486596883).cljs$core$IFn$_invoke$arity$1(grid),viz.grid.adjacent_points.call(null,grid,point));
});
viz.grid.empty_adjacent_points.call(null,viz.grid.add_point.call(null,viz.grid.add_point.call(null,viz.grid.new_grid.call(null,viz.grid.isometric),new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [(0),(0)], null)),new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [(0),(1)], null)),new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [(0),(1)], null));

//# sourceMappingURL=grid.js.map