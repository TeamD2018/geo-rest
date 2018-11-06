function create_resolver_cache()
    local s = box.schema.space.create('geo_cache', { if_not_exists = true, field_count = 2 })
    return s:create_index('address', { type = 'HASH', unique = true, if_not_exists = true, parts = { 1, 'string' } })
end

create_resolver_cache()

---save_to_cache
---@param address string
---@param point table
function save_to_cache(address, point)
    return box.space.geo_cache:replace { address, point }
end

---resolve
---@param address string
function resolve(address)
    return box.space.geo_cache:get { address }
end
