package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"

	"github.com/sirupsen/logrus"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	var (
		dataPath string
		dbURI    string
		table    string
	)

	flag.StringVar(&dataPath, "d", "", "data path")
	flag.StringVar(&dbURI, "u", "", "db uri")
	flag.StringVar(&table, "t", "", "table to load")

	flag.Parse()

	db, err := sqlx.Connect("postgres", dbURI)
	if err != nil {
		logrus.WithError(err).Fatal("failed to connect to DB")
	}

	defer func() {
		err := db.Close()
		if err != nil {
			logrus.WithError(err).Fatal("failed to close DB")
		}
	}()

	switch table {
	case "":
		err = loadPlant(dataPath, db)
		if err != nil {
			break
		}

		logrus.Info("plant loaded")

		err = loadStockingPoint(dataPath, db)
		if err != nil {
			break
		}

		logrus.Info("stocking_point loaded")

		err = loadResourceGroup(dataPath, db)
		if err != nil {
			break
		}

		logrus.Info("resource_group loaded")

		err = loadResource(dataPath, db)
		if err != nil {
			break
		}

		logrus.Info("resource loaded")

		err = loadProduct(dataPath, db)
		if err != nil {
			break
		}

		logrus.Info("product loaded")

		err = loadResourceGroupPeriod(dataPath, db)
		if err != nil {
			break
		}

		logrus.Info("resource_group_period loaded")

		err = loadRouting(dataPath, db)
		if err != nil {
			break
		}

		logrus.Info("routing loaded")

		err = loadRoutingStep(dataPath, db)
		if err != nil {
			break
		}

		logrus.Info("routing_step loaded")

		err = loadCol(dataPath, db)
		if err != nil {
			break
		}

		logrus.Info("col loaded")

		err = loadSupplyOrder(dataPath, db)
		if err != nil {
			break
		}

		logrus.Info("supply_order loaded")

		err = loadSupplyOrderOperation(dataPath, db)
		if err != nil {
			break
		}

		logrus.Info("supply_order_operation loaded")

	case "plant":
		err = loadPlant(dataPath, db)
	case "stocking_point":
		err = loadStockingPoint(dataPath, db)
	case "resource_group":
		err = loadResourceGroup(dataPath, db)
	case "resource":
		err = loadResource(dataPath, db)
	case "product":
		err = loadProduct(dataPath, db)
	case "resource_group_period":
		err = loadResourceGroupPeriod(dataPath, db)
	case "routing":
		err = loadRouting(dataPath, db)
	case "routing_step":
		err = loadRoutingStep(dataPath, db)
	case "col":
		err = loadCol(dataPath, db)
	case "supply_order":
		err = loadSupplyOrder(dataPath, db)
	case "supply_order_operation":
		err = loadSupplyOrderOperation(dataPath, db)
	default:
		logrus.Fatal("unknown table")
	}

	if err != nil {
		logrus.WithError(err).Fatal("failed to load data")
	}
}

func loadPlant(dataPath string, db *sqlx.DB) error {
	f, err := os.Open(path.Join(dataPath, "plant.csv"))
	if err != nil {
		return fmt.Errorf("failed to open data file: %w", err)
	}

	defer func() {
		err := f.Close()
		if err != nil {
			logrus.WithError(err).Error("failed to close data file")
		}
	}()

	r := csv.NewReader(f)
	r.FieldsPerRecord = 3

	// skip header
	_, err = r.Read()
	if err != nil {
		return fmt.Errorf("failed to skip header: %w", err)
	}

	for {
		l, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to read line: %w", err)
		}

		_, err = db.Exec(`
			insert into plant (id, name, description) values ($1, $2, $3)
		`, l[0], l[1], l[2])
		if err != nil {
			return fmt.Errorf("failed to insert line to DB: %w", err)
		}
	}

	return nil
}

func loadStockingPoint(dataPath string, db *sqlx.DB) error {
	f, err := os.Open(path.Join(dataPath, "stocking-point.csv"))
	if err != nil {
		return fmt.Errorf("failed to open data file: %w", err)
	}

	defer func() {
		err := f.Close()
		if err != nil {
			logrus.WithError(err).Error("failed to close data file")
		}
	}()

	r := csv.NewReader(f)
	r.FieldsPerRecord = 2

	// skip header
	_, err = r.Read()
	if err != nil {
		return fmt.Errorf("failed to skip header: %w", err)
	}

	for {
		l, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to read line: %w", err)
		}

		_, err = db.Exec(`
			insert into stocking_point (id, name) values ($1, $2)
		`, l[0], l[1])
		if err != nil {
			return fmt.Errorf("failed to insert line to DB: %w", err)
		}
	}

	return nil
}

