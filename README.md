![](out.gif)

i2g
---
Basic utility to convert images to gifs.

Usage of gif
------------
```
-d int
```
Time delay between frames, in 100ths of a second. (default 4)
```
-lc int
```
Controls the number of times an animation will be restarted during display.
* 0 - loops forever(default),
* -1 - shows each frame once,
* n - shows each frame n+1 times.
```
-o string
```
Output file. (default "out.gif")
```
-p string
```
Relative path to sequence of images. (default "frames/")