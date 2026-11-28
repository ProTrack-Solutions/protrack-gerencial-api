package globaldomain

import (
	"encoding/json"
	"testing"
)

type paginatedItem struct {
	ID string `json:"id"`
}

func TestNewPaginatedResponse_NilDataBecomesEmptyArray(t *testing.T) {
	var data []paginatedItem // busca sem resultado: o service não faz nenhum append

	resp := NewPaginatedResponse(data, 0, PaginationParams{Page: 1, PerPage: 20})

	if resp.Data == nil {
		t.Fatal("Data não deveria ser nil")
	}

	body, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("erro ao serializar: %v", err)
	}

	var decoded map[string]json.RawMessage
	if err := json.Unmarshal(body, &decoded); err != nil {
		t.Fatalf("erro ao ler o JSON: %v", err)
	}
	if got := string(decoded["data"]); got != "[]" {
		t.Fatalf(`esperava "data": [], veio "data": %s`, got)
	}
}

func TestNewPaginatedResponse_KeepsDataAndTotals(t *testing.T) {
	data := []paginatedItem{{ID: "a"}, {ID: "b"}}

	resp := NewPaginatedResponse(data, 45, PaginationParams{Page: 2, PerPage: 20})

	if len(resp.Data) != 2 || resp.Data[0].ID != "a" {
		t.Fatalf("dados alterados: %+v", resp.Data)
	}
	if resp.TotalRows != 45 || resp.TotalPages != 3 || resp.Page != 2 || resp.PerPage != 20 {
		t.Fatalf("paginação incorreta: %+v", resp)
	}
}
