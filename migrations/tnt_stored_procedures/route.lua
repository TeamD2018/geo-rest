function create_couriers_route_space()
    local s = box.schema.space.create('couriers_route', { if_not_exists = true, field_count = 2 })
    return s:create_index('courier_id', { type = 'HASH', unique = true, if_not_exists = true, parts = { 1, 'string' } })
end

create_couriers_route_space()

function add_courier(courier_id)
    local s = box.space.couriers_route
    if type(courier_id) ~= 'string' then
        error('courier_id must be a string')
    end
    s:replace { courier_id, {} }
end

function add_point_to_route(courier_id, point)
    local courier_id_idx = box.space.couriers_route.index.courier_id
    local res = box.space.couriers_route.index.courier_id:get { courier_id }
    if res == nil then
        error('courier with ' .. courier_id .. ' not found')
    end
    local route = res[2]
    if point.lat == nil or point.lon == nil then
        error('lat or lon not found')
    end
    if type(point.lat) ~= 'number' or type(point.lon) ~= 'number' then
        error('lat and lon must be a numbers')
    end
    if point.lat < -90 or point.lat > 90 or point.lon < -180 or point.lon > 180 then
        error('lat or lon have invalid format (-90 < lat < 90, -180 < lon < 180')
    end
    table.insert(route, { lat = point.lat, lon = point.lon, ts = point.ts })
    courier_id_idx:update({ courier_id }, { { '=', 2, route } })
end

function delete_courier(courier_id)
    box.space.couriers_route:delete { courier_id }
end

function get_route(courier_id, since)
    if type(courier_id) ~= 'string' then
        error('courier_id must be a string')
    end
    local temp_res = box.space.couriers_route.index.courier_id:get { courier_id }
    if temp_res == nil then
        error('courier with ' .. courier_id .. ' not found')
    end
    local res = {}
    for _, v in pairs(temp_res[2]) do
        if v.ts >= since then
            table.insert(res, v)
        end
    end
    return res
end