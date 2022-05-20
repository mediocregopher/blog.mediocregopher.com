(ns viz.grid)

;; grid     set of points relative to a common origin

(def euclidean [       [0 -1]
                [-1 0]   ,    [1 0]
                       [0 1]       ])

(def isometric [[-1 -1] [0 -2] [1 -1]
                          ,
                [-1 1]  [0 2]  [1 1]])

(def hexagonal [        [0 -1]
                          ,
                [-1 1]         [1 1]])

(defn new-grid [grid-def]
  { :grid-def grid-def
    :points   #{} })

(defn add-point [grid point]
  (update-in grid [:points] conj point))

(def my-grid (-> (new-grid euclidean)
                 (add-point [0 1])))

;; TODO this could be useful, but it's not needed now, as long as all points we
;; use are generated from adjacent-points
;;(defn valid-point? [grid point]
;;  (letfn [(ordered-dim-points [dim order]
;;            (->> (:grid-def grid)
;;                 (map #(%1 dim))
;;                 (sort (if (= order :asc) < > ))
;;                 (filter (if (= order :asc) #(> %1 0) #(< %1 0)))
;;                 ))
;;          (closest-in-dim [dim-i dim-jj]
;;            (reduce (fn [curr dim-j]
;;                      (let [next (+ curr dim-j)]
;;                        (reduce #(if (= ;; TODO wat
;;                                   (if (> 0 dim-i) 
;;                                     (min dim-i next)
;;                                     (max dim-i next))))
;;                    0 dim-jj))
;;
;;          ]
;;    (closest-in-dim 4 [1])))
;;    ;;(ordered-dim 1 :asc)))
;;
;;(valid-point? my-grid [0 1])

(defn rm-point [grid point]
  (update-in grid [:points] disj point))

(defn adjacent-points [grid point]
  (map #(map + %1 point) (:grid-def grid)))

(defn empty-adjacent-points [grid point]
  (remove (:points grid) (adjacent-points grid point)))

(-> (new-grid isometric)
    (add-point [0 0])
    (add-point [0 1])
    (empty-adjacent-points [0 1]))
