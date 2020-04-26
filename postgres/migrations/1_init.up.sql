create table plant (
    id text primary key,
    name text not null,
    description text not null
);

create table stocking_point (
    id text primary key,
    name text not null
);

create table resource_group (
    id text primary key,
    plaint_id text not null,
    name text not null
);

create table resource (
    id text primary key,
    resource_group_id text not null,
    short_name text not null,
    long_name text not null
);

create table product (
    id text primary key,
    name text not null
);

create table resource_group_period (
    id text primary key,
    resource_group_id text not null,
    available_capacity bigint not null, -- duration
    free_capacity bigint not null, -- duration
    start_date timestamp with time zone not null,
    has_finate_capacity bool not null
);

create table routing (
    id text,
    input_product_id text,
    output_product_id text not null,
    input_stocking_point_id text not null,
    output_stocking_point_id text not null
);

create table routing_step (
    id text not null,
    plant_id text not null,
    routing_id text not null,
    resource_group_id text not null,
    sequence_number bigint not null,
    yield double precision not null
);

create table col (
    id text not null,
    routing_id text not null,
    product_id text not null,
    quantity double precision not null,
    min_quantity double precision not null,
    max_quantity double precision not null,
    has_sales_budget_reservation bool not null,
    requires_order_combination bool not null,
    number_of_active_routing_chain_upstream bigint not null,
    selected_shipping_shop bigint not null,
    result_product_type text not null,
    delivery_type text not null,
    planned_status text not null,
    name text not null,
    product_name text not null,
    latest_desired_delivery_date timestamp with time zone not null,
    product_specification_id text not null,
    resource_group_ids text[] not null
);

create table supply_order (
    id text primary key,
    col_id text not null,
    routing_id text not null,
    product_id text not null,
    stocking_point_id text not null,
    order_position text not null,
    product_name text not null,
    product_type text not null,
    quantity double precision not null,
    planned_status text not null,
    start_time timestamp with time zone,
    end_time timestamp with time zone,
    deadline_time timestamp with time zone,
    product_full_id text not null
);

create table supply_order_operation (
    id text primary key,
    supply_order_id text not null,
    resource_group_id text not null,
    routing_step_id text not null,
    description text not null,
    sequence_number bigint not null,
    allowed_standard_resources text not null,
    start_time timestamp with time zone not null,
    end_time timestamp with time zone not null,
    production_time bigint not null, -- duration
    input_quantity double precision not null,
    output_quantity double precision not null,
    scheduling_space bigint not null, -- duration
    operation_code bigint not null
);