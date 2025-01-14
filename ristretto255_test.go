package ristretto255

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"testing"

	"filippo.io/edwards25519/field"
)

// The encoding of the canonical generator.
var compressedRistrettoBasepoint, _ = hex.DecodeString("e2f2ae0a6abc4e71a884a961c500515f58e30b6aa582dd8db6a65945e08d2d76")

func TestRistrettoBasepointRoundTrip(t *testing.T) {
	decodedBasepoint := &Element{}
	err := decodedBasepoint.Decode(compressedRistrettoBasepoint)
	if err != nil {
		t.Fatal(err)
	}

	ristrettoBasepoint := (&Element{}).Base()
	if decodedBasepoint.Equal(ristrettoBasepoint) != 1 {
		t.Error("decode succeeded, but got wrong point")
	}

	roundtripBasepoint := decodedBasepoint.Encode(nil)
	if !bytes.Equal(compressedRistrettoBasepoint, roundtripBasepoint) {
		t.Error("decode<>encode roundtrip produced different results")
	}

	encodedBasepoint := ristrettoBasepoint.Encode(nil)
	if !bytes.Equal(compressedRistrettoBasepoint, encodedBasepoint) {
		t.Error("point encode produced different results")
	}
}

func TestRistrettoRandomRoundtrip(t *testing.T) {
	// TODO quickcheck
}

func TestRistrettoSmallMultiplesTestVectors(t *testing.T) {
	var testVectors = [16]string{
		// This is the identity point
		"0000000000000000000000000000000000000000000000000000000000000000",
		// This is the basepoint
		"e2f2ae0a6abc4e71a884a961c500515f58e30b6aa582dd8db6a65945e08d2d76",
		// These are small multiples of the basepoint
		"6a493210f7499cd17fecb510ae0cea23a110e8d5b901f8acadd3095c73a3b919",
		"94741f5d5d52755ece4f23f044ee27d5d1ea1e2bd196b462166b16152a9d0259",
		"da80862773358b466ffadfe0b3293ab3d9fd53c5ea6c955358f568322daf6a57",
		"e882b131016b52c1d3337080187cf768423efccbb517bb495ab812c4160ff44e",
		"f64746d3c92b13050ed8d80236a7f0007c3b3f962f5ba793d19a601ebb1df403",
		"44f53520926ec81fbd5a387845beb7df85a96a24ece18738bdcfa6a7822a176d",
		"903293d8f2287ebe10e2374dc1a53e0bc887e592699f02d077d5263cdd55601c",
		"02622ace8f7303a31cafc63f8fc48fdc16e1c8c8d234b2f0d6685282a9076031",
		"20706fd788b2720a1ed2a5dad4952b01f413bcf0e7564de8cdc816689e2db95f",
		"bce83f8ba5dd2fa572864c24ba1810f9522bc6004afe95877ac73241cafdab42",
		"e4549ee16b9aa03099ca208c67adafcafa4c3f3e4e5303de6026e3ca8ff84460",
		"aa52e000df2e16f55fb1032fc33bc42742dad6bd5a8fc0be0167436c5948501f",
		"46376b80f409b29dc2b5f6f0c52591990896e5716f41477cd30085ab7f10301e",
		"e0c418f7c8d9c4cdd7395b93ea124f3ad99021bb681dfc3302a9d99a2e53e64e",
	}

	basepointMultiple := (&Element{}).Zero()
	ristrettoBasepoint := (&Element{}).Base()

	for i := range testVectors {
		// Grab the bytes of the encoding
		encoding, err := hex.DecodeString(testVectors[i])
		if err != nil {
			t.Fatalf("#%d: bad hex encoding in test vector: %v", i, err)
		}

		// Decode the test vector to a ristretto255 element
		decodedPoint := Element{}
		err = decodedPoint.Decode(encoding)
		if err != nil {
			t.Fatalf("#%d: could not decode test vector: %v", i, err)
		}
		// Re-encode and check round trips
		roundtripEncoding := decodedPoint.Encode(nil)
		if !bytes.Equal(encoding, roundtripEncoding) {
			t.Errorf("#%d: decode<>encode roundtrip failed", i)
		}

		// Check that the test vector encodes i * B
		if basepointMultiple.Equal(&decodedPoint) != 1 {
			t.Errorf("decoded small multiple %d * B is not %d * B", i, i)
		}
		computedEncoding := basepointMultiple.Encode(nil)
		if !bytes.Equal(encoding, computedEncoding) {
			t.Errorf("#%d: encoding computed value did not match", i)
		}

		// Ensure basepointMultiple = i * B in the next iteration
		basepointMultiple.Add(basepointMultiple, ristrettoBasepoint)
	}
}

