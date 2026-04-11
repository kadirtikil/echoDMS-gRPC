package services

import (
	"context"
	"fmt"

	pb "github.com/echoDMS/proto/echodms"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type PageService struct {
	pb.UnimplementedPageServiceServer
	db *pgxpool.Pool
}

func NewPageService(pool *pgxpool.Pool) *PageService {
	return &PageService{db: pool}
}

func (s *PageService) CreatePage(ctx context.Context, req *pb.CreatePageRequest) (*pb.Page, error) {
	sql := `INSERT INTO pages (content, document_id) VALUES ($1, $2) RETURNING id, content`
	row := s.db.QueryRow(ctx, sql, req.Content, req.DocumentId)

	var page pb.Page
	if err := row.Scan(&page.Id, &page.Content); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create page: %v", err)
	}

	return &page, nil
}

func (s *PageService) ListPages(ctx context.Context, req *pb.ListPagesRequest) (*pb.ListPagesResponse, error) {
	sql := `SELECT id, content FROM pages WHERE document_id = $1`
	rows, err := s.db.Query(ctx, sql, req.DocumentId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get pages: %v", err)
	}
	defer rows.Close()

	fmt.Println(rows)

	var pages []*pb.Page
	for rows.Next() {
		var page pb.Page
		if err := rows.Scan(&page.Id, &page.Content); err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to scan page: %v", err)
		}
		pages = append(pages, &page)
	}

	return &pb.ListPagesResponse{Pages: pages}, nil
}

func (s *PageService) GetPage(ctx context.Context, req *pb.GetPageRequest) (*pb.Page, error) {
	sql := `SELECT id, content FROM pages WHERE id = $1`
	row := s.db.QueryRow(ctx, sql, req.Id)

	var page pb.Page
	if err := row.Scan(&page.Id, &page.Content); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get page: %v", err)
	}

	return &page, nil
}

func (s *PageService) UpdatePage(ctx context.Context, req *pb.UpdatePageRequest) (*pb.Page, error) {
	sql := `UPDATE pages SET content = $1 WHERE id = $2 RETURNING id, content`
	row := s.db.QueryRow(ctx, sql, req.Content, req.Id)

	var page pb.Page
	if err := row.Scan(&page.Id, &page.Content); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update page: %v", err)
	}

	return &page, nil
}

func (s *PageService) DeletePage(ctx context.Context, req *pb.DeletePageRequest) (*emptypb.Empty, error) {
	sql := `DELETE FROM pages WHERE id = $1`
	_, err := s.db.Exec(ctx, sql, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to delete page: %v", err)
	}

	return &emptypb.Empty{}, nil
}
