(ns viz.core
  (:require [quil.core :as q :include-macros true]
            [quil.middleware :as m]
            [viz.forest :as forest]
            [viz.grid :as grid]
            [viz.ghost :as ghost]
            [goog.string :as gstring]
            [goog.string.format]
            ;[gil.core :as gil]
            ))

(defn- debug [& args]
  (.log js/console (clojure.string/join " " (map str args))))

;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;

(defn- window-partial [k]
  (int (* (aget js/document "documentElement" k) 0.95)))

(def window-size [ (min 1025 (window-partial "clientWidth"))
                   (int (* (window-partial "clientHeight") 0.75))
                   ])
(def window-half-size (apply vector (map #(float (/ %1 2)) window-size)))

(defn- new-state []
  {:frame-rate 15
   :exit-wait-frames 40
   :tail-length 15
   :frame 0
   :gif-seconds 0
   :grid-width 30 ; from the center
   :ghost (-> (ghost/new-ghost grid/euclidean)
              (ghost/new-active-node [0 0])
              )
   })

(defn- curr-second [state]
  (float (/ (:frame state) (:frame-rate state))))

(defn- grid-size [state]
  (let [h (int (* (window-size 1)
                  (float (/ (:grid-width state) (window-size 0)))))]
             [(:grid-width state) h]))

(defn- positive [n] (if (> 0 n) (- n) n))

(defn- spawn-chance [state]
  (let [period-seconds 1
        period-frames (* (:frame-rate state) period-seconds)]
  (if (zero? (rem (:frame state) period-frames))
    1 100)
    ))

  ;(let [period-seconds 1
  ;      rad-per-second (float (/ (/ Math/PI 2) period-seconds))
  ;      rad (* rad-per-second (curr-second state))
  ;      chance-raw (positive (q/sin rad))
  ;      ]
  ;    (if (> chance-raw 0.97) 3 50)
  ;  ))


;(defn- mk-poss-fn [state]
;  (let [chance (spawn-chance state)]
;    (fn [pos adj-poss]
;      (if (zero? (rand-int chance))
;        adj-poss
;        (take 1 (shuffle adj-poss))))
;    ))

(defn- mk-poss-fn [state]
  (fn [pos adj-poss]
    (take 2 (random-sample 0.6 adj-poss))))

(defn setup []
  (let [state (new-state)]
    (q/frame-rate (:frame-rate state))
    state))

(defn- scale [state xy]
  (map-indexed #(* %2 (float (/ (window-half-size %1)
                                ((grid-size state) %1)))) xy))

; each bound is a position vector
(defn- in-bounds? [min-bound max-bound pos]
  (let [pos-k (keep-indexed #(let [mini (min-bound %1)
                                   maxi (max-bound %1)]
                               (when (and (>= %2 mini) (<= %2 maxi)) %2)) pos)]
    (= (count pos) (count pos-k))))

(defn- quil-bounds [state buffer]
  (let [[w h] (apply vector (map #(- % buffer) (grid-size state)))]
    [[(- w) (- h)] [w h]]))

(defn- ghost-incr [state]
  (assoc state :ghost
         (ghost/filter-active-nodes (ghost/incr (:ghost state) (mk-poss-fn state))
                                    #(let [[minb maxb] (quil-bounds state 2)]
                                       (in-bounds? minb maxb (:pos %1))))))

(defn- ghost-expire-roots [state]
  (if-not (< (:tail-length state) (:frame state)) state
    (update-in state [:ghost] ghost/remove-roots)))

(defn- maybe-exit [state]
  (if (empty? (get-in state [:ghost :active-node-ids]))
    (if (zero? (:exit-wait-frames state)) (new-state)
      (update-in state [:exit-wait-frames] dec))
    state))

(defn update-state [state]
  (-> state
      (ghost-incr)
      (ghost-expire-roots)
      (update-in [:frame] inc)
      (maybe-exit)))

(defn- draw-ellipse [state pos size] ; size is [w h]
  (let [scaled-pos (scale state pos)
        scaled-size (map int (scale state size))]
    (apply q/ellipse (concat scaled-pos scaled-size))))

(defn- in-line? [& nodes]
  (apply = (map #(apply map - %1)
                (partition 2 1 (map :pos nodes)))))

(defn draw-lines [state forest parent node]
  "Draws the lines of all children leading from the node, recursively"
  (q/stroke 0xFF000000)
  (q/fill 0xFFFFFFFF)
  (let [children (map #(forest/get-node forest %) (:child-ids node))]

    (if-not parent
      (doseq [child children] (draw-lines state forest node child))
      (let [in-line-child (some #(if (in-line? parent node %) %) children)
            ]
        (doseq [child children]
          (if (and in-line-child (= in-line-child child))
            (draw-lines state forest parent child)
            (draw-lines state forest node child)))
        (when-not in-line-child
          (apply q/line (apply concat
                               (map #(scale state %)
                                    (map :pos (list parent node))))))
        ))

    ; we also take the opportunity to draw the leaves
    (when (empty? children)
      (draw-ellipse state (:pos node) [0.3 0.3]))

    ))

(defn draw-state [state]
  ; Clear the sketch by filling it with light-grey color.
  (q/background 0xFFFFFFFF)
  (q/with-translation [(/ (window-size 0) 2)
                       (/ (window-size 1) 2)]
    (let [lines (forest/lines (get-in state [:ghost :forest]))
          leaves (forest/leaves (get-in state [:ghost :forest]))
          active (ghost/active-nodes (:ghost state))
          roots (forest/roots (get-in state [:ghost :forest]))
          ]

      (q/stroke 0xFF000000)
      (doseq [root roots]
        (draw-lines state (get-in state [:ghost :forest]) nil root))

      (q/stroke 0xFF000000)
      (q/fill 0xFF000000)
      (doseq [active-node active]
        (let [pos (:pos active-node)]
          (draw-ellipse state pos [0.35 0.35])
          ))

      ))

    ;(when-not (zero? (:gif-seconds state))
    ;  (let [anim-frames (* (:gif-seconds state) (:frame-rate state))]
    ;    (gil/save-animation "quil.gif" anim-frames 0)
    ;    (when (> (:frame state) anim-frames) (q/exit))))

    ;(q/text (clojure.string/join
    ;          "\n"
    ;          (list
    ;            (gstring/format "frame:%d" (:frame state))
    ;            (gstring/format "second:%f" (curr-second state))
    ;            (gstring/format "spawn-chance:%d" (spawn-chance state))))
    ;        30 30)
  )

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
