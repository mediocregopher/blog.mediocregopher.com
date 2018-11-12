// Compiled by ClojureScript 1.9.473 {}
goog.provide('quil.sketch');
goog.require('cljs.core');
goog.require('quil.util');
goog.require('quil.middlewares.deprecated_options');
goog.require('goog.dom');
goog.require('goog.events');
goog.require('goog.events.EventType');
quil.sketch._STAR_applet_STAR_ = null;
quil.sketch.current_applet = (function quil$sketch$current_applet(){
return quil.sketch._STAR_applet_STAR_;
});
quil.sketch.rendering_modes = new cljs.core.PersistentArrayMap(null, 4, [new cljs.core.Keyword(null,"java2d","java2d",166099237),(Processing.prototype.PConstants["JAVA2D"]),new cljs.core.Keyword(null,"p2d","p2d",-2106175755),(Processing.prototype.PConstants["P2D"]),new cljs.core.Keyword(null,"p3d","p3d",-850380194),(Processing.prototype.PConstants["P3D"]),new cljs.core.Keyword(null,"opengl","opengl",-614998103),(Processing.prototype.PConstants["OPENGL"])], null);
quil.sketch.resolve_renderer = (function quil$sketch$resolve_renderer(mode){
return quil.util.resolve_constant_key.call(null,mode,quil.sketch.rendering_modes);
});
quil.sketch.size = (function quil$sketch$size(var_args){
var args8190 = [];
var len__7927__auto___8193 = arguments.length;
var i__7928__auto___8194 = (0);
while(true){
if((i__7928__auto___8194 < len__7927__auto___8193)){
args8190.push((arguments[i__7928__auto___8194]));

var G__8195 = (i__7928__auto___8194 + (1));
i__7928__auto___8194 = G__8195;
continue;
} else {
}
break;
}

var G__8192 = args8190.length;
switch (G__8192) {
case 2:
return quil.sketch.size.cljs$core$IFn$_invoke$arity$2((arguments[(0)]),(arguments[(1)]));

break;
case 3:
return quil.sketch.size.cljs$core$IFn$_invoke$arity$3((arguments[(0)]),(arguments[(1)]),(arguments[(2)]));

break;
default:
throw (new Error([cljs.core.str.cljs$core$IFn$_invoke$arity$1("Invalid arity: "),cljs.core.str.cljs$core$IFn$_invoke$arity$1(args8190.length)].join('')));

}
});

quil.sketch.size.cljs$core$IFn$_invoke$arity$2 = (function (width,height){
return quil.sketch.current_applet.call(null).size((width | (0)),(height | (0)));
});

quil.sketch.size.cljs$core$IFn$_invoke$arity$3 = (function (width,height,mode){
return quil.sketch.current_applet.call(null).size((width | (0)),(height | (0)),quil.util.resolve_constant_key.call(null,mode,quil.sketch.rendering_modes));
});

quil.sketch.size.cljs$lang$maxFixedArity = 3;

