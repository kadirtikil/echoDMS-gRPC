package document_service_test

import (
	"context"
	"testing"

	"github.com/echoDMS/db"
	"github.com/echoDMS/proto/document"
	document_service "github.com/echoDMS/services/document"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

var pg *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()
	pool, err := db.NewPool(ctx, "postgres://postgres:postgres@localhost:5432/echo_test")
	if err != nil {
		panic(err)
	}
	pg = pool

	code := m.Run()

	pool.Close()
	// os.Exit(code) — omit so cleanup runs; test binary exits with code automatically
	_ = code
}

func TestPaginateDocuments(t *testing.T) {
	service := document_service.NewDocumentService(pg)

	tests := []struct {
		name      string
		page      int32
		pageSize  int32
		ownerID   string
		wantCount int
		wantErr   bool
	}{
		{
			name:      "first page with results",
			page:      1,
			pageSize:  10,
			ownerID:   "a1b2c3d4-1111-1111-1111-000000000002",
			wantCount: 3,
		},
		{
			name:      "page size smaller than total",
			page:      1,
			pageSize:  2,
			ownerID:   "a1b2c3d4-1111-1111-1111-000000000002",
			wantCount: 2,
		},
		{
			name:      "second page with remaining",
			page:      2,
			pageSize:  2,
			ownerID:   "a1b2c3d4-1111-1111-1111-000000000002",
			wantCount: 1,
		},
		{
			name:      "page beyond results",
			page:      99,
			pageSize:  10,
			ownerID:   "a1b2c3d4-1111-1111-1111-000000000002",
			wantCount: 0,
		},
		{
			name:      "owner with no documents",
			page:      1,
			pageSize:  10,
			ownerID:   "00000000-0000-0000-0000-000000000000",
			wantCount: 0,
		},
		{
			name:      "empty owner id",
			page:      1,
			pageSize:  10,
			ownerID:   "",
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.PaginateDocuments(context.Background(), &document.PaginateDocumentsRequest{
				PageNumber: tt.page,
				PageSize:   tt.pageSize,
				OwnerId:    tt.ownerID,
			})

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Len(t, resp.Documents, tt.wantCount)

			for _, doc := range resp.Documents {
				assert.Equal(t, tt.ownerID, doc.OwnerId)
			}
		})
	}
}

func TestGetDocument(t *testing.T) {
	service := document_service.NewDocumentService(pg)

	tests := []struct {
		name      string
		id        string
		wantTitle string
		wantErr   bool
	}{
		{
			name:      "existing document",
			id:        "a1b2c3d4-1111-1111-1111-000000000001",
			wantTitle: "Thesis Draft",
			wantErr:   false,
		},
		{
			name:    "non-existent document",
			id:      "00000000-0000-0000-0000-000000000000",
			wantErr: true,
		},
		{
			name:    "invalid uuid",
			id:      "not-a-uuid",
			wantErr: true,
		},
		{
			name:    "empty id",
			id:      "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.GetDocument(context.Background(), &document.GetDocumentRequest{
				Id: tt.id,
			})

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp.Document)
			assert.Equal(t, tt.id, resp.Document.Id)
			assert.Equal(t, tt.wantTitle, resp.Document.Title)
		})
	}
}

func TestCreateDocument(t *testing.T) {
	service := document_service.NewDocumentService(pg)

	t.Cleanup(func() {
		pg.Exec(context.Background(), "DELETE FROM documents WHERE id NOT IN ('a1b2c3d4-1111-1111-1111-000000000001', 'a1b2c3d4-1111-1111-1111-000000000002', 'a1b2c3d4-1111-1111-1111-000000000003')")
	})

	tests := []struct {
		name            string
		title           string
		description     string
		format          document.DocumentFormat
		compiler        document.CompilerName
		compilerVersion string
		outputFormat    document.OutputFormat
		ownerID         string
		wantErr         bool
	}{
		{
			name:            "valid document",
			title:           "New Document",
			description:     "A new document for testing",
			format:          document.DocumentFormat_TYPST,
			compiler:        document.CompilerName_TYPST_COMPILER,
			compilerVersion: "0.12",
			outputFormat:    document.OutputFormat_PDF,
			ownerID:         "a1b2c3d4-1111-1111-1111-000000000002",
			wantErr:         false,
		},
		{
			name:            "invalid format enum",
			title:           "Bad Format",
			description:     "Document with invalid format",
			format:          document.DocumentFormat(999),
			compiler:        document.CompilerName_TYPST_COMPILER,
			compilerVersion: "0.12",
			outputFormat:    document.OutputFormat_PDF,
			ownerID:         "a1b2c3d4-1111-1111-1111-000000000002",
			wantErr:         true,
		},
		{
			name:            "invalid compiler enum",
			title:           "Bad Compiler",
			description:     "Document with invalid compiler",
			format:          document.DocumentFormat_TYPST,
			compiler:        document.CompilerName(999),
			compilerVersion: "0.12",
			outputFormat:    document.OutputFormat_PDF,
			ownerID:         "a1b2c3d4-1111-1111-1111-000000000002",
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.CreateDocument(context.Background(), &document.CreateDocumentRequest{
				Title:           tt.title,
				Description:     tt.description,
				Format:          tt.format,
				Compiler:        tt.compiler,
				CompilerVersion: tt.compilerVersion,
				OutputFormat:    tt.outputFormat,
				OwnerId:         tt.ownerID,
			})

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, tt.title, resp.Title)
			assert.NotEmpty(t, resp.Id)
		})
	}
}