func loadResourceGroup(dataPath string, db *sqlx.DB) error {
	f, err := os.Open(path.Join(dataPath, "resource-group.csv"))
	if err != nil {
		return fmt.Errorf("failed to open data file: %w", err)
	}

	defer func() {
		err := f.Close()
		if err != nil {
			logrus.WithError(err).Error("failed to close data file")
		}
	}()

	r := csv.NewReader(f)
	r.FieldsPerRecord = 5

	// skip header
	_, err = r.Read()
	if err != nil {
		return fmt.Errorf("failed to skip header: %w", err)
	}

	resourceGroups := map[string]string{}

	for {
		l, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to read line: %w", err)
		}
		resourceGroups[l[0]] = l[1]
	}

	for id, name := range resourceGroups {
		_, err = db.Exec(`
			insert into resource_group (id, name) values ($1, $2)
		`, id, name)
		if err != nil {
			return fmt.Errorf("failed to insert line to DB: %w", err)
		}
	}

	return nil
}

func loadResource(dataPath string, db *sqlx.DB) error {
	f, err := os.Open(path.Join(dataPath, "resource-group.csv"))
	if err != nil {
		return fmt.Errorf("failed to open data file: %w", err)
	}

	defer func() {
		err := f.Close()
		if err != nil {
			logrus.WithError(err).Error("failed to close data file")
		}
	}()

	r := csv.NewReader(f)
	r.FieldsPerRecord = 5

	// skip header
	_, err = r.Read()
	if err != nil {
		return fmt.Errorf("failed to skip header: %w", err)
	}

	for {
		l, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to read line: %w", err)
		}

		_, err = db.Exec(`
			insert into resource (id, resource_group_id, short_name, long_name)
			values ($1, $2, $3, $4)
		`, l[2], l[0], l[3], l[4])
		if err != nil {
			return fmt.Errorf("failed to insert line to DB: %w", err)
		}
	}

	return nil
}

func loadProduct(dataPath string, db *sqlx.DB) error {
	f, err := os.Open(path.Join(dataPath, "product.csv"))
	if err != nil {
		return fmt.Errorf("failed to open data file: %w", err)
	}

	defer func() {
		err := f.Close()
		if err != nil {
			logrus.WithError(err).Error("failed to close data file")
		}
	}()

	r := csv.NewReader(f)
	r.FieldsPerRecord = 2

	// skip header
	_, err = r.Read()
	if err != nil {
		return fmt.Errorf("failed to skip header: %w", err)
	}

	for {
		l, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to read line: %w", err)
		}

		_, err = db.Exec(`
			insert into product (id, name) values ($1, $2)
		`, l[0], l[1])
		if err != nil {
			return fmt.Errorf("failed to insert line to DB: %w", err)
		}
	}

	return nil
}

var capacityRe = regexp.MustCompile(
	`^(?:(?P<days>\d+) days?|(?P<time>\d+:\d+:\d+)(?:\.(?P<ms>\d+))?|(?P<days2>\d+) days?, (?P<time2>\d+:\d+:\d+)(?:\.(?P<ms2>\d+))?|00:00.(?P<ms3>\d+)|(?P<part>0,\d+)||(?P<zero>0))$`)

func parseDuration(s string) (d time.Duration, err error) {
	matches := capacityRe.FindStringSubmatch(s)
	for i, k := range capacityRe.SubexpNames() {
		if i == 0 || k == "" {
			continue
		}
		if matches[i] == "" {
			continue
		}
		v := matches[i]
		switch k {
		case "days", "days2":
			var days int
			days, err = strconv.Atoi(v)
			if err != nil {
				err = fmt.Errorf("failed to parse days: %w", err)
				return
			}
			d += time.Duration(days) * 24 * time.Hour
		case "time", "time2":
			parts := strings.Split(v, ":")
			if len(parts) != 3 {
				err = fmt.Errorf("time doesn't have 3 parts")
				return
			}
			var hour, minute, second int
			hour, err = strconv.Atoi(parts[0])
			if err != nil {
				err = fmt.Errorf("failed to parse time hour: %w", err)
				return
			}
			minute, err = strconv.Atoi(parts[1])
			if err != nil {
				err = fmt.Errorf("failed to parse time minute: %w", err)
				return
			}
			second, err = strconv.Atoi(parts[2])
			if err != nil {
				err = fmt.Errorf("failed to parse time second: %w", err)
				return
			}
			d += time.Duration(hour)*time.Hour +
				time.Duration(minute)*time.Minute +
				time.Duration(second)*time.Second
		case "ms", "ms2", "ms3":
			var ms int
			ms, err = strconv.Atoi(v)
			if err != nil {
				err = fmt.Errorf("failed to parse ms: %w", err)
				return
			}
			d += time.Duration(ms) * time.Millisecond
		case "part":
			v = strings.Replace(v, ",", ".", 1)
			var part float64
			part, err = strconv.ParseFloat(v, 64)
			if err != nil {
				err = fmt.Errorf("failed to parse part: %w", err)
				return
			}
			d += time.Duration(math.Round(float64(24*time.Hour) * part))
			return
		case "zero":
			return
		}
	}
	return
}

