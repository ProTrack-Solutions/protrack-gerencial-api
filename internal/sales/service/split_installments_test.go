package service

import (
	"fmt"
	"math"
	"reflect"
	"testing"
)

func TestSplitInstallments(t *testing.T) {
	tests := []struct {
		amount float64
		count  int32
		want   []float64
	}{
		{100, 3, []float64{33.34, 33.33, 33.33}},
		{200, 3, []float64{66.67, 66.67, 66.66}},
		{90, 3, []float64{30, 30, 30}},
		{10, 4, []float64{2.5, 2.5, 2.5, 2.5}},
		{0.05, 3, []float64{0.02, 0.02, 0.01}},
		{150 - 49.9, 7, []float64{14.30, 14.30, 14.30, 14.30, 14.30, 14.30, 14.30}},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%.2f_em_%dx", tt.amount, tt.count), func(t *testing.T) {
			got := splitInstallments(tt.amount, tt.count)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("splitInstallments(%v, %d) = %v, want %v", tt.amount, tt.count, got, tt.want)
			}
		})
	}
}

// A soma das parcelas precisa bater, centavo a centavo, com o valor lançado no saldo do cliente.
func TestSplitInstallmentsSumMatchesAmount(t *testing.T) {
	for cents := int64(1); cents <= 100000; cents += 7 {
		amount := float64(cents) / 100
		for count := int32(1); count <= 24; count++ {
			var sum int64
			for _, v := range splitInstallments(amount, count) {
				sum += int64(math.Round(v * 100))
			}
			if sum != cents {
				t.Fatalf("R$ %.2f em %dx: soma das parcelas = %d centavos, esperado %d", amount, count, sum, cents)
			}
		}
	}
}
