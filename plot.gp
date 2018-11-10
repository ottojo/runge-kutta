set datafile separator comma

set key off

unset autoscale
set xrange [-2e8:2e8]
set yrange [-6e8:6e8]
set zrange [-3e8:3e8]

set xyplane 0
set tics front
set ytics offset 0,-2,0 2e8
unset xtics

set terminal pngcairo size 1920,1080
set view equal xyz
set view 90,90,3
do for [ii=startindex:endindex] {
    set output sprintf("animation/%6d.png",ii)
    splot  \
    filename using 0:1:2 every ::1::ii w l, '' using 3:4:5 every ::1::ii w l, \
        '' using 6:7:8 every ::1::ii w l, '' using 9:10:11 every ::1::ii w l, '' using 12:13:14 every ::1::ii w l, \
        '' using 15:16:17 every ::1::ii w l, '' using 18:19:20 every ::1::ii w l, '' using 21:22:23 every ::1::ii w l, \
        '' using 24:25:26 every ::1::ii w l, '' using 27:28:29 every ::1::ii w l, \
    filename using 0:1:2 every ::ii::ii w p, '' using 3:4:5 every ::ii::ii w p, \
        '' using 6:7:8 every ::ii::ii w p, '' using 9:10:11 every ::ii::ii w p, '' using 12:13:14 every ::ii::ii w p, \
        '' using 15:16:17 every ::ii::ii w p, '' using 18:19:20 every ::ii::ii w p, '' using 21:22:23 every ::ii::ii w p, \
        '' using 24:25:26 every ::ii::ii w p, '' using 27:28:29 every ::ii::ii w p,
}
