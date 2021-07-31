(ns viz.dial
  (:require [quil.core :as q]))

(defn new-dial []
  {:val 0
   :min -1
   :max 1
   })

(defn- scale [v old-min old-max new-min new-max]
  (+ new-min (* (- new-max new-min)
                (/ (- v old-min) (- old-max old-min)))))

(defn scaled [dial min max]
  (let [new-val (scale (:val dial) (:min dial) (:max dial) min max)]
    (assoc dial :min min :max max :val new-val)))

(defn floored [dial at]
  (if (< (:val dial) at)
    (assoc dial :val at)
    dial))

(defn invert [dial]
  (assoc dial :val (* -1 (:val dial))))

;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;

; plot is a sequence of [t val], where t is a normalized time value between 0
; and 1, and val is the value the plot should become at that point.
(defn new-plot [frame-rate period-seconds plot]
  {:frame-rate frame-rate
   :period period-seconds
   :plot plot})

(defn by-plot [dial plot curr-frame]
  (let [dial-t (/ (mod (/ curr-frame (:frame-rate plot)) (:period plot)) (:period plot))
        ]
    (assoc dial :val
           (reduce
             (fn [curr-v [t v]] (if (<= t dial-t) v (reduced curr-v)))
             0 (:plot plot)))
    ))