quil.sketch.bind_handlers = (function quil$sketch$bind_handlers(prc,opts){
var seq__8209 = cljs.core.seq.call(null,cljs.core.PersistentHashMap.fromArrays([new cljs.core.Keyword(null,"keyPressed","keyPressed",1791025256),new cljs.core.Keyword(null,"mouseOut","mouseOut",-386669045),new cljs.core.Keyword(null,"mouseScrolled","mouseScrolled",31878252),new cljs.core.Keyword(null,"mouseDragged","mouseDragged",129975181),new cljs.core.Keyword(null,"setup","setup",1987730512),new cljs.core.Keyword(null,"keyReleased","keyReleased",541714964),new cljs.core.Keyword(null,"mouseClicked","mouseClicked",1764302965),new cljs.core.Keyword(null,"mouseReleased","mouseReleased",1116234838),new cljs.core.Keyword(null,"mousePressed","mousePressed",1776186454),new cljs.core.Keyword(null,"mouseMoved","mouseMoved",-1936954058),new cljs.core.Keyword(null,"mouseOver","mouseOver",-1334461930),new cljs.core.Keyword(null,"keyTyped","keyTyped",1437329399),new cljs.core.Keyword(null,"draw","draw",1358331674)],[new cljs.core.Keyword(null,"key-pressed","key-pressed",-757100364),new cljs.core.Keyword(null,"mouse-exited","mouse-exited",-483205244),new cljs.core.Keyword(null,"mouse-wheel","mouse-wheel",1811662439),new cljs.core.Keyword(null,"mouse-dragged","mouse-dragged",-1220073441),new cljs.core.Keyword(null,"setup","setup",1987730512),new cljs.core.Keyword(null,"key-released","key-released",215919828),new cljs.core.Keyword(null,"mouse-clicked","mouse-clicked",-199339421),new cljs.core.Keyword(null,"mouse-released","mouse-released",-664480061),new cljs.core.Keyword(null,"mouse-pressed","mouse-pressed",736955536),new cljs.core.Keyword(null,"mouse-moved","mouse-moved",-1918152310),new cljs.core.Keyword(null,"mouse-entered","mouse-entered",811350322),new cljs.core.Keyword(null,"key-typed","key-typed",-876037597),new cljs.core.Keyword(null,"draw","draw",1358331674)]));
var chunk__8210 = null;
var count__8211 = (0);
var i__8212 = (0);
while(true){
if((i__8212 < count__8211)){
var vec__8213 = cljs.core._nth.call(null,chunk__8210,i__8212);
var processing_name = cljs.core.nth.call(null,vec__8213,(0),null);
var quil_name = cljs.core.nth.call(null,vec__8213,(1),null);
var temp__4657__auto___8221 = opts.call(null,quil_name);
if(cljs.core.truth_(temp__4657__auto___8221)){
var handler_8222 = temp__4657__auto___8221;
(prc[cljs.core.name.call(null,processing_name)] = ((function (seq__8209,chunk__8210,count__8211,i__8212,handler_8222,temp__4657__auto___8221,vec__8213,processing_name,quil_name){
return (function (){
var _STAR_applet_STAR_8216 = quil.sketch._STAR_applet_STAR_;
quil.sketch._STAR_applet_STAR_ = prc;

try{return handler_8222.call(null);
}finally {quil.sketch._STAR_applet_STAR_ = _STAR_applet_STAR_8216;
}});})(seq__8209,chunk__8210,count__8211,i__8212,handler_8222,temp__4657__auto___8221,vec__8213,processing_name,quil_name))
);
} else {
}

var G__8223 = seq__8209;
var G__8224 = chunk__8210;
var G__8225 = count__8211;
var G__8226 = (i__8212 + (1));
seq__8209 = G__8223;
chunk__8210 = G__8224;
count__8211 = G__8225;
i__8212 = G__8226;
continue;
} else {
var temp__4657__auto__ = cljs.core.seq.call(null,seq__8209);
if(temp__4657__auto__){
var seq__8209__$1 = temp__4657__auto__;
if(cljs.core.chunked_seq_QMARK_.call(null,seq__8209__$1)){
var c__7633__auto__ = cljs.core.chunk_first.call(null,seq__8209__$1);
var G__8227 = cljs.core.chunk_rest.call(null,seq__8209__$1);
var G__8228 = c__7633__auto__;
var G__8229 = cljs.core.count.call(null,c__7633__auto__);
var G__8230 = (0);
seq__8209 = G__8227;
chunk__8210 = G__8228;
count__8211 = G__8229;
i__8212 = G__8230;
continue;
} else {
var vec__8217 = cljs.core.first.call(null,seq__8209__$1);
var processing_name = cljs.core.nth.call(null,vec__8217,(0),null);
var quil_name = cljs.core.nth.call(null,vec__8217,(1),null);
var temp__4657__auto___8231__$1 = opts.call(null,quil_name);
if(cljs.core.truth_(temp__4657__auto___8231__$1)){
var handler_8232 = temp__4657__auto___8231__$1;
(prc[cljs.core.name.call(null,processing_name)] = ((function (seq__8209,chunk__8210,count__8211,i__8212,handler_8232,temp__4657__auto___8231__$1,vec__8217,processing_name,quil_name,seq__8209__$1,temp__4657__auto__){
return (function (){
var _STAR_applet_STAR_8220 = quil.sketch._STAR_applet_STAR_;
quil.sketch._STAR_applet_STAR_ = prc;

try{return handler_8232.call(null);
}finally {quil.sketch._STAR_applet_STAR_ = _STAR_applet_STAR_8220;
}});})(seq__8209,chunk__8210,count__8211,i__8212,handler_8232,temp__4657__auto___8231__$1,vec__8217,processing_name,quil_name,seq__8209__$1,temp__4657__auto__))
);
} else {
}

var G__8233 = cljs.core.next.call(null,seq__8209__$1);
var G__8234 = null;
var G__8235 = (0);
var G__8236 = (0);
seq__8209 = G__8233;
chunk__8210 = G__8234;
count__8211 = G__8235;
i__8212 = G__8236;
continue;
}
} else {
return null;
}
}
break;
}
});
quil.sketch.make_sketch = (function quil$sketch$make_sketch(options){
var opts = cljs.core.merge.call(null,new cljs.core.PersistentArrayMap(null, 1, [new cljs.core.Keyword(null,"size","size",1098693007),new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [(500),(300)], null)], null),(function (p1__8237_SHARP_){
return p1__8237_SHARP_.call(null,options);
}).call(null,cljs.core.apply.call(null,cljs.core.comp,cljs.core.cons.call(null,quil.middlewares.deprecated_options.deprecated_options,new cljs.core.Keyword(null,"middleware","middleware",1462115504).cljs$core$IFn$_invoke$arity$2(options,cljs.core.PersistentVector.EMPTY)))));
var sketch_size = (function (){var or__6814__auto__ = new cljs.core.Keyword(null,"size","size",1098693007).cljs$core$IFn$_invoke$arity$1(opts);
if(cljs.core.truth_(or__6814__auto__)){
return or__6814__auto__;
} else {
return new cljs.core.PersistentVector(null, 2, 5, cljs.core.PersistentVector.EMPTY_NODE, [(200),(200)], null);
}
})();
var renderer = new cljs.core.Keyword(null,"renderer","renderer",336841071).cljs$core$IFn$_invoke$arity$1(opts);
var features = cljs.core.set.call(null,new cljs.core.Keyword(null,"features","features",-1146962336).cljs$core$IFn$_invoke$arity$1(opts));
var setup = ((function (opts,sketch_size,renderer,features){
return (function (){
cljs.core.apply.call(null,quil.sketch.size,cljs.core.concat.call(null,sketch_size,(cljs.core.truth_(renderer)?new cljs.core.PersistentVector(null, 1, 5, cljs.core.PersistentVector.EMPTY_NODE, [renderer], null):cljs.core.PersistentVector.EMPTY)));

if(cljs.core.truth_(new cljs.core.Keyword(null,"settings","settings",1556144875).cljs$core$IFn$_invoke$arity$1(opts))){
new cljs.core.Keyword(null,"settings","settings",1556144875).cljs$core$IFn$_invoke$arity$1(opts).call(null);
} else {
}

if(cljs.core.truth_(new cljs.core.Keyword(null,"setup","setup",1987730512).cljs$core$IFn$_invoke$arity$1(opts))){
return new cljs.core.Keyword(null,"setup","setup",1987730512).cljs$core$IFn$_invoke$arity$1(opts).call(null);
} else {
return null;
}
});})(opts,sketch_size,renderer,features))
;
var mouse_wheel = (cljs.core.truth_(new cljs.core.Keyword(null,"mouse-wheel","mouse-wheel",1811662439).cljs$core$IFn$_invoke$arity$1(opts))?((function (opts,sketch_size,renderer,features,setup){
return (function (){
return new cljs.core.Keyword(null,"mouse-wheel","mouse-wheel",1811662439).cljs$core$IFn$_invoke$arity$1(opts).call(null,((-1) * quil.sketch._STAR_applet_STAR_.mouseScroll));
});})(opts,sketch_size,renderer,features,setup))
:null);
var opts__$1 = cljs.core.assoc.call(null,opts,new cljs.core.Keyword(null,"setup","setup",1987730512),setup,new cljs.core.Keyword(null,"mouse-wheel","mouse-wheel",1811662439),mouse_wheel);
var attach_function = ((function (opts,sketch_size,renderer,features,setup,mouse_wheel,opts__$1){
return (function (prc){
quil.sketch.bind_handlers.call(null,prc,opts__$1);

prc.quil = cljs.core.atom.call(null,null);

return prc.target_frame_rate = cljs.core.atom.call(null,(60));
});})(opts,sketch_size,renderer,features,setup,mouse_wheel,opts__$1))
;
var sketch = (new Processing.Sketch(attach_function));
if(cljs.core.contains_QMARK_.call(null,features,new cljs.core.Keyword(null,"global-key-events","global-key-events",335064944))){
((sketch["options"])["globalKeyEvents"] = true);
} else {
}

return sketch;
});
quil.sketch.destroy_previous_sketch = (function quil$sketch$destroy_previous_sketch(host_elem){
var temp__4657__auto__ = host_elem.processing_obj;
if(cljs.core.truth_(temp__4657__auto__)){
var proc_obj = temp__4657__auto__;
return proc_obj.exit();
} else {
return null;
}
});
quil.sketch.sketch = (function quil$sketch$sketch(var_args){
var args__7934__auto__ = [];
var len__7927__auto___8239 = arguments.length;
var i__7928__auto___8240 = (0);
while(true){
if((i__7928__auto___8240 < len__7927__auto___8239)){
args__7934__auto__.push((arguments[i__7928__auto___8240]));

var G__8241 = (i__7928__auto___8240 + (1));
i__7928__auto___8240 = G__8241;
continue;
} else {
}
break;
}

var argseq__7935__auto__ = ((((0) < args__7934__auto__.length))?(new cljs.core.IndexedSeq(args__7934__auto__.slice((0)),(0),null)):null);
return quil.sketch.sketch.cljs$core$IFn$_invoke$arity$variadic(argseq__7935__auto__);
});

