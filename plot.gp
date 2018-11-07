set datafile separator comma
set key off
set term png
do for [ii=1:365*10] {
     set output sprintf("animation/%04d.png",ii)
    splot [-1e10:1e10] 'sonnensystem.csv' using 0:1:2 every ::1::ii with l, '' using 3:4:5 every ::1::ii with l, \
     '' using 6:7:8 every ::1::ii w l, '' using 9:10:11 every ::1::ii w l, '' using 12:13:14 every ::1::ii w l, \
     '' using 15:16:17 every ::1::ii with l, '' using 18:19:20 every ::1::ii w l, '' using 21:22:23 every ::1::ii w l, \
     '' using 24:25:26 every ::1::ii w l, '' using 27:28:29 every ::1::ii with l, \
     'sonnensystem.csv' using 0:1:2 every ::ii::ii w p, '' using 3:4:5 every ::ii::ii w p, \
                   '' using 6:7:8 every ::ii::ii w p, '' using 9:10:11 every ::ii::ii w p, '' using 12:13:14 every ::ii::ii w p, \
                   '' using 15:16:17 every ::ii::ii w p, '' using 18:19:20 every ::ii::ii w p, '' using 21:22:23 every ::ii::ii w p, \
                   '' using 24:25:26 every ::ii::ii w p, '' using 27:28:29 every ::ii::ii w p,
     replot
}
system('ffmpeg -f image2 -framerate 150 -i %04d.png out.mp4')