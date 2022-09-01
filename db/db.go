package db

import (
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //import postgres

	"github.com/gcinnovate/go-fcapp/config"
)

var regionDistricts map[string]map[string]interface{}
var subcountyFacilities map[string]map[string]interface{}
var db *sqlx.DB

func init() {
	psqlInfo := fmt.Sprintf("%s", config.FcAppConf.Database.URI)

	var err error
	db, err = ConnectDB(psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	regionDistricts, err = LoadRegionDistricts(db)
	if err != nil {
		log.Fatal(err)
	}
	districtSubcounties, err = LoadDistrictSubcounties(db)
	if err != nil {
		log.Fatal(err)
	}

	subcountyFacilities, err = LoadSubcountyFacilities(db)
	if err != nil {
		log.Fatal(err)
	}

}

// ConnectDB ...
func ConnectDB(dataSourceName string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	return db, nil
}

//GetDB ...
func GetDB() *sqlx.DB {
	return db
}

// LoadRegionDistricts simply loads our regions and their underlyng districts
func LoadRegionDistricts(db *sqlx.DB) (map[string]map[string]interface{}, error) {
	regDistricts := make(map[string]map[string]interface{})
	rows, err := db.Queryx("SELECT name FROM fcapp_orgunits WHERE hierarchylevel = 2")
	if err != nil {
		log.Fatal("Failed to load regions")
	}

	defer rows.Close()
	for rows.Next() {
		var region string
		err = rows.Scan(&region)
		if err != nil {
			log.Fatal("Failed to load regions")
			return nil, err
		}
		records, err := db.NamedQuery(
			`
			SELECT 
				name 
			FROM fcapp_orgunits 
			WHERE 
				parentid = (select id FROM fcapp_orgunits WHERE hierarchylevel = 2 AND name = :region) 
			ORDER BY name
			`, map[string]interface{}{"region": region})
		if err != nil {
			log.Fatal("Failed to load districts in region: ", region)
			return nil, err
		}
		defer records.Close()

		payload := make(map[string]interface{})
		districts := make(map[string]string)
		screen1 := ""
		screen2 := ""
		var dlist []string
		idx := 0
		for records.Next() {
			idx++
			var district string
			err = records.Scan(&district)
			if err != nil {
				log.Fatal("Failed to read districts from region: ", region)
				return nil, err
			}
			dlist = append(dlist, district)
			switch {
			case idx <= 10:
				screen1 += fmt.Sprintf("%d. %s\n", idx, district)
				districts[fmt.Sprintf("%d", idx)] = district
			case idx > 10 && idx < 20:
				screen2 += fmt.Sprintf("%d. %s\n", idx+1, district)
				districts[fmt.Sprintf("%d", idx)] = district

			}
		}

		if len(dlist) > 10 {
			dlist = append(dlist, "")
			insert(dlist, "#", 10)
		}
		payload["district_list"] = strings.Join(dlist, ",")
		if len(screen2) > 0 {
			screen1 += "11. More\n"
			screen2 += "0. Back"
		} else {
			screen1 += "0. Back"
		}
		payload["districts"] = districts
		payload["screen_1"] = screen1
		payload["screen_2"] = screen2
		regDistricts[fmt.Sprintf("%s", region)] = payload

	}
	// log.Println(regDistricts)
	return regDistricts, nil
}

// GetRegionDistricts ....
func GetRegionDistricts() map[string]map[string]interface{} {
	return regionDistricts
}

var districtSubcounties map[string]map[string]interface{}

// LoadDistrictSubcounties ...
func LoadDistrictSubcounties(db *sqlx.DB) (map[string]map[string]interface{}, error) {
	districtSubs := make(map[string]map[string]interface{})

	rows, err := db.Queryx("SELECT id, name FROM fcapp_orgunits WHERE hierarchylevel = 3 ORDER BY name")
	if err != nil {
		log.Fatal("Failed to load districts")
	}
	defer rows.Close()
	for rows.Next() {
		var district string
		var districtID int64
		err = rows.Scan(&districtID, &district)
		if err != nil {
			log.Fatal("Failed to load district:", err)
			return nil, err
		}
		records, err := db.NamedQuery(
			`
			SELECT name FROM fcapp_orgunits WHERE parentid = :district_id	
			`, map[string]interface{}{"district_id": districtID})
		if err != nil {
			log.Fatal("Failed to load subcounties in district: ", district, " ", districtID, err)
			return nil, err
		}
		defer records.Close()

		payload := make(map[string]interface{})
		subcounties := make(map[string]string)
		screen1 := ""
		screen2 := ""
		screen3 := ""
		var slist []string
		idx := 0
		for records.Next() {
			idx++
			var subcounty string
			err = records.Scan(&subcounty)
			if err != nil {
				log.Fatal("Failed to read subcounties from district: ", district)
				return nil, err
			}
			if idx == 11 || idx == 21 {
				slist = append(slist, "#")
				slist = append(slist, subcounty)
			} else {
				slist = append(slist, subcounty)
			}
			switch {
			case idx < 11:
				screen1 += fmt.Sprintf("%d. %s\n", idx, subcounty)
				subcounties[fmt.Sprintf("%d", idx)] = subcounty
			case idx > 10 && idx < 20:
				screen2 += fmt.Sprintf("%d. %s\n", idx+1, subcounty)
				subcounties[fmt.Sprintf("%d", idx+1)] = subcounty
			case idx > 20 && idx < 31:
				screen3 += fmt.Sprintf("%d. %s\n", idx+2, subcounty)
				subcounties[fmt.Sprintf("%d", idx+2)] = subcounty

			}
		}

		if len(screen2) > 0 {
			screen1 += "11. More\n"
			screen2 += "0. Back\n"
		}
		if len(screen3) > 0 {
			screen2 += "21. More\n"
			screen3 += "0. Back\n"
		}
		if len(screen2) < 1 {
			screen1 += "0. Back\n"
		}

		payload["subcounty_list"] = strings.Join(slist, ",")
		payload["s_screen_1"] = screen1
		payload["s_screen_2"] = screen2
		payload["s_screen_3"] = screen3

		districtSubs[fmt.Sprintf("%s", district)] = payload

	}
	// fmt.Println(districtSubs)
	return districtSubs, nil
}

func LoadSubcountyFacilities(db *sqlx.DB) (map[string]map[string]interface{}, error) {
	subFacilities := make(map[string]map[string]interface{})

	rows, err := db.Queryx(`
		SELECT id, name FROM fcapp_orgunits WHERE hierarchylevel = 4 
		    AND parentid NOT IN (SELECT id FROM fcapp_orgunits WHERE hierarchylevel= 4 
		        AND name IN('Wakiso', 'Kampala')) ORDER BY name
		`)
	if err != nil {
		log.Fatal("Failed to load subcounties")
	}
	defer rows.Close()

	for rows.Next() {
		var subcounty string
		var subcountyID int64
		err = rows.Scan(&subcountyID, &subcounty)
		if err != nil {
			log.Fatal("Failed to load subcounty:", err)
			return nil, err
		}

		records, err := db.NamedQuery(
			`
			SELECT name FROM fcapp_orgunits WHERE parentid = :subcounty_id	
			`, map[string]interface{}{"subcounty_id": subcountyID})
		if err != nil {
			log.Fatal("Failed to load facilities in subcounty: ", subcounty, " ", subcountyID, err)
			return nil, err
		}
		defer records.Close()

		payload := make(map[string]interface{})
		facilities := make(map[string]string)
		screen1 := ""
		screen2 := ""
		screen3 := ""
		var flist []string
		idx := 0

		for records.Next() {
			idx++
			var facility string
			err = records.Scan(&facility)
			if err != nil {
				log.Fatal("Failed to read facilities from subcounty: ", subcounty)
				return nil, err
			}
			if idx == 11 || idx == 21 {
				flist = append(flist, "#")
				flist = append(flist, facility)
			} else {
				flist = append(flist, facility)
			}
			switch {
			case idx < 11:
				screen1 += fmt.Sprintf("%d. %s\n", idx, facility)
				facilities[fmt.Sprintf("%d", idx)] = facility
			case idx > 10 && idx < 20:
				screen2 += fmt.Sprintf("%d. %s\n", idx+1, facility)
				facilities[fmt.Sprintf("%d", idx+1)] = facility
			case idx > 20 && idx < 31:
				screen3 += fmt.Sprintf("%d. %s\n", idx+2, facility)
				facilities[fmt.Sprintf("%d", idx+2)] = facility

			}
		}
		if len(screen2) > 0 {
			screen1 += "11. More\n"
			screen2 += "0. Back\n"
		}
		if len(screen3) > 0 {
			screen2 += "21. More\n"
			screen3 += "0. Back\n"
		}
		if len(screen2) < 1 {
			screen1 += "0. Back\n"
		}

		payload["facility_list"] = strings.Join(flist, ",")
		payload["s_screen_1"] = screen1
		payload["s_screen_2"] = screen2
		payload["s_screen_3"] = screen3

		subFacilities[fmt.Sprintf("%s", subcounty)] = payload
	}

	return subFacilities, nil
}

// GetDistrictSubcounties ...
func GetDistrictSubcounties() map[string]map[string]interface{} {
	return districtSubcounties
}

// GetSubcountyFacilities ...
func GetSubcountyFacilities() map[string]map[string]interface{} {
	return subcountyFacilities
}
func insert(a []string, s string, i int) []string {
	return append(a[:i], append([]string{s}, a[i:]...)...)
}
