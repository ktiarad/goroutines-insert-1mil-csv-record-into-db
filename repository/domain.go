package repository

import (
	"context"
	"database/sql"
	"fmt"
	"goinsertmil/model"
	"log"
)

type DomainRepository interface {
	InsertDomain(ctx context.Context, request model.Domain) (err error)
}

type domainRepository struct {
	db *sql.DB
}

func NewDomainRepository(db *sql.DB) *domainRepository {
	return &domainRepository{
		db: db,
	}
}

func (d *domainRepository) InsertDomain(ctx context.Context, request model.Domain) (err error) {
	data := []interface{}{
		request.GlobalRank,
		request.TldRank,
		request.Domain,
		request.TLD,
		request.RefSubNets,
		request.RefIPs,
		request.IDN_Domain,
		request.IDN_TLD,
		request.PrevGlobalRank,
		request.TldRank,
		request.PrevRefSubNets,
		request.PrevRefIPs,
	}

	// Create new connection since this function
	conn, err := d.db.Conn(ctx)
	if err != nil {
		log.Fatalf("error connect to db: %v", err.Error())
	}

	query := fmt.Sprintf("INSERT INTO domains (global_rank, tld_rank, domain, tld, ref_sub_nets, ref_ips, idn_domain, idn_tld, prev_global_rank, prev_tld_rank, prev_ref_sub_nets, prev_ref_ips) VALUES (%s)",
		GenerateDollarsMark(data),
	)

	_, err = conn.ExecContext(ctx, query, data...)
	if err != nil {
		log.Fatalf("error when execute: %v; data: %+v, len: %d", err.Error(), data, len(data))
	}

	err = conn.Close()
	if err != nil {
		log.Fatalf("error when close conn: %v", err.Error())
	}

	return err
}
