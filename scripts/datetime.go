package scripts

import (
	"strings"
	"time"
	"fmt"
	"html/template"

	"pgpartitiontable/models"

	"github.com/go-xweb/log"
)

// 按时间 进行分区
func DateTime(d int) {
	dateList := make([]string, 0)
	shortDateList := make([]string, 0)
	for i := 0; i < d; i++ {
		date := time.Now().AddDate(0, 0, i).Format("2006-01-02")
		dateList = append(dateList, date)
		shortDate := time.Now().AddDate(0, 0, i).Format("20060102")
		shortDateList = append(shortDateList, shortDate)
	}
	tableList := map[string]string{
		"test_datetime":      "create_time",
	}
	for tableName, checkItem := range tableList {
		data := map[string]string{
			"Dates":     strings.Join(dateList, ","),
			"Sdates":    strings.Join(shortDateList, ","),
			"TableName": tableName,
			"CheckItem": checkItem,
		}
		// 创建分区表
		crePartitionTable := `
	do language plpgsql $$ 
	declare
		dates varchar[] := '{ {{.Dates}} }';
		sdates varchar[] := array [{{.Sdates}}];
		i integer;
	begin
		for i in 1 .. array_upper(dates, 1) loop
			IF NOT EXISTS (select 1 from information_schema.tables where table_name='{{.TableName}}_' || sdates[i] ) THEN
				execute format('create table  {{.TableName}}_%s (like {{.TableName}} including all) inherits ({{.TableName}})', sdates[i]);
				execute format('alter table {{.TableName}}_%s add constraint ck check({{.CheckItem}}::date::text = ''%s'')', sdates[i], dates[i]);
			END IF;
		end loop;
	end;  
	$$;  
	`

		// 创建触发器函数
		creTableTriggerFuncStart := fmt.Sprintf(`
	create or replace function %s_insert_trigger() returns trigger as $$
	declare
        tablename varchar :='%s';
		checkitem varchar :='%s';
		check_date varchar;
        sdate varchar;
	begin`, tableName, tableName, checkItem)
		creTableTriggerFuncCase := fmt.Sprintf(`
	  check_date := NEW.%s::date::text;
      sdate := to_char(NEW.%s::date,'yyyymmdd');
      case check_date
			%s
	  end case;
	  return null;`, checkItem,checkItem, getDateTimeCaseStr(tableName,dateList, shortDateList))
		creTableTriggerFuncEnd := `
	end;
	$$ language plpgsql;
	`

		creTableTriggerFunc := creTableTriggerFuncStart + creTableTriggerFuncCase + creTableTriggerFuncEnd
		fmt.Println(creTableTriggerFunc)
		// 创建触发器
		creTableTrigger := fmt.Sprintf(`
	do language plpgsql $$ 
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'insert_%s_trigger') THEN
			create trigger insert_%s_trigger before insert on %s for each row execute procedure %s_insert_trigger();
		END IF;
	END
	$$;`, tableName, tableName, tableName, tableName)
		t := template.Must(template.New("test").Parse(crePartitionTable))
		builder := &strings.Builder{}
		if err := t.Execute(builder, data); err != nil {
			log.Error(fmt.Sprintf("PgPartitionTable %s template parse err:%v", tableName, err))
			return
		}
		crePartitionTable = builder.String()
		if err := models.CreatePartitionTable(crePartitionTable, creTableTriggerFunc, creTableTrigger); err != nil {
			log.Error(fmt.Sprintf("PgPartitionTable %s CreateHashTable err:%v", tableName, err))
			return
		}

	}

}

func getDateTimeCaseStr(tableName string,DateList, shortDateList []string) string {
	var ceseStr string
	for i, sdate := range shortDateList {
		ceseStr += fmt.Sprintf(`
		when '%s' then
           insert into %s_%s values (NEW.*);`, DateList[i], tableName, sdate)
	}
	// 如果不存在则创建分区表
	ceseStr += `
		else  
          IF NOT EXISTS (select 1 from information_schema.tables where table_name=tablename||'_'||sdate)
          THEN
				execute format('create table  %s_%s (like %s including all) inherits (%s)', tablename,sdate,tablename,tablename);
				execute format('alter table %s_%s add constraint ck check(%s::date::text = ''%s'')', tablename,sdate,checkitem, check_date);
		  END IF;
		  execute format('insert into %s_%s values ($1.*)', tablename,sdate) USING NEW;
	`
	return ceseStr
}
