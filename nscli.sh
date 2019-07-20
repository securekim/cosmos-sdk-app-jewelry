#/bin/bash
code=code12
bin/nscli query account $(nscli keys show jack -a)
bin/nscli tx nameservice buy-code $code 5nametoken --from jack
sleep 3
bin/nscli tx nameservice set-code $code carat1 cut1 clarity1 color1 fluorescence1 --from jack
#bin/nscli query nameservice resolve $code
sleep 5
bin/nscli query nameservice whichis $code