func loadResourceGroupPeriod(dataPath string, db *sqlx.DB) error {
	f, err := os.Open(path.Join(dataPath, "resource-group-period.csv"))
	if err != nil {
		return fmt.Errorf("failed to open data file: %w", err)
	}

	defer func() {
		err := f.Close()
		if err != nil {
			logrus.WithError(err).Error("failed to close data file")
		}
	}()

	r := csv.NewReader(f)
	r.FieldsPerRecord = 7

	// skip header
	_, err = r.Read()
	if err != nil {
		return fmt.Errorf("failed to skip header: %w", err)
	}

	for {
		l, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to read line: %w", err)
		}

		id := l[2]
		resourceGroupID := l[1]

		availableCapacity, err := parseDuration(l[3])
		if err != nil {
			return fmt.Errorf("parse available capacity `%s`: %w", l[3], err)
		}

		freeCapacity, err := parseDuration(l[4])
		if err != nil {
			return fmt.Errorf("parse free capacity `%s`: %w", l[4], err)
		}

		startDate, err := time.Parse("2006-01-02 15:04:05", l[5])
		if err != nil {
			return fmt.Errorf("parse start date `%s`: %w", l[5], err)
		}

		hasFinateCapacity, err := strconv.ParseBool(l[6])
		if err != nil {
			return fmt.Errorf("parse has finate capacity `%s`: %w", l[6], err)
		}

		_, err = db.Exec(`
			insert into resource_group_period (
				id,
			    resource_group_id,
			    available_capacity,
			    free_capacity,
			    start_date,
			    has_finate_capacity
			) values ($1, $2, $3, $4, $5, $6)
		`, id, resourceGroupID, availableCapacity, freeCapacity, startDate,
			hasFinateCapacity)
		if err != nil {
			return fmt.Errorf("insert line to DB: %w", err)
		}
	}

	return nil
}

func loadRouting(dataPath string, db *sqlx.DB) error {
	f, err := os.Open(path.Join(dataPath, "routing.csv"))
	if err != nil {
		return fmt.Errorf("failed to open data file: %w", err)
	}

	defer func() {
		err := f.Close()
		if err != nil {
			logrus.WithError(err).Error("failed to close data file")
		}
	}()

	r := csv.NewReader(f)
	r.FieldsPerRecord = 6

	// skip header
	_, err = r.Read()
	if err != nil {
		return fmt.Errorf("failed to skip header: %w", err)
	}

	for {
		l, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to read line: %w", err)
		}

		_, err = db.Exec(`
			insert into routing (id, input_product_id, output_product_id, 
				input_stocking_point_id, output_stocking_point_id)
			values ($1, $2, $3, $4, $5)
		`, l[1], l[2], l[3], l[4], l[5])
		if err != nil {
			return fmt.Errorf("failed to insert line to DB: %w", err)
		}
	}

	return nil
}

func loadRoutingStep(dataPath string, db *sqlx.DB) error {
	f, err := os.Open(path.Join(dataPath, "routing-step.csv"))
	if err != nil {
		return fmt.Errorf("failed to open data file: %w", err)
	}

	defer func() {
		err := f.Close()
		if err != nil {
			logrus.WithError(err).Error("failed to close data file")
		}
	}()

	r := csv.NewReader(f)
	r.FieldsPerRecord = 7

	// skip header
	_, err = r.Read()
	if err != nil {
		return fmt.Errorf("failed to skip header: %w", err)
	}

	for {
		l, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to read line: %w", err)
		}

		yield, err := strconv.ParseFloat(
			strings.Replace(l[5], ",", ".", 1), 64)
		if err != nil {
			return fmt.Errorf("parse yield `%s`: %w", l[5], err)
		}

		_, err = db.Exec(`
			insert into routing_step (id, sequence_number, routing_id, 
			    resource_group_id, yield, plant_id)
			values ($1, $2, $3, $4, $5, $6)
		`, l[1], l[2], l[3], l[4], yield, l[6])
		if err != nil {
			return fmt.Errorf("failed to insert line to DB: %w", err)
		}
	}

	return nil
}

