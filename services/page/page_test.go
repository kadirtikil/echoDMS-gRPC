package page_service_test

import (
	"context"
	"testing"

	"github.com/echoDMS/db"
	"github.com/echoDMS/proto/page"
	page_service "github.com/echoDMS/services/page"
	db_utils "github.com/echoDMS/utils/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pg *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()
	pool, err := db.NewPool(ctx, "postgresql://postgres:postgres@db:5432/echo_test?sslmode=disable")
	if err != nil {
		panic(err)
	}
	pg = pool

	db_utils.ReseedTestDB(context.Background(), pg)
	m.Run()

	pool.Close()
}

func TestPaginatePages(t *testing.T) {
	service := page_service.NewPageService(pg)

	tests := []struct {
		name       string
		page       int32
		pageSize   int32
		documentID string
		wantCount  int
		wantErr    bool
	}{
		{
			name:       "first page with results",
			page:       1,
			pageSize:   10,
			documentID: "a1b2c3d4-1111-1111-1111-000000000001",
			wantCount:  3,
		},
		{
			name:       "page size smaller than total",
			page:       1,
			pageSize:   2,
			documentID: "a1b2c3d4-1111-1111-1111-000000000001",
			wantCount:  2,
		},
		{
			name:       "second page with remaining",
			page:       2,
			pageSize:   2,
			documentID: "a1b2c3d4-1111-1111-1111-000000000001",
			wantCount:  1,
		},
		{
			name:       "page with no results",
			page:       3,
			pageSize:   2,
			documentID: "a1b2c3d4-1111-1111-1111-000000000001",
			wantCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.PaginatePages(context.Background(), &page.PaginatePagesRequest{
				DocumentId: tt.documentID,
				PageNumber: tt.page,
				PageSize:   tt.pageSize,
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("PaginatePages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(resp.Pages) != tt.wantCount {
				t.Errorf("PaginatePages() got %d pages, want %d", len(resp.Pages), tt.wantCount)
			}
		})
	}
}

func TestGetPage(t *testing.T) {
	service := page_service.NewPageService(pg)

	tests := []struct {
		name      string
		pageID    string
		wantDocID string
		wantErr   bool
	}{
		{
			name:      "existing page",
			pageID:    "b1b2c3d4-2222-2222-2222-000000000001",
			wantDocID: "a1b2c3d4-1111-1111-1111-000000000001",
			wantErr:   false,
		},
		{
			name:    "non-existent page",
			pageID:  "00000000-0000-0000-0000-000000000000",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.GetPage(context.Background(), &page.GetPageRequest{Id: tt.pageID})
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && resp.DocumentId != tt.wantDocID {
				t.Errorf("GetPage() got document ID %s, want %s", resp.DocumentId, tt.wantDocID)
			}
		})
	}
}

func TestCreatePage(t *testing.T) {
	service := page_service.NewPageService(pg)

	tests := []struct {
		name       string
		documentID string
		content    string
		pageNumber int32
		wantErr    bool
	}{
		{
			name:       "valid page creation",
			documentID: "s1d2f3g4-1111-1111-1111-000000000001",
			content:    "This is a new page.",
			pageNumber: 3,
			wantErr:    false,
		},
		{
			name:       "invalid document ID",
			documentID: "00000000-0000-0000-0000-000000000000",
			content:    "This page should fail.",
			pageNumber: 1,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.CreatePage(context.Background(), &page.CreatePageRequest{
				DocumentId: tt.documentID,
				Content:    tt.content,
				PageNumber: tt.pageNumber,
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("CreatePage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if resp.DocumentId != tt.documentID {
					t.Errorf("CreatePage() got document ID %s, want %s", resp.DocumentId, tt.documentID)
				}
				if resp.Content != tt.content {
					t.Errorf("CreatePage() got content %s, want %s", resp.Content, tt.content)
				}
				if resp.PageNumber != tt.pageNumber {
					t.Errorf("CreatePage() got page number %d, want %d", resp.PageNumber, tt.pageNumber)
				}
			}
		})
	}
}

func TestDeletePage(t *testing.T) {
	service := page_service.NewPageService(pg)

	tests := []struct {
		name    string
		pageID  string
		wantErr bool
	}{
		{
			name:    "existing page deletion",
			pageID:  "b1b2c3d4-2222-2222-2222-000000000001",
			wantErr: false,
		},
		{
			name:    "non-existent page deletion",
			pageID:  "00000000-0000-0000-0000-000000000000",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.DeletePage(context.Background(), &page.DeletePageRequest{Id: tt.pageID})
			if (err != nil) != tt.wantErr {
				t.Errorf("DeletePage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}

}
