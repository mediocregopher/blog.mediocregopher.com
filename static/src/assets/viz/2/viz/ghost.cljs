(ns viz.ghost
  (:require [quil.core :as q]
            [quil.middleware :as m]
            [viz.forest :as forest]
            [viz.grid   :as grid]
            clojure.set))

(defn new-ghost []
  {:active-node-ids #{}
   :color 0xFF000000
   })

(defn add-active-node [ghost id]
  (update-in ghost [:active-node-ids] conj id))

(defn rm-active-node [ghost id]
  (update-in ghost [:active-node-ids] disj id))

(defn- gen-new-poss [forest poss-fn id]
  "generates new positions branching from the given node"
  (let [pos (:pos (forest/get-node forest id))
        adj-poss (forest/empty-adjacent-points forest pos)
        new-poss (poss-fn pos adj-poss)]
    new-poss))

(defn- spawn-children [forest poss-fn id]
  (reduce (fn [[forest new-ids] pos]
            (let [[forest new-id] (forest/spawn-child forest id pos)]
              [forest (conj new-ids new-id)]))
          [forest #{}]
          (gen-new-poss forest poss-fn id)))

(defn- spawn-children-multi [forest poss-fn ids]
  (reduce (fn [[forest new-ids] id]
            (let [[forest this-new-ids] (spawn-children forest poss-fn id)]
              [forest (clojure.set/union new-ids this-new-ids)]))
          [forest #{}]
          ids))

(defn incr [ghost forest poss-fn]
  (let [[forest new-ids] (spawn-children-multi forest poss-fn (:active-node-ids ghost))]
    [(assoc ghost :active-node-ids new-ids)
     (reduce (fn [forest id]
               (forest/update-node-meta forest id
                  (fn [m] (assoc m :color (:color ghost)))))
             forest new-ids)]))

(defn- eg-poss-fn [pos adj-poss]
  (take 2 (random-sample 0.6 adj-poss)))
