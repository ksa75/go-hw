go build -o go-cp.exe

./go-cp.exe -from testdata/input.txt -to out.txt
diff out.txt testdata/out_offset0_limit0.txt

./go-cp.exe -from testdata/input.txt -to out.txt -limit 10
diff out.txt testdata/out_offset0_limit10.txt

./go-cp.exe -from testdata/input.txt -to out.txt -limit 1000
diff out.txt testdata/out_offset0_limit1000.txt

./go-cp.exe -from testdata/input.txt -to out.txt -limit 10000
diff out.txt testdata/out_offset0_limit10000.txt

./go-cp.exe -from testdata/input.txt -to out.txt -offset 100 -limit 1000
diff out.txt testdata/out_offset100_limit1000.txt

./go-cp.exe -from testdata/input.txt -to out.txt -offset 6000 -limit 1000
diff out.txt testdata/out_offset6000_limit1000.txt

rm -force go-cp.exe
rm -force out.txt

