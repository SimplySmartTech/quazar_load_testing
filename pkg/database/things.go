package database

import (
	"context"
	"fmt"
	"log"
	"simply_smart_mqtt_load_testing/pkg/model"
)

var (
	FetchTemplateKey = "SELECT template_key FROM templates"
	RemoveThingsData = "DELETE FROM %v WHERE 1=1 RETURNING *;"
)

func (db *databaseOps) FetchThingies(ctx context.Context, limit int32, sql string) (model.Thingies, error) {
	var returnThingies model.Thingies
	if sql == "" {
		sql = "SELECT id, name,template,thing_key,extract(epoch from last_activity) as last_activity,site_id,user_defined_property,templates_id FROM things t order by id"
	}
	rows, err := db.db.QueryContext(ctx, sql)
	if err != nil {
		log.Println(err)
		return returnThingies, err
	}

	// Fetch thingies row by row
	for rows.Next() {
		var t model.Thingie
		err = rows.Scan(&t.Id, &t.Name, &t.Template, &t.ThingKey, &t.LastActivity, &t.SiteId, &t.UserDefinedProperty, &t.TemplateId)
		if err != nil {
			log.Println(err)
		}
		returnThingies = append(returnThingies, t)
	}
	return returnThingies, err
}

func (db *databaseOps) TotalGateway(ctx context.Context, meterpergateway int32) (total int32) {
	row := db.db.QueryRowContext(ctx, "select COUNT(*) as total_g from things t")
	if row.Err() != nil {
		log.Println("things.go-TotalGateway", row.Err().Error())
		return 0
	}
	var total_g int32
	err := row.Scan(&total_g)
	if err != nil {
		log.Println(err)
		return 0
	}

	// Check for decimal point values and increase 1 if reminder is remaining
	total_g_float := float32(total_g)
	total_g_float = total_g_float / float32(meterpergateway)
	total_g = total_g / meterpergateway
	if total_g_float-float32(int32(total_g_float)) > 0 {
		total_g++
	}
	return total_g
}
func (db *databaseOps) GetTemplateById(ctx context.Context, id int) ([]string, error) {
	rows, err := db.db.QueryContext(ctx, fmt.Sprintf("SELECT template_key FROM templates where id = %d", id))
	if err != nil {
		fmt.Println("error during to fetching template key")
		return nil, err
	}
	result := []string{}
	for rows.Next() {
		templateKey := ""
		rows.Scan(&templateKey)
		result = append(result, templateKey)
	}
	return result, nil
}
func (db *databaseOps) FetchCountMeterReading(ctx context.Context, table string) (int, error) {
	query := fmt.Sprintf("select count(*) from %v", table)
	row, err := db.db.QueryContext(ctx, query)
	if err != nil {
		log.Println("things.go-TotalGateway", err)
		return 0, err
	}
	// fmt.Print("working\n")
	count := -1
	row.Next()
	err = row.Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}

func (db *databaseOps) FetchTemplateKey(ctx context.Context) ([]string, error) {
	rows, err := db.db.QueryContext(ctx, FetchTemplateKey)
	if err != nil {
		fmt.Println("error during to fetching template key")
		return nil, err
	}
	result := []string{}
	for rows.Next() {
		templateKey := ""
		rows.Scan(&templateKey)
		result = append(result, templateKey)
	}
	return result, nil
}

func (db *databaseOps) RemoveTemplateThingData(ctx context.Context, templateKey string) (int, error) {
	query := fmt.Sprintf(RemoveThingsData, templateKey)
	rows, err := db.db.ExecContext(ctx, query)
	if err != nil {
		fmt.Println("error during to fetching template key", query)
		return 0, err
	}

	val, err := rows.RowsAffected()
	if err != nil {
		fmt.Println("error during to fetching template key", err.Error())
		return 0, err
	}
	return int(val), nil
}
