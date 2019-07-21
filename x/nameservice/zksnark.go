package nameservice
​
import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/arnaucube/go-snark"
	"github.com/arnaucube/go-snark/circuitcompiler"
	"github.com/arnaucube/go-snark/r1csqap"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
)
​
// ============= Cryptonian =============== //
​
func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}
​
func CompileCircuit(circuitCode string ,  privateInputs []*big.Int, publicInputs []*big.Int) error {
	fmt.Println("cli")
​
	// parse circuit code
	//parser := circuitcompiler.NewParser(bufio.NewReader(circuit))
	parser := circuitcompiler.NewParser( strings.NewReader(circuitCode))
	circuit, err := parser.Parse()
	panicErr(err)
	fmt.Println("\ncircuit data:", circuit)
​
	// calculate wittness
	w, err := circuit.CalculateWitness(privateInputs, publicInputs)
	panicErr(err)
	fmt.Println("\nwitness", w)
​
	// flat code to R1CS
	fmt.Println("\ngenerating R1CS from flat code")
	a, b, c := circuit.GenerateR1CS()
	fmt.Println("\nR1CS:")
	fmt.Println("a:", a)
	fmt.Println("b:", b)
	fmt.Println("c:", c)
​
	// R1CS to QAP
	alphas, betas, gammas, zx := snark.Utils.PF.R1CSToQAP(a, b, c)
	fmt.Println("qap")
	fmt.Println(alphas)
	fmt.Println(betas)
	fmt.Println(gammas)
​
	ax, bx, cx, px := snark.Utils.PF.CombinePolynomials(w, alphas, betas, gammas)
​
	hx := snark.Utils.PF.DivisorPolynomial(px, zx)
​
	// hx==px/zx so px==hx*zx
	// assert.Equal(t, px, snark.Utils.PF.Mul(hx, zx))
	if !r1csqap.BigArraysEqual(px, snark.Utils.PF.Mul(hx, zx)) {
		panic(errors.New("px != hx*zx"))
	}
​
	// p(x) = a(x) * b(x) - c(x) == h(x) * z(x)
	abc := snark.Utils.PF.Sub(snark.Utils.PF.Mul(ax, bx), cx)
	// assert.Equal(t, abc, px)
	if !r1csqap.BigArraysEqual(abc, px) {
		panic(errors.New("abc != px"))
	}
	hz := snark.Utils.PF.Mul(hx, zx)
	if !r1csqap.BigArraysEqual(abc, hz) {
		panic(errors.New("abc != hz"))
	}
	// assert.Equal(t, abc, hz)
​
	div, rem := snark.Utils.PF.Div(px, zx)
	if !r1csqap.BigArraysEqual(hx, div) {
		panic(errors.New("hx != div"))
	}
	// assert.Equal(t, hx, div)
	// assert.Equal(t, rem, r1csqap.ArrayOfBigZeros(4))
	for _, r := range rem {
		if !bytes.Equal(r.Bytes(), big.NewInt(int64(0)).Bytes()) {
			panic(errors.New("error:error:  px/zx rem not equal to zeros"))
		}
	}
​
	// store circuit to json
	jsonData, err := json.Marshal(circuit)
	panicErr(err)
	// store setup into file
	jsonFile, err := os.Create("compiledcircuit.json")
	panicErr(err)
	defer jsonFile.Close()
	jsonFile.Write(jsonData)
	jsonFile.Close()
	fmt.Println("Compiled Circuit data written to ", jsonFile.Name())
​
​
	// store px
	jsonData, err = json.Marshal(px)
	panicErr(err)
	// store setup into file
	jsonFile, err = os.Create("px.json")
	panicErr(err)
	defer jsonFile.Close()
	jsonFile.Write(jsonData)
	jsonFile.Close()
	fmt.Println("Px data written to ", jsonFile.Name())
​
	return nil
}
​
func CompileCircuitOnly(circuitCode string) error {
​
	// parse the code
	parser := circuitcompiler.NewParser(strings.NewReader(circuitCode))
	circuit, _ := parser.Parse()
​
	// store circuit to json
	jsonData, err := json.Marshal(circuit)
	panicErr(err)
	// store setup into file
	jsonFile, err := os.Create("compiledcircuit.json")
	panicErr(err)
	defer jsonFile.Close()
	jsonFile.Write(jsonData)
	jsonFile.Close()
​
	//fmt.Println("Compiled Circuit data written to ", jsonFile.Name())
​
	return nil
}
​
func TrustedSetup(privateInputs []*big.Int, publicInputs []*big.Int) error {
​
	// open compiledcircuit.json
	compiledcircuitFile, err := ioutil.ReadFile("compiledcircuit.json")
	panicErr(err)
	var circuit circuitcompiler.Circuit
	json.Unmarshal([]byte(string(compiledcircuitFile)), &circuit)
	panicErr(err)
​
	// calculate wittness
	w, err := circuit.CalculateWitness(privateInputs, publicInputs)
	panicErr(err)
​
	// R1CS to QAP
	alphas, betas, gammas, _ := snark.Utils.PF.R1CSToQAP(circuit.R1CS.A, circuit.R1CS.B, circuit.R1CS.C)
	fmt.Println("qap")
	fmt.Println(alphas)
	fmt.Println(betas)
	fmt.Println(gammas)
​
	// calculate trusted setup
	setup, err := snark.GenerateTrustedSetup(len(w), circuit, alphas, betas, gammas)
	panicErr(err)
	fmt.Println("\nt:", setup.Toxic.T)
​
	// remove setup.Toxic
	var tsetup snark.Setup
	tsetup.Pk = setup.Pk
	tsetup.Vk = setup.Vk
	tsetup.G1T = setup.G1T
	tsetup.G2T = setup.G2T
​
	// store setup to json
	jsonData, err := json.Marshal(tsetup)
	panicErr(err)
	// store setup into file
	jsonFile, err := os.Create("trustedsetup.json")
	panicErr(err)
	defer jsonFile.Close()
	jsonFile.Write(jsonData)
	jsonFile.Close()
	fmt.Println("Trusted Setup data written to ", jsonFile.Name())
​
	return nil
}
​
func TrustedSetupOnly() error {
	compiledcircuitFile, err := ioutil.ReadFile("compiledcircuit.json")
	panicErr(err)
	var circuit circuitcompiler.Circuit
	json.Unmarshal([]byte(string(compiledcircuitFile)), &circuit)
	panicErr(err)
​
	// flat code to R1CS
	//fmt.Println("generating R1CS from flat code")
	a, b, c := circuit.GenerateR1CS()
​
	/*
	now we have the R1CS from the circuit:
	a: [[0 0 1 0 0 0 0 0] [0 0 1 0 0 0 0 0] [0 0 1 0 1 0 0 0] [5 0 0 0 0 1 0 0] [0 0 0 0 0 0 1 0] [0 1 0 0 0 0 0 0] [1 0 0 0 0 0 0 0]]
	b: [[0 0 1 0 0 0 0 0] [0 0 0 1 0 0 0 0] [1 0 0 0 0 0 0 0] [1 0 0 0 0 0 0 0] [1 0 0 0 0 0 0 0] [1 0 0 0 0 0 0 0] [1 0 0 0 0 0 0 0]]
	c: [[0 0 0 1 0 0 0 0] [0 0 0 0 1 0 0 0] [0 0 0 0 0 1 0 0] [0 0 0 0 0 0 1 0] [0 1 0 0 0 0 0 0] [0 0 0 0 0 0 1 0] [0 0 0 0 0 0 0 1]]
	*/
	alphas, betas, gammas, _ := snark.Utils.PF.R1CSToQAP(a,b,c)
​
	setup, _ := snark.GenerateTrustedSetup(8, circuit, alphas, betas, gammas )
​
	// store setup to json
	jsonData, err := json.Marshal(setup)
	panicErr(err)
	// store setup into file
	jsonFile, err := os.Create("trustedsetup.json")
	panicErr(err)
	defer jsonFile.Close()
	jsonFile.Write(jsonData)
	jsonFile.Close()
	//fmt.Println("Trusted Setup data written to ", jsonFile.Name())
​
	return nil
}
// privateInputs : private Key .. publicInputs : public Key
func GenerateProofs(privateInputs []*big.Int, publicInputs []*big.Int) error {
	// open compiledcircuit.json
	compiledcircuitFile, err := ioutil.ReadFile("compiledcircuit.json")
	panicErr(err)
	var circuit circuitcompiler.Circuit
	json.Unmarshal([]byte(string(compiledcircuitFile)), &circuit)
	panicErr(err)
​
	// open trustedsetup.json
	trustedsetupFile, err := ioutil.ReadFile("trustedsetup.json")
	panicErr(err)
	var trustedsetup snark.Setup
	json.Unmarshal([]byte(string(trustedsetupFile)), &trustedsetup)
	panicErr(err)
​
	// calculate wittness
	w, err := circuit.CalculateWitness(privateInputs, publicInputs)
	panicErr(err)
	//fmt.Println("witness", w)
​
	// flat code to R1CS
	a,b,c := circuit.GenerateR1CS()
/*	a := circuit.R1CS.A
	b := circuit.R1CS.B
	c := circuit.R1CS.C */
​
	// R1CS to QAP
	alphas, betas, gammas, _ := snark.Utils.PF.R1CSToQAP(a, b, c)
	_, _, _, px := snark.Utils.PF.CombinePolynomials(w, alphas, betas, gammas)
	hx := snark.Utils.PF.DivisorPolynomial(px, trustedsetup.Pk.Z)
	fmt.Println(hx)
	/*
	fmt.Println(circuit)
	fmt.Println(trustedsetup.G1T)
	fmt.Println(hx)
	fmt.Println(w)
	*/
	proof, err := snark.GenerateProofs(circuit, trustedsetup, w, px)
	panicErr(err)
​
	fmt.Println("\n proofs:")
	fmt.Println(proof)
​
	// store proofs to json
	jsonData, err := json.Marshal(proof)
	panicErr(err)
	// store proof into file
	jsonFile, err := os.Create("proofs.json")
	panicErr(err)
	defer jsonFile.Close()
	jsonFile.Write(jsonData)
	jsonFile.Close()
	fmt.Println("Proofs data written to ", jsonFile.Name())
	return nil
}
​
func VerifyProofs(publicInputs []*big.Int) error {
	// open proofs.json
	proofsFile, err := ioutil.ReadFile("proofs.json")
	panicErr(err)
	var proof snark.Proof
	json.Unmarshal([]byte(string(proofsFile)), &proof)
	panicErr(err)
​
	// open compiledcircuit.json
	compiledcircuitFile, err := ioutil.ReadFile("compiledcircuit.json")
	panicErr(err)
	var circuit circuitcompiler.Circuit
	json.Unmarshal([]byte(string(compiledcircuitFile)), &circuit)
	panicErr(err)
​
	// open trustedsetup.json
	trustedsetupFile, err := ioutil.ReadFile("trustedsetup.json")
	panicErr(err)
	var trustedsetup snark.Setup
	json.Unmarshal([]byte(string(trustedsetupFile)), &trustedsetup)
	panicErr(err)
​
	verified := snark.VerifyProof(circuit, trustedsetup, proof, publicInputs, true)
	if !verified {
		fmt.Println("ERROR: proofs not verified")
	} else {
		fmt.Println("Proofs verified")
	}
	return nil
}
​
func Onecode_zkSNARK () {
	flatCode := `
func exp3(private a):
	b = a * a
	c = a * b
	return c
​
func main(private s0, public s1):
	s3 = exp3(s0)
	s4 = s3 + s0
	s5 = s4 + 5
	equals(s1, s5)
	out = 1 * 1
	`
	// parse the code
	parser := circuitcompiler.NewParser(strings.NewReader(flatCode))
	circuit, _ := parser.Parse()
	//assert.Nil(t, err)
	fmt.Println(circuit)
​
	// Public / Private Value
	b3 := big.NewInt(int64(3))
	privateInputs := []*big.Int{b3}
	b35 := big.NewInt (int64(35))
	publicSignals := []*big.Int{b35}
​
	// y = x^3 + x + 5
​
	// Witness
	w, _ := circuit.CalculateWitness(privateInputs, publicSignals)
	//assert.Nil(t,err)
	fmt.Println("witness", w)
​
	// now we have the witness :
	// w = [1 35 3 9 27 30 35 1]
​
	// flat code to R1CS
	fmt.Println("generating R1CS from flat code")
	a, b, c := circuit.GenerateR1CS()
	/*
	now we have the R1CS from the circuit:
	a: [[0 0 1 0 0 0 0 0] [0 0 1 0 0 0 0 0] [0 0 1 0 1 0 0 0] [5 0 0 0 0 1 0 0] [0 0 0 0 0 0 1 0] [0 1 0 0 0 0 0 0] [1 0 0 0 0 0 0 0]]
	b: [[0 0 1 0 0 0 0 0] [0 0 0 1 0 0 0 0] [1 0 0 0 0 0 0 0] [1 0 0 0 0 0 0 0] [1 0 0 0 0 0 0 0] [1 0 0 0 0 0 0 0] [1 0 0 0 0 0 0 0]]
	c: [[0 0 0 1 0 0 0 0] [0 0 0 0 1 0 0 0] [0 0 0 0 0 1 0 0] [0 0 0 0 0 0 1 0] [0 1 0 0 0 0 0 0] [0 0 0 0 0 0 1 0] [0 0 0 0 0 0 0 1]]
	*/
​
	alphas, betas, gammas, _ := snark.Utils.PF.R1CSToQAP(a,b,c)
​
	_, _, _, px := snark.Utils.PF.CombinePolynomials(w,alphas, betas, gammas)
​
	// calculate trusted setup
	// func GenerateTrustedSetup(witnessLength int,
	// 							circuit circuitcompiler.Circuit,
	// 							alphas, betas, gammas [][]*big.Int)
	setup, _ := snark.GenerateTrustedSetup(len(w), *circuit, alphas, betas, gammas )
​
	//_ := snark.Utils.PF.DivisorPolynomial(px, setup.Pk.Z)
​
	proof, _ := snark.GenerateProofs(*circuit, setup, w, px)
​
	b35Verif := big.NewInt(int64(35))
	publicSignalsVerif := []*big.Int{b35Verif}
​
	//assert.True()
	// func VerifyProof(circuit circuitcompiler.Circuit,
	// 					setup Setup, proof Proof,
	// 					publicSignals []*big.Int, debug bool) bool {
	verified := snark.VerifyProof(*circuit, setup, proof, publicSignalsVerif, true)
	if !verified {
		fmt.Println("Error in Verifying Proof")
	} else {
		fmt.Println("Verifying OK!!!")
	}
​
}
