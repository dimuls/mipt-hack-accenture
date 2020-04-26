package entity

import "time"

type Plant struct {
	ID          string `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
}

type StockingPoint struct {
	ID   string `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type ResourceGroup struct {
	ID   string `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type Resource struct {
	ID              string `db:"id" json:"id"`
	ResourceGroupID string `db:"resource_group_id" json:"resource_group_id"`
	ShortName       string `db:"short_name" json:"short_name"`
	LongName        string `db:"long_name" json:"long_name"`
}

type Product struct {
	ID   string `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type ResourceGroupPeriod struct {
	ID                string        `db:"id" json:"id"`
	ResourceGroupID   string        `db:"resource_group_id" json:"resource_group_id"`
	AvailableCapacity time.Duration `db:"available_capacity" json:"available_capacity"`
	FreeCapacity      time.Duration `db:"free_capacity" json:"free_capacity"`
	StartDate         time.Time     `db:"start_date" json:"start_date"`
	HasFinateCapacity bool          `db:"has_finate_capacity" json:"has_finate_capacity"`
}

type Routing struct {
	ID                    string `db:"id" json:"id"`
	InputProductID        string `db:"input_product_id" json:"input_product_id"`
	OutputProductID       string `db:"output_product_id" json:"output_product_id"`
	InputStockingPointID  string `db:"input_stocking_point_id" json:"input_stocking_point_id"`
	OutputStockingPointID string `db:"output_stocking_point_id" json:"output_stocking_point_id"`
}

type RoutingStep struct {
	ID              string  `db:"id" json:"id"`
	PlantID         string  `db:"plant_id" json:"plant_id"`
	RoutingID       string  `db:"routing_id" json:"routing_id"`
	ResourceGroupID string  `db:"resource_group_id" json:"resource_group_id"`
	SequenceNumber  int     `db:"sequence_number" json:"sequence_number"`
	Yield           float64 `db:"yield" json:"yield"`
}

type Col struct {
	ID                                 string    `db:"id" json:"id"`
	RoutingID                          string    `db:"routing_id" json:"routing_id"`
	ProductID                          string    `db:"product_id" json:"product_id"`
	Quantity                           float64   `db:"quantity" json:"quantity"`
	MinQuantity                        float64   `db:"min_quantity" json:"min_quantity"`
	MaxQuantity                        float64   `db:"max_quantity" json:"max_quantity"`
	HasSalesBudgetReservation          bool      `db:"has_sales_budget_reservation" json:"has_sales_budget_reservation"`
	RequiresOrderCombination           bool      `db:"requires_order_combination" json:"requires_order_combination"`
	NumberOfActiveRoutingChainUpstream int       `db:"number_of_active_routing_chain_upstream" json:"number_of_active_routing_chain_upstream"`
	SelectedShippingShop               int       `db:"selected_shipping_shop" json:"selected_shipping_shop"`
	ResultProductType                  string    `db:"result_product_type" json:"result_product_type"`
	DeliveryType                       string    `db:"delivery_type" json:"delivery_type"`
	PlannedStatus                      string    `db:"planned_status" json:"planned_status"`
	Name                               string    `db:"name" json:"name"`
	ProductName                        string    `db:"product_name" json:"product_name"`
	LatestDesiredDeliveryDate          time.Time `db:"latest_desired_delivery_date" json:"latest_desired_delivery_date"`
	ProductSpecificationID             string    `db:"product_specification_id" json:"product_specification_id"`
	ResourceGroupIDs                   []string  `db:"resource_group_ids" json:"resource_group_ids"`
}

type SupplyOrder struct {
	ID              string    `db:"id" json:"id"`
	ColID           string    `db:"col_id" json:"col_id"`
	RoutingID       string    `db:"routing_id" json:"routing_id"`
	ProductID       string    `db:"product_id" json:"product_id"`
	StockingPointID string    `db:"stocking_point_id" json:"stocking_point_id"`
	OrderPosition   string    `db:"order_position" json:"order_position"`
	ProductName     string    `db:"product_name" json:"product_name"`
	ProductType     string    `db:"product_type" json:"product_type"`
	Quantity        float64   `db:"quantity" json:"quantity"`
	PlannedStatus   string    `db:"planned_status" json:"planned_status"`
	StartTime       time.Time `db:"start_time" json:"start_time"`
	EndTime         time.Time `db:"end_time" json:"end_time"`
	DeadlineTime    time.Time `db:"deadline_time" json:"deadline_time"`
	ProductFullID   string    `db:"product_full_id" json:"product_full_id"`
}

type SupplyOrderOperation struct {
	ID                       string        `db:"id" json:"id"`
	SupplyOrderID            string        `db:"supply_order_id" json:"supply_order_id"`
	ResourceGroupID          string        `db:"resource_group_id" json:"resource_group_id"`
	RoutingStepID            string        `db:"routing_step_id" json:"routing_step_id"`
	Description              string        `db:"description" json:"description"`
	SequenceNumber           int           `db:"sequence_number" json:"sequence_number"`
	AllowedStandardResources string        `db:"allowed_standard_resources" json:"allowed_standard_resources"`
	StartTime                time.Time     `db:"start_time" json:"start_time"`
	EndTime                  time.Time     `db:"end_time" json:"end_time"`
	ProductTime              time.Duration `db:"production_time" json:"production_time"`
	InputQuantity            float64       `db:"input_quantity" json:"input_quantity"`
	OutputQuantity           float64       `db:"output_quantity" json:"output_quantity"`
	SchedulingSpace          time.Duration `db:"scheduling_space" json:"scheduling_space"`
	OperationCode            int           `db:"operation_code" json:"operation_code"`
}
