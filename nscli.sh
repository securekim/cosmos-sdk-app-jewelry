#/bin/bash
code=code13
_PATH=`pwd`
cd bin
./nscli query account $($_PATH/bin/nscli keys show gemologist -a)
./nscli tx nameservice buy-code $code 5nametoken --from gemologist
sleep 5
./nscli tx nameservice set-code $code carat1 cut1 clarity1 color1 fluorescence1 --from gemologist
#bin/nscli query nameservice resolve $code
sleep 5
./nscli query nameservice whichis $code
./nscli tx nameservice buy-code $code 100nametoken --from gemologist
sleep 3
./nscli query nameservice whichis $code
sleep 2
./nscli tx nameservice buy-code $code 100nametoken --from wholesaler
sleep 3
./nscli query nameservice whichis $code
#bin/nscli query nameservice resolve $code