quil.sketch.sketch.cljs$core$IFn$_invoke$arity$variadic = (function (opts){
var opts_map = cljs.core.apply.call(null,cljs.core.hash_map,opts);
var host_elem = goog.dom.getElement(new cljs.core.Keyword(null,"host","host",-1558485167).cljs$core$IFn$_invoke$arity$1(opts_map));
var renderer = (function (){var or__6814__auto__ = new cljs.core.Keyword(null,"renderer","renderer",336841071).cljs$core$IFn$_invoke$arity$1(opts_map);
if(cljs.core.truth_(or__6814__auto__)){
return or__6814__auto__;
} else {
return new cljs.core.Keyword(null,"p2d","p2d",-2106175755);
}
})();
if(cljs.core.truth_(host_elem)){
if(cljs.core.truth_(host_elem.processing_context)){
if(cljs.core._EQ_.call(null,renderer,host_elem.processing_context)){
} else {
console.warn("WARNING: Using different context on one canvas!");
}
} else {
host_elem.processing_context = renderer;
}

quil.sketch.destroy_previous_sketch.call(null,host_elem);

return host_elem.processing_obj = (new Processing(host_elem,quil.sketch.make_sketch.call(null,opts_map)));
} else {
return console.error("ERROR: Cannot create sketch. :host is not specified.");
}
});

