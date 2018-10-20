function create_resolver_cache()
    local s = box.schema.space.create('geo_cache', { if_not_exists = true, field_count = 2 })
    s:create_index('address', { type = 'HASH', unique = true, if_not_exists = true, parts = { 1, 'string' } })
end

create_resolver_cache()