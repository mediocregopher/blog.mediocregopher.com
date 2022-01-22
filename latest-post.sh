#!/bin/bash

postsDir=static/src/_posts

echo $postsDir/$(ls -1 static/src/_posts | sort -n | tail -n1)