func loadCol(dataPath string, db *sqlx.DB) error {
	f, err := os.Open(path.Join(dataPath, "col.csv"))
	if err != nil {
		return fmt.Errorf("failed to open data file: %w", err)
	}

	defer func() {
		err := f.Close()
		if err != nil {
			logrus.WithError(err).Error("failed to close data file")
		}
	}()

	r := csv.NewReader(f)
	r.FieldsPerRecord = 19

	// skip header
	_, err = r.Read()
	if err != nil {
		return fmt.Errorf("failed to skip header: %w", err)
	}

	for {
		l, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to read line: %w", err)
		}

		id := l[1]

		quantity, err := strconv.ParseFloat(
			strings.Replace(l[2], ",", ".", 1), 64)
		if err != nil {
			return fmt.Errorf("parse quantity `%s`: %w", l[2], err)
		}

		minQuantity, err := strconv.ParseFloat(
			strings.Replace(l[3], ",", ".", 1), 64)
		if err != nil {
			return fmt.Errorf("parse min quantity `%s`: %w", l[3], err)
		}

		maxQuantity, err := strconv.ParseFloat(
			strings.Replace(l[4], ",", ".", 1), 64)
		if err != nil {
			return fmt.Errorf("parse max quantity `%s`: %w", l[4], err)
		}

		hasSalesBudgetReservation, err := strconv.ParseBool(l[5])
		if err != nil {
			return fmt.Errorf("parse has_sales_budget_reservation "+
				"`%s`: %w", l[5], err)
		}

		requiresOrderCombination, err := strconv.ParseBool(l[6])
		if err != nil {
			return fmt.Errorf("parse requires_order_combination "+
				"`%s`: %w", l[6], err)
		}

		numberOfActiveRoutingChainUpstream, err := strconv.Atoi(l[7])
		if err != nil {
			return fmt.Errorf("parse number_of_active_routing_chain_upstream "+
				"`%s`: %w", l[7], err)
		}

		selectedShippingShop, err := strconv.Atoi(l[8])
		if err != nil {
			return fmt.Errorf("parse selected_shipping_shop `%s`: %w",
				l[8], err)
		}

		resultProductType := l[9]
		deliveryType := l[10]
		plannedStatus := l[11]
		routingID := l[12]
		name := l[13]
		productID := l[14]
		productName := l[15]

		latestDesiredDeliveryDate, err := time.Parse("2-01-2006", l[16])
		if err != nil {
			return fmt.Errorf("parse latest_desired_delivery_date `%s`: %w",
				l[14], err)
		}

		productSpecificationID := l[17]
		resourceGroupIDs := strings.Split(l[18], ", ")

		_, err = db.Exec(`
			insert into col (
				id,
				routing_id,
			    product_id,
			    quantity,
				min_quantity,
				max_quantity,
				has_sales_budget_reservation,
				requires_order_combination,
				number_of_active_routing_chain_upstream,
				selected_shipping_shop,
				result_product_type,
				delivery_type,
				planned_status,
				name,
				product_name,
				latest_desired_delivery_date,
				product_specification_id,
				resource_group_ids
			) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13,
			          $14, $15, $16, $17, $18)
		`, id, routingID, productID, quantity, minQuantity, maxQuantity,
			hasSalesBudgetReservation, requiresOrderCombination,
			numberOfActiveRoutingChainUpstream, selectedShippingShop,
			resultProductType, deliveryType, plannedStatus, name,
			productName, latestDesiredDeliveryDate, productSpecificationID,
			pq.Array(resourceGroupIDs))
		if err != nil {
			return fmt.Errorf("failed to insert line to DB: %w", err)
		}
	}

	return nil
}

