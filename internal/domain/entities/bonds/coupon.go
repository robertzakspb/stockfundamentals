package bonds

import (
	"time"

	"github.com/google/uuid"
)

/*
Figi            string                 `protobuf:"bytes,1,opt,name=figi,proto3" json:"figi,omitempty"`                                                                                      //FIGI-идентификатор инструмента.
	CouponDate      *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=coupon_date,json=couponDate,proto3" json:"coupon_date,omitempty"`                                                        //Дата выплаты купона.
	CouponNumber    int64                  `protobuf:"varint,3,opt,name=coupon_number,json=couponNumber,proto3" json:"coupon_number,omitempty"`                                                 //Номер купона.
	FixDate         *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=fix_date,json=fixDate,proto3" json:"fix_date,omitempty"`                                                                 //Дата фиксации реестра для выплаты купона — опционально.
	PayOneBond      *MoneyValue            `protobuf:"bytes,5,opt,name=pay_one_bond,json=payOneBond,proto3" json:"pay_one_bond,omitempty"`                                                      //Выплата на одну облигацию.
	CouponType      CouponType             `protobuf:"varint,6,opt,name=coupon_type,json=couponType,proto3,enum=tinkoff.public.invest.api.contract.v1.CouponType" json:"coupon_type,omitempty"` //Тип купона.
	CouponStartDate *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=coupon_start_date,json=couponStartDate,proto3" json:"coupon_start_date,omitempty"`                                       //Начало купонного периода.
	CouponEndDate   *timestamppb.Timestamp `protobuf:"bytes,8,opt,name=coupon_end_date,json=couponEndDate,proto3" json:"coupon_end_date,omitempty"`                                             //Окончание купонного периода.
	CouponPeriod    int32
*/

type Coupon struct {
	Id         uuid.UUID
	Figi       string
	CouponDate time.Time
}
