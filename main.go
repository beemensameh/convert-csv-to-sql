package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gocarina/gocsv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

const (
	host     = ""
	user     = ""
	password = ""
	dbname   = ""

	csvFileName = ""
)

func getRequirementIdentitfersIndices(requirementIdentifers []string, startingIndex int) (requirementIdentifersStrings []interface{}, requirementIdentifierIndex []string) {
	for index, requirementIdentitfer := range requirementIdentifers {
		requirementIdentifersStrings = append(requirementIdentifersStrings, requirementIdentitfer)
		requirementIdentifierIndex = append(requirementIdentifierIndex, "$"+strconv.Itoa(startingIndex+index))
	}

	return requirementIdentifersStrings, requirementIdentifierIndex
}

func main() {
	var (
		usersUUID                  = []string{}
		updateQueryTemplate string = "UPDATE public.requirements SET value='{\"residence_country_code\":\"%v\",\"nationality_countries\":[\"%v\"]}'::json::json WHERE id=%v; -- user_uuid = '%v'\n"

		requirements []*Requirement
	)

	in, err := os.Open(csvFileName)
	if err != nil {
		panic(err)
	}
	defer in.Close()

	users := []*CSVSheet{}

	if err := gocsv.UnmarshalFile(in, &users); err != nil {
		panic(err)
	}

	for _, user := range users {
		usersUUID = append(usersUUID, user.UUID)
	}

	connStr := fmt.Sprintf("user=%s port=%d dbname=%s password=%s host=%s", user, 5432, dbname, password, host)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logrus.Errorln("Error in open connection =>", err)
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		logrus.Errorln("Error in ping =>", err)
		log.Fatal(err)
	}

	// Query logic
	requirementIdentitfersStrings, requirementIdentifiersIndices := getRequirementIdentitfersIndices(usersUUID, 1)
	userSql := fmt.Sprintf("select users.uuid as user_uuid, requirements.id, requirements.requirement_type, requirements.value from users join applications on (users.id = applications.user_id) join requirements on (applications.id= requirements.application_id) where users.uuid IN (%v) and requirement_type='nationality'", strings.Join(requirementIdentifiersIndices, ","))

	rows, err := db.Query(userSql, requirementIdentitfersStrings...)
	if err != nil {
		logrus.Errorln("Error in execut query =>", err)
		log.Fatal("Failed to execute query: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		requirement := &Requirement{}

		err = rows.Scan(
			&requirement.UserUUID,
			&requirement.ID,
			&requirement.Value,
			&requirement.RequirementType,
		)
		if err != nil {
			log.Fatal(err)
		}

		requirements = append(requirements, requirement)
	}
	err = rows.Err()
	if err != nil {
		logrus.Errorln("Error in scan query =>", err)
		log.Fatal("Failed to scan query: ", err)
	}

	f, err := os.Create("queries.sql")
	if err != nil {
		log.Fatal("Error in create sql file: ", err)
	}

	defer f.Close()

	for _, requirement := range requirements {
		for _, user := range users {
			if user.UUID == requirement.UserUUID {
				insertQuery := fmt.Sprintf(updateQueryTemplate, user.Country, user.Residence, requirement.ID, user.UUID)
				w := bufio.NewWriter(f)
				_, err = w.WriteString(insertQuery)
				if err != nil {
					log.Fatal("Error for writing a line in file: ", err)
				}

				w.Flush()
			}
		}
	}
}
