function create_courier_orders_space()
    local s = box.schema.space.create('courier_orders', { if_not_exists = true, field_count = 2 })
    return s:create_index('courier_id', { type = 'HASH', unique = true, if_not_exists = true, parts = { 1, 'string' } })
end

create_courier_orders_space()

function inc_courier_orders_counter(courier_id)
    local s = box.space.courier_orders
    if s then
        return s:upsert({ courier_id, 1 }, { { '+', 2, 1 } })
    end
    return nil, error('space "courier_orders" not exist')
end

function dec_courier_orders_counter(courier_id)
    local s = box.space.courier_orders
    if s then
        return s:upsert({ courier_id, 0 }, { { '-', 2, 1 } })
    end
    return nil, error('space "courier_orders" not exist')
end

function get_or_create_counter(courier_id)
    local s = box.space.courier_orders
    if not s then
        return nil, error('space "courier_orders" not exist')
    end
    local counter = s:get(courier_id)
    if not counter then
        return s:insert { courier_id, 0 }
    end
    return counter
end

function drop_courier_orders_counter(courier_id)
    local s = box.space.courier_orders
    if not s then
        return nil, error('space "courier_orders" not exist')
    end
    return s:delete(courier_id)
end

function get_counters(ids)
    local counters = {}
    for _, id in ipairs(ids) do
        table.insert(counters, get_or_create_counter(id))
    end
    return counters
end
