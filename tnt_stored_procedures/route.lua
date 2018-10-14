function create_couriers_route_space()
    s = box.schema.space.create('couriers_route', { if_not_exists = true, field_count = 3 })
    s:create_index('order_id', { type = 'HASH', if_not_exists = true, parts = { 2, 'string' } })
    s:create_index('courier_id', { type = 'TREE', unique = false, if_not_exists = true, parts = { 1, 'string' } })
    s:create_index('courier_order', { type = 'HASH', if_not_exists = true, parts = { 1, 'string', 2, 'string' } })
end

function add_courier_with_order(courier_id, order_id)
    if type(courier_id) ~= 'string' or type(order_id) ~= 'string' then
        error('courier_id or order_id must be a string')
    end
    s:replace { courier_id, order_id, {} }
end

function add_point_to_route(courier_id, point)
    local res = box.space.couriers_route.index.courier_id:select { courier_id }
    for i, v in pairs(res) do
        local order_id = v[2]
        local r = v[3]
        if point.lat == nil or point.lon == nil then
            error('lat or lon not found')
        end
        if type(point.lat) ~= 'number' or type(point.lon) ~= 'number' then
            error('lat and lon must be a numbers')
        end
        if point.lat < -90 or point.lat > 90 or point.lon < -180 or point.lon > 180 then
            error('lat or lon have invalid format (-90 < lat < 90, -180 < lon < 180')
        end
        table.insert(r, { lat = point.lat, lon = point.lon })
        box.space.couriers_route.index.courier_order:update({ courier_id, order_id }, { { '=', 3, r } })
    end
end

function delete_courier(courier_id)
    local couriers = box.space.couriers_route.index.courier_id:select { courier_id }
    for i, v in pairs(couriers) do
        local res = box.space.couriers_route:delete { v[2] }
        if res == nil then
            return error("courier_id not found")
        end
    end
end

function get_route(courier_id, order_id)
    if type(courier_id) ~= 'string' or type(order_id) ~= 'string' then
        error('courier_id or order_id must be a string')
    end
    return box.space.couriers_route.index.courier_order:get { courier_id, order_id }[3]
end