func TestRistrettoBadEncodingsTestVectors(t *testing.T) {
	var testVectors = []string{
		// These are all bad because they're non-canonical field encodings.
		"00ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff7f",
		"f3ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff7f",
		"edffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff7f",
		// These are all bad because they're negative field elements.
		"0100000000000000000000000000000000000000000000000000000000000000",
		"01ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff7f",
		"ed57ffd8c914fb201471d1c3d245ce3c746fcbe63a3679d51b6a516ebebe0e20",
		"c34c4e1826e5d403b78e246e88aa051c36ccf0aafebffe137d148a2bf9104562",
		"c940e5a4404157cfb1628b108db051a8d439e1a421394ec4ebccb9ec92a8ac78",
		"47cfc5497c53dc8e61c91d17fd626ffb1c49e2bca94eed052281b510b1117a24",
		"f1c6165d33367351b0da8f6e4511010c68174a03b6581212c71c0e1d026c3c72",
		"87260f7a2f12495118360f02c26a470f450dadf34a413d21042b43b9d93e1309",
		// These are all bad because they give a nonsquare x^2.
		"26948d35ca62e643e26a83177332e6b6afeb9d08e4268b650f1f5bbd8d81d371",
		"4eac077a713c57b4f4397629a4145982c661f48044dd3f96427d40b147d9742f",
		"de6a7b00deadc788eb6b6c8d20c0ae96c2f2019078fa604fee5b87d6e989ad7b",
		"bcab477be20861e01e4a0e295284146a510150d9817763caf1a6f4b422d67042",
		"2a292df7e32cababbd9de088d1d1abec9fc0440f637ed2fba145094dc14bea08",
		"f4a9e534fc0d216c44b218fa0c42d99635a0127ee2e53c712f70609649fdff22",
		"8268436f8c4126196cf64b3c7ddbda90746a378625f9813dd9b8457077256731",
		"2810e5cbc2cc4d4eece54f61c6f69758e289aa7ab440b3cbeaa21995c2f4232b",
		// These are all bad because they give a negative xy value.
		"3eb858e78f5a7254d8c9731174a94f76755fd3941c0ac93735c07ba14579630e",
		"a45fdc55c76448c049a1ab33f17023edfb2be3581e9c7aade8a6125215e04220",
		"d483fe813c6ba647ebbfd3ec41adca1c6130c2beeee9d9bf065c8d151c5f396e",
		"8a2e1d30050198c65a54483123960ccc38aef6848e1ec8f5f780e8523769ba32",
		"32888462f8b486c68ad7dd9610be5192bbeaf3b443951ac1a8118419d9fa097b",
		"227142501b9d4355ccba290404bde41575b037693cef1f438c47f8fbf35d1165",
		"5c37cc491da847cfeb9281d407efc41e15144c876e0170b499a96a22ed31e01e",
		"445425117cb8c90edcbc7c1cc0e74f747f2c1efa5630a967c64f287792a48a4b",
		// This is s = -1, which causes y = 0.
		"ecffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff7f",
	}

	basepointMultiple := Element{}
	basepointMultiple.Zero()

	for i := range testVectors {
		// Grab the bytes of the encoding
		encoding, err := hex.DecodeString(testVectors[i])
		if err != nil {
			t.Fatalf("#%d: bad hex encoding in test vector: %v", i, err)
		}

		// Attempt decoding
		decodedPoint := Element{}
		err = decodedPoint.Decode(encoding)
		if err == nil {
			t.Fatalf("#%d: did not fail on bad encoding", i)
		}
	}
}

