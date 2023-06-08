package main

import (
	"bufio"
	"fmt"
	"github.com/xuri/excelize/v2"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	// 读取excel生成对应的sql
	file, err := excelize.OpenFile("test/file/南坪C（护栏不拆分编码）20230606.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err = file.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	sheetName := "多类别明细表"
	rows, err := file.GetRows(sheetName)
	if err != nil {
		fmt.Println(err)
		return
	}

	tunnelMaps := map[string]any{
		"BYTD12": "26",
	}
	// 线路映射,线路名称+空间类型+线路类型作为key，id作为value
	line_space_id := map[string]string{
		"BHZ-18-15":   "79",
		"BHZ-116-15":  "80",
		"BHZ-159-15":  "81",
		"BHY-18-16":   "82",
		"BHY-116-16":  "83",
		"BHY-159-16":  "84",
		"BHA-18-17":   "85",
		"BHA-116-17":  "86",
		"BHA-159-17":  "87",
		"BHB-18-18":   "88",
		"BHB-116-18":  "89",
		"BHB-159-18":  "90",
		"BHC-18-19":   "91",
		"BHC-116-19":  "92",
		"BHC-159-19":  "93",
		"BHD-18-20":   "94",
		"BHD-116-20":  "95",
		"BHD-159-20":  "96",
		"NPL2-18-21":  "97",
		"NPL2-116-21": "98",
		"NPL2-159-21": "99",
		"NPR2-18-22":  "100",
		"NPR2-116-22": "101",
		"NPR2-159-22": "102",
		"NPL1-18-23":  "103",
		"NPL1-116-23": "104",
		"NPL1-159-23": "105",
		"NPR1-18-24":  "106",
		"NPR1-116-24": "107",
		"NPR1-159-24": "108",
		"NPA-18-25":   "109",
		"NPA-116-25":  "110",
		"NPA-159-25":  "111",
		"NPC-18-26":   "112",
		"NPC-116-26":  "113",
		"NPC-159-26":  "114",
		"NPD-18-27":   "115",
		"NPD-116-27":  "116",
		"NPD-159-27":  "117",
		"NPH-18-28":   "118",
		"NPH-116-28":  "119",
		"NPH-159-28":  "120",
		"NPK-18-29":   "121",
		"NPK-116-29":  "122",
		"NPK-159-29":  "123",
		"RXTQ-18-30":  "124",
		"RXTQ-116-30": "125",
		"RXTQ-159-30": "126",
		"CFDD-18-31":  "127",
		"CFDD-116-31": "128",
		"CFDD-159-31": "129",
		"BHI-18-32":   "130",
		"BHI-116-32":  "131",
		"BHI-159-32":  "132",
		"BHL1-18-33":  "133",
		"BHL1-116-33": "134",
		"BHL1-159-33": "135",
		"BHR2-18-34":  "136",
		"BHR2-116-34": "137",
		"BHR2-159-34": "138",
		"NBY-18-35":   "139",
		"NBY-116-35":  "140",
		"NBY-159-35":  "141",
		"NBZ-18-36":   "142",
		"NBZ-116-36":  "143",
		"NBZ-159-36":  "144",
		"NGD-18-37":   "145",
		"NGD-116-37":  "146",
		"NGD-159-37":  "147",
		"NGE-18-38":   "148",
		"NGE-116-38":  "149",
		"NGE-159-38":  "150",
		"NGF-18-39":   "151",
		"NGF-116-39":  "152",
		"NGF-159-39":  "153",
		"NGR1-18-40":  "154",
		"NGR1-116-40": "155",
		"NGR1-159-40": "156",
		"NGS1-18-41":  "157",
		"NGS1-116-41": "158",
		"NGS1-159-41": "159",
		"BHC-159-26":  "236",
		"BHC-18-26":   "237",
		"BHC-116-26":  "238",
	}

	component_id := map[string]any{
		"QM":   "1",
		"SB":   "2",
		"XB":   "3",
		"FS":   "4",
		"PZ":   "71",
		"SSF":  "73",
		"RXD":  "77",
		"ZL":   "79",
		"CZS":  "170",
		"NZS":  "173",
		"DZ":   "215",
		"JK":   "232",
		"QSX":  "233",
		"YSX":  "234",
		"LM":   "240",
		"JTZS": "244",
		"QBB":  "245",
		"JKJC": "246",
		"PLJC": "247",
		"JGYB": "248",
		"WSD":  "249",
		"ZLND": "250",
		"HXL":  "252",
		"GL":   "253",
		"JC":   "254",
		"TS":   "255",
		"ZZ":   "256",
		"NHL":  "262",
		"BX":   "283",
		"TM":   "284",
		"SPZ":  "285",
		"BZP":  "286",
		"LMJ":  "287",
	}
	object_id := map[string]string{
		"S": "34", // 设施
		"E": "35", // 设备
	}
	running_water_id := map[string]string{
		"1":  "21",
		"2":  "36",
		"3":  "37",
		"4":  "38",
		"5":  "39",
		"6":  "40",
		"7":  "41",
		"8":  "42",
		"9":  "43",
		"10": "44",
		"11": "66",
		"12": "67",
		"13": "68",
		"14": "69",
		"15": "142",
		"16": "143",
		"17": "144",
		"18": "145",
		"19": "146",
		"20": "147",
		"21": "148",
		"22": "149",
		"23": "150",
		"24": "151",
		"25": "152",
		"26": "153",
		"27": "154",
		"28": "155",
		"29": "156",
		"30": "157",
		"31": "158",
		"32": "163",
		"33": "164",
		"34": "165",
		"35": "166",
		"36": "167",
		"37": "168",
		"38": "169",
		"39": "170",
		"40": "171",
		"41": "172",
		"42": "173",
		"43": "174",
		"44": "175",
	}

	space_type_id := map[string]string{
		"SB": "18",
		"TD": "24",
		"FJ": "26",
		"QM": "116",
		"RH": "140",
		"CH": "141",
		"XB": "159",
		"LM": "160",
		"DQ": "161",
		"BP": "162",
	}

	// 线路类型
	line_type := map[string]string{
		"BYTD12": "26",
	}

	space_unit_id := map[string]string{
		"1":     "22",
		"2":     "27",
		"3":     "28",
		"4":     "29",
		"5":     "30",
		"6":     "31",
		"7":     "32",
		"8":     "33",
		"9":     "34",
		"10":    "35",
		"11":    "45",
		"12":    "46",
		"13":    "47",
		"14":    "48",
		"15":    "49",
		"16":    "50",
		"17":    "51",
		"18":    "52",
		"19":    "53",
		"20":    "54",
		"21":    "55",
		"22":    "56",
		"23":    "57",
		"24":    "58",
		"25":    "59",
		"26":    "60",
		"27":    "61",
		"28":    "62",
		"29":    "63",
		"30":    "64",
		"31":    "87",
		"32":    "70",
		"33":    "137",
		"34":    "71",
		"35":    "117",
		"36":    "72",
		"37":    "73",
		"38":    "118",
		"39":    "119",
		"40":    "89",
		"41":    "88",
		"42":    "122",
		"43":    "123",
		"44":    "124",
		"45":    "90",
		"46":    "126",
		"47":    "74",
		"48":    "127",
		"49":    "75",
		"50":    "76",
		"51":    "83",
		"52":    "77",
		"53":    "84",
		"54":    "78",
		"55":    "85",
		"56":    "79",
		"57":    "131",
		"58":    "80",
		"59":    "81",
		"60":    "86",
		"61":    "82",
		"62":    "133",
		"63":    "92",
		"64":    "91",
		"65":    "138",
		"66":    "139",
		"67":    "176",
		"68":    "177",
		"69":    "178",
		"70":    "179",
		"71":    "180",
		"A0":    "181",
		"A1":    "182",
		"A12":   "183",
		"A13":   "184",
		"A15":   "185",
		"A7":    "186",
		"P1":    "187",
		"P10":   "188",
		"P11":   "189",
		"P12":   "190",
		"P13":   "191",
		"P14":   "192",
		"P2":    "193",
		"P3":    "194",
		"P4":    "195",
		"P5":    "196",
		"P6":    "197",
		"P7":    "198",
		"P8":    "199",
		"P9":    "200",
		"0":     "201",
		"BA00":  "202",
		"LL04":  "203",
		"LL11":  "204",
		"R1-11": "205",
		"LL12":  "206",
		"RR12":  "207",
		"LL13":  "208",
		"RR13":  "209",
		"LL14":  "210",
		"RR14":  "211",
		"LL15":  "212",
		"RR15":  "213",
		"LL16":  "214",
		"RR16":  "215",
		"LL17":  "216",
		"RR17":  "217",
		"LL18":  "218",
		"RR18":  "219",
		"L1-04": "220",
		"LL19":  "221",
		"RR19":  "222",
		"R2-00": "225",
		"LL20":  "226",
		"RR20":  "227",
		"LL21":  "228",
		"RR21":  "229",
		"LL22":  "230",
		"RR22":  "231",
		"R2-03": "232",
		"L1-00": "233",
		"LL23":  "234",
		"RR23":  "235",
		"LL24":  "236",
		"RR24":  "237",
		"LL25":  "238",
		"RR25":  "239",
		"LL26":  "240",
		"RR26":  "241",
		"LL27":  "242",
		"RR27":  "243",
		"LL30":  "244",
		"LL31":  "245",
		"RR30":  "246",
		"LL02":  "247",
		"LL03":  "248",
		"LL05":  "249",
		"LL06":  "250",
		"LL07":  "251",
		"LL08":  "252",
		"LL09":  "253",
		"LL28":  "254",
		"LL29":  "255",
		"RR29":  "256",
		"RR28":  "257",
		"NE01":  "258",
		"NE02":  "259",
		"NE03":  "260",
		"NE04":  "261",
		"NE05":  "262",
		"NE06":  "263",
		"NE07":  "264",
		"NE00":  "265",
		"S1-00": "266",
		"S1-01": "267",
		"S1-02": "268",
		"S1-03": "269",
		"S1-04": "270",
		"NF01":  "271",
		"NF02":  "272",
		"NF03":  "273",
		"NF04":  "274",
		"NF05":  "275",
		"ND03":  "276",
		"ND05":  "277",
		"ND06":  "278",
		"ND07":  "279",
		"ND01":  "280",
		"ND02":  "281",
		"ND08":  "282",
		"ND09":  "283",
		"ND10":  "284",
		"ND11":  "285",
		"ND12":  "286",
		"ND13":  "287",
		"ND04":  "288",
		"R1-07": "289",
		"R1-08": "290",
		"R1-09": "291",
		"R1-01": "292",
		"R1-02": "293",
		"R1-03": "294",
		"R1-04": "295",
		"R1-05": "296",
		"R1-06": "297",
		"LL01":  "298",
		"L1-01": "299",
		"L1-02": "300",
		"R2-01": "301",
		"R2-02": "302",
		"L1-03": "303",
		"LL32":  "304",
		"LL33":  "305",
		"LL34":  "306",
		"RR32":  "307",
		"RR33":  "308",
		"RR34":  "309",
		"RR35":  "310",
		"BB02":  "311",
		"BA01":  "312",
		"BA02":  "313",
		"BB01":  "314",
		"BB03":  "315",
		"BB04":  "316",
		"BC00":  "317",
		"BC01":  "318",
		"BC02":  "319",
		"BA03":  "320",
		"LL36":  "321",
		"LL35":  "322",
		"RR36":  "323",
		"LL37":  "324",
		"RR37":  "325",
		"BD03":  "326",
		"BD04":  "327",
		"BD01":  "328",
		"BD02":  "329",
		"BD05":  "330",
		"BD06":  "331",
		"BD07":  "332",
		"BD08":  "333",
		"BD00":  "334",
		"BD09":  "335",
		"ND14":  "336",
		"ND00":  "337",
		"BC04":  "338",
		"BC05":  "339",
		"BI01":  "340",
		"BI02":  "341",
		"BC03":  "342",
		"BC06":  "343",
		"BC07":  "344",
		"BB00":  "345",
		"NE08":  "346",
		"BI00":  "347",
		"BI03":  "348",
		"R1-00": "349",
		"LL00":  "350",
		"LL10":  "351",
		"R1-10": "352",
		"NF00":  "353",
		"NE09":  "354",
		"RR11":  "355",
		"R2-04": "356",
		"A14":   "377",
	}
	// INSERT INTO `stec_bytd`.`bt_component_code_manager` (`record_create_date`,`record_update_date`,`code`,`bridge_tunnel_id`,`city_id`,`component_id`,
	//`end_number`,`line_space_id`,`object_id`,`running_water_id`,`space_type_id`,`space_unit_id`,`start_number`,`sub_component_id`,`device_id`,`psn`,`flag`,
	//`bit`,`alarm_status`,`model_id`,`model_ids`) VALUES (
	//'2023-05-26 09:53:36.067000','2023-05-26 09:53:36.067000','SZS-BYTD69-S-TJ-NZS-TD-BF-017-0002',83,
	//135,172,'LK1+800',6,34,36,24,51,'LK1+788',173,'745204',20230526095336189,0,0,0,'745204','745204,745205,745206,745207,745208,745209,745210,745211,745212,
	//753621,753622');

	psn := fmt.Sprintf("%s%s", time.Now().Format("20060102150405"), "001")
	fmt.Println(psn)
	f, err := os.Create("test/save/t.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	writer := bufio.NewWriter(f)
	defer writer.Flush()
	for i, row := range rows {
		if i > 0 {
			sql := "INSERT INTO `stec_bytd`.`bt_component_code_manager` (`record_create_date`,`record_update_date`,`code`,`bridge_tunnel_id`,`city_id`," +
				"`component_id`, `end_number`,`line_space_id`,`object_id`,`running_water_id`,`space_type_id`,`space_unit_id`,`start_number`,`sub_component_id`," +
				"`device_id`,`psn`,`flag`, `bit`,`alarm_status`,`model_id`,`model_ids`) VALUES ('2023-05-30 09:53:36.067000','2023-05-30 09:53:36.067000',"
			sql += "'" + row[22] + "',"
			sql += tunnelMaps[row[7]].(string) + ","
			sql += "135" + ","
			sql += component_id[row[11]].(string) + ",'" // 构建大类
			sql += row[21] + "',"
			sql += line_space_id[row[17]+"-"+space_type_id[row[15]]+"-"+line_type[row[7]]] + "," // 线路空间, key是空间名称+类型+线路，值是线路id
			fmt.Println(row[17] + "-" + space_type_id[row[15]] + "-" + line_type[row[7]])
			sql += object_id[row[9]] + ","               // 对象属性
			sql += running_water_id[row[19]] + ","       // 流水编号
			sql += space_type_id[row[15]] + ","          // 空间类型
			sql += space_unit_id[row[18]] + ",'"         // 空间单元
			sql += row[20] + "',"                        // 开始里程
			sql += component_id[row[13]].(string) + ",'" // 构建小类
			sql += row[23] + "',"                        // deviceId
			at, _ := strconv.Atoi(psn)
			sql += strconv.Itoa(at+i-1) + "," // psn
			sql += "0" + ","                  //
			sql += "0" + ","                  //
			sql += "0" + ",'"                 //
			sql += row[23] + "','"            // 模型id
			split := strings.Split(row[25], "/")
			models := strings.Join(split, ",")
			sql += models + "'" // 模型集合id
			sql += ");\n"
			fmt.Println(sql)
			_, err = writer.WriteString(sql)
			if err != nil {
				panic(err)
			}
		}
	}

}