func TestDeleteDocument(t *testing.T) {
	service := document_service.NewDocumentService(pg)

	created, err := service.CreateDocument(context.Background(), &document.CreateDocumentRequest{
		Title:           "To Be Deleted",
		Description:     "This will be deleted",
		Format:          document.DocumentFormat_TYPST,
		Compiler:        document.CompilerName_TYPST_COMPILER,
		CompilerVersion: "0.12",
		OutputFormat:    document.OutputFormat_PDF,
		OwnerId:         "a1b2c3d4-1111-1111-1111-000000000002",
	})
	if err != nil {
		t.Fatalf("failed to create test document: %v", err)
	}

	tests := []struct {
		name        string
		id          string
		wantDeleted bool
		wantErr     bool
	}{
		{
			name:        "delete existing document",
			id:          created.Id,
			wantDeleted: true,
		},
		{
			name:        "delete non-existent document",
			id:          "00000000-0000-0000-0000-000000000000",
			wantDeleted: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.DeleteDocument(context.Background(), &document.DeleteDocumentRequest{
				Id: tt.id,
			})

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantDeleted, resp.Success)
		})
	}
}

func TestUpdateDocument(t *testing.T) {
	service := document_service.NewDocumentService(pg)

	// Create a document to update
	created, err := service.CreateDocument(context.Background(), &document.CreateDocumentRequest{
		Title:           "Before Update",
		Description:     "Original description",
		Format:          document.DocumentFormat_TYPST,
		Compiler:        document.CompilerName_TYPST_COMPILER,
		CompilerVersion: "0.12",
		OutputFormat:    document.OutputFormat_PDF,
		OwnerId:         "a1b2c3d4-1111-1111-1111-000000000002",
	})
	if err != nil {
		t.Fatalf("failed to create test document: %v", err)
	}

	t.Cleanup(func() {
		pg.Exec(context.Background(), "DELETE FROM documents WHERE id = $1", created.Id)
	})

	tests := []struct {
		name            string
		id              string
		title           string
		description     string
		format          document.DocumentFormat
		compiler        document.CompilerName
		compilerVersion string
		outputFormat    document.OutputFormat
		wantErr         bool
	}{
		{
			name:            "valid update",
			id:              created.Id,
			title:           "After Update",
			description:     "Updated description",
			format:          document.DocumentFormat_LATEX,
			compiler:        document.CompilerName_PDFLATEX,
			compilerVersion: "1.0",
			outputFormat:    document.OutputFormat_PDF,
			wantErr:         false,
		},
		{
			name:            "non-existent document",
			id:              "00000000-0000-0000-0000-000000000000",
			title:           "Should Fail",
			description:     "This update should fail",
			format:          document.DocumentFormat_TYPST,
			compiler:        document.CompilerName_TYPST_COMPILER,
			compilerVersion: "0.12",
			outputFormat:    document.OutputFormat_PDF,
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.UpdateDocument(context.Background(), &document.UpdateDocumentRequest{
				Id:              tt.id,
				Title:           tt.title,
				Description:     tt.description,
				Format:          tt.format,
				Compiler:        tt.compiler,
				CompilerVersion: tt.compilerVersion,
				OutputFormat:    tt.outputFormat,
			})

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, tt.id, resp.Id)
			assert.Equal(t, tt.title, resp.Title)
		})
	}
}