func TestRistrettoFromUniformBytesTestVectors(t *testing.T) {
	inputs := []string{
		"Ristretto is traditionally a short shot of espresso coffee",
		"made with the normal amount of ground coffee but extracted with",
		"about half the amount of water in the same amount of time",
		"by using a finer grind.",
		"This produces a concentrated shot of coffee per volume.",
		"Just pulling a normal shot short will produce a weaker shot",
		"and is not a Ristretto as some believe.",
	}
	elements := []string{
		"3066f82a1a747d45120d1740f14358531a8f04bbffe6a819f86dfe50f44a0a46",
		"f26e5b6f7d362d2d2a94c5d0e7602cb4773c95a2e5c31a64f133189fa76ed61b",
		"006ccd2a9e6867e6a2c5cea83d3302cc9de128dd2a9a57dd8ee7b9d7ffe02826",
		"f8f0c87cf237953c5890aec3998169005dae3eca1fbb04548c635953c817f92a",
		"ae81e7dedf20a497e10c304a765c1767a42d6e06029758d2d7e8ef7cc4c41179",
		"e2705652ff9f5e44d3e841bf1c251cf7dddb77d140870d1ab2ed64f1a9ce8628",
		"80bd07262511cdde4863f8a7434cef696750681cb9510eea557088f76d9e5065",
	}

	var element Element
	for i, input := range inputs {
		hash := sha512.Sum512([]byte(input))
		element.FromUniformBytes(hash[:])
		if encoding := hex.EncodeToString(element.Encode(nil)); encoding != elements[i] {
			t.Errorf("#%d: expected %q, got %q", i, elements[i], encoding)
		}
	}
}

func TestEquivalentFromUniformBytes(t *testing.T) {
	inputs := []string{
		"edffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
			"1200000000000000000000000000000000000000000000000000000000000000",
		"edffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff7f" +
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		"0000000000000000000000000000000000000000000000000000000000000080" +
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff7f",
		"0000000000000000000000000000000000000000000000000000000000000000" +
			"1200000000000000000000000000000000000000000000000000000000000080",
	}
	expected := "304282791023b73128d277bdcb5c7746ef2eac08dde9f2983379cb8e5ef0517f"

	var element Element
	for i, input := range inputs {
		h, err := hex.DecodeString(input)
		if err != nil {
			t.Fatalf("#%d: bad hex encoding in test vector: %v", i, err)
		}

		element.FromUniformBytes(h[:])
		if encoding := hex.EncodeToString(element.Encode(nil)); encoding != expected {
			t.Errorf("#%d: expected %q, got %q", i, expected, encoding)
		}
	}
}

func TestMarshalScalar(t *testing.T) {
	x := new(Scalar)
	// generate an arbitrary scalar
	xbytes := sha512.Sum512([]byte("Hello World"))
	x.SetUniformBytes(xbytes[:])
	text, err := json.Marshal(x)
	if err != nil {
		t.Fatalf("Could not marshal json: %v", err)
	}
	t.Logf("json: %s", text)
	y := new(Scalar)
	err = json.Unmarshal(text, y)
	if err != nil || y.Equal(x) == 0 {
		t.Fatalf("Error unmarshaling scalar from json: %s %v", text, err)
	}
}

func TestMarshalElement(t *testing.T) {
	x := new(Element)
	// generate an arbitrary element
	xbytes := sha512.Sum512([]byte("Hello World"))
	x.FromUniformBytes(xbytes[:])
	text, err := json.Marshal(x)
	if err != nil {
		t.Fatalf("Could not marshal json: %v", err)
	}
	t.Logf("json: %s", text)
	y := new(Element)
	err = json.Unmarshal(text, y)
	if err != nil || y.Equal(x) == 0 {
		t.Fatalf("Error unmarshaling element from json: %s %v", text, err)
	}
}

