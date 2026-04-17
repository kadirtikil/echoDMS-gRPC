package page_service

import (
	"context"

	"github.com/echoDMS/proto/page"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PageService struct {
	page.UnimplementedPageServiceServer
	pool *pgxpool.Pool
}

func NewPageService(pool *pgxpool.Pool) *PageService {
	return &PageService{pool: pool}
}

func (ps *PageService) PaginatePages(ctx context.Context, req *page.PaginatePagesRequest) (*page.PaginatePagesResponse, error) {
	sql := `SELECT id, document_id, content, page_number FROM pages WHERE document_id = $1 LIMIT $2 OFFSET $3`
	offset := (req.PageNumber - 1) * req.PageSize
	rows, err := ps.pool.Query(ctx, sql, req.DocumentId, req.PageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pages []*page.Page
	for rows.Next() {
		var p page.Page
		if err := rows.Scan(&p.Id, &p.DocumentId, &p.Content, &p.PageNumber); err != nil {
			return nil, err
		}
		pages = append(pages, &p)
	}

	return &page.PaginatePagesResponse{Pages: pages}, nil
}

func (ps *PageService) GetPage(ctx context.Context, req *page.GetPageRequest) (*page.Page, error) {
	sql := `SELECT id, document_id, content, page_number FROM pages WHERE id = $1`
	row := ps.pool.QueryRow(ctx, sql, req.Id)

	var p page.Page
	if err := row.Scan(&p.Id, &p.DocumentId, &p.Content, &p.PageNumber); err != nil {
		return nil, err
	}
	return &p, nil
}

func (ps *PageService) CreatePage(ctx context.Context, req *page.CreatePageRequest) (*page.Page, error) {
	sql := `INSERT INTO pages (document_id, content, page_number) VALUES ($1, $2, $3) RETURNING id`
	var id string
	err := ps.pool.QueryRow(ctx, sql, req.DocumentId, req.Content, req.PageNumber).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &page.Page{Id: id, DocumentId: req.DocumentId, Content: req.Content, PageNumber: req.PageNumber}, nil
}

func (ps *PageService) DeletePage(ctx context.Context, req *page.DeletePageRequest) (*page.DeletePageResponse, error) {
	sql := `DELETE FROM pages WHERE id = $1`
	result, err := ps.pool.Exec(ctx, sql, req.Id)
	if err != nil {
		return nil, err
	}
	return &page.DeletePageResponse{Id: req.Id, Success: result.RowsAffected() > 0}, nil
}
