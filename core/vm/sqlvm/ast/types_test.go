package ast

import (
	"testing"

	"github.com/dexon-foundation/dexon/common"
	dec "github.com/dexon-foundation/dexon/core/vm/sqlvm/common/decimal"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
)

type TypesTestSuite struct{ suite.Suite }

func (s *TypesTestSuite) requireEncodeAndDecodeDecimalNoError(
	d DataType, t decimal.Decimal, bs int) {
	encode, err := DecimalEncode(d, t)
	s.Require().NoError(err)
	s.Require().Len(encode, bs)
	decode, err := DecimalDecode(d, encode)
	s.Require().NoError(err)
	s.Require().Equal(t.String(), decode.String())
}

func (s *TypesTestSuite) requireEncodeAndDecodeNoError(
	d DataType, t TypeNode) {
	encode, err := DataTypeEncode(t)
	s.Require().NoError(err)
	s.Require().Equal(d, encode)
	decode, err := DataTypeDecode(d)
	s.Require().NoError(err)
	s.Require().Equal(t, decode)
}

func (s *TypesTestSuite) requireEncodeError(input TypeNode) {
	_, err := DataTypeEncode(input)
	s.Require().Error(err)
}

func (s *TypesTestSuite) requireDecodeError(input DataType) {
	_, err := DataTypeDecode(input)
	s.Require().Error(err)
}

func (s *TypesTestSuite) TestEncodeAndDecode() {
	s.requireEncodeAndDecodeNoError(
		ComposeDataType(DataTypeMajorBool, 0),
		&BoolTypeNode{})
	s.requireEncodeAndDecodeNoError(
		ComposeDataType(DataTypeMajorAddress, 0),
		&AddressTypeNode{})
	s.requireEncodeAndDecodeNoError(
		ComposeDataType(DataTypeMajorInt, 1),
		&IntTypeNode{Size: 16})
	s.requireEncodeAndDecodeNoError(
		ComposeDataType(DataTypeMajorUint, 2),
		&IntTypeNode{Unsigned: true, Size: 24})
	s.requireEncodeAndDecodeNoError(
		ComposeDataType(DataTypeMajorFixedBytes, 3),
		&FixedBytesTypeNode{Size: 4})
	s.requireEncodeAndDecodeNoError(
		ComposeDataType(DataTypeMajorDynamicBytes, 0),
		&DynamicBytesTypeNode{})
	s.requireEncodeAndDecodeNoError(
		ComposeDataType(DataTypeMajorFixed, 1),
		&FixedTypeNode{Size: 8, FractionalDigits: 1})
	s.requireEncodeAndDecodeNoError(
		ComposeDataType(DataTypeMajorUfixed+1, 2),
		&FixedTypeNode{Unsigned: true, Size: 16, FractionalDigits: 2})
}

func (s *TypesTestSuite) TestEncodeError() {
	s.requireEncodeError(nil)
	s.requireEncodeError(&IntTypeNode{Size: 1})
	s.requireEncodeError(&IntTypeNode{Size: 257})
	s.requireEncodeError(&FixedBytesTypeNode{Size: 0})
	s.requireEncodeError(&FixedBytesTypeNode{Size: 257})
	s.requireEncodeError(&FixedTypeNode{Size: 1, FractionalDigits: 0})
	s.requireEncodeError(&FixedTypeNode{Size: 257, FractionalDigits: 0})
	s.requireEncodeError(&FixedTypeNode{Size: 8, FractionalDigits: 81})
}

func (s *TypesTestSuite) TestDecodeError() {
	s.requireDecodeError(DataTypeUnknown)
	s.requireDecodeError(ComposeDataType(DataTypeMajorBool, 1))
	s.requireDecodeError(ComposeDataType(DataTypeMajorAddress, 1))
	s.requireDecodeError(ComposeDataType(DataTypeMajorInt, 0x20))
	s.requireDecodeError(ComposeDataType(DataTypeMajorUint, 0x20))
	s.requireDecodeError(ComposeDataType(DataTypeMajorFixedBytes, 0x20))
	s.requireDecodeError(ComposeDataType(DataTypeMajorDynamicBytes, 1))
	s.requireDecodeError(ComposeDataType(DataTypeMajorFixed, 81))
	s.requireDecodeError(ComposeDataType(DataTypeMajorUfixed, 81))
	s.requireDecodeError(ComposeDataType(DataTypeMajorUfixed+0x20, 80))
}

