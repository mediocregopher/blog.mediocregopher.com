// Compiled by ClojureScript 1.9.473 {}
goog.provide('viz.dial');
goog.require('cljs.core');
goog.require('quil.core');
viz.dial.new_dial = (function viz$dial$new_dial(){
return new cljs.core.PersistentArrayMap(null, 3, [new cljs.core.Keyword(null,"val","val",128701612),(0),new cljs.core.Keyword(null,"min","min",444991522),(-1),new cljs.core.Keyword(null,"max","max",61366548),(1)], null);
});
viz.dial.scale = (function viz$dial$scale(v,old_min,old_max,new_min,new_max){
return (new_min + ((new_max - new_min) * ((v - old_min) / (old_max - old_min))));
});
viz.dial.scaled = (function viz$dial$scaled(dial,min,max){
var new_val = viz.dial.scale.call(null,new cljs.core.Keyword(null,"val","val",128701612).cljs$core$IFn$_invoke$arity$1(dial),new cljs.core.Keyword(null,"min","min",444991522).cljs$core$IFn$_invoke$arity$1(dial),new cljs.core.Keyword(null,"max","max",61366548).cljs$core$IFn$_invoke$arity$1(dial),min,max);
return cljs.core.assoc.call(null,dial,new cljs.core.Keyword(null,"min","min",444991522),min,new cljs.core.Keyword(null,"max","max",61366548),max,new cljs.core.Keyword(null,"val","val",128701612),new_val);
});
viz.dial.floored = (function viz$dial$floored(dial,at){
if((new cljs.core.Keyword(null,"val","val",128701612).cljs$core$IFn$_invoke$arity$1(dial) < at)){
return cljs.core.assoc.call(null,dial,new cljs.core.Keyword(null,"val","val",128701612),at);
} else {
return dial;
}
});
viz.dial.invert = (function viz$dial$invert(dial){
return cljs.core.assoc.call(null,dial,new cljs.core.Keyword(null,"val","val",128701612),((-1) * new cljs.core.Keyword(null,"val","val",128701612).cljs$core$IFn$_invoke$arity$1(dial)));
});
viz.dial.new_plot = (function viz$dial$new_plot(frame_rate,period_seconds,plot){
return new cljs.core.PersistentArrayMap(null, 3, [new cljs.core.Keyword(null,"frame-rate","frame-rate",-994918942),frame_rate,new cljs.core.Keyword(null,"period","period",-352129191),period_seconds,new cljs.core.Keyword(null,"plot","plot",2086832225),plot], null);
});
viz.dial.by_plot = (function viz$dial$by_plot(dial,plot,curr_frame){
var dial_t = (cljs.core.mod.call(null,(curr_frame / new cljs.core.Keyword(null,"frame-rate","frame-rate",-994918942).cljs$core$IFn$_invoke$arity$1(plot)),new cljs.core.Keyword(null,"period","period",-352129191).cljs$core$IFn$_invoke$arity$1(plot)) / new cljs.core.Keyword(null,"period","period",-352129191).cljs$core$IFn$_invoke$arity$1(plot));
return cljs.core.assoc.call(null,dial,new cljs.core.Keyword(null,"val","val",128701612),cljs.core.reduce.call(null,((function (dial_t){
return (function (curr_v,p__8119){
var vec__8120 = p__8119;
var t = cljs.core.nth.call(null,vec__8120,(0),null);
var v = cljs.core.nth.call(null,vec__8120,(1),null);
if((t <= dial_t)){
return v;
} else {
return cljs.core.reduced.call(null,curr_v);
}
});})(dial_t))
,(0),new cljs.core.Keyword(null,"plot","plot",2086832225).cljs$core$IFn$_invoke$arity$1(plot)));
});

//# sourceMappingURL=dial.js.map