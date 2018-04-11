local dbId = ARGV[1]
local boxes = { "ss7box", "smppbox", "httpbox", "router" }
local kpilistid = "kpis"
local cntpattern = ":cnt:"
local kpiValues = {}
local kpiKeys = {}
local nextKpiPos = 1
local selectId = redis.call('SELECT', dbId)
for pos, box in ipairs(boxes) do
    local kpilistname = kpilistid .. ":" .. box
    kpiKeys = redis.call('LRANGE', kpilistname, 0, -1)
    if #kpiKeys > 0 then
        for keypos = 1, #kpiKeys do
            kpiValues[nextKpiPos] = {}
            kpiValues[nextKpiPos][1] = kpiKeys[keypos]
            local cntpos = string.find(kpiKeys[keypos], cntpattern)
            if cntpos ~= nil then
                kpiValues[nextKpiPos][2] = redis.call('GETSET', kpiKeys[keypos], 0)
            else
                kpiValues[nextKpiPos][2] = redis.call('GET', kpiKeys[keypos])
            end
            nextKpiPos = nextKpiPos + 1
        end
    end
end
return kpiValues