func loadSupplyOrder(dataPath string, db *sqlx.DB) error {
	f, err := os.Open(path.Join(dataPath, "supply-order.csv"))
	if err != nil {
		return fmt.Errorf("failed to open data file: %w", err)
	}

	defer func() {
		err := f.Close()
		if err != nil {
			logrus.WithError(err).Error("failed to close data file")
		}
	}()

	r := csv.NewReader(f)
	r.FieldsPerRecord = 15

	// skip header
	_, err = r.Read()
	if err != nil {
		return fmt.Errorf("failed to skip header: %w", err)
	}

	for {
		l, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to read line: %w", err)
		}

		quantity, err := strconv.ParseFloat(
			strings.Replace(l[6], ",", ".", 1), 64)
		if err != nil {
			return fmt.Errorf("parse quantity `%s`: %w", l[6], err)
		}

		startTime, err := time.Parse("2006-01-02 15:04:05", l[9])
		if err != nil {
			return fmt.Errorf("parse start_time `%s`: %w", l[9], err)
		}

		endTime, err := time.Parse("2006-01-02 15:04:05", l[10])
		if err != nil {
			return fmt.Errorf("parse end_time `%s`: %w", l[10], err)
		}

		deadlineTime, err := time.Parse("2006-01-02 15:04:05", l[11])
		if err != nil {
			return fmt.Errorf("parse deadline_time `%s`: %w", l[11], err)
		}

		_, err = db.Exec(`
			insert into supply_order (id, product_id, order_position,
				product_name, product_type, quantity, stocking_point_id,
			    planned_status, start_time, end_time, deadline_time,
			    product_full_id, routing_id, col_id)
			    values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12,
			            $13, $14)
		`, l[1], l[2], l[3], l[4], l[5], quantity, l[7], l[8],
			startTime, endTime, deadlineTime, l[12], l[13], l[14])
		if err != nil {
			return fmt.Errorf("failed to insert line to DB: %w", err)
		}
	}

	return nil
}

func loadSupplyOrderOperation(dataPath string, db *sqlx.DB) error {
	f, err := os.Open(path.Join(dataPath, "supply-order-operation.csv"))
	if err != nil {
		return fmt.Errorf("failed to open data file: %w", err)
	}

	defer func() {
		err := f.Close()
		if err != nil {
			logrus.WithError(err).Error("failed to close data file")
		}
	}()

	r := csv.NewReader(f)
	r.FieldsPerRecord = 14

	// skip header
	_, err = r.Read()
	if err != nil {
		return fmt.Errorf("failed to skip header: %w", err)
	}

	for {
		l, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to read line: %w", err)
		}

		sequenceNumber, err := strconv.Atoi(l[3])
		if err != nil {
			return fmt.Errorf("parse sequence_number `%s`: %w",
				l[3], err)
		}

		startTimeParts := strings.Split(l[5], ".")
		startTime, err := time.Parse("Jan-2-2006 15:04:05",
			startTimeParts[0])
		if err != nil {
			return fmt.Errorf("parse start_time `%s`: %w", l[5], err)
		}
		if len(startTimeParts) > 1 {
			ms, err := strconv.Atoi(startTimeParts[1])
			if err != nil {
				return fmt.Errorf("parse start_time ms part `%s`: %w",
					l[5], err)
			}
			startTime.Add(time.Duration(ms) * time.Millisecond)
		}

		endTime, err := time.Parse("2006-01-02 15:04:05", l[6])
		if err != nil {
			return fmt.Errorf("parse end_time `%s`: %w", l[6], err)
		}

		productionTime, err := parseDuration(l[7])
		if err != nil {
			return fmt.Errorf("parse production_time `%s`: %w", l[7], err)
		}

		inputQuantity, err := strconv.ParseFloat(strings.Replace(l[8], ",", ".", 1), 64)
		if err != nil {
			return fmt.Errorf("parse input_quantity `%s: %w", l[8], err)
		}

		outputQuantity, err := strconv.ParseFloat(strings.Replace(l[9], ",", ".", 1), 64)
		if err != nil {
			return fmt.Errorf("parse output_quantity `%s: %w", l[9], err)
		}

		schedulingSpace, err := parseDuration(l[10])
		if err != nil {
			return fmt.Errorf("parse scheduling_space `%s`: %w", l[10], err)
		}

		operationCode, err := strconv.Atoi(l[12])
		if err != nil {
			return fmt.Errorf("parse operation_code `%s`: %w", l[12], err)
		}

		_, err = db.Exec(`
			insert into supply_order_operation (id, description,
			    sequence_number, allowed_standard_resources,
			    start_time, end_time, production_time, input_quantity,
			    output_quantity, scheduling_space, resource_group_id,
				operation_code, routing_step_id)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		`, l[1], l[2], sequenceNumber, l[4], startTime, endTime, productionTime,
			inputQuantity, outputQuantity, schedulingSpace, l[11],
			operationCode, l[13])
		if err != nil {
			return fmt.Errorf("failed to insert line to DB: %w", err)
		}
	}

	return nil
}
