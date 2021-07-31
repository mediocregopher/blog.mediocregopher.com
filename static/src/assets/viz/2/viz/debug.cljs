(ns viz.debug)

(defn- log [& args]
  (.log js/console (clojure.string/join " " (map str args))))