quil.sketch.sketch.cljs$lang$maxFixedArity = (0);

quil.sketch.sketch.cljs$lang$applyTo = (function (seq8238){
return quil.sketch.sketch.cljs$core$IFn$_invoke$arity$variadic(cljs.core.seq.call(null,seq8238));
});

quil.sketch.sketch_init_list = cljs.core.atom.call(null,cljs.core.List.EMPTY);
quil.sketch.empty_body_QMARK_ = (function quil$sketch$empty_body_QMARK_(){
var child = document.body.childNodes;
return (child.length <= (1));
});
quil.sketch.add_canvas = (function quil$sketch$add_canvas(canvas_id){
var canvas = document.createElement("canvas");
canvas.setAttribute("id",canvas_id);

return document.body.appendChild(canvas);
});
quil.sketch.init_sketches = (function quil$sketch$init_sketches(){
var add_elem_QMARK__8250 = quil.sketch.empty_body_QMARK_.call(null);
var seq__8246_8251 = cljs.core.seq.call(null,cljs.core.deref.call(null,quil.sketch.sketch_init_list));
var chunk__8247_8252 = null;
var count__8248_8253 = (0);
var i__8249_8254 = (0);
while(true){
if((i__8249_8254 < count__8248_8253)){
var sk_8255 = cljs.core._nth.call(null,chunk__8247_8252,i__8249_8254);
if(cljs.core.truth_(add_elem_QMARK__8250)){
quil.sketch.add_canvas.call(null,new cljs.core.Keyword(null,"host-id","host-id",742376279).cljs$core$IFn$_invoke$arity$1(sk_8255));
} else {
}

new cljs.core.Keyword(null,"fn","fn",-1175266204).cljs$core$IFn$_invoke$arity$1(sk_8255).call(null);

var G__8256 = seq__8246_8251;
var G__8257 = chunk__8247_8252;
var G__8258 = count__8248_8253;
var G__8259 = (i__8249_8254 + (1));
seq__8246_8251 = G__8256;
chunk__8247_8252 = G__8257;
count__8248_8253 = G__8258;
i__8249_8254 = G__8259;
continue;
} else {
var temp__4657__auto___8260 = cljs.core.seq.call(null,seq__8246_8251);
if(temp__4657__auto___8260){
var seq__8246_8261__$1 = temp__4657__auto___8260;
if(cljs.core.chunked_seq_QMARK_.call(null,seq__8246_8261__$1)){
var c__7633__auto___8262 = cljs.core.chunk_first.call(null,seq__8246_8261__$1);
var G__8263 = cljs.core.chunk_rest.call(null,seq__8246_8261__$1);
var G__8264 = c__7633__auto___8262;
var G__8265 = cljs.core.count.call(null,c__7633__auto___8262);
var G__8266 = (0);
seq__8246_8251 = G__8263;
chunk__8247_8252 = G__8264;
count__8248_8253 = G__8265;
i__8249_8254 = G__8266;
continue;
} else {
var sk_8267 = cljs.core.first.call(null,seq__8246_8261__$1);
if(cljs.core.truth_(add_elem_QMARK__8250)){
quil.sketch.add_canvas.call(null,new cljs.core.Keyword(null,"host-id","host-id",742376279).cljs$core$IFn$_invoke$arity$1(sk_8267));
} else {
}

new cljs.core.Keyword(null,"fn","fn",-1175266204).cljs$core$IFn$_invoke$arity$1(sk_8267).call(null);

var G__8268 = cljs.core.next.call(null,seq__8246_8261__$1);
var G__8269 = null;
var G__8270 = (0);
var G__8271 = (0);
seq__8246_8251 = G__8268;
chunk__8247_8252 = G__8269;
count__8248_8253 = G__8270;
i__8249_8254 = G__8271;
continue;
}
} else {
}
}
break;
}

return cljs.core.reset_BANG_.call(null,quil.sketch.sketch_init_list,cljs.core.PersistentVector.EMPTY);
});
quil.sketch.add_sketch_to_init_list = (function quil$sketch$add_sketch_to_init_list(sk){
cljs.core.swap_BANG_.call(null,quil.sketch.sketch_init_list,cljs.core.conj,sk);

if(cljs.core._EQ_.call(null,document.readyState,"complete")){
return quil.sketch.init_sketches.call(null);
} else {
return null;
}
});
goog.events.listenOnce(window,goog.events.EventType.LOAD,quil.sketch.init_sketches);

//# sourceMappingURL=sketch.js.map