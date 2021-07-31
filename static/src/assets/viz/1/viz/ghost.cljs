(ns viz.ghost
  (:require [viz.forest :as forest]
            [viz.grid   :as grid]
            clojure.set))

(defn new-ghost [grid-def]
  { :grid (grid/new-grid grid-def)
   :forest (forest/new-forest)
   :active-node-ids #{}
   })

(defn new-active-node [ghost pos]
  (let [[forest id] (forest/add-node (:forest ghost) pos)
        grid (grid/add-point (:grid ghost) pos)]
    (-> ghost
        (assoc :grid grid :forest forest)
        (update-in [:active-node-ids] conj id))))

(defn- gen-new-poss [ghost poss-fn id]
  "generates new positions branching from the given node"
  (let [pos (:pos (forest/get-node (:forest ghost) id))
        adj-poss (grid/empty-adjacent-points (:grid ghost) pos)]
    (poss-fn pos adj-poss)))

(defn- spawn-children [ghost poss-fn id]
  (reduce (fn [[ghost new-ids] pos]
            (let [[forest new-id] (forest/spawn-child (:forest ghost) id pos)
                  grid (grid/add-point (:grid ghost) pos)]
              [(assoc ghost :forest forest :grid grid) (conj new-ids new-id)]))
          [ghost #{}]
          (gen-new-poss ghost poss-fn id)))

(defn- spawn-children-multi [ghost poss-fn ids]
  (reduce (fn [[ghost new-ids] id]
            (let [[ghost this-new-ids] (spawn-children ghost poss-fn id)]
              [ghost (clojure.set/union new-ids this-new-ids)]))
          [ghost #{}]
          ids))

(defn incr [ghost poss-fn]
  (let [[ghost new-ids] (spawn-children-multi ghost poss-fn (:active-node-ids ghost))]
    (assoc ghost :active-node-ids new-ids)))

(defn active-nodes [ghost]
  (map #(get-in ghost [:forest :nodes %]) (:active-node-ids ghost)))

(defn filter-active-nodes [ghost pred]
  (assoc ghost :active-node-ids
         (reduce #(if (pred %2) (conj %1 (:id %2)) %1) #{}
                 (active-nodes ghost))))

(defn remove-roots [ghost]
  (let [roots (forest/roots (:forest ghost))
        root-ids (map :id roots)
        root-poss (map :pos roots)
        ]
    (-> ghost
        (update-in [:active-node-ids] #(reduce disj %1 root-ids))
        (update-in [:forest] #(reduce forest/remove-node %1 root-ids))
        (update-in [:grid] #(reduce grid/rm-point %1 root-poss))
        )))

(defn- eg-poss-fn [pos adj-poss]
  (take 2 (random-sample 0.6 adj-poss)))

(-> (new-ghost grid/euclidean)
    (new-active-node [0 0])
    (incr eg-poss-fn)
    (incr eg-poss-fn)
    (incr eg-poss-fn)
    (remove-roots)
    )
