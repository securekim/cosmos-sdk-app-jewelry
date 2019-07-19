#/bin/bash
bin/nscli query account $(nscli keys show jack -a)
bin/nscli tx nameservice buy-code code1 5nametoken --from jack
bin/nscli tx nameservice set-code code1 carat1 cut1 clarity1 color1 fluorescence1 --from jack
bin/nscli query nameservice resolve code1
bin/nscli query nameservice whichis code1
