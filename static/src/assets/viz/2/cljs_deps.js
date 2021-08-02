goog.addDependency("base.js", ['goog'], []);
goog.addDependency("../cljs/core.js", ['cljs.core'], ['goog.string', 'goog.Uri', 'goog.object', 'goog.math.Integer', 'goog.string.StringBuffer', 'goog.array', 'goog.math.Long']);
goog.addDependency("../process/env.js", ['process.env'], ['cljs.core']);
goog.addDependency("../viz/grid.js", ['viz.grid'], ['cljs.core']);
goog.addDependency("../viz/forest.js", ['viz.forest'], ['cljs.core', 'viz.grid']);
goog.addDependency("../processing.js", ['org.processingjs.Processing'], [], {'foreign-lib': true});
goog.addDependency("../quil/middlewares/deprecated_options.js", ['quil.middlewares.deprecated_options'], ['cljs.core']);
goog.addDependency("../clojure/string.js", ['clojure.string'], ['goog.string', 'cljs.core', 'goog.string.StringBuffer']);
goog.addDependency("../quil/util.js", ['quil.util'], ['cljs.core', 'clojure.string']);
goog.addDependency("../quil/sketch.js", ['quil.sketch'], ['goog.dom', 'cljs.core', 'quil.middlewares.deprecated_options', 'goog.events.EventType', 'goog.events', 'quil.util']);
goog.addDependency("../quil/core.js", ['quil.core'], ['org.processingjs.Processing', 'quil.sketch', 'cljs.core', 'clojure.string', 'quil.util']);
goog.addDependency("../quil/middlewares/navigation_3d.js", ['quil.middlewares.navigation_3d'], ['cljs.core', 'quil.core']);
goog.addDependency("../quil/middlewares/navigation_2d.js", ['quil.middlewares.navigation_2d'], ['cljs.core', 'quil.core']);
goog.addDependency("../quil/middlewares/fun_mode.js", ['quil.middlewares.fun_mode'], ['cljs.core', 'quil.core']);
goog.addDependency("../quil/middleware.js", ['quil.middleware'], ['cljs.core', 'quil.middlewares.navigation_3d', 'quil.middlewares.navigation_2d', 'quil.middlewares.fun_mode']);
goog.addDependency("../clojure/set.js", ['clojure.set'], ['cljs.core']);
goog.addDependency("../viz/ghost.js", ['viz.ghost'], ['cljs.core', 'viz.forest', 'clojure.set', 'quil.core', 'quil.middleware', 'viz.grid']);
goog.addDependency("../viz/dial.js", ['viz.dial'], ['cljs.core', 'quil.core']);
goog.addDependency("../viz/core.js", ['viz.core'], ['goog.string', 'cljs.core', 'viz.forest', 'quil.core', 'quil.middleware', 'goog.string.format', 'viz.ghost', 'viz.grid', 'viz.dial']);