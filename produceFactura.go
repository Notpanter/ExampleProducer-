package service

import (
	"app/internal/config"
	"app/internal/database"
	"app/internal/domain/producer"
	"app/internal/routers/common"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
)

func ProduceFactura() error {
	log.Println("Start produceFACTURA")
	dbLists, err := database.GetDbList()
	if err != nil {
		return fmt.Errorf("ping Database error: %w", err)
	}
	ctx := context.Background()
	var aList producer.TfFactura
	aDictionary := []producer.TfFactura{}
	for i, rec := range dbLists {
		fmt.Println("DB Name", i, rec.DbName, rec.Ip) //
		connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
			"10.1.233.5", "ForTest", "Aa123654789",
			1433, "ALA50_MarketSQL")

		//connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;connection timeout=%d;",
		//	rec.Ip, config.Config.DBConfig.User, config.Config.DBConfig.Password,
		//	config.Config.DBConfig.Port, rec.DbName, 5)
		//fmt.Printf("connection: ", connString)
		sourceDb, err := sql.Open("sqlserver", connString)
		if err != nil {
			log.Println("open SQL Connection: ", err, " DbName: ", rec.DbName)
			continue
		} // Check if database is alive.
		err = sourceDb.PingContext(ctx)
		if err != nil {
			log.Println("sourceDb.PingContext error: ", err, " DbName: ", rec.DbName)
			continue
		}
		sqlStatement := `  SELECT     100000 AS [USER], WH, FACTURA, CASH_N, CHECK_N, [DATE], CUSTOMER, RNN, IIN_BIN, LEGAL, FIRM_ID, FACTURA_YEAR, ACCOUNT, BANK_NAME, BIK, ADDRESS, 
                     DOVERENNOST_NOMER, DOVERENNOST_DATA, DOVERENNOST_KOMU, N_DOG, DATA_DOG, BASE_CASH_N, BASE_CHECK_N, SAVE_DATE, CONTACTS, 
                      ISNULL(APPROVED, 0) AS APPROVED
                      FROM         tf_FACTURA
                       WHERE     TRANSFERRED_TO_KAFKA=0`
		rows, err := sourceDb.QueryContext(ctx, sqlStatement)
		if err != nil {
			log.Println("sqlStatementGetDB sourceDb.QueryContext: ", err, "DbName: ", rec.DbName)
			continue
		}
		if err != nil {
			return fmt.Errorf("get information Factura error: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			//var IdWarehouseGroup *int
			err := rows.Scan(
				&aList.USER,
				&aList.WH,
				&aList.FACTURA,
				&aList.CashN,
				&aList.CheckN,
				&aList.Date,
				&aList.Customer,
				&aList.Rnn,
				&aList.IinBin,
				&aList.Legal,
				&aList.FirmId,
				&aList.FacturaYear,
				&aList.Account,
				&aList.BankName,
				&aList.Bik,
				&aList.Address,
				&aList.DoverennostNomer,
				&aList.DoverennostData,
				&aList.DoverennostKomu,
				&aList.NDog,
				&aList.DataDog,
				&aList.BaseCashN,
				&aList.BaseCheckN,
				&aList.SaveDate,
				&aList.Contacts,
				&aList.Approved,
			)
			if err != nil {
				return fmt.Errorf("loop information about Factura: %w", err)
			}
			aDictionary = append(aDictionary, aList)
		}

		for _, perem := range aDictionary {
			sqlStatement := `UPDATE tf_FACTURA SET TRANSFERRED_TO_KAFKA = 1 WHERE FACTURA = @Factura AND WH = @WH AND CASH_N=@CASH AND CHECK_N=@CHK`
			rows, err = sourceDb.QueryContext(ctx, sqlStatement,
				sql.Named("Factura", perem.FACTURA),
				sql.Named("WH", perem.WH),
				sql.Named("CASH", perem.CashN),
				sql.Named("CHK", perem.CheckN),
			)
			//msg := producer.TfFacturaLoad{Count: len(aDictionary), Items: aDictionary}

		}
		//TODO: поставить ограничение на количество записей в сообщении

		msg := producer.TfFacturaLoad{Count: len(aDictionary), Items: aDictionary}

		kafkaMsg, err := json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("marshal data for topic queue MSG, IdWarehouse: %v, %w", err)
		}

		if err := common.PushCommentToQueue(config.Config.Kafka.TopicFactura,
			config.Config.Kafka.TopicFacturaMaxSize,
			kafkaMsg); err != nil {
			return fmt.Errorf("push Comment to topic queue MSG, IdWarehouse: %v, %w", err)
		}

	}
	log.Println("Finish produceFactura")
	return nil
}
