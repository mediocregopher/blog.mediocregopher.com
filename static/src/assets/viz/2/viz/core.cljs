(ns viz.core
  (:require [quil.core :as q]
            [quil.middleware :as m]
            [viz.forest :as forest]
            [viz.grid :as grid]
            [viz.ghost :as ghost]
            [viz.dial :as dial]
            [goog.string :as gstring]
            [goog.string.format]
            ))

(defn- debug [& args]
  (.log js/console (clojure.string/join " " (map str args))))
(defn- observe [v] (debug v) v)

(defn- positive [n] (if (> 0 n) (- n) n))

;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;; initialization

(defn- window-partial [k]
  (int (aget js/document "documentElement" k)))

(def window-size
  (let [w (int (min 1024 (window-partial "clientWidth")))]
    [w (int (min (* w 0.75) (window-partial "clientHeight")))]))

(def window-half-size (apply vector (map #(float (/ %1 2)) window-size)))

(defn- set-grid-size [state]
  (let [h (int (* (window-size 1)
                  (float (/ (:grid-width state) (window-size 0)))))]
    (assoc state :grid-size [(:grid-width state) h])))

(defn- add-ghost [state ghost-def]
  (let [[forest id] (forest/add-node (:forest state) (:start-pos ghost-def))
        ghost       (-> (ghost/new-ghost)
                        (ghost/add-active-node id)
                        (assoc :ghost-def ghost-def))
        ]
    (assoc state
           :forest forest
           :ghosts (cons ghost (:ghosts state)))))

(defn- new-state []
  (-> {:frame-rate 15
       :color-cycle-period 8
       :tail-length 7
       :frame 0
       :grid-width 45 ; from the center
       :forest (forest/new-forest grid/isometric)
       }
      (set-grid-size)
      (add-ghost {:start-pos [-10 -10]
                  :color-fn (fn [state]
                              (let [frames-per-color-cycle
                                    (* (:color-cycle-period state) (:frame-rate state))]
                                (q/color
                                  (/ (mod (:frame state) frames-per-color-cycle)
                                     frames-per-color-cycle)
                                  1 1)))
                  })
      ))

(defn setup []
  (q/color-mode :hsb 1 1 1)
  (new-state))

;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;; scaling and unit conversion related

(defn- curr-second [state]
  (float (/ (:frame state) (:frame-rate state))))

(defn- scale [grid-size xy]
  (map-indexed #(* %2 (float (/ (window-half-size %1)
                                (grid-size %1)))) xy))

;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;; poss-fn

(def bounds-buffer 1)

