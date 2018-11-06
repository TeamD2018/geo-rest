function create_resolver_cache()
    local s = box.schema.space.create('geo_cache', { if_not_exists = true, field_count = 2 })
    s:create_index('address', { type = 'HASH', unique = true, if_not_exists = true, parts = { 1, 'string' } })
    s:create_index('geo_point', { type = 'RTREE', unique = false, if_not_exists = true, parts = { 2, 'array' } })
    return
end

create_resolver_cache()

---save_to_cache
---@param address string
---@param point array
function save_to_cache(address, point)
    return box.space.geo_cache:replace { address, point }
end

---resolve
---@param address string
function revers_resolve(address)
    return box.space.geo_cache:get { address }
end

---resolve
---@param geo_point array
function resolve(geo_point)
    return box.space.geo_cache.index.geo_point:select(geo_point)
end

---clear_cache
function clear_cache()
    return box.space.geo_cache:truncate()
end
