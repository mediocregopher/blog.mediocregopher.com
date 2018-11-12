(ns viz.forest)

(defn new-forest []
  {:nodes {}
   :roots #{}
   :leaves #{}
   :next-id 0})

(defn- new-id [forest]
  (let [id (:next-id forest)]
    [(assoc forest :next-id (inc id))
     id]))

(defn- unset-parent [forest id parent-id]
  (-> forest
      (update-in [:nodes id] dissoc :parent-id :parent-pos)
      (update-in [:nodes parent-id :child-ids] disj id)
      (update-in [:roots] conj id)
      (update-in [:leaves] conj parent-id)
      ))

(defn- set-parent [forest id parent-id]
  (let [parent-pos (get-in forest [:nodes parent-id :pos])
        prev-parent-id (get-in forest [:nodes id :parent-id])
        ]
    (-> forest
        (assoc-in [:nodes id :parent-id] parent-id)
        (assoc-in [:nodes id :parent-pos] parent-pos)
        (update-in [:nodes parent-id :child-ids] #(if %1 (conj %1 id) #{id}))
        (update-in [:roots] disj id)
        (update-in [:leaves] disj parent-id)
        ;; If there was a previous parent of the child, unset that shit
        (#(if prev-parent-id (unset-parent %1 id prev-parent-id) %1))
        )))

(defn add-node [forest pos]
  (let [[forest id] (new-id forest)
        forest (-> forest
                   (assoc-in [:nodes id] {:id id :pos pos})
                   (update-in [:roots] conj id)
                   (update-in [:leaves] conj id)
                   )
        ]
    [forest id]))

(defn remove-node [forest id]
  (let [child-ids (get-in forest [:nodes id :child-ids])
        parent-id (get-in forest [:nodes id :parent-id])]
    (-> forest
        ;; unset this node's parent, if it has one
        (#(if parent-id (unset-parent %1 id parent-id) %1))
        ;; unset this node's children, if it has any
        ((fn [forest] (reduce #(unset-parent %1 %2 id) forest child-ids)))
        ;; remove from all top-level sets
        (update-in [:nodes] dissoc id)
        (update-in [:roots] disj id)
        (update-in [:leaves] disj id)
        )))

(defn get-node [forest id]
  (get-in forest [:nodes id]))

(defn spawn-child [forest parent-id pos]
  (let [[forest id] (add-node forest pos)
        forest (-> forest
                   (set-parent id parent-id)
                   )
        ]
    [forest id]))

(defn roots [forest] (-> forest :nodes (select-keys (:roots forest)) (vals)))
(defn root? [node] (not (boolean (:parent-id node))))

(defn leaves [forest] (-> forest :nodes (select-keys (:leaves forest)) (vals)))
(defn leaf? [node] (empty? (:child-ids node)))

(defn lines [forest]
  (->> forest
       (:nodes)
       (vals)
       (remove #(empty? (:parent-pos %)))
       (map #(vector (:pos %) (:parent-pos %)))
       ))

(def my-forest
  (let [forest (new-forest)
        [forest id0] (add-node forest [0 0])
        [forest id1] (spawn-child forest id0 [1 1])
        [forest id2] (spawn-child forest id0 [-1 -1])
        [forest id3] (spawn-child forest id1 [2 2])
        forest (remove-node forest id1)
        ]
    forest))

(identity my-forest)
(lines my-forest)
