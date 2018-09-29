#!/bin/sh

# requires imagemagick and perl-image-exiftool

set -e

widths="500 1000 1500 2000 2500 3000"

for img in $@; do
    echo $img

    # make target directories
    dir=$(dirname "$img") # gets directory
    for targetWidth in $widths; do
        mkdir -p $dir/${targetWidth}px
    done

    # get width
    width=$(identify "$img" | awk '{print $3}' | cut -dx -f1)
    echo -e "\toriginal width: $width"

    echo -e "\tremoving metadata"
    exiftool -all= "$img"
    rm -f "${img}_original" # exiftool makes a copy of the original, delete it

    for targetWidth in $widths; do
        targetFile=$dir/${targetWidth}px/$(basename "$img")
        echo -en "\tresizing into $targetFile... "
        if [ "$targetWidth" -ge "$width" ]; then
            echo "skipping, original image too small"
            continue
        elif [ -e "$targetFile" ]; then
            echo "skipping, target file exists"
            continue
        fi
        convert "$img" -resize $targetWidth "$targetFile"
        echo "done"
    done
done
