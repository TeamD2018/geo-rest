function create_couriers_route_history_space()
    s = box.schema.space.create('route_history', { if_not_exists = true, field_count = 3 })
    s:create_index('courier_id', { type = 'HASH', if_not_exists = true, parts = { 1, 'string' } })
    s:create_index('order_id', { type = 'HASH', if_not_exists = true, parts = { 2, 'string' } })
end