(defn- in-bounds? [grid-size pos]
  (let [[w h] (apply vector (map #(- % bounds-buffer) grid-size))]
    (every?
      #(and (>= (% 1) (- (% 0))) (<= (% 1) (% 0)))
      (map vector [w h] pos))))

(defn- dist-from-sqr [pos1 pos2]
  (reduce + (map #(* % %) (map - pos1 pos2))))

(defn- dist-from [pos1 pos2]
  (q/sqrt (dist-from-sqr pos1 pos2)))

(defn take-adj-poss [grid-width pos adj-poss]
  (let [dist-from-center (dist-from [0 0] pos)
        width grid-width
        dist-ratio (/ (- width dist-from-center) width)
        ]
    (take
      (int (* (q/map-range (rand) 0 1 0.75 1)
              dist-ratio
              (count adj-poss)))
      adj-poss)))

(defn- mk-poss-fn [state]
  (let [grid-size (:grid-size state)]
    (fn [pos adj-poss]
      (->> adj-poss
           (filter #(in-bounds? grid-size %))
           (sort-by #(dist-from-sqr % [0 0]))
           (take-adj-poss (grid-size 0) pos)))))

;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;; update

(defn- update-ghost-forest [state update-fn]
  (let [[ghosts forest]
        (reduce (fn [[ghosts forest] ghost]
                  (let [[ghost forest] (update-fn ghost forest)]
                    [(cons ghost ghosts) forest]))
                [nil (:forest state)]
                (:ghosts state))]
    (assoc state :ghosts (reverse ghosts) :forest forest)))

(defn- ghost-incr [state poss-fn]
  (update-ghost-forest state #(ghost/incr %1 %2 poss-fn)))

(defn rm-nodes [state node-ids]
  (update-ghost-forest state (fn [ghost forest]
                               [(reduce ghost/rm-active-node ghost node-ids)
                                (reduce forest/remove-node forest node-ids)])))

(defn- maybe-remove-roots [state]
  (if (>= (:tail-length state) (:frame state))
    state
    (rm-nodes state (map :id (forest/roots (:forest state))))))

(defn- ghost-set-color [state]
  (update-ghost-forest state (fn [ghost forest]
                               (let [color ((get-in ghost [:ghost-def :color-fn]) state)]
                                 [(assoc ghost :color color) forest]))))

(defn update-state [state]
  (let [poss-fn (mk-poss-fn state)]
    (-> state
        (ghost-set-color)
        (ghost-incr poss-fn)
        (maybe-remove-roots)
        (update-in [:frame] inc)
        )))

;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;; draw

(defn- draw-ellipse [pos size scale-fn] ; size is [w h]
  (let [scaled-pos (scale-fn pos)
        scaled-size (map int (scale-fn size))]
    (apply q/ellipse (concat scaled-pos scaled-size))))

(defn- in-line? [& nodes]
  (apply = (map #(apply map - %1)
                (partition 2 1 (map :pos nodes)))))

(defn- draw-node [node active? scale-fn]
  (let [pos (:pos node)
        stroke (get-in node [:meta :color])
        fill   (if active? stroke 0xFFFFFFFF)
        ]
    (q/stroke stroke)
    (q/fill fill)
    (draw-ellipse pos [0.30 0.30] scale-fn)))

(defn- draw-line [node parent scale-fn]
  (let [node-color (get-in node [:meta :color])
        parent-color (get-in node [:meta :color])
        color (q/lerp-color node-color parent-color 0.5)
        ]
    (q/stroke color)
    (q/stroke-weight 1)
    (apply q/line (map scale-fn (map :pos (list parent node))))))

(defn- draw-lines [forest parent node scale-fn]
  "Draws the lines of all children leading from the node, recursively"
  (let [children (map #(forest/get-node forest %) (:child-ids node))]

    (if-not parent
      (doseq [child children] (draw-lines forest node child scale-fn))
      (let [in-line-child (some #(if (in-line? parent node %) %) children)
            ]
        (doseq [child children]
          (if (and in-line-child (= in-line-child child))
            (draw-lines forest parent child scale-fn)
            (draw-lines forest node child scale-fn)))
        (when-not in-line-child
          (draw-line node parent scale-fn))
        ))

    ; we also take the opportunity to draw the leaves
    (when (empty? children)
      (draw-node node false scale-fn))
    ))

(defn draw-dial [state dial posL posR]
  (let [dial-norm (q/norm (:val dial) (:min dial) (:max dial))
        dial-pos (map #(q/lerp %1 %2 dial-norm) posL posR)]
    (q/stroke 0xFF000000)
    (q/stroke-weight 1)
    (q/fill   0xFF000000)
    (apply q/line (concat posL posR))
    (apply q/ellipse (concat dial-pos [5 5]))
    ))

(defn draw-state [state]
  ; Clear the sketch by filling it with light-grey color.
  (q/background 0xFFFFFFFF)
  (q/with-translation window-half-size

    (let [grid-size (:grid-size state)
          scale-fn #(scale grid-size %)
          ghost (:ghost state)
          forest (:forest state)
          roots (forest/roots forest)
          ]

      (doseq [root roots]
        (draw-lines forest nil root scale-fn))

      (doseq [ghost (:ghosts state)]
        (doseq [active-node (map #(forest/get-node forest %)
                                 (:active-node-ids ghost))]
          (draw-node active-node true scale-fn)))

      ))

    ;(draw-dial state (:dial state) [30 30] [100 30])

    ;(q/text (clojure.string/join
    ;          "\n"
    ;          (list
    ;            (gstring/format "frame:%d" (:frame state))
    ;            (gstring/format "second:%f" (curr-second state))
    ;            (gstring/format "spawn-chance:%d" (spawn-chance state))))
    ;        30 30)
  )

;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;; def

(q/defsketch viz
  :title ""
  :host "viz"
  :size window-size
  ; setup function called only once, during sketch initialization.
  :setup setup
  ; update-state is called on each iteration before draw-state.
  :update update-state
  :draw draw-state
  :features [:keep-on-top]
  ; This sketch uses functional-mode middleware.
  ; Check quil wiki for more info about middlewares and particularly
  ; fun-mode.
  :middleware [m/fun-mode])