func (s *TypesTestSuite) TestEncodeAndDecodeDecimal() {
	pos := decimal.New(15, 0)
	zero := decimal.Zero
	neg := decimal.New(-15, 0)

	s.requireEncodeAndDecodeDecimalNoError(
		ComposeDataType(DataTypeMajorInt, 2),
		pos,
		3)
	s.requireEncodeAndDecodeDecimalNoError(
		ComposeDataType(DataTypeMajorInt, 2),
		zero,
		3)
	s.requireEncodeAndDecodeDecimalNoError(
		ComposeDataType(DataTypeMajorInt, 2),
		neg,
		3)

	s.requireEncodeAndDecodeDecimalNoError(
		ComposeDataType(DataTypeMajorUint, 2),
		pos,
		3)
	s.requireEncodeAndDecodeDecimalNoError(
		ComposeDataType(DataTypeMajorUint, 2),
		zero,
		3)

	s.requireEncodeAndDecodeDecimalNoError(
		ComposeDataType(DataTypeMajorAddress, 0),
		pos,
		20)
	s.requireEncodeAndDecodeDecimalNoError(
		ComposeDataType(DataTypeMajorAddress, 0),
		zero,
		20)

	pos = decimal.New(15, -2)
	neg = decimal.New(-15, -2)

	s.requireEncodeAndDecodeDecimalNoError(
		ComposeDataType(DataTypeMajorFixed+2, 2),
		pos,
		3)
	s.requireEncodeAndDecodeDecimalNoError(
		ComposeDataType(DataTypeMajorFixed+2, 2),
		zero,
		3)
	s.requireEncodeAndDecodeDecimalNoError(
		ComposeDataType(DataTypeMajorFixed+2, 2),
		neg,
		3)

	s.requireEncodeAndDecodeDecimalNoError(
		ComposeDataType(DataTypeMajorUfixed+2, 2),
		pos,
		3)
	s.requireEncodeAndDecodeDecimalNoError(
		ComposeDataType(DataTypeMajorUfixed+2, 2),
		zero,
		3)
}

func (s *TypesTestSuite) TestGetMinMax() {
	decAddressMax := decimal.New(2, 0).Pow(decimal.New(common.AddressLength*8, 0)).Sub(dec.One)
	testcases := []struct {
		Name     string
		In       DataType
		Min, Max decimal.Decimal
		Err      error
	}{
		{"Bool", ComposeDataType(DataTypeMajorBool, 0), dec.False, dec.True, nil},
		{"Address", ComposeDataType(DataTypeMajorAddress, 0), decimal.Zero, decAddressMax, nil},
		{"Int8", ComposeDataType(DataTypeMajorInt, 0), decimal.New(-128, 0), decimal.New(127, 0), nil},
		{"Int16", ComposeDataType(DataTypeMajorInt, 1), decimal.New(-32768, 0), decimal.New(32767, 0), nil},
		{"UInt8", ComposeDataType(DataTypeMajorUint, 0), decimal.Zero, decimal.New(255, 0), nil},
		{"UInt16", ComposeDataType(DataTypeMajorUint, 1), decimal.Zero, decimal.New(65535, 0), nil},
		{"Bytes1", ComposeDataType(DataTypeMajorFixedBytes, 0), decimal.Zero, decimal.New(255, 0), nil},
		{"Bytes2", ComposeDataType(DataTypeMajorFixedBytes, 1), decimal.Zero, decimal.New(65535, 0), nil},
		{"Dynamic Bytes", ComposeDataType(DataTypeMajorDynamicBytes, 0), decimal.Zero, decimal.Zero, ErrGetMinMax},
	}

	var (
		min, max decimal.Decimal
		err      error
	)
	for _, t := range testcases {
		min, max, err = GetMinMax(t.In)
		s.Require().Equal(t.Err, err, "Case: %v. Error not equal: %v != %v", t.Name, t.Err, err)
		if t.Err != nil {
			continue
		}

		s.Require().True(t.Min.Equal(min), "Case: %v. Min not equal: %v != %v", t.Name, t.Min, min)
		s.Require().True(t.Max.Equal(max), "Case: %v. Max not equal: %v != %v", t.Name, t.Max, max)
	}
}

func TestTypes(t *testing.T) {
	suite.Run(t, new(TypesTestSuite))
}