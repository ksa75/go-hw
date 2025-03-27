go build -o go-cp.exe

./go-cp.exe -from testdata/input.txt -to out.txt
diff (gc out.txt) (gc testdata/out_offset0_limit0.txt)

./go-cp.exe -from testdata/input.txt -to out.txt -limit 10
diff  (gc out.txt) (gc testdata/out_offset0_limit10.txt)

./go-cp.exe -from testdata/input.txt -to out.txt -limit 1000
diff (gc out.txt) (gc testdata/out_offset0_limit1000.txt)

./go-cp.exe -from testdata/input.txt -to out.txt -limit 10000
diff (gc out.txt) (gc testdata/out_offset0_limit10000.txt)

./go-cp.exe -from testdata/input.txt -to out.txt -offset 100 -limit 1000
diff (gc out.txt) (gc testdata/out_offset100_limit1000.txt)

./go-cp.exe -from testdata/input.txt -to out.txt -offset 6000 -limit 1000
diff (gc out.txt) (gc testdata/out_offset6000_limit1000.txt)


./go-cp.exe -from testdata/input.txt -to out.txt -offset 6000 -limit 0
./go-cp.exe -from testdata/input.txt -to out.txt -offset 10000 -limit 10000
./go-cp.exe -from testdata/input.txt -to out.txt -offset 6616 -limit 100
./go-cp.exe -from testdata/input.txt -to out.txt -offset 6717 -limit 100
./go-cp.exe -from testdata/input.txt -to out.txt -offset 6917 -limit 100
./go-cp.exe -from not_existing -to out.txt
./go-cp.exe -from testdata/input.txt -to ?*.txt

rm -force go-cp.exe
rm -force out.txt

