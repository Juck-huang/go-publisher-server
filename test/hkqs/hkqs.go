package main

import (
	"bufio"
	"fmt"
	"github.com/xuri/excelize/v2"
	"os"
	"time"
)

func main() {
	// 读取excel生成对应的sql
	file, err := excelize.OpenFile("test/hkqs/file/文明东隧道编码表20230606.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err = file.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	sheetName := "Sheet1"
	rows, err := file.GetRows(sheetName)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 隧道名称
	tunnelMaps := map[string]any{
		"WMDSD": "263",
	}
	// 线路映射,线路名称+空间类型+线路类型作为key，id作为value
	line_space_id := map[string]string{
		"NZXCK-67-263":       "385",
		"NXZRK-67-263":       "386",
		"BZXYSBF-67-263":     "387",
		"AYSBF-67-263":       "388",
		"BYSBF-67-263":       "389",
		"CYSBF-67-263":       "390",
		"DYSBF-67-263":       "391",
		"ZXFS-68-263":        "392",
		"CFSBF-68-263":       "393",
		"NZX-22-263":         "394",
		"BZX-22-263":         "395",
		"AZD-22-263":         "396",
		"BZD-22-263":         "397",
		"CZD-22-263":         "398",
		"DZD-22-263":         "399",
		"XSBFQPS-64-263":     "400",
		"XSBFQPJ-64-263":     "401",
		"XSBFFJ-64-263":      "402",
		"XSBFPYJF-64-263":    "403",
		"XSBFPFJF-64-263":    "404",
		"XSBFBYJ-64-263":     "405",
		"XSBF10KVPDJ-64-263": "406",
		"XSBF400VPDJ-64-263": "407",
		"XSBFXFBF-64-263":    "408",
		"XSBFXFSC-64-263":    "409",
		"XSBFKZS-64-263":     "410",
		"XSBFXFJF-64-263":    "411",
		"XSBFZBS-64-263":     "412",
		"XSBFYJZMDYS-64-263": "413",
		"XSBFZMPDJ-64-263":   "414",
		"XSBFRDJKS-64-263":   "415",
		"XSBFDLJ-64-263":     "416",
		"DSBFPYJF-23-263":    "417",
		"DSBFPFJF-23-263":    "418",
		"DSBFLTJ-23-263":     "419",
		"DSBFQS-23-263":      "420",
		"DSBFZMPDJ-23-263":   "421",
		"DSBFQPJ-23-263":     "422",
		"DSBFKZS-23-263":     "423",
		"DSBF10KVPDS-23-263": "424",
		"DSBF100VPDS-23-263": "425",
		"DSBFXFJ-23-263":     "426",
		"DSBFBFJF-23-263":    "427",
		"DSBFRDKZS-23-263":   "428",
		"DSBFYJZMDYS-23-263": "429",
		"DSBFZBS-23-263":     "430",
		"DSBFDLJ-23-263":     "431",
		"DSBFBYJ-23-263":     "432",
		"DSBFXFBF-23-263":    "433",
		"DSBFPFJ-23-263":     "434",
		"DSBFDLJC-23-263":    "435",
	}

	component_id := map[string]any{
		"JBZM":  "143",
		"YJZM":  "144",
		"JKSXT": "146",
		"LM":    "189",
		"CS":    "191",
		"ZSMB":  "192",
		"HJG":   "193",
		"BG":    "194",
		"JQZM":  "202",
		"YDZM":  "203",
		"FZQ":   "209",
		"TD":    "235",
		"DGY":   "236",
		"TCGB":  "352",
		"SXT":   "355",
		"SD":    "357",
		"TJ":    "6",
		"ZM":    "102",
		"JK":    "103",
	}
	object_id := map[string]string{
		"S": "114", // 设施
		"E": "115", // 设备
	}
	running_water_id := map[string]string{
		"0001": "10",
		"0002": "11",
	}

	space_type_id := map[string]string{
		"XCD": "22",
	}

	// 线路类型
	line_type := map[string]string{
		"WMDSD": "263",
	}

	space_unit_id := map[string]string{
		"001": "12",
		"002": "13",
		"003": "24",
		"004": "25",
		"005": "26",
		"006": "27",
		"007": "31",
		"008": "32",
		"009": "33",
		"010": "34",
		"011": "35",
		"012": "36",
		"013": "37",
		"014": "38",
		"015": "39",
		"016": "40",
		"017": "41",
		"018": "42",
		"019": "43",
		"020": "44",
		"021": "60",
		"022": "70",
		"023": "71",
		"024": "72",
		"025": "73",
		"026": "74",
		"027": "75",
		"028": "76",
		"029": "77",
		"030": "78",
		"031": "79",
		"032": "80",
		"033": "81",
		"034": "82",
		"035": "83",
		"036": "84",
		"037": "85",
		"038": "86",
		"039": "87",
		"040": "88",
		"041": "89",
		"042": "90",
		"043": "91",
		"044": "92",
		"045": "93",
		"046": "94",
		"047": "95",
		"048": "96",
		"049": "97",
		"050": "98",
		"051": "99",
		"052": "100",
		"053": "101",
	}
	// INSERT INTO `stec_bytd`.`bt_component_code_manager` (`record_create_date`,`record_update_date`,`code`,`bridge_tunnel_id`,`city_id`,`component_id`,
	//`end_number`,`line_space_id`,`object_id`,`running_water_id`,`space_type_id`,`space_unit_id`,`start_number`,`sub_component_id`,`device_id`,`psn`,`flag`,
	//`bit`,`alarm_status`,`model_id`,`model_ids`) VALUES (
	//'2023-05-26 09:53:36.067000','2023-05-26 09:53:36.067000','SZS-BYTD69-S-TJ-NZS-TD-BF-017-0002',83,
	//135,172,'LK1+800',6,34,36,24,51,'LK1+788',173,'745204',20230526095336189,0,0,0,'745204','745204,745205,745206,745207,745208,745209,745210,745211,745212,
	//753621,753622');

	psn := fmt.Sprintf("%s%s", time.Now().Format("20060102150405"), "001")
	fmt.Println(psn)
	f, err := os.Create("test/hkqs/save/t.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	writer := bufio.NewWriter(f)
	defer writer.Flush()
	for i, row := range rows {
		if i > 0 {
			sql := "INSERT INTO `stec_hkqs`.`bt_component_code_manager` (`record_create_date`,`record_update_date`,`code`,`bridge_tunnel_id`,`city_id`," +
				"`component_id`, `end_number`,`line_space_id`,`object_id`,`running_water_id`,`space_type_id`,`space_unit_id`,`start_number`,`sub_component_id`) VALUES ('2023-06-07 09:53:36.067000','2023-06-07 09:53:36.067000',"
			sql += "'" + row[18] + "',"              // code
			sql += tunnelMaps[row[3]].(string) + "," // 桥梁隧道id
			sql += "117" + ","
			sql += component_id[row[7]].(string) + ",'"                                          // 构建大类
			sql += row[16] + "',"                                                                // 结束里程
			sql += line_space_id[row[13]+"-"+space_type_id[row[11]]+"-"+line_type[row[3]]] + "," // 线路空间, key是线路名称+空间类型+线路类型（隧道编号），值是线路id
			sql += object_id[row[5]] + ","                                                       // 对象属性
			sql += running_water_id[row[17]] + ","                                               // 流水编号
			sql += space_type_id[row[11]] + ","                                                  // 空间类型
			sql += space_unit_id[row[14]] + ",'"                                                 // 空间单元
			sql += row[15] + "',"                                                                // 开始里程
			sql += component_id[row[9]].(string)                                                 // 构建小类
			sql += ");\n"
			fmt.Println(sql)
			_, err = writer.WriteString(sql)
			if err != nil {
				panic(err)
			}
		}
	}

}