func TestElementSet(t *testing.T) {
	// Test this, because the internal point type being hard-copyable isn't part of the spec.

	el1 := NewIdentityElement()
	el2 := NewGeneratorElement()

	if el1.Equal(el2) == 1 {
		t.Error("shouldn't be the same")
	}

	// Check new value
	el1.Set(el2)

	if el1.Equal(el2) == 0 {
		t.Error("failed to set the value")
	}

	// Mutate source var
	el2.Add(el2, el2)

	if el1.Equal(el2) == 1 {
		t.Error("shouldn't have changed")
	}
}

func TestScalarSet(t *testing.T) {
	// Test this, because the internal scalar representation being hard-copyable isn't part of the spec.

	// 32-byte little endian value of "1"
	scOne := make([]byte, 32)
	scOne[0] = 0x01

	sc1, sc2 := NewScalar(), NewScalar()

	// sc1 <- 1
	sc1.SetCanonicalBytes(scOne)

	// 1 != 0
	if sc1.Equal(sc2) == 1 {
		t.Error("shouldn't be the same")
	}

	// sc2 <- sc1
	sc2.Set(sc1)

	// 1 == 1
	if sc1.Equal(sc2) == 0 {
		t.Error("failed to set the value")
	}

	// sc1 <- 1 + 1
	sc1.Add(sc1, sc1)

	// 2 != 1
	if sc1.Equal(sc2) == 1 {
		t.Error("shouldn't have changed")
	}
}

func TestConstants(t *testing.T) {
	// From draft-hdevalence-cfrg-ristretto-01, Section 3.1.
	t.Run("d", func(t *testing.T) {
		testConstant(t, d,
			"37095705934669439343138083508754565189542113879843219016388785533085940283555")
	})
	t.Run("sqrtM1", func(t *testing.T) {
		testConstant(t, sqrtM1,
			"19681161376707505956807079304988542015446066515923890162744021073123829784752")
	})
	t.Run("sqrtADMinusOne", func(t *testing.T) {
		testConstant(t, sqrtADMinusOne,
			"25063068953384623474111414158702152701244531502492656460079210482610430750235")
	})
	t.Run("invSqrtAMinusD", func(t *testing.T) {
		testConstant(t, invSqrtAMinusD,
			"54469307008909316920995813868745141605393597292927456921205312896311721017578")
	})
	t.Run("oneMinusDSQ", func(t *testing.T) {
		testConstant(t, oneMinusDSQ,
			"1159843021668779879193775521855586647937357759715417654439879720876111806838")
	})
	t.Run("dMinusOneSQ", func(t *testing.T) {
		testConstant(t, dMinusOneSQ,
			"40440834346308536858101042469323190826248399146238708352240133220865137265952")
	})
}

func testConstant(t *testing.T, f *field.Element, decimal string) {
	b, ok := new(big.Int).SetString(decimal, 10)
	if !ok {
		t.Fatal("invalid decimal")
	}
	buf := b.FillBytes(make([]byte, 32))
	for i := 0; i < len(buf)/2; i++ {
		buf[i], buf[len(buf)-i-1] = buf[len(buf)-i-1], buf[i]
	}
	if !bytes.Equal(f.Bytes(), buf) {
		t.Errorf("expected %x", buf)
	}
}

func TestAffine(t *testing.T) {
	ristrettoBasepoint := (&Element{}).Base()
	x, y := ristrettoBasepoint.ToAffineCoordinate()
	result, err := ToExtendedProjectveCoordinate(x, y)
	if err != nil {
		t.Errorf(err.Error())
	}
	if ristrettoBasepoint.Equal(result) != 1 {
		t.Fatal("invalid point")
	}
}
