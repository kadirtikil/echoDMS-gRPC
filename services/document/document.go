package document_service

import (
	"context"
	"strings"

	"github.com/echoDMS/proto/document"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DocumentService struct {
	document.UnimplementedDocumentServiceServer
	pool *pgxpool.Pool
}

func NewDocumentService(pool *pgxpool.Pool) *DocumentService {
	return &DocumentService{pool: pool}
}

var formatToDB = map[document.DocumentFormat]string{
	document.DocumentFormat_TYPST:    "typst",
	document.DocumentFormat_LATEX:    "latex",
	document.DocumentFormat_MARKDOWN: "markdown",
}

var compilerToDB = map[document.CompilerName]string{
	document.CompilerName_TYPST_COMPILER: "typst",
	document.CompilerName_PDFLATEX:       "pdflatex",
	document.CompilerName_XELATEX:        "xelatex",
	document.CompilerName_LUALATEX:       "lualatex",
	document.CompilerName_PANDOC:         "pandoc",
}

var outputFormatToDB = map[document.OutputFormat]string{
	document.OutputFormat_PDF:  "pdf",
	document.OutputFormat_PNG:  "png",
	document.OutputFormat_SVG:  "svg",
	document.OutputFormat_HTML: "html",
}

func (ds *DocumentService) PaginateDocuments(ctx context.Context, req *document.PaginateDocumentsRequest) (*document.PaginateDocumentsResponse, error) {
	offset := (req.PageNumber - 1) * req.PageSize
	sql := `SELECT id, title, description, format, owner_id, is_archived FROM documents WHERE owner_id = $1 LIMIT $2 OFFSET $3`
	rows, err := ds.pool.Query(ctx, sql, req.OwnerId, req.PageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documents []*document.Document
	for rows.Next() {
		var doc document.Document
		var formatStr string
		if err := rows.Scan(&doc.Id, &doc.Title, &doc.Description, &formatStr, &doc.OwnerId, &doc.IsArchived); err != nil {
			return nil, err
		}
		doc.Format = document.DocumentFormat(document.DocumentFormat_value[strings.ToUpper(formatStr)])
		documents = append(documents, &doc)
	}

	return &document.PaginateDocumentsResponse{
		Documents: documents,
	}, nil
}

func (ds *DocumentService) GetDocument(ctx context.Context, req *document.GetDocumentRequest) (*document.GetDocumentResponse, error) {
	sql := `SELECT id, title, description, format, owner_id, is_archived FROM documents WHERE id = $1`
	row := ds.pool.QueryRow(ctx, sql, req.Id)

	var doc document.Document
	var formatStr string
	if err := row.Scan(&doc.Id, &doc.Title, &doc.Description, &formatStr, &doc.OwnerId, &doc.IsArchived); err != nil {
		return nil, err
	}
	doc.Format = document.DocumentFormat(document.DocumentFormat_value[strings.ToUpper(formatStr)])

	return &document.GetDocumentResponse{
		Document: &doc,
	}, nil
}

func (ds *DocumentService) CreateDocument(ctx context.Context, req *document.CreateDocumentRequest) (*document.Document, error) {
	sql := `INSERT INTO documents (title, description, format, compiler, compiler_version, output_format, owner_id)
              VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	formatStr := formatToDB[req.Format]
	compilerStr := compilerToDB[req.Compiler]
	outputStr := outputFormatToDB[req.OutputFormat]

	var id string
	err := ds.pool.QueryRow(ctx, sql, req.Title, req.Description, formatStr, compilerStr, req.CompilerVersion, outputStr, req.OwnerId).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &document.Document{
		Id:          id,
		Title:       req.Title,
		Description: req.Description,
		Format:      req.Format,
		OwnerId:     req.OwnerId,
	}, nil
}
