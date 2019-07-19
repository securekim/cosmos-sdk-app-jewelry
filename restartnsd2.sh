#/bin/bash
./nscli query account $(nscli keys show jack -a)
./nscli tx nameservice set-code code1 carat1 cut1 clarity1 color1 fluorescence1 --from jack
./nscli query nameservice whichis code1
